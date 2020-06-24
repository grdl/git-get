package git

import (
	"fmt"
	"git-get/pkg/file"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

const (
	dotgit    = ".git"
	untracked = "??" // Untracked files are marked as "??" in git status output.
	master    = "master"
	head      = "HEAD"
)

// Repo represents a git repository on disk.
type Repo struct {
	path string
}

// CloneOpts specify detail about repository to clone.
type CloneOpts struct {
	URL            *url.URL
	Path           string // TODO: should Path be a part of clone opts?
	Branch         string
	Quiet          bool
	IgnoreExisting bool
}

// Open checks if given path can be accessed and returns a Repo instance pointing to it.
func Open(path string) (*Repo, error) {
	_, err := file.Exists(path)
	if err != nil {
		return nil, err
	}

	return &Repo{
		path: path,
	}, nil
}

// Clone clones repository specified with CloneOpts.
func Clone(opts *CloneOpts) (*Repo, error) {
	// TODO: not sure if this check should be here
	if opts.IgnoreExisting {
		return nil, nil
	}

	args := []string{"clone", "--progress", "-v"}

	if opts.Branch != "" {
		args = append(args, "--branch", opts.Branch, "--single-branch")
	}

	if opts.Quiet {
		args = append(args, "--quiet")
	}

	args = append(args, opts.URL.String())
	args = append(args, opts.Path)

	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return nil, errors.Wrapf(err, "git clone failed")
	}

	repo, err := Open(opts.Path)
	return repo, err
}

// Fetch preforms a git fetch on all remotes
func (r *Repo) Fetch() error {
	cmd := gitCmd(r.path, "fetch", "--all", "--quiet")
	return cmd.Run()
}

// Uncommitted returns the number of uncommitted files in the repository.
// Only tracked files are not counted.
func (r *Repo) Uncommitted() (int, error) {
	cmd := gitCmd(r.path, "status", "--ignore-submodules", "--porcelain")

	out, err := cmd.Output()
	if err != nil {
		return 0, cmdError(cmd, err)
	}

	lines := lines(out)
	count := 0
	for _, line := range lines {
		// Don't count lines with untracked files and empty lines.
		if !strings.HasPrefix(line, untracked) && strings.TrimSpace(line) != "" {
			count++
		}
	}

	return count, nil
}

// Untracked returns the number of untracked files in the repository.
func (r *Repo) Untracked() (int, error) {
	cmd := gitCmd(r.path, "status", "--ignore-submodules", "--porcelain")

	out, err := cmd.Output()
	if err != nil {
		return 0, cmdError(cmd, err)
	}

	lines := lines(out)
	count := 0
	for _, line := range lines {
		if strings.HasPrefix(line, untracked) {
			count++
		}
	}

	return count, nil
}

// CurrentBranch returns the short name currently checked-out branch for the repository.
// If repo is in a detached head state, it will return "HEAD".
func (r *Repo) CurrentBranch() (string, error) {
	cmd := gitCmd(r.path, "rev-parse", "--symbolic-full-name", "--abbrev-ref", "HEAD")

	out, err := cmd.Output()
	if err != nil {
		return "", cmdError(cmd, err)
	}

	lines := lines(out)
	return lines[0], nil
}

// Branches returns a list of local branches in the repository.
func (r *Repo) Branches() ([]string, error) {
	cmd := gitCmd(r.path, "branch", "--format=%(refname:short)")

	out, err := cmd.Output()
	if err != nil {
		return nil, cmdError(cmd, err)
	}

	lines := lines(out)

	// TODO: Is detached head shown always on the first line? Maybe we don't need to iterate over everything.
	// Remove the line containing detached head.
	for i, line := range lines {
		if strings.Contains(line, "HEAD detached") {
			lines = append(lines[:i], lines[i+1:]...)
		}
	}

	return lines, nil
}

// Upstream returns the name of an upstream branch if a given branch is tracking one.
// Otherwise it returns an empty string.
func (r *Repo) Upstream(branch string) (string, error) {
	cmd := gitCmd(r.path, "rev-parse", "--abbrev-ref", "--symbolic-full-name", fmt.Sprintf("%s@{upstream}", branch))

	// TODO: no upstream will also throw an error.
	out, err := cmd.Output()
	if err != nil {
		return "", cmdError(cmd, err)
	}

	lines := lines(out)
	return lines[0], nil
}

// AheadBehind returns the number of commits a given branch is ahead and/or behind the upstream.
func (r *Repo) AheadBehind(branch string, upstream string) (int, int, error) {
	cmd := gitCmd(r.path, "rev-list", "--left-right", "--count", fmt.Sprintf("%s...%s", branch, upstream))

	out, err := cmd.Output()
	if err != nil {
		return 0, 0, cmdError(cmd, err)
	}

	lines := lines(out)

	// rev-list --left-right --count output is separated by a tab
	lr := strings.Split(lines[0], "\t")

	ahead, err := strconv.Atoi(lr[0])
	if err != nil {
		return 0, 0, err
	}

	behind, err := strconv.Atoi(lr[1])
	if err != nil {
		return 0, 0, err
	}

	return ahead, behind, nil
}

func gitCmd(repoPath string, args ...string) *exec.Cmd {
	args = append([]string{"--work-tree", repoPath, "--git-dir", path.Join(repoPath, dotgit)}, args...)
	return exec.Command("git", args...)
}

func lines(output []byte) []string {
	lines := strings.TrimSuffix(string(output), "\n")
	return strings.Split(lines, "\n")
}

func cmdError(cmd *exec.Cmd, err error) error {
	if err != nil {
		return errors.Wrapf(err, "git %s failed", cmd.Args[4]) // Show which git command failed (skip "--work-tree and --gitdir flags")
	}
	return nil
}
