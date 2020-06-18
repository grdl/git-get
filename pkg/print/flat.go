package print

import (
	"fmt"
	"git-get/pkg/repo"
	"path/filepath"
	"strings"
)

// FlatPrinter implements Printer interface and provides method for printing list of repos in flat format.
type FlatPrinter struct{}

// Print generates a flat list of repositories and their statuses - each repo in new line with full path.
func (p *FlatPrinter) Print(root string, repos []*repo.Repo) string {
	var str strings.Builder

	for _, r := range repos {
		path := strings.TrimPrefix(r.Path, root)
		path = strings.Trim(path, string(filepath.Separator))

		str.WriteString(fmt.Sprintf("\n%s %s", path, printWorktreeStatus(r)))

		for _, branch := range r.Status.Branches {
			// Don't print the status of the current branch. It was already printed above.
			if branch.Name == r.Status.CurrentBranch {
				continue
			}

			indent := strings.Repeat(" ", len(path))
			str.WriteString(fmt.Sprintf("\n%s %s", indent, printBranchStatus(branch)))
		}
	}

	return str.String()
}
