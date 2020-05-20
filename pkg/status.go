package pkg

import (
	git "github.com/libgit2/git2go/v30"
	"github.com/pkg/errors"
)

type RepoStatus struct {
	HasUntrackedFiles     bool
	HasUncommittedChanges bool
	Branches              map[string]BranchStatus
}

type BranchStatus struct {
	Name        string
	IsRemote    bool
	HasUpstream bool
	NeedsPull   bool
	NeedsPush   bool
	Ahead       int
	Behind      int
}

func loadStatus(r *git.Repository) (*RepoStatus, error) {
	entries, err := statusEntries(r)
	if err != nil {
		return nil, err
	}

	branches, err := branches(r)
	if err != nil {
		return nil, err
	}

	status := &RepoStatus{
		Branches: branches,
	}

	for _, entry := range entries {
		switch entry.Status {
		case git.StatusWtNew:
			status.HasUntrackedFiles = true
		case git.StatusIndexNew:
			status.HasUncommittedChanges = true
		}
	}

	return status, nil
}

func statusEntries(r *git.Repository) ([]git.StatusEntry, error) {
	opts := &git.StatusOptions{
		Show:  git.StatusShowIndexAndWorkdir,
		Flags: git.StatusOptIncludeUntracked,
	}

	status, err := r.StatusList(opts)
	if err != nil {
		return nil, errors.Wrap(err, "Failed getting repository status list")
	}

	entryCount, err := status.EntryCount()
	if err != nil {
		return nil, errors.Wrap(err, "Failed getting repository status list count")
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

func branches(r *git.Repository) (map[string]BranchStatus, error) {
	iter, err := r.NewBranchIterator(git.BranchAll)
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

	statuses := make(map[string]BranchStatus)
	for _, branch := range branches {
		status, err := branchStatus(branch)
		if err != nil {
			// TODO: Handle error. We should tell user that we couldn't read status of that branch but probably shouldn't exit
			continue
		}
		statuses[status.Name] = status
	}

	return statuses, nil
}

func branchStatus(branch *git.Branch) (BranchStatus, error) {
	var status BranchStatus

	name, err := branch.Name()
	if err != nil {
		return status, errors.Wrap(err, "Failed getting branch name")
	}
	status.Name = name

	// If branch is a remote one, return immediately. Upstream can only be found for local branches.
	if branch.IsRemote() {
		status.IsRemote = true
		return status, nil
	}

	upstream, err := branch.Upstream()
	if err != nil && !git.IsErrorCode(err, git.ErrNotFound) {
		return status, errors.Wrap(err, "Failed getting branch upstream")
	}

	// If there's no upstream, return immediately. Ahead/Behind can only be found when upstream exists.
	if upstream == nil {
		return status, nil
	}

	status.HasUpstream = true

	ahead, behind, err := branch.Owner().AheadBehind(branch.Target(), upstream.Target())
	if err != nil {
		return status, errors.Wrap(err, "Failed getting ahead/behind information")
	}

	status.Ahead = ahead
	status.Behind = behind

	if ahead > 0 {
		status.NeedsPush = true
	}

	if behind > 0 {
		status.NeedsPull = true
	}

	return status, nil
}
