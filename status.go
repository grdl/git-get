package main

import (
	git "github.com/libgit2/git2go/v30"
	"github.com/pkg/errors"
)

type RepoStatus int

const (
	StatusOk RepoStatus = iota
	StatusUntrackedFiles
	StatusUncommittedChanges
	StatusUnknown
)

func GetStatus(path string) ([]RepoStatus, error) {
	repo, err := git.OpenRepository(path)
	if err != nil {
		return nil, errors.Wrap(err, "Failed opening repository")
	}

	entries, err := statusEntries(repo)
	if err != nil {
		return nil, errors.Wrap(err, "Failed opening repository")
	}

	statusSet := make(map[RepoStatus]bool)

	statusSet[StatusOk] = true
	for _, entry := range entries {
		switch entry.Status {
		case git.StatusWtNew:
			statusSet[StatusUntrackedFiles] = true
			statusSet[StatusOk] = false
		case git.StatusIndexNew:
			statusSet[StatusUncommittedChanges] = true
			statusSet[StatusOk] = false
		default:
			statusSet[StatusUnknown] = true
			statusSet[StatusOk] = false
		}
	}

	var statusSlice []RepoStatus
	for k, v := range statusSet {
		if v {
			statusSlice = append(statusSlice, k)
		}
	}
	return statusSlice, nil
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
