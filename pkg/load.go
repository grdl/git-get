package pkg

import (
	"fmt"
	"git-get/pkg/git"
	"strings"
)

// Loaded represents a repository which status is Loaded from disk and stored for printing.
type Loaded struct {
	path     string
	current  string
	branches map[string]string // key: branch name, value: branch status
	worktree string
	remote   string
	errors   []string
}

// Load reads status of a repository at a given path.
func Load(path string, fetch bool) *Loaded {
	loaded := &Loaded{
		path:     path,
		branches: make(map[string]string),
		errors:   make([]string, 0),
	}

	repo, err := git.Open(path)
	if err != nil {
		loaded.errors = append(loaded.errors, err.Error())
		return loaded
	}

	if fetch {
		err = repo.Fetch()
		if err != nil {
			loaded.errors = append(loaded.errors, err.Error())
		}
	}

	loaded.current, err = repo.CurrentBranch()
	if err != nil {
		loaded.errors = append(loaded.errors, err.Error())
	}

	var errs []error
	loaded.branches, errs = loadBranches(repo)
	for _, err := range errs {
		loaded.errors = append(loaded.errors, err.Error())
	}

	loaded.worktree, err = loadWorkTree(repo)
	if err != nil {
		loaded.errors = append(loaded.errors, err.Error())
	}

	loaded.remote, err = repo.Remote()
	if err != nil {
		loaded.errors = append(loaded.errors, err.Error())
	}

	return loaded
}

func loadBranches(r git.Repo) (map[string]string, []error) {
	statuses := make(map[string]string)
	errors := make([]error, 0)

	branches, err := r.Branches()
	if err != nil {
		errors = append(errors, err)
		return statuses, errors
	}

	for _, branch := range branches {
		status, err := loadBranchStatus(r, branch)
		statuses[branch] = status
		if err != nil {
			errors = append(errors, err)
		}
	}

	return statuses, errors
}

func loadBranchStatus(r git.Repo, branch string) (string, error) {
	upstream, err := r.Upstream(branch)
	if err != nil {
		return "", err
	}

	if upstream == "" {
		return "no upstream", nil
	}

	ahead, behind, err := r.AheadBehind(branch, upstream)
	if err != nil {
		return "", err
	}

	if ahead == 0 && behind == 0 {
		return "", nil
	}

	var res []string
	if ahead != 0 {
		res = append(res, fmt.Sprintf("%d ahead", ahead))
	}
	if behind != 0 {
		res = append(res, fmt.Sprintf("%d behind", behind))
	}

	return strings.Join(res, " "), nil
}

func loadWorkTree(r git.Repo) (string, error) {
	uncommitted, err := r.Uncommitted()
	if err != nil {
		return "", err
	}

	untracked, err := r.Untracked()
	if err != nil {
		return "", err
	}

	if uncommitted == 0 && untracked == 0 {
		return "", nil
	}

	var res []string
	if uncommitted != 0 {
		res = append(res, fmt.Sprintf("%d uncommitted", uncommitted))
	}
	if untracked != 0 {
		res = append(res, fmt.Sprintf("%d untracked", untracked))
	}

	return strings.Join(res, " "), nil
}

// Path returns path to a repository.
func (r *Loaded) Path() string {
	return r.path
}

// Current returns the name of currently checked out branch (or tag or detached HEAD).
func (r *Loaded) Current() string {
	return r.current
}

// Branches returns a list of all branches names except the currently checked out one. Use Current() to get its name.
func (r *Loaded) Branches() []string {
	var branches []string
	for b := range r.branches {
		if b != r.current {
			branches = append(branches, b)
		}
	}
	return branches
}

// BranchStatus returns status of a given branch
func (r *Loaded) BranchStatus(branch string) string {
	return r.branches[branch]
}

// WorkTreeStatus returns status of a worktree
func (r *Loaded) WorkTreeStatus() string {
	return r.worktree
}

// Remote returns URL to remote repository
func (r *Loaded) Remote() string {
	return r.remote
}

// Errors is a slice of errors that occurred when loading repo status
func (r *Loaded) Errors() []string {
	return r.errors
}
