package repo

import (
	"fmt"
	"git-get/pkg/cfg"

	"github.com/go-git/go-git/v5/plumbing"

	"io"
	"io/ioutil"
	"net/url"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	go_git_ssh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

type Repo struct {
	*git.Repository
	Path   string
	Status *RepoStatus
}

// CloneOpts specify details about repository to clone.
type CloneOpts struct {
	URL            *url.URL
	Path           string // TODO: should Path be a part of clone opts?
	Branch         string
	Quiet          bool
	IgnoreExisting bool
}

// Clone clones repository specified in CloneOpts.
func Clone(opts *CloneOpts) (*Repo, error) {
	var progress io.Writer
	if !opts.Quiet {
		progress = os.Stdout
		fmt.Printf("Cloning into '%s'...\n", opts.Path)
	}

	// TODO: can this be cleaner?
	var auth transport.AuthMethod
	var err error
	if opts.URL.Scheme == "ssh" {
		if auth, err = sshKeyAuth(); err != nil {
			return nil, err
		}
	}

	if opts.Branch == "" {
		opts.Branch = cfg.DefBranch
	}

	// If branch name is actually a tag (ie. is prefixed with refs/tags) - check out that tag.
	// Otherwise, assume it's a branch name and check it out.
	refName := plumbing.ReferenceName(opts.Branch)
	if !refName.IsTag() {
		refName = plumbing.NewBranchReferenceName(opts.Branch)
	}

	gitOpts := &git.CloneOptions{
		URL:               opts.URL.String(),
		Auth:              auth,
		RemoteName:        git.DefaultRemoteName,
		ReferenceName:     refName,
		SingleBranch:      false,
		NoCheckout:        false,
		Depth:             0,
		RecurseSubmodules: git.NoRecurseSubmodules,
		Progress:          progress,
		Tags:              git.AllTags,
	}

	repo, err := git.PlainClone(opts.Path, false, gitOpts)
	if err != nil {

		if opts.IgnoreExisting && errors.Is(err, git.ErrRepositoryAlreadyExists) {
			return nil, nil
		}

		return nil, errors.Wrapf(err, "failed cloning %s", opts.URL.String())
	}

	return New(repo, opts.Path), nil
}

// Open opens a repository on a given path.
func Open(path string) (*Repo, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed opening repo %s", path)
	}

	return New(repo, path), nil
}

// New returns a new Repo instance from a given go-git Repository.
func New(repo *git.Repository, path string) *Repo {
	return &Repo{
		Repository: repo,
		Path:       path,
		Status:     &RepoStatus{},
	}
}

// Fetch performs a git fetch on all remotes
func (r *Repo) Fetch() error {
	remotes, err := r.Remotes()
	if err != nil {
		return errors.Wrapf(err, "failed getting remotes of repo %s", r.Path)
	}

	for _, remote := range remotes {
		err = remote.Fetch(&git.FetchOptions{})
		if err != nil {
			return errors.Wrapf(err, "failed fetching remote %s", remote.Config().Name)
		}
	}

	return nil
}

func sshKeyAuth() (transport.AuthMethod, error) {
	privateKey := viper.GetString(cfg.KeyPrivateKey)
	sshKey, err := ioutil.ReadFile(privateKey)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open ssh private key %s", privateKey)
	}

	signer, err := ssh.ParsePrivateKey([]byte(sshKey))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse ssh private key %s", privateKey)
	}

	// TODO: can it ba a different user
	auth := &go_git_ssh.PublicKeys{User: "git", Signer: signer}
	return auth, nil
}

// CurrentBranchStatus returns the BranchStatus of a currently checked out branch.
func (r *Repo) CurrentBranchStatus() *BranchStatus {
	if r.Status.CurrentBranch == StatusDetached || r.Status.CurrentBranch == StatusUnknown {
		return nil
	}

	for _, b := range r.Status.Branches {
		if b.Name == r.Status.CurrentBranch {
			return b
		}
	}

	return nil
}
