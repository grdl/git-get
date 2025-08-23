package main

import (
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// This program behaves as a git subcommand (see https://git.github.io/htmldocs/howto/new-command.html)
	// When added to PATH, git recognizes it as its subcommand and it can be invoked as "git get..." or "git list..."
	// It can also be invoked as a regular binary with subcommands: "git-get get..." or "git-get list"
	// The following flow detects the invokation method and runs the appropriate command.

	programName := filepath.Base(os.Args[0])

	// Remove common executable extensions
	programName = strings.TrimSuffix(programName, ".exe")

	// Determine which command to run based on program name or first argument
	var command string
	var args []string

	switch programName {
	case "git-get":
		// Check if first argument is a subcommand
		if len(os.Args) > 1 && (os.Args[1] == "get" || os.Args[1] == "list") {
			// Called as: git-get get <repo> or git-get list
			command = os.Args[1]
			args = os.Args[2:]
		} else {
			// Called as: git-get <repo> (default to get command)
			command = "get"
			args = os.Args[1:]
		}
	case "git-list":
		// Called as: git-list (symlinked binary)
		command = "list"
		args = os.Args[1:]
	default:
		// Fallback: use first argument as command
		if len(os.Args) > 1 {
			command = os.Args[1]
			args = os.Args[2:]
		} else {
			command = "get"
			args = []string{}
		}
	}

	// Execute the appropriate command
	switch command {
	case "get":
		runGet(args)
	case "list":
		runList(args)
	default:
		// Default to get command for unknown commands
		runGet(os.Args[1:])
	}
}
