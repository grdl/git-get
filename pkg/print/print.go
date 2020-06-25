package print

import "fmt"

const (
	head = "HEAD"
)

// Printable represents a repository which status can be printed
type Printable interface {
	Path() string
	Current() string
	Branches() []string
	BranchStatus(string) string
	WorkTreeStatus() string
	Remote() string
	Errors() []string
}

// TODO: not sure if this works on windows. See https://github.com/mattn/go-colorable
func red(str string) string {
	return fmt.Sprintf("\033[1;31m%s\033[0m", str)
}

func green(str string) string {
	return fmt.Sprintf("\033[1;32m%s\033[0m", str)
}

func blue(str string) string {
	return fmt.Sprintf("\033[1;34m%s\033[0m", str)
}

func yellow(str string) string {
	return fmt.Sprintf("\033[1;33m%s\033[0m", str)
}
