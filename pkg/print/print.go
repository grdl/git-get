package print

import (
	"fmt"
	"git-get/pkg/repo"
	"strings"
)

type Printer interface {
	Print(root string, repos []*repo.Repo) string
}

const (
	ColorRed    = "\033[1;31m%s\033[0m"
	ColorGreen  = "\033[1;32m%s\033[0m"
	ColorBlue   = "\033[1;34m%s\033[0m"
	ColorYellow = "\033[1;33m%s\033[0m"
)

func printWorktreeStatus(r *repo.Repo) string {
	clean := true
	var status []string

	// if current branch status can't be found it's probably a detached head
	// TODO: what if current HEAD points to a tag?
	if current := r.CurrentBranchStatus(); current == nil {
		status = append(status, fmt.Sprintf(ColorYellow, r.Status.CurrentBranch))
	} else {
		status = append(status, printBranchStatus(current))
	}

	// TODO: this is ugly
	// unset clean flag to use it to render braces around worktree status and remove "ok" from branch status if it's there
	if r.Status.HasUncommittedChanges || r.Status.HasUntrackedFiles {
		clean = false
	}

	if !clean {
		status[len(status)-1] = strings.TrimSuffix(status[len(status)-1], repo.StatusOk)
		status = append(status, "[")
	}

	if r.Status.HasUntrackedFiles {
		status = append(status, fmt.Sprintf(ColorRed, repo.StatusUntracked))
	}

	if r.Status.HasUncommittedChanges {
		status = append(status, fmt.Sprintf(ColorRed, repo.StatusUncommitted))
	}

	if !clean {
		status = append(status, "]")
	}

	return strings.Join(status, " ")
}

func printBranchStatus(branch *repo.BranchStatus) string {
	// ok indicates that the branch has upstream and is not ahead or behind it
	ok := true
	var status []string

	status = append(status, fmt.Sprintf(ColorBlue, branch.Name))

	if branch.Upstream == "" {
		ok = false
		status = append(status, fmt.Sprintf(ColorYellow, repo.StatusNoUpstream))
	}

	if branch.Behind != 0 {
		ok = false
		status = append(status, fmt.Sprintf(ColorYellow, fmt.Sprintf("%d %s", branch.Behind, repo.StatusBehind)))
	}

	if branch.Ahead != 0 {
		ok = false
		status = append(status, fmt.Sprintf(ColorYellow, fmt.Sprintf("%d %s", branch.Ahead, repo.StatusAhead)))
	}

	if ok {
		status = append(status, fmt.Sprintf(ColorGreen, repo.StatusOk))
	}

	return strings.Join(status, " ")
}
