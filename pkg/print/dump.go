package print

import (
	"strings"
)

// DumpRepo is a git repository printable into a dump file.
type DumpRepo interface {
	Path() string
	Remote() (string, error)
	CurrentBranch() (string, error)
}

// DumpPrinter prints a list of repos in a dump file format.
type DumpPrinter struct{}

// NewDumpPrinter creates a DumpPrinter.
func NewDumpPrinter() *DumpPrinter {
	return &DumpPrinter{}
}

// Print generates a list of repos URLs. Each line contains a URL and, if applicable, a currently checked out branch name.
// It's a way to dump all repositories managed by git-get and is supposed to be consumed by `git get --dump`.
func (p *DumpPrinter) Print(repos []DumpRepo) string {
	var str strings.Builder

	for i, r := range repos {
		url, err := r.Remote()
		if err != nil {
			continue
			// TODO: handle error?
		}

		str.WriteString(url)

		current, err := r.CurrentBranch()
		if err != nil || current != detached {
			str.WriteString(" " + current)
		}

		if i < len(repos)-1 {
			str.WriteString("\n")
		}
	}

	return str.String()
}
