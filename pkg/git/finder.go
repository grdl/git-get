package git

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"syscall"

	"github.com/karrick/godirwalk"
	"github.com/pkg/errors"
)

// errSkipNode is used as an error indicating that .git directory has been found.
// It's handled by ErrorsCallback to tell the WalkCallback to skip this dir.
var errSkipNode = errors.New(".git directory found, skipping this node")

// errDirectoryAccess indicates a directory doesn't exists or can't be accessed
var errDirectoryAccess = errors.New("directory doesn't exist or can't be accessed")

// Exists returns true if a directory exists. If it doesn't or the directory can't be accessed it returns an error.
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)

	if err == nil {
		return true, nil
	}

	if err != nil {
		if os.IsNotExist(err) {
			return false, errDirectoryAccess
		}
	}

	// Directory exists but can't be accessed
	return true, errDirectoryAccess
}

// RepoFinder finds git repositories inside a given path.
type RepoFinder struct {
	root   string
	repos  []*Repo
	errors []error
}

// NewRepoFinder returns a RepoFinder pointed at given root path.
func NewRepoFinder(root string) *RepoFinder {
	return &RepoFinder{
		root: root,
	}
}

// Find finds git repositories inside a given root path.
// It doesn't add repositories nested inside other git repos.
// Returns error if root repo path can't be found or accessed.
func (f *RepoFinder) Find() error {
	if _, err := Exists(f.root); err != nil {
		return err
	}

	walkOpts := &godirwalk.Options{
		ErrorCallback: f.errorCb,
		Callback:      f.walkCb,
		// Use Unsorted to improve speed because repos will be processed by goroutines in a random order anyway.
		Unsorted: true,
	}

	err := godirwalk.Walk(f.root, walkOpts)
	if err != nil {
		return err
	}

	if len(f.repos) == 0 {
		return fmt.Errorf("no git repos found in root path %s", f.root)
	}

	return nil
}

// LoadAll loads and returns sorted slice of statuses of all repositories found by RepoFinder.
// If fetch equals true, it first fetches from the remote repo before loading the status.
// Each repo is loaded concurrently in its own goroutine, with max 100 repos being loaded at the same time.
func (f *RepoFinder) LoadAll(fetch bool) []*Status {
	var ss []*Status

	loadedChan := make(chan *Status)

	for _, repo := range f.repos {
		go func(repo *Repo) {
			loadedChan <- repo.LoadStatus(fetch)
		}(repo)
	}

	for l := range loadedChan {
		ss = append(ss, l)

		// Close the channel when all repos are loaded.
		if len(ss) == len(f.repos) {
			close(loadedChan)
		}
	}

	// Sort the status slice by path
	sort.Slice(ss, func(i, j int) bool {
		return strings.Compare(ss[i].path, ss[j].path) < 0
	})

	return ss
}

func (f *RepoFinder) walkCb(path string, ent *godirwalk.Dirent) error {
	// Do not traverse .git directories
	if ent.IsDir() && ent.Name() == dotgit {
		f.addIfOk(path)
		return errSkipNode
	}

	// Do not traverse directories containing a .git directory
	if ent.IsDir() {
		_, err := os.Stat(filepath.Join(path, dotgit))
		if err == nil {
			f.addIfOk(path)
			return errSkipNode
		}
	}
	return nil
}

// addIfOk adds the found repo to the repos slice if it can be opened.
// If repo path can't be accessed it will add an error to the errors slice.
func (f *RepoFinder) addIfOk(path string) {
	repo, err := Open(strings.TrimSuffix(path, dotgit))
	if err != nil {
		f.errors = append(f.errors, err)
		return
	}

	f.repos = append(f.repos, repo)
}

func (f *RepoFinder) errorCb(_ string, err error) godirwalk.ErrorAction {
	// Skip .git directory and directories we don't have permissions to access
	// TODO: Will syscall.EACCES work on windows?
	if errors.Is(err, errSkipNode) || errors.Is(err, syscall.EACCES) {
		return godirwalk.SkipNode
	}
	return godirwalk.Halt
}
