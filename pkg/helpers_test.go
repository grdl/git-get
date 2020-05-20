package pkg

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"

	"github.com/pkg/errors"

	git "github.com/libgit2/git2go/v30"
)

func checkFatal(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("%+v", err)
	}
}

func newTempDir(t *testing.T) string {
	dir, err := ioutil.TempDir("", "git-get-repo-")
	checkFatal(t, errors.Wrap(err, "Failed creating test repo directory"))

	// Automatically remove repo when test is over
	t.Cleanup(func() {
		err := os.RemoveAll(dir)
		if err != nil {
			t.Errorf("failed cleaning up repo")
		}
	})

	return dir
}

func newTestRepo(t *testing.T) *git.Repository {
	dir := newTempDir(t)
	repo, err := git.InitRepository(dir, false)
	checkFatal(t, errors.Wrap(err, "Failed initializing a temp repo"))

	return repo
}

func createFile(t *testing.T, repo *git.Repository, name string) {
	err := ioutil.WriteFile(path.Join(repo.Workdir(), name), []byte("I'm a file"), 0644)
	checkFatal(t, errors.Wrap(err, "Failed writing a file"))
}

func stageFile(t *testing.T, repo *git.Repository, name string) {
	index, err := repo.Index()
	checkFatal(t, errors.Wrap(err, "Failed getting repo index"))

	err = index.AddByPath(name)
	checkFatal(t, errors.Wrap(err, "Failed adding file to index"))

	err = index.Write()
	checkFatal(t, errors.Wrap(err, "Failed writing index"))
}

func createCommit(t *testing.T, repo *git.Repository, message string) *git.Commit {
	index, err := repo.Index()
	checkFatal(t, errors.Wrap(err, "Failed getting repo index"))

	treeId, err := index.WriteTree()
	checkFatal(t, errors.Wrap(err, "Failed building tree from index"))

	tree, err := repo.LookupTree(treeId)
	checkFatal(t, errors.Wrap(err, "Failed looking up tree id"))

	signature := &git.Signature{
		Name:  "Some Guy",
		Email: "someguy@example.com",
		When:  time.Date(2000, 01, 01, 16, 00, 00, 0, time.UTC),
	}

	empty, err := repo.IsEmpty()
	checkFatal(t, errors.Wrap(err, "Failed checking if repo is empty"))

	var commitId *git.Oid
	if !empty {
		currentBranch, err := repo.Head()
		checkFatal(t, errors.Wrap(err, "Failed getting current branch"))

		currentTip, err := repo.LookupCommit(currentBranch.Target())
		checkFatal(t, errors.Wrap(err, "Failed getting current tip"))

		commitId, err = repo.CreateCommit("HEAD", signature, signature, message, tree, currentTip)
	} else {
		commitId, err = repo.CreateCommit("HEAD", signature, signature, message, tree)
	}

	commit, err := repo.LookupCommit(commitId)
	checkFatal(t, errors.Wrap(err, "Failed looking up a commit"))

	return commit
}

func createBranch(t *testing.T, repo *git.Repository, name string) *git.Branch {
	head, err := repo.Head()
	checkFatal(t, errors.Wrap(err, "Failed getting repo head"))

	commit, err := repo.LookupCommit(head.Target())
	checkFatal(t, errors.Wrap(err, "Failed getting commit id from head"))

	branch, err := repo.CreateBranch(name, commit, false)
	checkFatal(t, errors.Wrap(err, "Failed creating branch"))

	return branch
}

func checkoutBranch(t *testing.T, repo *git.Repository, name string) {
	branch, err := repo.LookupBranch(name, git.BranchAll)

	// If branch can't be found, let's check if it's a remote branch
	if branch == nil {
		branch, err = repo.LookupBranch("origin/"+name, git.BranchAll)
	}
	checkFatal(t, errors.Wrap(err, "Failed looking up branch"))

	// If branch is remote, we need to create a local one first
	if branch.IsRemote() {
		commit, err := repo.LookupCommit(branch.Target())
		checkFatal(t, errors.Wrap(err, "Failed looking up commit"))

		localBranch, err := repo.CreateBranch(name, commit, false)
		checkFatal(t, errors.Wrap(err, "Failed creating local branch"))

		err = localBranch.SetUpstream("origin/" + name)
		checkFatal(t, errors.Wrap(err, "Failed setting upstream"))
	}

	err = repo.SetHead("refs/heads/" + name)
	checkFatal(t, errors.Wrap(err, "Failed setting head"))

	options := &git.CheckoutOpts{
		Strategy: git.CheckoutForce,
	}
	err = repo.CheckoutHead(options)
	checkFatal(t, errors.Wrap(err, "Failed checking out tree"))
}
