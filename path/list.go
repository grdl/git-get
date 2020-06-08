package path

import (
	"fmt"
	"git-get/cfg"
	"git-get/git"
	"os"
	"sort"
	"strings"
	"syscall"

	"github.com/karrick/godirwalk"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// skipNode is used as an error indicating that .git directory has been found.
// It's handled by ErrorsCallback to tell the WalkCallback to skip this dir.
var skipNode = errors.New(".git directory found, skipping this node")

var repos []string

func FindRepos() ([]string, error) {
	repos = []string{}

	root := viper.GetString(cfg.KeyReposRoot)

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
	if ent.IsDir() && ent.Name() == ".git" {
		repos = append(repos, strings.TrimSuffix(path, ".git"))
		return skipNode
	}
	return nil
}

func ErrorCb(_ string, err error) godirwalk.ErrorAction {
	// Skip .git directory and directories we don't have permissions to access
	// TODO: Will syscall.EACCES work on windows?
	if errors.Is(err, skipNode) || errors.Is(err, syscall.EACCES) {
		return godirwalk.SkipNode
	}
	return godirwalk.Halt
}

func OpenAll(paths []string) ([]*git.Repo, error) {
	var repos []*git.Repo
	reposChan := make(chan *git.Repo)

	for _, path := range paths {
		go func(path string) {
			repo, err := git.OpenRepo(path)

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
		return strings.Compare(repos[i].Path, repos[j].Path) < 0
	})

	return repos, nil
}
