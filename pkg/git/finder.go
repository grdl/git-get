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

// Max number of concurrently running status loading workers.
const maxWorkers = 100

// errSkipNode is used as an error indicating that .git directory has been found.
// It's handled by ErrorsCallback to tell the WalkCallback to skip this dir.
var errSkipNode = errors.New(".git directory found, skipping this node")

var errDirNoAccess = errors.New("directory can't be accessed")
var errDirNotExist = errors.New("directory doesn't exist")

// Exists returns true if a directory exists. If it doesn't or the directory can't be accessed it returns an error.
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)

	if err == nil {
		return true, nil
	}

	if err != nil {
		if os.IsNotExist(err) {
			return false, errors.Wrapf(errDirNotExist, "can't access %s", path)
		}
	}

	// Directory exists but can't be accessed
	return true, errors.Wrapf(errDirNoAccess, "can't access %s", path)
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
func (f *RepoFinder) addIfOk(path string) {
	// TODO: is the case below really correct? What if there's a race condition and the dir becomes unaccessible between finding it and opening?

	// Open() should never return an error here. If a finder found a .git inside this dir, it means it could open and access it.
	// If the dir was unaccessible, then it would have been skipped by the check in errorCb().
	repo, err := Open(strings.TrimSuffix(path, dotgit))
	if err == nil {
		f.repos = append(f.repos, repo)
	}
}

func (f *RepoFinder) errorCb(_ string, err error) godirwalk.ErrorAction {
	// Skip .git directory and directories we don't have permissions to access
	// TODO: Will syscall.EACCES work on windows?
	if errors.Is(err, errSkipNode) || errors.Is(err, syscall.EACCES) {
		return godirwalk.SkipNode
	}
	return godirwalk.Halt
}
