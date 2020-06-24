package print

import (
	"fmt"
	"strings"
)

// FlatPrinter prints a list of repos in a flat format.
type FlatPrinter struct{}

// NewFlatPrinter creates a FlatPrinter.
func NewFlatPrinter() *FlatPrinter {
	return &FlatPrinter{}
}

// Print generates a flat list of repositories and their statuses - each repo in new line with full path.
func (p *FlatPrinter) Print(repos []Repo) string {
	var str strings.Builder

	for _, r := range repos {
		str.WriteString(fmt.Sprintf("\n%s %s", r.Path(), printCurrentBranchLine(r)))

		branches, err := r.Branches()
		if err != nil {
			str.WriteString(printErr(err))
			continue
		}

		current, err := r.CurrentBranch()
		if err != nil {
			str.WriteString(printErr(err))
			continue
		}

		for _, branch := range branches {
			// Don't print the status of the current branch. It was already printed above.
			if branch == current {
				continue
			}

			status, err := printBranchStatus(r, branch)
			if err != nil {
				status = printErr(err)
			}

			indent := strings.Repeat(" ", len(r.Path()))
			str.WriteString(fmt.Sprintf("\n%s %s %s", indent, printBranchName(branch), status))
		}
	}

	return str.String()
}
