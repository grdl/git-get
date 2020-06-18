package print

import (
	"fmt"
	"git-get/pkg/repo"
	"path/filepath"
	"strings"
)

type FlatPrinter struct{}

func (p *FlatPrinter) Print(root string, repos []*repo.Repo) string {
	val := root

	for _, r := range repos {
		path := strings.TrimPrefix(r.Path, root)
		path = strings.Trim(path, string(filepath.Separator))

		val += fmt.Sprintf("\n%s %s", path, printWorktreeStatus(r))

		for _, branch := range r.Status.Branches {
			// Don't print the status of the current branch. It was already printed above.
			if branch.Name == r.Status.CurrentBranch {
				continue
			}

			indent := strings.Repeat(" ", len(path))
			val += fmt.Sprintf("\n%s %s", indent, printBranchStatus(branch))
		}
	}

	return val
}
