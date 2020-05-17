package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"

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
		t.Errorf("failed cleaning up repo %s", err.Error())
	}
}

func newTempRepo() (*git.Repository, error) {
	dir, err := ioutil.TempDir("", "test-repo-")
	if err != nil {
		return nil, fmt.Errorf("failed creating a temp repo %s", err.Error())
	}

	repo, err := git.InitRepository(dir, false)
	if err != nil {
		return nil, fmt.Errorf("failed initializing a temp repo %s", err.Error())
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
		return nil, fmt.Errorf("failed writing a file %s", err.Error())
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
		return nil, fmt.Errorf("failed getting repo index %s", err.Error())
	}

	err = index.AddByPath(ReadmeFile)
	if err != nil {
		return nil, fmt.Errorf("failed adding file to index %s", err.Error())
	}

	err = index.Write()
	if err != nil {
		return nil, fmt.Errorf("failed writing index %s", err.Error())
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
		return nil, fmt.Errorf("failed creating a temp repo %s", err.Error())
	}

	treeId, err := index.WriteTree()
	if err != nil {
		return nil, fmt.Errorf("failed building tree from index %s", err.Error())
	}

	tree, err := repo.LookupTree(treeId)
	if err != nil {
		return nil, fmt.Errorf("failed looking up tree %s", err.Error())
	}

	signature := &git.Signature{
		Name:  CommitterName,
		Email: CommitterEmail,
		When:  time.Date(2000, 01, 01, 16, 00, 00, 0, time.UTC),
	}
	message := "Initial commit"

	_, err = repo.CreateCommit("HEAD", signature, signature, message, tree)
	if err != nil {
		return nil, fmt.Errorf("failed creating commit %s", err.Error())
	}

	return repo, nil
}
func TestStatus(t *testing.T) {
	repo, err := git.OpenRepository("/tmp/testgit")
	if err != nil {
		t.Fatalf("error: %s", err.Error())
	}

	defer cleanupRepo(t, repo)

	opts := &git.StatusOptions{
		Show:  git.StatusShowIndexAndWorkdir,
		Flags: git.StatusOptIncludeUntracked | git.StatusOptIncludeIgnored,
	}

	status, _ := repo.StatusList(opts)
	entryCount, _ := status.EntryCount()
	for i := 0; i < entryCount; i++ {
		entry, _ := status.ByIndex(i)
		fmt.Println(entry)
	}
}

func TestCreatingRepoWithCommit(t *testing.T) {
	repo, err := newTempRepoWithCommit()
	if err != nil {
		t.Fatalf("failed creating test repository %s", err.Error())
	}
	defer cleanupRepo(t, repo)
	fmt.Println(repo.Path())
}
