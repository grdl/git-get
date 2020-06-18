package print

import (
	"git-get/pkg/git"
	"strings"
)

type DumpPrinter struct{}

// Print generates a list of repos URLs. Each line contains a URL and, if applicable, a currently checked out branch name.
// It's a way to dump all repositories managed by git-get and is supposed to be consumed by `git get --dump`.
func (p *DumpPrinter) Print(_ string, repos []*git.Repo) string {
	var str strings.Builder

	for i, repo := range repos {
		remotes, err := repo.Remotes()
		if err != nil || len(remotes) == 0 {
			continue
		}

		// TODO: Needs work. Right now we're just assuming the first remote is the origin one and the one from which the current branch is checked out.
		url := remotes[0].Config().URLs[0]
		current := repo.Status.CurrentBranch

		str.WriteString(url)

		if current != git.StatusDetached && current != git.StatusUnknown {
			str.WriteString(" " + current)
		}

		if i < len(repos)-1 {
			str.WriteString("\n")
		}
	}

	return str.String()
}
