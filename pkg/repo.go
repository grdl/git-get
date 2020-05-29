package pkg

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path"

	"github.com/pkg/errors"

	"github.com/go-git/go-git/v5"
)

type Repo struct {
	repo   *git.Repository
	path   string
	Status *RepoStatus
}

func CloneRepo(url *url.URL, reposRoot string, quiet bool) (*Repo, error) {
	repoSubPath := URLToPath(url)
	repoPath := path.Join(reposRoot, repoSubPath)

	var output io.Writer
	if !quiet {
		output = os.Stdout
		fmt.Printf("Cloning into '%s'...\n", repoPath)
	}

	opts := &git.CloneOptions{
		URL:               url.String(),
		Auth:              nil,
		RemoteName:        git.DefaultRemoteName,
		ReferenceName:     "",
		SingleBranch:      false,
		NoCheckout:        false,
		Depth:             0,
		RecurseSubmodules: git.NoRecurseSubmodules,
		Progress:          output,
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
