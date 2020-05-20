package pkg

import (
	"github.com/pkg/errors"

	git "github.com/libgit2/git2go/v30"
)

func CloneRepo(url string, path string) (*git.Repository, error) {
	options := &git.CloneOptions{
		CheckoutOpts:         nil,
		FetchOptions:         nil,
		Bare:                 false,
		CheckoutBranch:       "",
		RemoteCreateCallback: nil,
	}

	repo, err := git.Clone(url, path, options)
	return repo, errors.Wrap(err, "Failed cloning repo")
}

func Fetch(repo *git.Repository) error {
	remotes, err := repo.Remotes.List()
	if err != nil {
		return errors.Wrap(err, "Failed listing remotes")
	}

	for _, r := range remotes {
		remote, err := repo.Remotes.Lookup(r)
		if err != nil {
			return errors.Wrap(err, "Failed looking up remote")
		}

		err = remote.Fetch(nil, nil, "")
		if err != nil {
			return errors.Wrap(err, "Failed fetching remote")
		}
	}

	return nil
}
