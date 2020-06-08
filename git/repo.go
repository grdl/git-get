package git

import (
	"fmt"
	"git-get/cfg"

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

func CloneRepo(url *url.URL, path string, quiet bool) (*Repo, error) {
	var progress io.Writer
	if !quiet {
		progress = os.Stdout
		fmt.Printf("Cloning into '%s'...\n", path)
	}

	// TODO: can this be cleaner?
	var auth transport.AuthMethod
	var err error
	if url.Scheme == "ssh" {
		if auth, err = sshKeyAuth(); err != nil {
			return nil, err
		}
	}

	opts := &git.CloneOptions{
		URL:               url.String(),
		Auth:              auth,
		RemoteName:        git.DefaultRemoteName,
		ReferenceName:     "",
		SingleBranch:      false,
		NoCheckout:        false,
		Depth:             0,
		RecurseSubmodules: git.NoRecurseSubmodules,
		Progress:          progress,
		Tags:              git.AllTags,
	}

	repo, err := git.PlainClone(path, false, opts)
	if err != nil {
		return nil, errors.Wrap(err, "Failed cloning repo")
	}

	return NewRepo(repo, path), nil
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
