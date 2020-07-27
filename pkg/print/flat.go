package print

import (
	"fmt"
	"os"
	"strings"
)

// FlatPrinter prints a list of repos in a flat format.
type FlatPrinter struct{}

// NewFlatPrinter creates a FlatPrinter.
func NewFlatPrinter() *FlatPrinter {
	return &FlatPrinter{}
}

// Print generates a flat list of repositories and their statuses - each repo in new line with full path.
func (p *FlatPrinter) Print(repos []Printable) string {
	var str strings.Builder

	for _, r := range repos {
		str.WriteString(strings.TrimSuffix(r.Path(), string(os.PathSeparator)))

		if len(r.Errors()) > 0 {
			str.WriteString(" " + red("error") + "\n")
			continue
		}

		str.WriteString(" " + blue(r.Current()))

		current := r.BranchStatus(r.Current())
		worktree := r.WorkTreeStatus()

		if worktree != "" {
			worktree = fmt.Sprintf("[ %s ]", worktree)
		}

		if worktree == "" && current == "" {
			str.WriteString(" " + green("ok"))
		} else {
			str.WriteString(" " + strings.Join([]string{yellow(current), red(worktree)}, " "))
		}

		for _, branch := range r.Branches() {
			status := r.BranchStatus(branch)
			if status == "" {
				status = green("ok")
			}

			indent := strings.Repeat(" ", len(r.Path())-1)
			str.WriteString(fmt.Sprintf("\n%s %s %s", indent, blue(branch), yellow(status)))
		}

		str.WriteString("\n")
	}

	return str.String() + Errors(repos)
}
