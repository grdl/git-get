package pkg

import (
	urlpkg "net/url"
	"os"
	"testing"
)

func TestFetch(t *testing.T) {
	// Create origin repo with a single commit in master
	origin := newTestRepo(t)
	createFile(t, origin, "file")
	stageFile(t, origin, "file")
	createCommit(t, origin, "Initial commit")

	// Clone the origin repo
	repoRoot := newTempDir(t)
	url, err := urlpkg.Parse(origin.Path())
	checkFatal(t, err)
	path, err := CloneRepo(url, repoRoot)
	checkFatal(t, err)

	// Open cloned repo and load its status
	repo, err := OpenRepo(path)
	checkFatal(t, err)

	// Check cloned status. It should not be behind origin
	if repo.Status.Branches["master"].Behind != 0 {
		t.Errorf("Master should not be behind")
	}

	// Add another commit to origin
	createFile(t, origin, "anotherFile")
	stageFile(t, origin, "anotherFile")
	createCommit(t, origin, "Second commit")

	// Fetch cloned repo and check the status again
	err = repo.Fetch()
	checkFatal(t, err)
	err = repo.Reload()
	checkFatal(t, err)

	// Cloned master should now be 1 commit behind origin
	if repo.Status.Branches["master"].Behind != 1 {
		t.Errorf("Master should be 1 commit behind")
	}

	if repo.Status.Branches["master"].Ahead != 0 {
		t.Errorf("Master should not be ahead")
	}
}

func TestMakeDir(t *testing.T) {
	repoRoot := newTempDir(t)
	repoPath := "github.com/grdl/git-get"

	dir, err := MakeDir(repoRoot, repoPath)
	checkFatal(t, err)

	stat, err := os.Stat(dir)
	checkFatal(t, err)

	if !stat.IsDir() {
		t.Errorf("Path is not a directory: %s", dir)
	}
}
