// Package run provides methods for running git command and capturing their output and errors
package run

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	pathpkg "path"
	"strings"
)

// Cmd represents a git command.
// The command is executed by chaining functions: Git() + optional OnRepo() + output specifier.
// This way the function chain reads more naturally.
//
// Examples of different compositions:
//
// - run.Git("clone", <URL>).AndShow()
//   means running "git clone <URL>" and printing the progress into stdout
//
// - run.Git("branch","-a").OnRepo(<REPO>).AndCaptureLines()
//   means running "git branch -a" inside <REPO> and returning a slice of branch names
//
// - run.Git("pull").OnRepo(<REPO>).AndShutUp()
//   means running "git pull" inside <REPO> and not printing any output
type Cmd struct {
	cmd *exec.Cmd
}

// Git creates a git command with given arguments.
func Git(args ...string) *Cmd {
	return &Cmd{
		cmd: exec.Command("git", args...),
	}
}

// OnRepo makes the command run inside a given repository path. Otherwise the command is run outside of any repository.
// Commands like "git clone" or "git config --global" don't have to (or shouldn't in some cases) be run inside a repo.
func (c *Cmd) OnRepo(path string) *Cmd {
	if strings.TrimSpace(path) == "" {
		return c
	}

	insert := []string{"--work-tree", path, "--git-dir", pathpkg.Join(path, ".git")}
	// Insert into the args slice after the 1st element (https://github.com/golang/go/wiki/SliceTricks#insert)
	c.cmd.Args = append(c.cmd.Args[:1], append(insert, c.cmd.Args[1:]...)...)

	return c
}

// AndCaptureLines executes the command and returns its output as a slice of lines.
func (c *Cmd) AndCaptureLines() ([]string, error) {
	errStream := &bytes.Buffer{}
	c.cmd.Stderr = errStream

	out, err := c.cmd.Output()
	if err != nil {
		return nil, &GitError{errStream, c.cmd.Args, err}
	}

	lines := lines(out)
	if len(lines) == 0 {
		return []string{""}, nil
	}

	return lines, nil
}

// AndCaptureLine executes the command and returns the first line of its output.
func (c *Cmd) AndCaptureLine() (string, error) {
	lines, err := c.AndCaptureLines()
	if err != nil {
		return "", err
	}
	return lines[0], nil
}

// AndShow executes the command and prints its output into standard output.
func (c *Cmd) AndShow() error {
	c.cmd.Stdout = os.Stdout

	errStream := &bytes.Buffer{}
	c.cmd.Stderr = errStream

	err := c.cmd.Run()
	if err != nil {
		return &GitError{errStream, c.cmd.Args, err}
	}
	return nil
}

// AndShutUp executes the command and doesn't return or show any output.
func (c *Cmd) AndShutUp() error {
	c.cmd.Stdout = nil

	errStream := &bytes.Buffer{}
	c.cmd.Stderr = errStream

	err := c.cmd.Run()
	if err != nil {
		return &GitError{errStream, c.cmd.Args, err}
	}
	return nil
}

// GitError provides more visibility into why an git command had failed.
type GitError struct {
	Stderr *bytes.Buffer
	Args   []string
	Err    error
}

func (e GitError) Error() string {
	msg := e.Stderr.String()
	if msg != "" && !strings.HasSuffix(msg, "\n") {
		msg += "\n"
	}
	return fmt.Sprintf("%s%q: %s", msg, strings.Join(e.Args, " "), e.Err)
}

func lines(output []byte) []string {
	lines := strings.TrimSuffix(string(output), "\n")
	return strings.Split(lines, "\n")
}
