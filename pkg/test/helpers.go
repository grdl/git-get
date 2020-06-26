package test

import (
	"fmt"
	"git-get/pkg/io"
	"git-get/pkg/run"
	"path"
	"testing"
)

func (r *Repo) writeFile(filename string, content string) {
	path := path.Join(r.path, filename)
	err := io.Write(path, content)
	checkFatal(r.t, err)
}

func (r *Repo) init() {
	err := run.Git("init", "--quiet", r.path).AndShutUp()
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
	dir, err := io.TempDir()
	checkFatal(r.t, err)

	url := fmt.Sprintf("file://%s/.git", r.path)
	err = run.Git("clone", url, dir).AndShutUp()
	checkFatal(r.t, err)

	clone := &Repo{
		path: dir,
		t:    r.t,
	}

	clone.t.Cleanup(r.cleanup)
	return clone
}

func (r *Repo) fetch() {
	err := run.Git("fetch", "--all").OnRepo(r.path).AndShutUp()
	checkFatal(r.t, err)
}

func checkFatal(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("failed making test repo: %+v", err)
	}
}
