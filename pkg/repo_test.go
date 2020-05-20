package pkg

import "testing"

func TestFetch(t *testing.T) {
	// Create origin repo with a single commit in master
	origin := newTestRepo(t)
	createFile(t, origin, "file")
	stageFile(t, origin, "file")
	createCommit(t, origin, "Initial commit")

	// Clone the origin repo
	dir := newTempDir(t)
	err := CloneRepo(origin.Path(), dir)
	checkFatal(t, err)

	// Open cloned repo and load its status
	repo, err := OpenRepo(dir)
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
