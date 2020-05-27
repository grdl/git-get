package new

import (
	"io/ioutil"
	pkgurl "net/url"
	"os"
	"testing"
	"time"

	"github.com/go-git/go-git/v5/plumbing"

	"github.com/go-git/go-git/v5/plumbing/object"

	"github.com/go-git/go-git/v5"
	"github.com/pkg/errors"
)

type TestRepo struct {
	Repo *git.Repository
	Path string
	URL  *pkgurl.URL
	t    *testing.T
}

func NewRepoEmpty(t *testing.T) *TestRepo {
	dir := NewTempDir(t)

	repo, err := git.PlainInit(dir, false)
	checkFatal(t, err)

	url, err := ParseURL("file://" + dir)
	checkFatal(t, err)

	return &TestRepo{
		Repo: repo,
		Path: dir,
		URL:  url,
		t:    t,
	}
}

func NewRepoWithUntracked(t *testing.T) *TestRepo {
	repo := NewRepoEmpty(t)
	repo.NewFile("README", "I'm a README file")

	return repo
}

func NewRepoWithStaged(t *testing.T) *TestRepo {
	repo := NewRepoEmpty(t)
	repo.NewFile("README", "I'm a README file")
	repo.AddFile("README")

	return repo
}
func NewRepoWithCommit(t *testing.T) *TestRepo {
	repo := NewRepoEmpty(t)
	repo.NewFile("README", "I'm a README file")
	repo.AddFile("README")
	repo.NewCommit("Initial commit")

	return repo
}

func NewTempDir(t *testing.T) string {
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

func (r *TestRepo) NewFile(name string, content string) {
	wt, err := r.Repo.Worktree()
	checkFatal(r.t, errors.Wrap(err, "Failed getting worktree"))

	file, err := wt.Filesystem.Create(name)
	checkFatal(r.t, errors.Wrap(err, "Failed creating a file"))

	_, err = file.Write([]byte(content))
	checkFatal(r.t, errors.Wrap(err, "Failed writing a file"))
}

func (r *TestRepo) AddFile(name string) {
	wt, err := r.Repo.Worktree()
	checkFatal(r.t, errors.Wrap(err, "Failed getting worktree"))

	_, err = wt.Add(name)
	checkFatal(r.t, errors.Wrap(err, "Failed adding file to index"))
}

func (r *TestRepo) NewCommit(msg string) {
	wt, err := r.Repo.Worktree()
	checkFatal(r.t, errors.Wrap(err, "Failed getting worktree"))

	opts := &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Some Guy",
			Email: "someguy@example.com",
			When:  time.Date(2000, 01, 01, 16, 00, 00, 0, time.UTC),
		},
	}

	_, err = wt.Commit(msg, opts)
	checkFatal(r.t, errors.Wrap(err, "Failed creating commit"))
}

func (r *TestRepo) NewBranch(name string) {
	head, err := r.Repo.Head()
	checkFatal(r.t, err)

	ref := plumbing.NewHashReference(plumbing.NewBranchReferenceName(name), head.Hash())

	err = r.Repo.Storer.SetReference(ref)
	checkFatal(r.t, err)
}

func (r *TestRepo) Clone() *TestRepo {
	dir := NewTempDir(r.t)

	repo, err := CloneRepo(r.URL, dir, true)
	checkFatal(r.t, err)

	url, err := ParseURL("file://" + dir)
	checkFatal(r.t, err)

	return &TestRepo{
		Repo: repo.repo,
		Path: dir,
		URL:  url,
		t:    r.t,
	}
}

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

func checkFatal(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("%+v", err)
	}
}
