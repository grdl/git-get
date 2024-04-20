package git

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"git-get/pkg/run"
)

const (
	dotgit    = ".git"
	untracked = "??" // Untracked files are marked as "??" in git status output.
	master    = "master"
	head      = "HEAD"
)

// Repo represents a git Repository cloned or initialized on disk.
type Repo struct {
	path string
}

// CloneOpts specify detail about Repository to clone.
type CloneOpts struct {
	URL    *url.URL
	Path   string // TODO: should Path be a part of clone opts?
	Branch string
	Quiet  bool
}

// Open checks if given path can be accessed and returns a Repo instance pointing to it.
func Open(path string) (*Repo, error) {
	if _, err := Exists(path); err != nil {
		return nil, err
	}

	return &Repo{
		path: path,
	}, nil
}

// Clone clones Repository specified with CloneOpts.
func Clone(opts *CloneOpts) (*Repo, error) {
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
		cleanupFailedClone(opts.Path)
		return nil, err
	}

	Repo, err := Open(opts.Path)
	return Repo, err
}

// Fetch preforms a git fetch on all remotes
func (r *Repo) Fetch() error {
	err := run.Git("fetch", "--all").OnRepo(r.path).AndShutUp()
	return err
}

// Uncommitted returns the number of uncommitted files in the Repository.
// Only tracked files are not counted.
func (r *Repo) Uncommitted() (int, error) {
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

// Untracked returns the number of untracked files in the Repository.
func (r *Repo) Untracked() (int, error) {
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

// CurrentBranch returns the short name currently checked-out branch for the Repository.
// If Repo is in a detached head state, it will return "HEAD".
func (r *Repo) CurrentBranch() (string, error) {
	out, err := run.Git("rev-parse", "--symbolic-full-name", "--abbrev-ref", "HEAD").OnRepo(r.path).AndCaptureLine()
	if err != nil {
		return "", err
	}

	return out, nil
}

// Branches returns a list of local branches in the Repository.
func (r *Repo) Branches() ([]string, error) {
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
func (r *Repo) Upstream(branch string) (string, error) {
	out, err := run.Git("rev-parse", "--abbrev-ref", "--symbolic-full-name", fmt.Sprintf("%s@{upstream}", branch)).OnRepo(r.path).AndCaptureLine()
	if err != nil {
		// TODO: no upstream will also throw an error.
		return "", nil
	}

	return out, nil
}

// Description returns description of a branch if a given branch has one set with "git branch --edit-description".
// Otherwise it returns an empty slice.
func (r *Repo) Description(branch string) ([]string, error) {
	out, err := run.Git("config", fmt.Sprintf("branch.%s.description", branch)).OnRepo(r.path).AndCaptureLines()
	if err != nil {
		return nil, nil
	}
	return out, nil
}

// AheadBehind returns the number of commits a given branch is ahead and/or behind the upstream.
func (r *Repo) AheadBehind(branch string, upstream string) (int, int, error) {
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

// Remote returns URL of remote Repository.
func (r *Repo) Remote() (string, error) {
	// https://stackoverflow.com/a/16880000/1085632
	out, err := run.Git("ls-remote", "--get-url").OnRepo(r.path).AndCaptureLine()
	if err != nil {
		return "", err
	}

	// TODO: needs testing. What happens when there are more than 1 remotes?
	return out, nil
}

// Path returns path to the Repository.
func (r *Repo) Path() string {
	return r.path
}

// cleanupFailedClone removes empty directories created by a failed git clone.
// Git itself will delete the final repo directory if a clone has failed,
// but it won't delete all the parent dirs that it created when cloning.
// eg:
// When operation like `git clone https://github.com/grdl/git-get /tmp/some/temp/dir/git-get` fails,
// git will only delete the final `git-get` dir in the path, but will leave /tmp/some/temp/dir even if it just created them.
//
// os.Remove will only delete an empty dir so we traverse the path "upwards" and delete all directories
// until a non-empty one is reached.
func cleanupFailedClone(path string) {
	for {
		path = filepath.Dir(path)
		if err := os.Remove(path); err != nil {
			return
		}
	}
}
