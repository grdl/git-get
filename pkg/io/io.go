// Package io provides functions to read, write and search files and directories.
package io

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"syscall"

	"github.com/karrick/godirwalk"
	"github.com/pkg/errors"
)

// ErrSkipNode is used as an error indicating that .git directory has been found.
// It's handled by ErrorsCallback to tell the WalkCallback to skip this dir.
var ErrSkipNode = errors.New(".git directory found, skipping this node")

// ErrDirectoryAccess indicated a direcotry doesn't exists or can't be accessed
var ErrDirectoryAccess = errors.New("directory doesn't exist or can't be accessed")

// TempDir creates a temporary directory for test repos.
func TempDir() (string, error) {
	dir, err := ioutil.TempDir("", "git-get-repo-")
	if err != nil {
		return "", errors.Wrap(err, "failed creating test repo directory")
	}

	return dir, nil
}

// Write writes string content into a file. If file doesn't exists, it will create it.
func Write(path string, content string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return errors.Wrapf(err, "failed opening a file for writing %s", path)
	}

	_, err = file.Write([]byte(content))
	if err != nil {
		errors.Wrapf(err, "Failed writing to a file %s", path)
	}
	return nil
}

// Exists returns true if a directory exists. If it doesn't or the directory can't be accessed it returns an error.
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)

	if err == nil {
		return true, nil
	}

	if err != nil {
		if os.IsNotExist(err) {
			return false, ErrDirectoryAccess
		}
	}

	// Directory exists but can't be accessed
	return true, ErrDirectoryAccess
}

// RepoFinder finds paths to git repos inside given path.
type RepoFinder struct {
	root  string
	repos []string
}

// NewRepoFinder returns a RepoFinder pointed at given root path.
func NewRepoFinder(root string) *RepoFinder {
	return &RepoFinder{
		root: root,
	}
}

// Find returns paths to git repos found inside a given root path.
// Returns error if root repo path can't be found or accessed.
func (r *RepoFinder) Find(root string) ([]string, error) {
	if _, err := Exists(root); err != nil {
		return nil, err
	}

	walkOpts := &godirwalk.Options{
		ErrorCallback: r.errorCb,
		Callback:      r.walkCb,
		// Use Unsorted to improve speed because repos will be processed by goroutines in a random order anyway.
		Unsorted: true,
	}

	err := godirwalk.Walk(root, walkOpts)
	if err != nil {
		return nil, err
	}

	if len(r.repos) == 0 {
		return nil, fmt.Errorf("no git repos found in root path %s", root)
	}

	return r.repos, nil
}

func (r *RepoFinder) walkCb(path string, ent *godirwalk.Dirent) error {
	if ent.IsDir() && ent.Name() == ".git" {
		r.repos = append(r.repos, strings.TrimSuffix(path, ".git"))
		return ErrSkipNode
	}
	return nil
}

func (r *RepoFinder) errorCb(_ string, err error) godirwalk.ErrorAction {
	// Skip .git directory and directories we don't have permissions to access
	// TODO: Will syscall.EACCES work on windows?
	if errors.Is(err, ErrSkipNode) || errors.Is(err, syscall.EACCES) {
		return godirwalk.SkipNode
	}
	return godirwalk.Halt
}
