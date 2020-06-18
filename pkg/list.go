package pkg

import (
	"fmt"
	"git-get/pkg/cfg"
	"git-get/pkg/print"
	"git-get/pkg/repo"
	"os"
	"sort"
	"strings"
	"syscall"

	"github.com/karrick/godirwalk"
	"github.com/pkg/errors"
)

// errSkipNode is used as an error indicating that .git directory has been found.
// It's handled by ErrorsCallback to tell the WalkCallback to skip this dir.
var errSkipNode = errors.New(".git directory found, skipping this node")

var repos []string

// ListCfg provides configuration for the List command.
type ListCfg struct {
	Fetch      bool
	Output     string
	PrivateKey string
	Root       string
}

// List executes the "git list" command.
func List(c *ListCfg) error {
	paths, err := findRepos(c.Root)
	if err != nil {
		return err
	}

	repos, err := openAll(paths)
	if err != nil {
		return err
	}

	var printer print.Printer
	switch c.Output {
	case cfg.OutFlat:
		printer = &print.FlatPrinter{}
	case cfg.OutTree:
		printer = &print.TreePrinter{}
	case cfg.OutSmart:
		printer = &print.SmartPrinter{}
	case cfg.OutDump:
		printer = &print.DumpPrinter{}
	default:
		return fmt.Errorf("invalid --out flag; allowed values: %v", []string{cfg.OutFlat, cfg.OutTree, cfg.OutSmart})
	}

	fmt.Println(printer.Print(c.Root, repos))
	return nil
}

func findRepos(root string) ([]string, error) {
	repos = []string{}

	if _, err := os.Stat(root); err != nil {
		return nil, fmt.Errorf("Repos root %s does not exist or can't be accessed", root)
	}

	walkOpts := &godirwalk.Options{
		ErrorCallback: errorCb,
		Callback:      walkCb,
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

func walkCb(path string, ent *godirwalk.Dirent) error {
	if ent.IsDir() && ent.Name() == ".git" {
		repos = append(repos, strings.TrimSuffix(path, ".git"))
		return errSkipNode
	}
	return nil
}

func errorCb(_ string, err error) godirwalk.ErrorAction {
	// Skip .git directory and directories we don't have permissions to access
	// TODO: Will syscall.EACCES work on windows?
	if errors.Is(err, errSkipNode) || errors.Is(err, syscall.EACCES) {
		return godirwalk.SkipNode
	}
	return godirwalk.Halt
}

func openAll(paths []string) ([]*repo.Repo, error) {
	var repos []*repo.Repo
	reposChan := make(chan *repo.Repo)

	for _, path := range paths {
		go func(path string) {
			repo, err := repo.Open(path)

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
