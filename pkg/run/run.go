// Package run provides methods for running git command and capturing their output and errors
package run

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Cmd represents a git command.
// The command is executed by chaining functions: Git() + optional OnRepo() + output specifier.
// This way the function chain reads more naturally.
//
// Examples of different compositions:
//
//   - run.Git("clone", <URL>).AndShow()
//     means running "git clone <URL>" and printing the progress into stdout
//
//   - run.Git("branch","-a").OnRepo(<REPO>).AndCaptureLines()
//     means running "git branch -a" inside <REPO> and returning a slice of branch names
//
//   - run.Git("pull").OnRepo(<REPO>).AndShutUp()
//     means running "git pull" inside <REPO> and not printing any output
type Cmd struct {
	cmd  *exec.Cmd
	args string
	path string
}

// Git creates a git command with given arguments.
func Git(args ...string) *Cmd {
	ctx := context.Background()

	return &Cmd{
		cmd:  exec.CommandContext(ctx, "git", args...),
		args: strings.Join(args, " "),
	}
}

// OnRepo makes the command run inside a given repository path. Otherwise the command is run outside of any repository.
// Commands like "git clone" or "git config --global" don't have to (or shouldn't in some cases) be run inside a repo.
func (c *Cmd) OnRepo(path string) *Cmd {
	if strings.TrimSpace(path) == "" {
		return c
	}

	insert := []string{"--work-tree", path, "--git-dir", filepath.Join(path, ".git")}
	// Insert into the args slice after the 1st element (https://github.com/golang/go/wiki/SliceTricks#insert)
	c.cmd.Args = append(c.cmd.Args[:1], append(insert, c.cmd.Args[1:]...)...)

	c.path = path

	return c
}

// AndCaptureLines executes the command and returns its output as a slice of lines.
func (c *Cmd) AndCaptureLines() ([]string, error) {
	errStream := &bytes.Buffer{}
	c.cmd.Stderr = errStream

	out, err := c.cmd.Output()
	if err != nil {
		return nil, &GitError{errStream, c.args, c.path, err}
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

// AndShow executes the command and prints its stderr and stdout.
func (c *Cmd) AndShow() error {
	c.cmd.Stdout = os.Stdout
	c.cmd.Stderr = os.Stderr

	err := c.cmd.Run()
	if err != nil {
		return &GitError{&bytes.Buffer{}, c.args, c.path, err}
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
		return &GitError{errStream, c.args, c.path, err}
	}

	return nil
}

// GitError provides more visibility into why an git command had failed.
type GitError struct {
	Stderr *bytes.Buffer
	Args   string
	Path   string
	Err    error
}

func (e GitError) Error() string {
	msg := e.Stderr.String()

	if e.Path == "" {
		return fmt.Sprintf("git %s failed: %s", e.Args, msg)
	}

	return fmt.Sprintf("git %s failed on %s: %s", e.Args, e.Path, msg)
}

func lines(output []byte) []string {
	lines := strings.TrimSuffix(string(output), "\n")
	return strings.Split(lines, "\n")
}
