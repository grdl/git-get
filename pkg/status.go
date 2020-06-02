package pkg

import (
	"sort"
	"strings"

	"github.com/go-git/go-billy/v5/osfs"

	"github.com/go-git/go-git/v5/plumbing/format/gitignore"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/pkg/errors"
)

type RepoStatus struct {
	HasUntrackedFiles     bool
	HasUncommittedChanges bool
	Branches              []*BranchStatus
}

type BranchStatus struct {
	Name      string
	Upstream  string
	NeedsPull bool
	NeedsPush bool
}

func (r *Repo) LoadStatus() error {
	wt, err := r.repo.Worktree()
	if err != nil {
		return errors.Wrap(err, "Failed getting worktree")
	}

	// worktree.Status doesn't load gitignore patterns that may be defined outside of .gitignore file using excludesfile
	// We need to load them explicitly here
	// TODO: variables are not expanded so if excludesfile is declared like "~/gitignore_global" or "$HOME/gitignore_global", this will fail to open it
	globalPatterns, err := gitignore.LoadGlobalPatterns(osfs.New(""))
	if err != nil {
		return errors.Wrap(err, "Failed loading global gitignore patterns")
	}
	wt.Excludes = append(wt.Excludes, globalPatterns...)

	systemPatterns, err := gitignore.LoadSystemPatterns(osfs.New(""))
	if err != nil {
		return errors.Wrap(err, "Failed loading system gitignore patterns")
	}
	wt.Excludes = append(wt.Excludes, systemPatterns...)

	status, err := wt.Status()
	if err != nil {
		return errors.Wrap(err, "Failed getting worktree status")
	}

	r.Status.HasUncommittedChanges = hasUncommitted(status)
	r.Status.HasUntrackedFiles = hasUntracked(status)

	err = r.loadBranchesStatus()
	if err != nil {
		return err
	}

	return nil
}

// hasUntracked returns true if there are any untracked files in the worktree
func hasUntracked(status git.Status) bool {
	for _, fs := range status {
		if fs.Worktree == git.Untracked || fs.Staging == git.Untracked {
			return true
		}
	}

	return false
}

// hasUncommitted returns true if there are any uncommitted (but tracked) files in the worktree
func hasUncommitted(status git.Status) bool {
	// If repo is clean it means every file in worktree and staging has 'Unmodified' state
	if status.IsClean() {
		return false
	}

	// If repo is not clean, check if any file has state different than 'Untracked' - it means they are tracked and have uncommitted modifications
	for _, fs := range status {
		if fs.Worktree != git.Untracked || fs.Staging != git.Untracked {
			return true
		}
	}

	return false
}

func (r *Repo) loadBranchesStatus() error {
	iter, err := r.repo.Branches()
	if err != nil {
		return errors.Wrap(err, "Failed getting branches iterator")
	}

	err = iter.ForEach(func(reference *plumbing.Reference) error {
		bs, err := r.newBranchStatus(reference.Name().Short())
		if err != nil {
			return err
		}

		r.Status.Branches = append(r.Status.Branches, bs)
		return nil
	})
	if err != nil {
		return errors.Wrap(err, "Failed iterating over branches")
	}

	// Sort branches by name. It's useful to have them sorted for printing and testing.
	sort.Slice(r.Status.Branches, func(i, j int) bool {
		return strings.Compare(r.Status.Branches[i].Name, r.Status.Branches[j].Name) < 0
	})
	return nil
}

func (r *Repo) newBranchStatus(branch string) (*BranchStatus, error) {
	bs := &BranchStatus{
		Name: branch,
	}

	upstream, err := r.upstream(branch)
	if err != nil {
		return nil, err
	}

	if upstream == "" {
		return bs, nil
	}

	needsPull, needsPush, err := r.needsPullOrPush(branch, upstream)
	if err != nil {
		return nil, err
	}

	bs.Upstream = upstream
	bs.NeedsPush = needsPush
	bs.NeedsPull = needsPull

	return bs, nil
}

// upstream finds if a given branch tracks an upstream.
// Returns found upstream branch name (eg, origin/master) or empty string if branch doesn't track an upstream.
//
// Information about upstream is taken from .git/config file.
// If a branch has an upstream, there's a [branch] section in the file with two fields:
// "remote" - name of the remote containing upstream branch (or "." if upstream is a local branch)
// "merge" - full ref name of the upstream branch (eg, ref/heads/master)
func (r *Repo) upstream(branch string) (string, error) {
	cfg, err := r.repo.Config()
	if err != nil {
		return "", errors.Wrap(err, "Failed getting repo config")
	}

	// Check if our branch exists in "branch" config sections. If not, it doesn't have an upstream configured.
	bcfg := cfg.Branches[branch]
	if bcfg == nil {
		return "", nil
	}

	remote := bcfg.Remote
	if remote == "" {
		return "", nil
	}

	merge := bcfg.Merge.Short()
	if merge == "" {
		return "", nil
	}
	return remote + "/" + merge, nil
}

func (r *Repo) needsPullOrPush(localBranch string, upstreamBranch string) (needsPull bool, needsPush bool, err error) {
	localHash, err := r.repo.ResolveRevision(plumbing.Revision(localBranch))
	if err != nil {
		return false, false, errors.Wrapf(err, "Failed resolving revision %s", localBranch)
	}

	upstreamHash, err := r.repo.ResolveRevision(plumbing.Revision(upstreamBranch))
	if err != nil {
		return false, false, errors.Wrapf(err, "Failed resolving revision %s", upstreamBranch)
	}

	localCommit, err := r.repo.CommitObject(*localHash)
	if err != nil {
		return false, false, errors.Wrapf(err, "Failed finding a commit for hash %s", localHash.String())
	}

	upstreamCommit, err := r.repo.CommitObject(*upstreamHash)
	if err != nil {
		return false, false, errors.Wrapf(err, "Failed finding a commit for hash %s", upstreamHash.String())
	}

	// If local branch hash is the same as upstream, it means there is no difference between local and upstream
	if *localHash == *upstreamHash {
		return false, false, nil
	}

	commons, err := localCommit.MergeBase(upstreamCommit)
	if err != nil {
		return false, false, errors.Wrapf(err, "Failed finding common ancestors for branches %s & %s", localBranch, upstreamBranch)
	}

	if len(commons) == 0 {
		// TODO: No common ancestors. This should be an error
		return false, false, nil
	}

	if len(commons) > 1 {
		// TODO: multiple best ancestors. How to handle this?
		return false, false, nil
	}

	mergeBase := commons[0]

	// If merge base is the same as upstream branch, local branch is ahead and push is needed
	// If merge base is the same as local branch, local branch is behind and pull is needed
	// If merge base is something else, branches have diverged and merge is needed (both pull and push)
	// ref: https://stackoverflow.com/a/17723781/1085632

	if mergeBase.Hash == *upstreamHash {
		return false, true, nil
	}

	if mergeBase.Hash == *localHash {
		return true, false, nil
	}

	return true, true, nil
}
