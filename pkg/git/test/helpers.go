// Package test contains helper utilities and functions creating pre-configured test repositories for testing purposes
package test

import (
	"fmt"
	"git-get/pkg/run"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// TempDir creates a temporary directory inside the parent dir.
// If parent is empty, it will use a system default temp dir (usually /tmp).
func TempDir(t *testing.T, parent string) string {
	t.Helper()

	// t.TempDir() is not enough in this case, we need to be able to create dirs inside the parent dir
	//nolint:usetesting
	dir, err := os.MkdirTemp(parent, "git-get-repo-")
	checkFatal(t, err)

	// Automatically remove temp dir when the test is over.
	t.Cleanup(func() {
		removeTestDir(t, dir)
	})

	return dir
}

func (r *Repo) init() {
	err := run.Git("init", "--quiet", "--initial-branch=main", r.path).AndShutUp()
	checkFatal(r.t, err)

	r.setupGitConfig()
}

// setupGitConfig sets up local git config for test repository only.
func (r *Repo) setupGitConfig() {
	err := run.Git("config", "user.name", "Test User").OnRepo(r.path).AndShutUp()
	checkFatal(r.t, err)

	err = run.Git("config", "user.email", "test@example.com").OnRepo(r.path).AndShutUp()
	checkFatal(r.t, err)
}

// writeFile writes the content string into a file. If file doesn't exists, it will create it.
func (r *Repo) writeFile(filename string, content string) {
	path := filepath.Join(r.path, filename)

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	checkFatal(r.t, err)

	_, err = file.WriteString(content)
	checkFatal(r.t, err)
}

func (r *Repo) stageFile(path string) {
	err := run.Git("add", path).OnRepo(r.path).AndShutUp()
	checkFatal(r.t, err)
}

func (r *Repo) commit(msg string) {
	err := run.Git("commit", "-m", fmt.Sprintf("%q", msg)).OnRepo(r.path).AndShutUp()
	checkFatal(r.t, err)
}

func (r *Repo) branch(name string) {
	err := run.Git("branch", name).OnRepo(r.path).AndShutUp()
	checkFatal(r.t, err)
}

func (r *Repo) tag(name string) {
	err := run.Git("tag", "-a", name, "-m", name).OnRepo(r.path).AndShutUp()
	checkFatal(r.t, err)
}

func (r *Repo) checkout(name string) {
	err := run.Git("checkout", name).OnRepo(r.path).AndShutUp()
	checkFatal(r.t, err)
}

func (r *Repo) clone() *Repo {
	dir := TempDir(r.t, "")

	url := fmt.Sprintf("file://%s/.git", r.path)
	err := run.Git("clone", url, dir).AndShutUp()
	checkFatal(r.t, err)

	clone := &Repo{
		path: dir,
		t:    r.t,
	}

	// Set up git config in the cloned repository
	clone.setupGitConfig()

	return clone
}

func (r *Repo) fetch() {
	err := run.Git("fetch", "--all").OnRepo(r.path).AndShutUp()
	checkFatal(r.t, err)
}

func checkFatal(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("failed making test repo: %+v", err)
	}
}

// removeTestDir removes a test directory.
func removeTestDir(t *testing.T, dir string) {
	t.Helper()
	// Skip cleanup on Windows to avoid file locking issues in CI
	// The CI runner environment is destroyed after tests anyway
	if runtime.GOOS == "windows" {
		return
	}

	err := os.RemoveAll(dir)
	if err != nil {
		t.Logf("warning: failed removing test repo %s: %v", dir, err)
	}
}
