package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"

	"github.com/pkg/errors"

	git "github.com/libgit2/git2go/v30"
)

const (
	ReadmeFile     = "README.md"
	ReadmeContent  = "I'm a readme file\n"
	CommitterName  = "Some Guy"
	CommitterEmail = "someguy@example.com"
)

func cleanupRepo(t *testing.T, repo *git.Repository) {
	err := os.RemoveAll(repo.Workdir())
	if err != nil {
		t.Errorf("failed cleaning up repo")
	}
}

func newTempRepo() (*git.Repository, error) {
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

func newTempRepoWithUntracked() (*git.Repository, error) {
	repo, err := newTempRepo()
	if err != nil {
		return nil, err
	}

	err = ioutil.WriteFile(path.Join(repo.Workdir(), ReadmeFile), []byte(ReadmeContent), 0644)
	if err != nil {
		return nil, errors.Wrap(err, "failed writing a file")
	}

	return repo, nil
}

func newTempRepoWithStaged() (*git.Repository, error) {
	repo, err := newTempRepoWithUntracked()
	if err != nil {
		return nil, err
	}

	index, err := repo.Index()
	if err != nil {
		return nil, errors.Wrap(err, "failed getting repo index")
	}

	err = index.AddByPath(ReadmeFile)
	if err != nil {
		return nil, errors.Wrap(err, "failed adding file to index")
	}

	err = index.Write()
	if err != nil {
		return nil, errors.Wrap(err, "failed writing index")
	}

	return repo, nil
}

func newTempRepoWithCommit() (*git.Repository, error) {
	repo, err := newTempRepoWithStaged()
	if err != nil {
		return nil, err
	}

	index, err := repo.Index()
	if err != nil {
		return nil, errors.Wrap(err, "failed creating a temp repo")
	}

	treeId, err := index.WriteTree()
	if err != nil {
		return nil, errors.Wrap(err, "failed building tree from index")
	}

	tree, err := repo.LookupTree(treeId)
	if err != nil {
		return nil, errors.Wrap(err, "failed looking up tree")
	}

	signature := &git.Signature{
		Name:  CommitterName,
		Email: CommitterEmail,
		When:  time.Date(2000, 01, 01, 16, 00, 00, 0, time.UTC),
	}
	message := "Initial commit"

	_, err = repo.CreateCommit("HEAD", signature, signature, message, tree)
	if err != nil {
		return nil, errors.Wrap(err, "failed creating commit")
	}

	return repo, nil
}

func TestCreatingRepoWithCommit(t *testing.T) {
	repo, err := newTempRepoWithCommit()
	if err != nil {
		t.Fatalf("failed creating test repository")
	}
	defer cleanupRepo(t, repo)
	fmt.Println(repo.Path())
}
