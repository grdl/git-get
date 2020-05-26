package new

import (
	"testing"
	"time"

	"github.com/go-git/go-git/v5/plumbing/object"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/pkg/errors"
	//"github.com/go-git/go-git/v5"
)

func checkFatal(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("%+v", err)
	}
}

//func newTempDir(t *testing.T) string {
//	dir, err := ioutil.TempDir("", "git-get-repo-")
//	checkFatal(t, errors.Wrap(err, "Failed creating test repo directory"))
//
//	// Automatically remove repo when test is over
//	t.Cleanup(func() {
//		err := os.RemoveAll(dir)
//		if err != nil {
//			t.Errorf("failed cleaning up repo")
//		}
//	})
//
//	return dir
//}

func newTestRepo(t *testing.T) *git.Repository {
	fs := memfs.New()
	storage := memory.NewStorage()

	repo, err := git.Init(storage, fs)
	checkFatal(t, errors.Wrap(err, "Failed initializing a temp repo"))

	return repo
}

func createFile(t *testing.T, repo *git.Repository, name string) {
	wt, err := repo.Worktree()
	checkFatal(t, errors.Wrap(err, "Failed getting worktree"))

	file, err := wt.Filesystem.Create(name)
	checkFatal(t, errors.Wrap(err, "Failed creating a file"))

	_, err = file.Write([]byte("I'm a file"))
	checkFatal(t, errors.Wrap(err, "Failed writing a file"))
}

func stageFile(t *testing.T, repo *git.Repository, name string) {
	wt, err := repo.Worktree()
	checkFatal(t, errors.Wrap(err, "Failed getting worktree"))

	_, err = wt.Add(name)
	checkFatal(t, errors.Wrap(err, "Failed adding file to index"))
}

func createCommit(t *testing.T, repo *git.Repository, msg string) {
	wt, err := repo.Worktree()
	checkFatal(t, errors.Wrap(err, "Failed getting worktree"))

	opts := &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Some Guy",
			Email: "someguy@example.com",
			When:  time.Date(2000, 01, 01, 16, 00, 00, 0, time.UTC),
		},
	}

	_, err = wt.Commit(msg, opts)
	checkFatal(t, errors.Wrap(err, "Failed creating commit"))
}

//
//func createBranch(t *testing.T, repo *git.Repository, name string) *git.Branch {
//	head, err := repo.Head()
//	checkFatal(t, errors.Wrap(err, "Failed getting repo head"))
//
//	commit, err := repo.LookupCommit(head.Target())
//	checkFatal(t, errors.Wrap(err, "Failed getting commit id from head"))
//
//	branch, err := repo.CreateBranch(name, commit, false)
//	checkFatal(t, errors.Wrap(err, "Failed creating branch"))
//
//	return branch
//}
//
//func checkoutBranch(t *testing.T, repo *git.Repository, name string) {
//	branch, err := repo.LookupBranch(name, git.BranchAll)
//
//	// If branch can't be found, let's check if it's a remote branch
//	if branch == nil {
//		branch, err = repo.LookupBranch("origin/"+name, git.BranchAll)
//	}
//	checkFatal(t, errors.Wrap(err, "Failed looking up branch"))
//
//	// If branch is remote, we need to create a local one first
//	if branch.IsRemote() {
//		commit, err := repo.LookupCommit(branch.Target())
//		checkFatal(t, errors.Wrap(err, "Failed looking up commit"))
//
//		localBranch, err := repo.CreateBranch(name, commit, false)
//		checkFatal(t, errors.Wrap(err, "Failed creating local branch"))
//
//		err = localBranch.SetUpstream("origin/" + name)
//		checkFatal(t, errors.Wrap(err, "Failed setting upstream"))
//	}
//
//	err = repo.SetHead("refs/heads/" + name)
//	checkFatal(t, errors.Wrap(err, "Failed setting head"))
//
//	options := &git.CheckoutOpts{
//		Strategy: git.CheckoutForce,
//	}
//	err = repo.CheckoutHead(options)
//	checkFatal(t, errors.Wrap(err, "Failed checking out tree"))
//}
