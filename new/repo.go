package new

import (
	"io"
	"net/url"
	"os"

	"github.com/pkg/errors"

	"github.com/go-git/go-git/v5"
)

type Repo struct {
	repo   *git.Repository
	Status *RepoStatus
}

func CloneRepo(url *url.URL, path string, quiet bool) (r *Repo, err error) {
	var output io.Writer
	if !quiet {
		output = os.Stdout
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

	repo, err := git.PlainClone(path, false, opts)
	if err != nil {
		return nil, errors.Wrap(err, "Failed cloning repo")
	}

	return newRepo(repo), nil
}

func OpenRepo(path string) (r *Repo, err error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, errors.Wrap(err, "Failed opening repo")
	}

	return newRepo(repo), nil
}

func newRepo(repo *git.Repository) *Repo {
	return &Repo{
		repo:   repo,
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
