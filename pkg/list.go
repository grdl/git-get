package pkg

import (
	"errors"
	"fmt"
	"strings"

	"github.com/grdl/git-get/pkg/cfg"
	"github.com/grdl/git-get/pkg/git"
	"github.com/grdl/git-get/pkg/out"
)

var ErrInvalidOutput = errors.New("invalid output format")

// ListCfg provides configuration for the List command.
type ListCfg struct {
	Fetch  bool
	Output string
	Root   string
}

// List executes the "git list" command.
func List(conf *ListCfg) error {
	finder := git.NewRepoFinder(conf.Root)
	if err := finder.Find(); err != nil {
		return err
	}

	statuses := finder.LoadAll(conf.Fetch)

	printables := make([]out.Printable, len(statuses))

	for i := range statuses {
		printables[i] = statuses[i]
	}

	switch conf.Output {
	case cfg.OutFlat:
		fmt.Print(out.NewFlatPrinter().Print(printables))
	case cfg.OutTree:
		fmt.Print(out.NewTreePrinter().Print(conf.Root, printables))
	case cfg.OutDump:
		fmt.Print(out.NewDumpPrinter().Print(printables))
	default:
		return fmt.Errorf("%w, allowed values: [%s]", ErrInvalidOutput, strings.Join(cfg.AllowedOut, ", "))
	}

	return nil
}
