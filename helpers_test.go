package main

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

func cleanupRepo(t *testing.T, repo *git.Repository) {
	err := os.RemoveAll(repo.Workdir())
	if err != nil {
		t.Errorf("failed cleaning up repo")
	}
}

func newTestRepo() (*git.Repository, error) {
	dir, err := ioutil.TempDir("", "test-repo-")
	if err != nil {
		return nil, errors.Wrap(err, "failed creating a temp repo")
	}

	repo, err := git.InitRepository(dir, false)
	if err != nil {
		return nil, errors.Wrap(err, "failed initializing a temp repo")
	}

	return repo, nil
}

func createFile(repo *git.Repository, name string) error {
	err := ioutil.WriteFile(path.Join(repo.Workdir(), name), []byte("I'm a file"), 0644)
	if err != nil {
		return errors.Wrap(err, "failed writing a file")
	}

	return nil
}

func stageFile(repo *git.Repository, name string) error {
	index, err := repo.Index()
	if err != nil {
		return errors.Wrap(err, "failed getting repo index")
	}

	err = index.AddByPath(name)
	if err != nil {
		return errors.Wrap(err, "failed adding file to index")
	}

	err = index.Write()
	if err != nil {
		return errors.Wrap(err, "failed writing index")
	}

	return nil
}

func createCommit(repo *git.Repository, message string) error {
	index, err := repo.Index()
	if err != nil {
		return errors.Wrap(err, "failed creating a temp repo")
	}

	treeId, err := index.WriteTree()
	if err != nil {
		return errors.Wrap(err, "failed building tree from index")
	}

	tree, err := repo.LookupTree(treeId)
	if err != nil {
		return errors.Wrap(err, "failed looking up tree")
	}

	signature := &git.Signature{
		Name:  "Some Guy",
		Email: "someguy@example.com",
		When:  time.Date(2000, 01, 01, 16, 00, 00, 0, time.UTC),
	}

	_, err = repo.CreateCommit("HEAD", signature, signature, message, tree)
	if err != nil {
		return errors.Wrap(err, "failed creating commit")
	}

	return nil
}
