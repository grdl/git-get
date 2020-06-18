package repo

import (
	"net/url"

	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/go-git/go-git/v5/plumbing"

	"github.com/go-git/go-git/v5/plumbing/object"

	"github.com/go-git/go-git/v5"
	"github.com/pkg/errors"
)

const (
	testUser  = "Test User"
	testEmail = "testuser@example.com"
)

func newRepoEmpty(t *testing.T) *Repo {
	dir := newTempDir(t)

	repo, err := git.PlainInit(dir, false)
	checkFatal(t, err)

	return New(repo, dir)
}

func newRepoWithUntracked(t *testing.T) *Repo {
	r := newRepoEmpty(t)
	r.writeFile(t, "README", "I'm a README file")

	return r
}

func newRepoWithStaged(t *testing.T) *Repo {
	r := newRepoEmpty(t)
	r.writeFile(t, "README", "I'm a README file")
	r.addFile(t, "README")

	return r
}

func newRepoWithCommit(t *testing.T) *Repo {
	r := newRepoEmpty(t)
	r.writeFile(t, "README", "I'm a README file")
	r.addFile(t, "README")
	r.newCommit(t, "Initial commit")

	return r
}

func newRepoWithModified(t *testing.T) *Repo {
	r := newRepoEmpty(t)
	r.writeFile(t, "README", "I'm a README file")
	r.addFile(t, "README")
	r.newCommit(t, "Initial commit")
	r.writeFile(t, "README", "I'm modified")

	return r
}

func newRepoWithIgnored(t *testing.T) *Repo {
	r := newRepoEmpty(t)
	r.writeFile(t, ".gitignore", "ignoreme")
	r.addFile(t, ".gitignore")
	r.newCommit(t, "Initial commit")
	r.writeFile(t, "ignoreme", "I'm being ignored")

	return r
}

func newRepoWithLocalBranch(t *testing.T) *Repo {
	r := newRepoWithCommit(t)
	r.newBranch(t, "local")
	return r
}

func newRepoWithClonedBranch(t *testing.T) *Repo {
	origin := newRepoWithCommit(t)

	r := origin.clone(t, "master")
	r.newBranch(t, "local")
	r.checkoutBranch(t, "local")

	return r
}

func newRepoWithDetachedHead(t *testing.T) *Repo {
	r := newRepoWithCommit(t)

	r.writeFile(t, "new", "I'm a new file")
	r.addFile(t, "new")
	hash := r.newCommit(t, "new commit")

	r.checkoutHash(t, hash)

	return r
}

func newRepoWithBranchAhead(t *testing.T) *Repo {
	origin := newRepoWithCommit(t)

	r := origin.clone(t, "master")
	r.writeFile(t, "new", "I'm a new file")
	r.addFile(t, "new")
	r.newCommit(t, "new commit")

	return r
}

func newRepoWithBranchBehind(t *testing.T) *Repo {
	origin := newRepoWithCommit(t)

	r := origin.clone(t, "master")

	origin.writeFile(t, "origin.new", "I'm a new file on origin")
	origin.addFile(t, "origin.new")
	origin.newCommit(t, "new origin commit")

	r.fetch(t)
	return r
}

// generate repo with 2 commits ahead and 3 behind the origin
func newRepoWithBranchAheadAndBehind(t *testing.T) *Repo {
	origin := newRepoWithCommit(t)

	r := origin.clone(t, "master")
	r.writeFile(t, "local.new", "local 1")
	r.addFile(t, "local.new")
	r.newCommit(t, "1st local commit")

	r.writeFile(t, "local.new", "local 2")
	r.addFile(t, "local.new")
	r.newCommit(t, "2nd local commit")

	origin.writeFile(t, "origin.new", "origin 1")
	origin.addFile(t, "origin.new")
	origin.newCommit(t, "1st origin commit")

	origin.writeFile(t, "origin.new", "origin 2")
	origin.addFile(t, "origin.new")
	origin.newCommit(t, "2nd origin commit")

	origin.writeFile(t, "origin.new", "origin 3")
	origin.addFile(t, "origin.new")
	origin.newCommit(t, "3rd origin commit")

	r.fetch(t)
	return r
}

func newRepoWithCheckedOutBranch(t *testing.T) *Repo {
	origin := newRepoWithCommit(t)
	origin.newBranch(t, "feature/branch1")

	r := origin.clone(t, "feature/branch1")
	return r
}

func newRepoWithCheckedOutTag(t *testing.T) *Repo {
	origin := newRepoWithCommit(t)
	origin.newTag(t, "v1.0.0")

	r := origin.clone(t, "refs/tags/v1.0.0")
	return r
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

func (r *Repo) writeFile(t *testing.T, name string, content string) {
	wt, err := r.Worktree()
	checkFatal(t, errors.Wrap(err, "Failed getting worktree"))

	file, err := wt.Filesystem.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	checkFatal(t, errors.Wrap(err, "Failed opening a file"))

	_, err = file.Write([]byte(content))
	checkFatal(t, errors.Wrap(err, "Failed writing a file"))
}

func (r *Repo) addFile(t *testing.T, name string) {
	wt, err := r.Worktree()
	checkFatal(t, errors.Wrap(err, "Failed getting worktree"))

	_, err = wt.Add(name)
	checkFatal(t, errors.Wrap(err, "Failed adding file to index"))
}

func (r *Repo) newCommit(t *testing.T, msg string) plumbing.Hash {
	wt, err := r.Worktree()
	checkFatal(t, errors.Wrap(err, "Failed getting worktree"))

	opts := &git.CommitOptions{
		Author: &object.Signature{
			Name:  testUser,
			Email: testEmail,
			When:  time.Date(2000, 01, 01, 16, 00, 00, 0, time.UTC),
		},
	}

	hash, err := wt.Commit(msg, opts)
	checkFatal(t, errors.Wrap(err, "Failed creating commit"))
	return hash
}

func (r *Repo) newBranch(t *testing.T, name string) {
	head, err := r.Head()
	checkFatal(t, err)

	ref := plumbing.NewHashReference(plumbing.NewBranchReferenceName(name), head.Hash())

	err = r.Storer.SetReference(ref)
	checkFatal(t, err)
}

func (r *Repo) newTag(t *testing.T, name string) {
	head, err := r.Head()
	checkFatal(t, err)

	ref := plumbing.NewHashReference(plumbing.NewTagReferenceName(name), head.Hash())

	err = r.Storer.SetReference(ref)
	checkFatal(t, err)
}

func (r *Repo) clone(t *testing.T, branch string) *Repo {
	dir := newTempDir(t)
	repoURL, err := url.Parse("file://" + r.Path)
	checkFatal(t, err)

	cloneOpts := &CloneOpts{
		URL:    repoURL,
		Path:   dir,
		Branch: branch,
		Quiet:  true,
	}

	repo, err := Clone(cloneOpts)
	checkFatal(t, err)

	return repo
}

func (r *Repo) fetch(t *testing.T) {
	err := r.Fetch()
	checkFatal(t, err)
}

func (r *Repo) checkoutBranch(t *testing.T, name string) {
	wt, err := r.Worktree()
	checkFatal(t, errors.Wrap(err, "Failed getting worktree"))

	opts := &git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(name),
	}
	err = wt.Checkout(opts)
	checkFatal(t, errors.Wrap(err, "Failed checking out branch"))
}

func (r *Repo) checkoutHash(t *testing.T, hash plumbing.Hash) {
	wt, err := r.Worktree()
	checkFatal(t, errors.Wrap(err, "Failed getting worktree"))

	opts := &git.CheckoutOptions{
		Hash: hash,
	}
	err = wt.Checkout(opts)
	checkFatal(t, errors.Wrap(err, "Failed checking out hash"))
}

func checkFatal(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("%+v", err)
	}
}
