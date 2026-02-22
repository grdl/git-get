package pkg

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/grdl/git-get/pkg/git"
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/sahilm/fuzzy"
)

var (
	ErrNoMatchFound         = errors.New("no repository matching found")
	ErrMissingQueryTerminal = errors.New("missing query for non-interactive terminal")
)

// CdCfg provides configuration for the Cd command.
type CdCfg struct {
	Roots []string
	Query string
}

// Cd executes the "git cd" command.
func Cd(conf *CdCfg) error {
	finder := git.NewRepoFinder(conf.Roots)
	if err := finder.Find(); err != nil {
		return err
	}

	repos := finder.LoadAll(false)

	absPaths := make([]string, len(repos))
	for i, r := range repos {
		absPaths[i] = r.Path()
	}

	searchPaths := getRelativePaths(conf.Roots, absPaths)

	var selectedPath string

	var err error

	if conf.Query != "" {
		selectedPath, err = findSelectedPathWithQuery(conf.Query, searchPaths, absPaths)
	} else {
		selectedPath, err = findSelectedPathNoQuery(absPaths)
	}

	if err != nil {
		return err
	}

	if isStdoutTerminal() {
		return startSubshell(selectedPath)
	}

	fmt.Println(selectedPath)

	return nil
}

// Cmd interface defines methods for exec.Cmd that we use.
type Cmd interface {
	Run() error
	SetDir(dir string)
	SetStdin(stdin *os.File)
	SetStdout(stdout *os.File)
	SetStderr(stderr *os.File)
}

type commander interface {
	CommandContext(ctx context.Context, name string, arg ...string) Cmd
}

type osCmd struct {
	*exec.Cmd
}

func (o *osCmd) SetDir(dir string) {
	o.Dir = dir
}

func (o *osCmd) SetStdin(stdin *os.File) {
	o.Stdin = stdin
}

func (o *osCmd) SetStdout(stdout *os.File) {
	o.Stdout = stdout
}

func (o *osCmd) SetStderr(stderr *os.File) {
	o.Stderr = stderr
}

type osCommander struct{}

//nolint:ireturn
func (osCommander) CommandContext(ctx context.Context, name string, arg ...string) Cmd {
	return &osCmd{exec.CommandContext(ctx, name, arg...)}
}

var defaultCommander commander = &osCommander{}

func startSubshell(path string) error {
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/sh"
	}

	cmd := defaultCommander.CommandContext(context.Background(), shell)
	cmd.SetDir(path)
	cmd.SetStdin(os.Stdin)
	cmd.SetStdout(os.Stdout)
	cmd.SetStderr(os.Stderr)

	fmt.Fprintf(os.Stderr, "Starting subshell in %s (type 'exit' to return)\n", path)

	return cmd.Run()
}

func isStdoutTerminal() bool {
	fileInfo, _ := os.Stdout.Stat()

	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

func findSelectedPathWithQuery(query string, searchPaths []string, absPaths []string) (string, error) {
	matches := fuzzy.Find(query, searchPaths)
	if len(matches) == 0 {
		return "", fmt.Errorf("%w: '%s'", ErrNoMatchFound, query)
	}

	if len(matches) == 1 {
		return absPaths[matches[0].Index], nil
	}

	// Multiple matches, if interactive use fuzzyfinder
	if isInteractive() {
		// Use relative paths for display in fuzzyfinder too, it's cleaner
		filteredSearchPaths := make([]string, len(matches))

		filteredAbsPaths := make([]string, len(matches))
		for i, m := range matches {
			filteredSearchPaths[i] = searchPaths[m.Index]
			filteredAbsPaths[i] = absPaths[m.Index]
		}

		return findInteractive(filteredSearchPaths, filteredAbsPaths)
	}

	// Non-interactive: pick the best match (the first one from fuzzy.Find).
	return absPaths[matches[0].Index], nil
}

func findSelectedPathNoQuery(paths []string) (string, error) {
	if !isInteractive() {
		return "", ErrMissingQueryTerminal
	}

	idx, err := fuzzyfinder.Find(
		paths,
		func(i int) string {
			return paths[i]
		},
	)
	if err != nil {
		return "", fmt.Errorf("fuzzyfinder failed: %w", err)
	}

	return paths[idx], nil
}

func findInteractive(searchPaths []string, absPaths []string) (string, error) {
	idx, err := fuzzyfinder.Find(
		searchPaths,
		func(i int) string {
			return searchPaths[i]
		},
	)
	if err != nil {
		return "", fmt.Errorf("fuzzyfinder failed: %w", err)
	}

	return absPaths[idx], nil
}

var isInteractive = func() bool {
	fileInfo, _ := os.Stdin.Stat()

	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}
