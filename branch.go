package main

import (
	"fmt"

	git "github.com/libgit2/git2go/v30"
	"github.com/pkg/errors"
)

type BranchStatus struct {
	Name        string
	IsRemote    bool
	HasUpstream bool
	NeedsPull   bool
	NeedsPush   bool
}

func Branches(repo *git.Repository) ([]*git.Branch, error) {
	it, err := repo.NewBranchIterator(git.BranchAll)
	if err != nil {
		return nil, errors.Wrap(err, "Failed creating branch iterator")
	}

	it.ForEach(func(branch *git.Branch, branchType git.BranchType) error {
		fmt.Print(branch.IsRemote())
		upstream, err := branch.Upstream()
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println(upstream.Name())
		}
		return nil
	})

	return nil, nil
}
