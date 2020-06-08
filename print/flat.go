package print

import (
	"fmt"
	"git-get/git"
	"path/filepath"
	"strings"
)

type FlatPrinter struct{}

func (p *FlatPrinter) Print(root string, repos []*git.Repo) string {
	val := root

	for _, repo := range repos {
		path := strings.TrimPrefix(repo.Path, root)
		path = strings.Trim(path, string(filepath.Separator))

		val += fmt.Sprintf("\n%s %s", path, printWorktreeStatus(repo))

		for _, branch := range repo.Status.Branches {
			// Don't print the status of the current branch. It was already printed above.
			if branch.Name == repo.Status.CurrentBranch {
				continue
			}

			indent := strings.Repeat(" ", len(path))
			val += fmt.Sprintf("\n%s %s", indent, printBranchStatus(branch))
		}
	}

	return val
}
