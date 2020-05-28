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
	tr := NewRepoEmpty(t)
	tr.WriteFile("README", "I'm a README file")

	return tr
}

func NewRepoWithStaged(t *testing.T) *TestRepo {
	tr := NewRepoEmpty(t)
	tr.WriteFile("README", "I'm a README file")
	tr.AddFile("README")

	return tr
}
func NewRepoWithCommit(t *testing.T) *TestRepo {
	tr := NewRepoEmpty(t)
	tr.WriteFile("README", "I'm a README file")
	tr.AddFile("README")
	tr.NewCommit("Initial commit")

	return tr
}

func NewRepoWithModified(t *testing.T) *TestRepo {
	tr := NewRepoEmpty(t)
	tr.WriteFile("README", "I'm a README file")
	tr.AddFile("README")
	tr.NewCommit("Initial commit")
	tr.WriteFile("README", "I'm modified")

	return tr
}

func NewRepoWithIgnored(t *testing.T) *TestRepo {
	tr := NewRepoEmpty(t)
	tr.WriteFile(".gitignore", "ignoreme")
	tr.AddFile(".gitignore")
	tr.NewCommit("Initial commit")
	tr.WriteFile("ignoreme", "I'm being ignored")

	return tr
}

func NewRepoWithLocalBranch(t *testing.T) *TestRepo {
	tr := NewRepoWithCommit(t)
	tr.NewBranch("local")
	return tr
}

func NewRepoWithClonedBranch(t *testing.T) *TestRepo {
	origin := NewRepoWithCommit(t)

	tr := origin.Clone()
	tr.NewBranch("local")

	return tr
}

func NewRepoWithBranchAhead(t *testing.T) *TestRepo {
	origin := NewRepoWithCommit(t)

	tr := origin.Clone()
	tr.WriteFile("new", "I'm a new file")
	tr.AddFile("new")
	tr.NewCommit("New commit")

	return tr
}

func NewRepoWithBranchBehind(t *testing.T) *TestRepo {
	origin := NewRepoWithCommit(t)

	tr := origin.Clone()

	origin.WriteFile("origin.new", "I'm a new file on origin")
	origin.AddFile("origin.new")
	origin.NewCommit("New origin commit")

	tr.Fetch()
	return tr
}

func NewRepoWithBranchAheadAndBehind(t *testing.T) *TestRepo {
	origin := NewRepoWithCommit(t)

	tr := origin.Clone()
	tr.WriteFile("local.new", "I'm a new file on local")
	tr.AddFile("local.new")
	tr.NewCommit("New local commit")

	origin.WriteFile("origin.new", "I'm a new file on origin")
	origin.AddFile("origin.new")
	origin.NewCommit("New origin commit")

	tr.Fetch()
	return tr
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

func (r *TestRepo) WriteFile(name string, content string) {
	wt, err := r.Repo.Worktree()
	checkFatal(r.t, errors.Wrap(err, "Failed getting worktree"))

	file, err := wt.Filesystem.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	checkFatal(r.t, errors.Wrap(err, "Failed opening a file"))

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

func (r *TestRepo) Fetch() {
	repo := &Repo{
		repo: r.Repo,
	}

	err := repo.Fetch()
	checkFatal(r.t, err)
}

func checkFatal(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("%+v", err)
	}
}
