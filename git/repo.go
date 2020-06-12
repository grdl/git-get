package git

import (
	"fmt"
	"git-get/cfg"

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
	IgnoreExisting bool // TODO: implement!
}

func CloneRepo(opts *CloneOpts) (*Repo, error) {
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
		return nil, errors.Wrap(err, "Failed cloning repo")
	}

	return NewRepo(repo, opts.Path), nil
}

func OpenRepo(repoPath string) (*Repo, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, errors.Wrap(err, "Failed opening repo")
	}

	return NewRepo(repo, repoPath), nil
}

func NewRepo(repo *git.Repository, repoPath string) *Repo {
	return &Repo{
		Repository: repo,
		Path:       repoPath,
		Status:     &RepoStatus{},
	}
}

// Fetch performs a git fetch on all remotes
func (r *Repo) Fetch() error {
	remotes, err := r.Remotes()
	if err != nil {
		return errors.Wrap(err, "Failed getting remotes")
	}

	for _, remote := range remotes {
		err = remote.Fetch(&git.FetchOptions{})
		if err != nil {
			return errors.Wrapf(err, "Failed fetching remote %s", remote.Config().Name)
		}
	}

	return nil
}

func sshKeyAuth() (transport.AuthMethod, error) {
	privateKey := viper.GetString(cfg.KeyPrivateKey)
	sshKey, err := ioutil.ReadFile(privateKey)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to open ssh private key %s", privateKey)
	}

	signer, err := ssh.ParsePrivateKey([]byte(sshKey))
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to parse ssh private key %s", privateKey)
	}

	// TODO: can it ba a different user
	auth := &go_git_ssh.PublicKeys{User: "git", Signer: signer}
	return auth, nil
}

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
