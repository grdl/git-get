package main

import (
	"os"
	"path/filepath"
	"strings"
)

// This program behaves as a git subcommand (see https://git.github.io/htmldocs/howto/new-command.html)
// When added to PATH, git recognizes it as its subcommand and it can be invoked as "git get..." or "git list..."
// It can also be invoked as a regular binary with subcommands: "git-get get..." or "git-get list"
// The following flow detects the invokation method and runs the appropriate command.

func main() {
	command, args := determineCommand()
	executeCommand(command, args)
}

func determineCommand() (string, []string) {
	programName := strings.TrimSuffix(filepath.Base(os.Args[0]), ".exe")

	switch programName {
	case "git-get":
		return handleGitGetInvocation()
	case "git-list":
		return handleGitListInvocation()
	default:
		return handleDefaultInvocation()
	}
}

func handleGitGetInvocation() (string, []string) {
	if len(os.Args) > 1 && (os.Args[1] == "get" || os.Args[1] == "list") {
		return os.Args[1], os.Args[2:]
	}

	return "get", os.Args[1:]
}

func handleGitListInvocation() (string, []string) {
	return "list", os.Args[1:]
}

func handleDefaultInvocation() (string, []string) {
	if len(os.Args) > 1 {
		return os.Args[1], os.Args[2:]
	}

	return "get", []string{}
}

func executeCommand(command string, args []string) {
	switch command {
	case "get":
		runGet(args)
	case "list":
		runList(args)
	default:
		runGet(os.Args[1:])
	}
}
