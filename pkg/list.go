package pkg

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/xlab/treeprint"

	"github.com/go-git/go-git/v5"

	"github.com/karrick/godirwalk"
)

// skipNode is used as an error indicating that .git directory has been found.
// It's handled by ErrorsCallback to tell the WalkCallback to skip this dir.
var skipNode = errors.New(".git directory found, skipping this node")

var repos []string

func FindRepos() ([]string, error) {
	repos = []string{}

	root := Cfg.ReposRoot()

	if _, err := os.Stat(root); err != nil {
		return nil, fmt.Errorf("Repos root %s does not exist or can't be accessed", root)
	}

	walkOpts := &godirwalk.Options{
		ErrorCallback: ErrorCb,
		Callback:      WalkCb,
		// Use Unsorted to improve speed because repos will be processed by goroutines in a random order anyway.
		Unsorted: true,
	}

	err := godirwalk.Walk(root, walkOpts)
	if err != nil {
		return nil, err
	}

	if len(repos) == 0 {
		return nil, fmt.Errorf("No git repos found in repos root %s", root)
	}

	return repos, nil
}

func WalkCb(path string, ent *godirwalk.Dirent) error {
	if ent.IsDir() && ent.Name() == git.GitDirName {
		repos = append(repos, strings.TrimSuffix(path, git.GitDirName))
		return skipNode
	}
	return nil
}

func ErrorCb(_ string, err error) godirwalk.ErrorAction {
	if errors.Is(err, skipNode) {
		return godirwalk.SkipNode
	}
	return godirwalk.Halt
}

func OpenAll(paths []string) ([]*Repo, error) {
	var repos []*Repo
	reposChan := make(chan *Repo)

	for _, path := range paths {
		go func(path string) {
			repo, err := OpenRepo(path)

			if err != nil {
				// TODO handle error
				fmt.Println(err)
			}

			err = repo.LoadStatus()
			if err != nil {
				// TODO handle error
				fmt.Println(err)
			}
			// when error happened we just sent a nil
			reposChan <- repo
		}(path)
	}

	for repo := range reposChan {
		repos = append(repos, repo)

		// TODO: is this the right way to close the channel? What if we have non-unique paths?
		if len(repos) == len(paths) {
			close(reposChan)
		}
	}

	// sort the final array to make printing easier
	sort.Slice(repos, func(i, j int) bool {
		return strings.Compare(repos[i].path, repos[j].path) < 0
	})

	return repos, nil
}

func PrintRepos(repos []*Repo) {
	root := Cfg.ReposRoot()

	seg := make([][]string, len(repos))

	t := treeprint.New()
	t.SetValue(root)

	for i, repo := range repos {
		path := strings.TrimPrefix(repo.path, root)
		path = strings.Trim(path, string(filepath.Separator))
		subpaths := strings.Split(path, string(filepath.Separator))

		seg[i] = make([]string, len(subpaths))

		//t.AddBranch(fmt.Sprintf("\033[1;31m%s\033[0m", path))

		branch := t
		for j, sub := range subpaths {
			seg[i][j] = sub

			if i > 0 && seg[i][j] == seg[i-1][j] {
				branch = branch.FindLastNode()
				continue
			}

			value := seg[i][j]

			// if this is the last segment, it means that's the name of the repository and we need to print its status
			if j == len(seg[i])-1 {
				value = value + PrintRepoStatus(repo)
			}

			branch = branch.AddBranch(value)
		}
	}

	fmt.Println(t.String())
}

const (
	ColorRed    = "\033[1;31m%s\033[0m"
	ColorGreen  = "\033[0;32m%s\033[0m"
	ColorBlue   = "\033[1;34m%s\033[0m"
	ColorYellow = "\033[1;33m%s\033[0m"
)

func PrintRepoStatus(repo *Repo) string {
	status := fmt.Sprintf(ColorGreen, StatusOk)

	if repo.Status.HasUntrackedFiles {
		status = fmt.Sprintf(ColorRed, StatusUntracked)
	}

	if repo.Status.HasUncommittedChanges {
		status = fmt.Sprintf(ColorRed, StatusUncommitted)
	}

	return " " + status
}
