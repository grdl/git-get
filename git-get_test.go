package main

import (
	"github.com/libgit2/git2go/v30"
	"io/ioutil"
	"os"
	"testing"
)

func createTempRepo(t *testing.T) *git.Repository {
	dir, err := ioutil.TempDir("", "test-repo-")
	if err != nil {
		t.Fatalf("Couldn't create a temp repo directory: %s", err.Error())
	}

	t.Cleanup(func() {
		_ = os.RemoveAll(dir)
	})

	repo, err := git.InitRepository(dir, false)
	if err != nil {
		t.Fatalf("Couldn't init a temp repo: %s", err.Error())
	}
	return repo
}

func TestTempRepo(t *testing.T) {
	repo := createTempRepo(t)

	if repo.IsBare() {
		t.Errorf("Repository %s should not be bare", repo.Path())
	}
}
