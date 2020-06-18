package repo

import (
	"git-get/pkg/cfg"
	"sort"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/revlist"

	"github.com/spf13/viper"

	"github.com/go-git/go-billy/v5/osfs"

	"github.com/go-git/go-git/v5/plumbing/format/gitignore"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/pkg/errors"
)

const (
	StatusUnknown     = "unknown"
	StatusDetached    = "detached HEAD"
	StatusNoUpstream  = "no upstream"
	StatusAhead       = "ahead"
	StatusBehind      = "behind"
	StatusOk          = "ok"
	StatusUncommitted = "uncommitted"
	StatusUntracked   = "untracked"
)

type RepoStatus struct {
	HasUntrackedFiles     bool
	HasUncommittedChanges bool
	CurrentBranch         string
	Branches              []*BranchStatus
}

type BranchStatus struct {
	Name     string
	Upstream string
	Ahead    int
	Behind   int
}

func (r *Repo) LoadStatus() error {
	// Fetch from remotes if executed with --fetch flag. Ignore the "already up-to-date" errors.
	if viper.GetBool(cfg.KeyFetch) {
		err := r.Fetch()
		if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
			return errors.Wrap(err, "Failed fetching from remotes")
		}
	}

	wt, err := r.Worktree()
	if err != nil {
		return errors.Wrap(err, "Failed getting worktree")
	}

	// worktree.Status doesn't load gitignore patterns that are defined outside of .gitignore file using excludesfile.
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
	r.Status.CurrentBranch = currentBranch(r)

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

func currentBranch(r *Repo) string {
	head, err := r.Head()
	if err != nil {
		return StatusUnknown
	}

	if head.Name().Short() == plumbing.HEAD.String() {
		return StatusDetached
	}

	return head.Name().Short()
}

func (r *Repo) loadBranchesStatus() error {
	iter, err := r.Branches()
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

	// Sort branches by name (but with "master" always at the top). It's useful to have them sorted for printing and testing.
	sort.Slice(r.Status.Branches, func(i, j int) bool {
		if r.Status.Branches[i].Name == "master" {
			return true
		}

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

	ahead, behind, err := r.aheadBehind(branch, upstream)
	if err != nil {
		return nil, err
	}

	bs.Upstream = upstream
	bs.Ahead = ahead
	bs.Behind = behind

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
	cfg, err := r.Config()
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

func (r *Repo) aheadBehind(localBranch string, upstreamBranch string) (ahead int, behind int, err error) {
	localHash, err := r.ResolveRevision(plumbing.Revision(localBranch))
	if err != nil {
		return 0, 0, errors.Wrapf(err, "Failed resolving revision %s", localBranch)
	}

	upstreamHash, err := r.ResolveRevision(plumbing.Revision(upstreamBranch))
	if err != nil {
		return 0, 0, errors.Wrapf(err, "Failed resolving revision %s", upstreamBranch)
	}

	behind, err = r.revlistCount(*localHash, *upstreamHash)
	if err != nil {
		return 0, 0, errors.Wrapf(err, "Failed counting commits behind %s", upstreamBranch)
	}

	ahead, err = r.revlistCount(*upstreamHash, *localHash)
	if err != nil {
		return 0, 0, errors.Wrapf(err, "Failed counting commits ahead of %s", upstreamBranch)
	}

	return ahead, behind, nil
}

// revlistCount counts the number of commits between two hashes.
// https://github.com/src-d/go-git/issues/757#issuecomment-452697701
// TODO: See if this can be optimized. Running the loop twice feels wrong.
func (r *Repo) revlistCount(hash1, hash2 plumbing.Hash) (int, error) {
	ref1hist, err := revlist.Objects(r.Storer, []plumbing.Hash{hash1}, nil)
	if err != nil {
		return 0, err
	}

	ref2hist, err := revlist.Objects(r.Storer, []plumbing.Hash{hash2}, ref1hist)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, h := range ref2hist {
		if _, err = r.CommitObject(h); err == nil {
			count++
		}
	}

	return count, nil
}
