package new

import (
	"github.com/go-git/go-git/v5"
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

func (r *Repo) LoadStatus() error {
	wt, err := r.repo.Worktree()
	if err != nil {
		return errors.Wrap(err, "Failed getting worktree")
	}

	status, err := wt.Status()
	if err != nil {
		return errors.Wrap(err, "Failed getting worktree status")
	}

	r.Status.HasUncommittedChanges = !status.IsClean()
	r.Status.HasUntrackedFiles = hasUntracked(status)
	return nil
}

// hasUntracked returns true if there's any untracked file in the worktree
func hasUntracked(status git.Status) bool {
	for _, fs := range status {
		if fs.Worktree == git.Untracked {
			return true
		}
	}
	return false
}
