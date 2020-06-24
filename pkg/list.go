package pkg

import (
	"fmt"
	"git-get/pkg/cfg"
	"git-get/pkg/git"
	"git-get/pkg/io"
	"git-get/pkg/print"
	"strings"
)

var repos []string

// ListCfg provides configuration for the List command.
type ListCfg struct {
	Fetch  bool
	Output string
	Root   string
}

// List executes the "git list" command.
func List(c *ListCfg) error {
	paths, err := io.NewRepoFinder(c.Root).Find()
	if err != nil {
		return err
	}

	// TODO: we should open, fetch and read status of each repo in separate goroutine
	var repos []git.Repo
	for _, path := range paths {
		repo, err := git.Open(path)
		if err != nil {
			// TODO: how should we handle it?
			continue
		}

		if c.Fetch {
			err := repo.Fetch()
			if err != nil {
				// TODO: handle error
			}
		}

		repos = append(repos, *repo)
	}

	switch c.Output {
	case cfg.OutFlat:
		printables := make([]print.Repo, len(repos))
		for i := range repos {
			printables[i] = &repos[i]
		}
		fmt.Println(print.NewFlatPrinter().Print(printables))

	case cfg.OutTree:
		printables := make([]print.Repo, len(repos))
		for i := range repos {
			printables[i] = &repos[i]
		}
		fmt.Println(print.NewTreePrinter().Print(c.Root, printables))

	case cfg.OutDump:
		printables := make([]print.DumpRepo, len(repos))
		for i := range repos {
			printables[i] = &repos[i]
		}
		fmt.Println(print.NewDumpPrinter().Print(printables))

	default:
		return fmt.Errorf("invalid --out flag; allowed values: [%s]", strings.Join(cfg.AllowedOut, ", "))
	}

	return nil
}
