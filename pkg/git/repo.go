package git

import (
	"fmt"
	"git-get/pkg/io"
	"git-get/pkg/run"
	"net/url"
	"strconv"
	"strings"
)

const (
	dotgit    = ".git"
	untracked = "??" // Untracked files are marked as "??" in git status output.
	master    = "master"
	head      = "HEAD"
)

// Repo represents a git repository on disk.
type Repo interface {
	Path() string
	Branches() ([]string, error)
	CurrentBranch() (string, error)
	Fetch() error
	Remote() (string, error)
	Uncommitted() (int, error)
	Untracked() (int, error)
	Upstream(string) (string, error)
	AheadBehind(string, string) (int, int, error)
}

type repo struct {
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
func Open(path string) (Repo, error) {
	_, err := io.Exists(path)
	if err != nil {
		return nil, err
	}

	return &repo{
		path: path,
	}, nil
}

// Clone clones repository specified with CloneOpts.
func Clone(opts *CloneOpts) (Repo, error) {
	// TODO: not sure if this check should be here
	if opts.IgnoreExisting {
		return nil, nil
	}

	runGit := run.Git("clone", opts.URL.String(), opts.Path)
	if opts.Branch != "" {
		runGit = run.Git("clone", "--branch", opts.Branch, "--single-branch", opts.URL.String(), opts.Path)
	}

	var err error
	if opts.Quiet {
		err = runGit.AndShutUp()
	} else {
		err = runGit.AndShow()
	}

	if err != nil {
		return nil, err
	}

	repo, err := Open(opts.Path)
	return repo, err
}

// Fetch preforms a git fetch on all remotes
func (r *repo) Fetch() error {
	err := run.Git("fetch", "--all").OnRepo(r.path).AndShutUp()
	return err
}

// Uncommitted returns the number of uncommitted files in the repository.
// Only tracked files are not counted.
func (r *repo) Uncommitted() (int, error) {
	out, err := run.Git("status", "--ignore-submodules", "--porcelain").OnRepo(r.path).AndCaptureLines()
	if err != nil {
		return 0, err
	}

	count := 0
	for _, line := range out {
		// Don't count lines with untracked files and empty lines.
		if !strings.HasPrefix(line, untracked) && strings.TrimSpace(line) != "" {
			count++
		}
	}

	return count, nil
}

// Untracked returns the number of untracked files in the repository.
func (r *repo) Untracked() (int, error) {
	out, err := run.Git("status", "--ignore-submodules", "--untracked-files=all", "--porcelain").OnRepo(r.path).AndCaptureLines()
	if err != nil {
		return 0, err
	}

	count := 0
	for _, line := range out {
		if strings.HasPrefix(line, untracked) {
			count++
		}
	}

	return count, nil
}

// CurrentBranch returns the short name currently checked-out branch for the repository.
// If repo is in a detached head state, it will return "HEAD".
func (r *repo) CurrentBranch() (string, error) {
	out, err := run.Git("rev-parse", "--symbolic-full-name", "--abbrev-ref", "HEAD").OnRepo(r.path).AndCaptureLine()
	if err != nil {
		return "", err
	}

	return out, nil
}

// Branches returns a list of local branches in the repository.
func (r *repo) Branches() ([]string, error) {
	out, err := run.Git("branch", "--format=%(refname:short)").OnRepo(r.path).AndCaptureLines()
	if err != nil {
		return nil, err
	}

	// TODO: Is detached head shown always on the first line? Maybe we don't need to iterate over everything.
	// Remove the line containing detached head.
	for i, line := range out {
		if strings.Contains(line, "HEAD detached") {
			out = append(out[:i], out[i+1:]...)
		}
	}

	return out, nil
}

// Upstream returns the name of an upstream branch if a given branch is tracking one.
// Otherwise it returns an empty string.
func (r *repo) Upstream(branch string) (string, error) {
	out, err := run.Git("rev-parse", "--abbrev-ref", "--symbolic-full-name", fmt.Sprintf("%s@{upstream}", branch)).OnRepo(r.path).AndCaptureLine()
	if err != nil {
		// TODO: no upstream will also throw an error.
		return "", nil
	}

	return out, nil
}

// AheadBehind returns the number of commits a given branch is ahead and/or behind the upstream.
func (r *repo) AheadBehind(branch string, upstream string) (int, int, error) {
	out, err := run.Git("rev-list", "--left-right", "--count", fmt.Sprintf("%s...%s", branch, upstream)).OnRepo(r.path).AndCaptureLine()
	if err != nil {
		return 0, 0, err
	}

	// rev-list --left-right --count output is separated by a tab
	lr := strings.Split(out, "\t")

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

// Remote returns URL of remote repository.
func (r *repo) Remote() (string, error) {
	// https://stackoverflow.com/a/16880000/1085632
	out, err := run.Git("ls-remote", "--get-url").OnRepo(r.path).AndCaptureLine()
	if err != nil {
		return "", err
	}

	// TODO: needs testing. What happens when there are more than 1 remotes?
	return out, nil
}

// Path returns path to the repository.
func (r *repo) Path() string {
	return r.path
}
