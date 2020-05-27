package new

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/pkg/errors"
)

type RepoStatus struct {
	HasUntrackedFiles     bool
	HasUncommittedChanges bool
	Branches              map[string]*BranchStatus
}

type BranchStatus struct {
	Name      string
	Upstream  *Upstream
	NeedsPull bool
	NeedsPush bool
	Ahead     int
	Behind    int
}

type Upstream struct {
	Remote string
	Branch string
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

	err = r.LoadBranchesStatus()
	if err != nil {
		return err
	}

	return nil
}

// hasUntracked returns true if there are any untracked files in the worktree
func hasUntracked(status git.Status) bool {
	for _, fs := range status {
		if fs.Worktree == git.Untracked {
			return true
		}
	}
	return false

}

func (r *Repo) LoadBranchesStatus() error {
	iter, err := r.repo.Branches()
	if err != nil {
		return errors.Wrap(err, "Failed getting branches iterator")
	}

	err = iter.ForEach(func(reference *plumbing.Reference) error {
		bs, err := r.newBranchStatus(reference.Name().Short())
		if err != nil {
			return err
		}

		r.Status.Branches[bs.Name] = bs
		return nil
	})
	if err != nil {
		return errors.Wrap(err, "Failed iterating over branches")
	}

	return nil
}

func (r *Repo) newBranchStatus(branch string) (*BranchStatus, error) {
	upstream, err := r.upstream(branch)
	if err != nil {
		return nil, err
	}

	return &BranchStatus{
		Name:     branch,
		Upstream: upstream,
	}, nil
}

// upstream finds if a given branch tracks an upstream.
// Returns found upstream or nil if branch doesn't track an upstream.
//
// Information about upstream is taken from .git/config file.
// If a branch has an upstream, there's a [branch] section in the file with two fields:
// "remote" - name of the remote containing upstreamn branch (or "." if upstream is a local branch)
// "merge" - full ref name of the upstream branch (eg, ref/heads/master)
func (r *Repo) upstream(branch string) (*Upstream, error) {
	cfg, err := r.repo.Config()
	if err != nil {
		return nil, errors.Wrap(err, "Failed getting repo config")
	}

	// Find our branch in "branch" config sections
	bcfg := cfg.Branches[branch]
	if bcfg == nil {
		return nil, nil
	}

	remote := bcfg.Remote
	if remote == "" {
		return nil, nil
	}

	// TODO: check if this should be short or full ref name
	merge := bcfg.Merge.Short()
	if merge == "" {
		return nil, nil
	}

	return &Upstream{
		Remote: remote,
		Branch: merge,
	}, nil
}
