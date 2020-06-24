package git

import (
	"fmt"
	"git-get/pkg/file"
	"net/url"
	"os"
	"os/exec"
	"path"
	"testing"
)

// testRepo embeds testing.T into a Repo instance to simplify creation of test repos.
// Any error thrown while creating a test repo will cause a t.Fatal call.
type testRepo struct {
	*Repo
	*testing.T
}

// TODO: this should be a method of a tempDir, not a repo
// Automatically remove test repo when the test is over
func (r *testRepo) cleanup() {
	err := os.RemoveAll(r.path)
	if err != nil {
		r.T.Errorf("failed removing test repo directory %s", r.path)
	}
}

func testRepoEmpty(t *testing.T) *testRepo {
	dir, err := file.TempDir()
	checkFatal(t, err)

	r, err := Open(dir)
	checkFatal(t, err)

	tr := &testRepo{
		Repo: r,
		T:    t,
	}

	t.Cleanup(tr.cleanup)

	tr.init()
	return tr
}

func testRepoWithUntracked(t *testing.T) *testRepo {
	r := testRepoEmpty(t)
	r.writeFile("README.md", "I'm a readme file")

	return r
}

func testRepoWithStaged(t *testing.T) *testRepo {
	r := testRepoEmpty(t)
	r.writeFile("README.md", "I'm a readme file")
	r.stageFile("README.md")

	return r
}

func testRepoWithCommit(t *testing.T) *testRepo {
	r := testRepoEmpty(t)
	r.writeFile("README.md", "I'm a readme file")
	r.stageFile("README.md")
	r.commit("Initial commit")

	return r
}

func testRepoWithUncommittedAndUntracked(t *testing.T) *testRepo {
	r := testRepoEmpty(t)
	r.writeFile("README.md", "I'm a readme file")
	r.stageFile("README.md")
	r.commit("Initial commit")
	r.writeFile("README.md", "These changes won't be committed")
	r.writeFile("untracked.txt", "I'm untracked")

	return r
}

func testRepoWithBranch(t *testing.T) *testRepo {
	r := testRepoWithCommit(t)
	r.branch("feature/branch")
	r.checkout("feature/branch")

	return r
}

func testRepoWithTag(t *testing.T) *testRepo {
	r := testRepoWithCommit(t)
	r.tag("v0.0.1")
	r.checkout("v0.0.1")

	return r
}

func testRepoWithBranchWithUpstream(t *testing.T) *testRepo {
	origin := testRepoWithCommit(t)
	origin.branch("feature/branch")

	r := origin.clone()
	r.checkout("feature/branch")
	return r
}

func testRepoWithBranchWithoutUpstream(t *testing.T) *testRepo {
	origin := testRepoWithCommit(t)

	r := origin.clone()
	r.branch("feature/branch")
	r.checkout("feature/branch")
	return r
}

func testRepoWithBranchAhead(t *testing.T) *testRepo {
	origin := testRepoWithCommit(t)
	origin.branch("feature/branch")

	r := origin.clone()
	r.checkout("feature/branch")

	r.writeFile("local.new", "local.new")
	r.stageFile("local.new")
	r.commit("local.new")

	return r
}

func testRepoWithBranchBehind(t *testing.T) *testRepo {
	origin := testRepoWithCommit(t)
	origin.branch("feature/branch")
	origin.checkout("feature/branch")

	r := origin.clone()
	r.checkout("feature/branch")

	origin.writeFile("origin.new", "origin.new")
	origin.stageFile("origin.new")
	origin.commit("origin.new")

	err := r.Fetch()
	checkFatal(r.T, err)

	return r
}

// returns a repo with 2 commits ahead and 1 behind
func testRepoWithBranchAheadAndBehind(t *testing.T) *testRepo {
	origin := testRepoWithCommit(t)
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

	err := r.Fetch()
	checkFatal(r.T, err)

	return r
}

func (r *testRepo) writeFile(filename string, content string) {
	path := path.Join(r.path, filename)
	err := file.Write(path, content)
	checkFatal(r.T, err)
}

func (r *testRepo) init() {
	cmd := exec.Command("git", "init", "--quiet", r.path)
	runGitCmd(r.T, cmd)
}

func (r *testRepo) stageFile(path string) {
	cmd := gitCmd(r.path, "add", path)
	runGitCmd(r.T, cmd)
}

func (r *testRepo) commit(msg string) {
	cmd := gitCmd(r.path, "commit", "-m", msg)
	runGitCmd(r.T, cmd)
}

func (r *testRepo) branch(name string) {
	cmd := gitCmd(r.path, "branch", name)
	runGitCmd(r.T, cmd)
}

func (r *testRepo) tag(name string) {
	cmd := gitCmd(r.path, "tag", "-a", name, "-m", name)
	runGitCmd(r.T, cmd)
}

func (r *testRepo) checkout(name string) {
	cmd := gitCmd(r.path, "checkout", name)
	runGitCmd(r.T, cmd)
}

func (r *testRepo) clone() *testRepo {
	dir, err := file.TempDir()
	checkFatal(r.T, err)

	url, err := url.Parse(fmt.Sprintf("file://%s/.git", r.path))
	checkFatal(r.T, err)

	opts := &CloneOpts{
		URL:   url,
		Quiet: true,
		Path:  dir,
	}

	repo, err := Clone(opts)
	checkFatal(r.T, err)

	tr := &testRepo{
		Repo: repo,
		T:    r.T,
	}

	tr.T.Cleanup(tr.cleanup)
	return tr
}

func runGitCmd(t *testing.T, cmd *exec.Cmd) {
	err := cmd.Run()
	checkFatal(t, cmdError(cmd, err))
}

func checkFatal(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("%+v", err)
	}
}
