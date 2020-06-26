package test

import (
	"git-get/pkg/io"
	"os"
	"testing"
)

// Repo represents a test repository.
// It embeds testing.T so that any error thrown while creating a test repo will cause a t.Fatal call.
type Repo struct {
	path string
	t    *testing.T
}

func (r *Repo) Path() string {
	return r.path
}

// TODO: this should be a method of a tempDir, not a repo
// Automatically remove test repo when the test is over
func (r *Repo) cleanup() {
	err := os.RemoveAll(r.path)
	if err != nil {
		r.t.Errorf("failed removing test repo directory %s", r.path)
	}
}

func RepoEmpty(t *testing.T) *Repo {
	dir, err := io.TempDir()
	checkFatal(t, err)

	r := &Repo{
		path: dir,
		t:    t,
	}

	t.Cleanup(r.cleanup)

	r.init()
	return r
}

func RepoWithUntracked(t *testing.T) *Repo {
	r := RepoEmpty(t)
	r.writeFile("README.md", "I'm a readme file")

	return r
}

func RepoWithStaged(t *testing.T) *Repo {
	r := RepoEmpty(t)
	r.writeFile("README.md", "I'm a readme file")
	r.stageFile("README.md")

	return r
}

func RepoWithCommit(t *testing.T) *Repo {
	r := RepoEmpty(t)
	r.writeFile("README.md", "I'm a readme file")
	r.stageFile("README.md")
	r.commit("Initial commit")

	return r
}

func RepoWithUncommittedAndUntracked(t *testing.T) *Repo {
	r := RepoEmpty(t)
	r.writeFile("README.md", "I'm a readme file")
	r.stageFile("README.md")
	r.commit("Initial commit")
	r.writeFile("README.md", "These changes won't be committed")
	r.writeFile("untracked.txt", "I'm untracked")

	return r
}

func RepoWithBranch(t *testing.T) *Repo {
	r := RepoWithCommit(t)
	r.branch("feature/branch")
	r.checkout("feature/branch")

	return r
}

func RepoWithTag(t *testing.T) *Repo {
	r := RepoWithCommit(t)
	r.tag("v0.0.1")
	r.checkout("v0.0.1")

	return r
}

func RepoWithBranchWithUpstream(t *testing.T) *Repo {
	origin := RepoWithCommit(t)
	origin.branch("feature/branch")

	r := origin.clone()
	r.checkout("feature/branch")
	return r
}

func RepoWithBranchWithoutUpstream(t *testing.T) *Repo {
	origin := RepoWithCommit(t)

	r := origin.clone()
	r.branch("feature/branch")
	r.checkout("feature/branch")
	return r
}

func RepoWithBranchAhead(t *testing.T) *Repo {
	origin := RepoWithCommit(t)
	origin.branch("feature/branch")

	r := origin.clone()
	r.checkout("feature/branch")

	r.writeFile("local.new", "local.new")
	r.stageFile("local.new")
	r.commit("local.new")

	return r
}

func RepoWithBranchBehind(t *testing.T) *Repo {
	origin := RepoWithCommit(t)
	origin.branch("feature/branch")
	origin.checkout("feature/branch")

	r := origin.clone()
	r.checkout("feature/branch")

	origin.writeFile("origin.new", "origin.new")
	origin.stageFile("origin.new")
	origin.commit("origin.new")

	r.fetch()

	return r
}

// returns a repo with 2 commits ahead and 1 behind
func RepoWithBranchAheadAndBehind(t *testing.T) *Repo {
	origin := RepoWithCommit(t)
	origin.branch("feature/branch")
	origin.checkout("feature/branch")

	r := origin.clone()
	r.checkout("feature/branch")

	origin.writeFile("origin.new", "origin.new")
	origin.stageFile("origin.new")
	origin.commit("origin.new")

	r.writeFile("local.new", "local.new")
	r.stageFile("local.new")
	r.commit("local.new")

	r.writeFile("local.new2", "local.new2")
	r.stageFile("local.new2")
	r.commit("local.new2")

	r.fetch()

	return r
}
