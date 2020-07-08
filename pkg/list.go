package pkg

import (
	"fmt"
	"git-get/pkg/cfg"
	"git-get/pkg/git"
	"git-get/pkg/print"
	"sort"
	"strings"
)

// ListCfg provides configuration for the List command.
type ListCfg struct {
	Fetch  bool
	Output string
	Root   string
}

// List executes the "git list" command.
func List(c *ListCfg) error {
	paths, err := git.NewRepoFinder(c.Root).Find()
	if err != nil {
		return err
	}

	loaded := loadAll(paths, c.Fetch)

	printables := make([]print.Printable, len(loaded))
	for i := range loaded {
		printables[i] = loaded[i]
	}

	switch c.Output {
	case cfg.OutFlat:
		fmt.Println(print.NewFlatPrinter().Print(printables))
	case cfg.OutTree:
		fmt.Println(print.NewTreePrinter().Print(c.Root, printables))
	case cfg.OutDump:
		fmt.Println(print.NewDumpPrinter().Print(printables))
	default:
		return fmt.Errorf("invalid --out flag; allowed values: [%s]", strings.Join(cfg.AllowedOut, ", "))
	}

	return nil
}

// loadAll runs a separate goroutine to open, fetch (if asked to) and load status of git repo
func loadAll(paths []string, fetch bool) []*Loaded {
	var ll []*Loaded

	loadedChan := make(chan *Loaded)

	for _, path := range paths {
		go func(path string) {

			loadedChan <- Load(path, fetch)
		}(path)
	}

	for l := range loadedChan {
		ll = append(ll, l)

		// Close the channell when loaded all paths
		if len(ll) == len(paths) {
			close(loadedChan)
		}
	}

	// sort the loaded slice by path
	sort.Slice(ll, func(i, j int) bool {
		return strings.Compare(ll[i].path, ll[j].path) < 0
	})

	return ll
}
