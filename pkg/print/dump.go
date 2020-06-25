package print

import (
	"strings"
)

// DumpPrinter prints a list of repos in a dump file format.
type DumpPrinter struct{}

// NewDumpPrinter creates a DumpPrinter.
func NewDumpPrinter() *DumpPrinter {
	return &DumpPrinter{}
}

// Print generates a list of repos URLs. Each line contains a URL and, if applicable, a currently checked out branch name.
// It's a way to dump all repositories managed by git-get and is supposed to be consumed by `git get --dump`.
func (p *DumpPrinter) Print(repos []Printable) string {
	var str strings.Builder

	for i, r := range repos {
		str.WriteString(r.Remote())

		// TODO: if head is detached maybe we should get the revision it points to in case it's a tag
		if current := r.Current(); current != "" && current != head {
			str.WriteString(" " + current)
		}

		if i < len(repos)-1 {
			str.WriteString("\n")
		}
	}

	return str.String()
}
