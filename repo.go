package main

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
