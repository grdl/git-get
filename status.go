package main

import (
	git "github.com/libgit2/git2go/v30"
	"github.com/pkg/errors"
)

type RepoStatus struct {
	HasUntrackedFiles     bool
	HasUncommittedChanges bool
	BranchStatuses        map[string]BranchStatus
}

func NewRepoStatus(path string) (RepoStatus, error) {
	var status RepoStatus

	repo, err := git.OpenRepository(path)
	if err != nil {
		return status, errors.Wrap(err, "Failed opening repository")
	}

	entries, err := statusEntries(repo)
	if err != nil {
		return status, errors.Wrap(err, "Failed getting repository status")
	}

	for _, entry := range entries {
		switch entry.Status {
		case git.StatusWtNew:
			status.HasUntrackedFiles = true
		case git.StatusIndexNew:
			status.HasUncommittedChanges = true
		}
	}

	branchStatuses, err := Branches(repo)
	if err != nil {
		return status, errors.Wrap(err, "Failed getting branches statuses")
	}
	status.BranchStatuses = branchStatuses

	return status, nil
}

func statusEntries(repo *git.Repository) ([]git.StatusEntry, error) {
	opts := &git.StatusOptions{
		Show:  git.StatusShowIndexAndWorkdir,
		Flags: git.StatusOptIncludeUntracked,
	}

	status, err := repo.StatusList(opts)
	if err != nil {
		return nil, errors.Wrap(err, "Failed getting repository status")
	}

	entryCount, err := status.EntryCount()
	if err != nil {
		return nil, errors.Wrap(err, "Failed getting repository status count")
	}

	var entries []git.StatusEntry
	for i := 0; i < entryCount; i++ {
		entry, err := status.ByIndex(i)
		if err != nil {
			return nil, errors.Wrap(err, "Failed getting repository status entry")
		}

		entries = append(entries, entry)
	}

	return entries, nil
}
