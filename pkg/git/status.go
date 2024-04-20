package git

import (
	"fmt"
	"strings"
)

// Status contains human readable (and printable) representation of a git repository status.
type Status struct {
	path               string
	current            string
	branches           map[string]string   // key: branch name, value: branch status
	branchDescriptions map[string][]string // key: branch name, value: branch description
	worktree           string
	remote             string
	errors             []string // Slice of errors which occurred when loading the status.
}

// LoadStatus reads status of a repository.
// If fetch equals true, it first fetches from the remote repo before loading the status.
// If errors occur during loading, they are stored in Status.errors slice.
func (r *Repo) LoadStatus(fetch bool) *Status {
	status := &Status{
		path:     r.path,
		branches: make(map[string]string),
		errors:   make([]string, 0),
	}

	if fetch {
		if err := r.Fetch(); err != nil {
			status.errors = append(status.errors, err.Error())
		}
	}

	var err error
	status.current, err = r.CurrentBranch()
	if err != nil {
		status.errors = append(status.errors, err.Error())
	}

	var errs []error
	status.branches, status.branchDescriptions, errs = r.loadBranches()
	for _, err := range errs {
		status.errors = append(status.errors, err.Error())
	}

	status.worktree, err = r.loadWorkTree()
	if err != nil {
		status.errors = append(status.errors, err.Error())
	}

	status.remote, err = r.Remote()
	if err != nil {
		status.errors = append(status.errors, err.Error())
	}

	return status
}

func (r *Repo) loadBranches() (map[string]string, map[string][]string, []error) {
	statuses := make(map[string]string)
	descriptions := make(map[string][]string)
	errors := make([]error, 0)

	branches, err := r.Branches()
	if err != nil {
		errors = append(errors, err)
		return statuses, descriptions, errors
	}

	for _, branch := range branches {
		status, err := r.loadBranchStatus(branch)
		statuses[branch] = status
		if err != nil {
			errors = append(errors, err)
		}
		description, err := r.loadBranchDescription(branch)
		descriptions[branch] = description
		if err != nil {
			errors = append(errors, err)
		}
	}

	return statuses, descriptions, errors
}

func (r *Repo) loadBranchStatus(branch string) (string, error) {
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

func (r *Repo) loadBranchDescription(branch string) ([]string, error) {
	return r.Description(branch)
}

func (r *Repo) loadWorkTree() (string, error) {
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
func (s *Status) Path() string {
	return s.path
}

// Current returns the name of currently checked out branch (or tag or detached HEAD).
func (s *Status) Current() string {
	return s.current
}

// Branches returns a list of all branches names except the currently checked out one. Use Current() to get its name.
func (s *Status) Branches() []string {
	var branches []string
	for b := range s.branches {
		if b != s.current {
			branches = append(branches, b)
		}
	}
	return branches
}

// BranchStatus returns status of a given branch
func (s *Status) BranchStatus(branch string) string {
	return s.branches[branch]
}

// BranchDescription returns description of a given branch
func (s *Status) BranchDescription(branch string) []string {
	return s.branchDescriptions[branch]
}

// WorkTreeStatus returns status of a worktree
func (s *Status) WorkTreeStatus() string {
	return s.worktree
}

// Remote returns URL to remote repository
func (s *Status) Remote() string {
	return s.remote
}

// Errors is a slice of errors that occurred when loading repo status
func (s *Status) Errors() []string {
	return s.errors
}
