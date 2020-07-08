package test

import (
	"fmt"
	"git-get/pkg/run"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func (r *Repo) init() {
	err := run.Git("init", "--quiet", r.path).AndShutUp()
	checkFatal(r.t, err)
}

// writeFile writes the content string into a file. If file doesn't exists, it will create it.
func (r *Repo) writeFile(filename string, content string) {
	path := filepath.Join(r.path, filename)

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	checkFatal(r.t, err)

	_, err = file.Write([]byte(content))
	checkFatal(r.t, err)
}

func (r *Repo) stageFile(path string) {
	err := run.Git("add", path).OnRepo(r.path).AndShutUp()
	checkFatal(r.t, err)
}

func (r *Repo) commit(msg string) {
	err := run.Git("commit", "-m", fmt.Sprintf("%q", msg), "--author=\"user <user@example.com>\"").OnRepo(r.path).AndShutUp()
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
	dir := tempDir(r.t, "")

	url := fmt.Sprintf("file://%s/.git", r.path)
	err := run.Git("clone", url, dir).AndShutUp()
	checkFatal(r.t, err)

	clone := &Repo{
		path: dir,
		t:    r.t,
	}

	return clone
}

func (r *Repo) fetch() {
	err := run.Git("fetch", "--all").OnRepo(r.path).AndShutUp()
	checkFatal(r.t, err)
}

// tempDir creates a temporary directory inside the parent dir.
// If parent is empty, it will use a system default temp dir (usually /tmp).
func tempDir(t *testing.T, parent string) string {
	dir, err := ioutil.TempDir(parent, "git-get-repo-")
	checkFatal(t, err)

	// Automatically remove temp dir when the test is over.
	t.Cleanup(func() {
		err := os.RemoveAll(dir)
		if err != nil {
			t.Errorf("failed removing test repo %s", dir)
		}
	})

	return dir
}

func checkFatal(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("failed making test repo: %+v", err)
	}
}
