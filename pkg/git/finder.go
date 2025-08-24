package git

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Max number of concurrently running status loading workers.
const maxWorkers = 100

var errDirNoAccess = fmt.Errorf("directory can't be accessed")
var errDirNotExist = fmt.Errorf("directory doesn't exist")

// Exists returns true if a directory exists. If it doesn't or the directory can't be accessed it returns an error.
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, fmt.Errorf("can't access %s: %w", path, errDirNotExist)
	}

	// Directory exists but can't be accessed
	return true, fmt.Errorf("can't access %s: %w", path, errDirNoAccess)
}

// RepoFinder finds git repositories inside a given path and loads their status.
type RepoFinder struct {
	root       string
	repos      []*Repo
	maxWorkers int
}

// NewRepoFinder returns a RepoFinder pointed at given root path.
func NewRepoFinder(root string) *RepoFinder {
	return &RepoFinder{
		root:       root,
		maxWorkers: maxWorkers,
	}
}

// Find finds git repositories inside a given root path.
// It doesn't add repositories nested inside other git repos.
// Returns error if root repo path can't be found or accessed.
func (f *RepoFinder) Find() error {
	if _, err := Exists(f.root); err != nil {
		return fmt.Errorf("failed to access root path: %w", err)
	}

	err := filepath.WalkDir(f.root, func(path string, d fs.DirEntry, err error) error {
		// Handle walk errors
		if err != nil {
			// Skip permission errors but continue walking
			if os.IsPermission(err) {
				return nil // Skip this path but continue
			}

			return fmt.Errorf("failed to walk %s: %w", path, err)
		}

		// Only process directories
		if !d.IsDir() {
			return nil
		}

		// Case 1: We're looking at a .git directory itself
		if d.Name() == dotgit {
			parentPath := filepath.Dir(path)
			f.addIfOk(parentPath)

			return fs.SkipDir // Skip the .git directory contents
		}

		// Case 2: Check if this directory contains a .git subdirectory
		gitPath := filepath.Join(path, dotgit)
		if _, err := os.Stat(gitPath); err == nil {
			f.addIfOk(path)
			return fs.SkipDir // Skip this directory's contents since it's a repo
		}

		return nil // Continue walking
	})
	if err != nil {
		return fmt.Errorf("failed to walk directory tree: %w", err)
	}

	if len(f.repos) == 0 {
		return fmt.Errorf("no git repos found in root path %s", f.root)
	}

	return nil
}

// LoadAll loads and returns sorted slice of statuses of all repositories found by RepoFinder.
// If fetch equals true, it first fetches from the remote repo before loading the status.
// Each repo is loaded concurrently by a separate worker, with max 100 workers being active at the same time.
func (f *RepoFinder) LoadAll(fetch bool) []*Status {
	var ss []*Status

	reposChan := make(chan *Repo, f.maxWorkers)
	statusChan := make(chan *Status, f.maxWorkers)

	// Fire up workers. They listen on reposChan, load status and send the result to statusChan.
	for i := 0; i < f.maxWorkers; i++ {
		go statusWorker(fetch, reposChan, statusChan)
	}

	// Start loading the slice of repos found by finder into the reposChan.
	// It runs in a goroutine so that as soon as repos appear on the channel they can be processed and sent to statusChan.
	go loadRepos(f.repos, reposChan)

	// Read statuses from the statusChan and add then to the result slice.
	// Close the channel when all repos are loaded.
	for status := range statusChan {
		ss = append(ss, status)
		if len(ss) == len(f.repos) {
			close(statusChan)
		}
	}

	// Sort the status slice by path
	sort.Slice(ss, func(i, j int) bool {
		return strings.Compare(ss[i].path, ss[j].path) < 0
	})

	return ss
}

func loadRepos(repos []*Repo, reposChan chan<- *Repo) {
	for _, repo := range repos {
		reposChan <- repo
	}

	close(reposChan)
}

func statusWorker(fetch bool, reposChan <-chan *Repo, statusChan chan<- *Status) {
	for repo := range reposChan {
		statusChan <- repo.LoadStatus(fetch)
	}
}

// addIfOk adds the found repo to the repos slice if it can be opened.
func (f *RepoFinder) addIfOk(path string) {
	// Open() should never return an error here since we already verified the .git directory exists.
	// The path should already be the repository root (not the .git subdirectory).
	repo, err := Open(path)
	if err == nil {
		f.repos = append(f.repos, repo)
	}
}
