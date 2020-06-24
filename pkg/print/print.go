package print

import (
	"fmt"
	"strings"
)

// TODO: not sure if this works on windows. See https://github.com/mattn/go-colorable
const (
	colorRed    = "\033[1;31m%s\033[0m"
	colorGreen  = "\033[1;32m%s\033[0m"
	colorBlue   = "\033[1;34m%s\033[0m"
	colorYellow = "\033[1;33m%s\033[0m"
)

const (
	untracked   = "untracked"
	uncommitted = "uncommitted"
	ahead       = "ahead"
	behind      = "behind"
	noUpstream  = "no upstream"
	ok          = "ok"
	detached    = "detached"
	head        = "HEAD"
)

// Repo is a git repository
// TODO: maybe branch should be a separate interface
type Repo interface {
	Path() string
	Branches() ([]string, error)
	CurrentBranch() (string, error)
	Upstream(branch string) (string, error)
	AheadBehind(branch string, upstream string) (int, int, error)
	Uncommitted() (int, error)
	Untracked() (int, error)
}

// // Printer provides a way to print a list of repos and their statuses
// type Printer interface {
// 	Print(root string, repos []Repo) string
// }

// prints status of currently checked out branch and the work tree.
// The format is: branch_name branch_status [ worktree_status ]
// Eg: master 1 head 2 behind [ 1 uncomitted ]
func printCurrentBranchLine(r Repo) string {
	var res []string

	current, err := r.CurrentBranch()
	if err != nil {
		return printErr(err)
	}

	// if current head is detached don't print its status
	if current == head {
		return fmt.Sprintf(colorYellow, detached)
	}

	status, err := printBranchStatus(r, current)
	if err != nil {
		return printErr(err)
	}

	worktree, err := printWorkTreeStatus(r)
	if err != nil {
		return printErr(err)
	}

	res = append(res, printBranchName(current))

	// if worktree is not clean and branch is ok then it shouldn't be ok
	if worktree != "" && strings.Contains(status, ok) {
		res = append(res, worktree)
	} else {
		res = append(res, status)
		res = append(res, worktree)
	}

	return strings.Join(res, " ")
}

func printBranchName(branch string) string {
	return fmt.Sprintf(colorBlue, branch)
}

func printBranchStatus(r Repo, branch string) (string, error) {
	var res []string
	upstream, err := r.Upstream(branch)
	if err != nil {
		return "", err
	}

	if upstream == "" {
		return fmt.Sprintf(colorYellow, noUpstream), nil
	}

	a, b, err := r.AheadBehind(branch, upstream)
	if err != nil {
		return printErr(err), nil
	}

	if a == 0 && b == 0 {
		return fmt.Sprintf(colorGreen, ok), nil
	}

	if a != 0 {
		res = append(res, fmt.Sprintf(colorYellow, fmt.Sprintf("%d %s", a, ahead)))
	}
	if b != 0 {
		res = append(res, fmt.Sprintf(colorYellow, fmt.Sprintf("%d %s", b, behind)))
	}

	return strings.Join(res, " "), nil
}

func printWorkTreeStatus(r Repo) (string, error) {
	uc, err := r.Uncommitted()
	if err != nil {
		return "", err
	}

	ut, err := r.Untracked()
	if err != nil {
		return "", err
	}

	if uc == 0 && ut == 0 {
		return "", nil
	}

	var res []string
	res = append(res, "[")
	if uc != 0 {
		res = append(res, fmt.Sprintf(colorRed, fmt.Sprintf("%d %s", uc, uncommitted)))
	}
	if ut != 0 {
		res = append(res, fmt.Sprintf(colorRed, fmt.Sprintf("%d %s", ut, untracked)))
	}

	res = append(res, "]")

	return strings.Join(res, " "), nil
}

func printErr(err error) string {
	return fmt.Sprintf(colorRed, err.Error())
}
