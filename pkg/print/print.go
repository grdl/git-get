package print

import (
	"fmt"
	"strings"
)

const (
	head = "HEAD"
)

// Printable represents a repository which status can be printed
type Printable interface {
	Path() string
	Current() string
	Branches() []string
	BranchStatus(string) string
	BranchDescription(string) []string
	WorkTreeStatus() string
	Remote() string
	Errors() []string
}

// Errors returns a printable list of errors from the slice of Printables or an empty string if there are no errors.
// It's meant to be appended at the end of Print() result.
func Errors(repos []Printable) string {
	errors := []string{}

	for _, repo := range repos {
		for _, err := range repo.Errors() {
			errors = append(errors, err)
		}
	}

	if len(errors) == 0 {
		return ""
	}

	var str strings.Builder
	str.WriteString(red("\nOops, errors happened when loading repository status:\n"))
	str.WriteString(strings.Join(errors, "\n"))

	return str.String()
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

func magenta(str string) string {
	return fmt.Sprintf("\033[1;35m%s\033[0m", str)
}
