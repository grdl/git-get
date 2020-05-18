package main

import (
	git "github.com/libgit2/git2go/v30"
	"github.com/pkg/errors"
)

type BranchStatus struct {
	Name        string
	IsRemote    bool
	HasUpstream bool
	NeedsPull   bool
	NeedsPush   bool
	Ahead       int
	Behind      int
}

func Branches(repo *git.Repository) ([]BranchStatus, error) {
	iter, err := repo.NewBranchIterator(git.BranchAll)
	if err != nil {
		return nil, errors.Wrap(err, "Failed creating branch iterator")
	}

	var branches []*git.Branch
	err = iter.ForEach(func(branch *git.Branch, branchType git.BranchType) error {
		branches = append(branches, branch)
		return nil
	})

	if err != nil {
		return nil, errors.Wrap(err, "Failed iterating over branches")
	}

	var statuses []BranchStatus
	for _, branch := range branches {
		status, err := NewBranchStatus(branch)
		if err != nil {
			// TODO: handle error
			continue
		}
		statuses = append(statuses, status)
	}

	return statuses, nil
}

func NewBranchStatus(branch *git.Branch) (BranchStatus, error) {
	var status BranchStatus

	name, err := branch.Name()
	if err != nil {
		return status, errors.Wrap(err, "Failed getting branch name")
	}
	status.Name = name

	status.IsRemote = branch.IsRemote()

	_, err = branch.Upstream()
	if err != nil {
		if git.IsErrorCode(err, git.ErrNotFound) {
			status.HasUpstream = false
		} else {
			return status, errors.Wrap(err, "Failed getting branch upstream")
		}
	}

	return status, nil
}
