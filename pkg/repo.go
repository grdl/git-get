package pkg

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	go_git_ssh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

type Repo struct {
	repo   *git.Repository
	path   string
	Status *RepoStatus
}

func CloneRepo(url *url.URL, reposRoot string, quiet bool) (*Repo, error) {
	repoPath := path.Join(reposRoot, URLToPath(url))

	var progress io.Writer
	if !quiet {
		progress = os.Stdout
		fmt.Printf("Cloning into '%s'...\n", repoPath)
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

	repo, err := git.PlainClone(repoPath, false, opts)
	if err != nil {
		return nil, errors.Wrap(err, "Failed cloning repo")
	}

	return newRepo(repo, repoPath), nil
}

func OpenRepo(repoPath string) (*Repo, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, errors.Wrap(err, "Failed opening repo")
	}

	return newRepo(repo, repoPath), nil
}

func newRepo(repo *git.Repository, repoPath string) *Repo {
	return &Repo{
		repo:   repo,
		path:   repoPath,
		Status: &RepoStatus{},
	}
}

// Fetch performs a git fetch on all remotes
func (r *Repo) Fetch() error {
	remotes, err := r.repo.Remotes()
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
	privateKey := viper.GetString(KeyPrivateKey)
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
