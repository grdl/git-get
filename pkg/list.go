package pkg

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/viper"

	"github.com/karrick/godirwalk"
)

// skipNode is used as an error indicating that .git directory has been found.
// It's handled by ErrorsCallback to tell the WalkCallback to skip this dir.
var skipNode = errors.New(".git directory found, skipping this node")

var repos []string

func FindRepos() ([]string, error) {
	repos = []string{}

	root := viper.GetString(KeyReposRoot)

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
	root := viper.GetString(KeyReposRoot)

	tree := BuildTree(root, repos)
	fmt.Println(RenderSmartTree(tree))
}

const (
	ColorRed    = "\033[1;31m%s\033[0m"
	ColorGreen  = "\033[1;32m%s\033[0m"
	ColorBlue   = "\033[1;34m%s\033[0m"
	ColorYellow = "\033[1;33m%s\033[0m"
)

func renderWorktreeStatus(repo *Repo) string {
	clean := true
	var status []string

	// if current branch status can't be found it's probably a detached head
	// TODO: what if current HEAD points to a tag?
	if current := repo.findCurrentBranchStatus(); current == nil {
		status = append(status, fmt.Sprintf(ColorYellow, repo.Status.CurrentBranch))
	} else {
		status = append(status, renderBranchStatus(current))
	}

	// TODO: this is ugly
	// unset clean flag to use it to render braces around worktree status and remove "ok" from branch status if it's there
	if repo.Status.HasUncommittedChanges || repo.Status.HasUntrackedFiles {
		clean = false
	}

	if !clean {
		status[len(status)-1] = strings.TrimSuffix(status[len(status)-1], StatusOk)
		status = append(status, "[")
	}

	if repo.Status.HasUntrackedFiles {
		status = append(status, fmt.Sprintf(ColorRed, StatusUntracked))
	}

	if repo.Status.HasUncommittedChanges {
		status = append(status, fmt.Sprintf(ColorRed, StatusUncommitted))
	}

	if !clean {
		status = append(status, "]")
	}

	return strings.Join(status, " ")
}

func renderBranchStatus(branch *BranchStatus) string {
	// ok indicates that the branch has upstream and is not ahead or behind it
	ok := true
	var status []string

	status = append(status, fmt.Sprintf(ColorBlue, branch.Name))

	if branch.Upstream == "" {
		ok = false
		status = append(status, fmt.Sprintf(ColorYellow, StatusNoUpstream))
	}

	if branch.NeedsPull {
		ok = false
		status = append(status, fmt.Sprintf(ColorYellow, StatusBehind))
	}

	if branch.NeedsPush {
		ok = false
		status = append(status, fmt.Sprintf(ColorYellow, StatusAhead))
	}

	if ok {
		status = append(status, fmt.Sprintf(ColorGreen, StatusOk))
	}

	return strings.Join(status, " ")
}

func (r *Repo) findCurrentBranchStatus() *BranchStatus {
	if r.Status.CurrentBranch == StatusDetached || r.Status.CurrentBranch == StatusUnknown {
		return nil
	}

	for _, b := range r.Status.Branches {
		if b.Name == r.Status.CurrentBranch {
			return b
		}
	}

	return nil
}
