package pkg

import (
	"fmt"
	"git-get/pkg/cfg"
	"git-get/pkg/git"
	"git-get/pkg/print"
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
	finder := git.NewRepoFinder(c.Root)
	if err := finder.Find(); err != nil {
		return err
	}

	statuses := finder.LoadAll(c.Fetch)
	printables := make([]print.Printable, len(statuses))
	for i := range statuses {
		printables[i] = statuses[i]
	}

	switch c.Output {
	case cfg.OutFlat:
		fmt.Print(print.NewFlatPrinter().Print(printables))
	case cfg.OutTree:
		fmt.Print(print.NewTreePrinter().Print(c.Root, printables))
	case cfg.OutDump:
		fmt.Print(print.NewDumpPrinter().Print(printables))
	default:
		return fmt.Errorf("invalid --out flag; allowed values: [%s]", strings.Join(cfg.AllowedOut, ", "))
	}

	return nil
}
