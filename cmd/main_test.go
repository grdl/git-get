package main

import (
	"os"
	"reflect"
	"testing"
)

//nolint:paralleltest // Tests modifies global state (os.Args) and cannot run in parallel
func TestDetermineCommand(t *testing.T) {
	tests := []struct {
		name        string
		programName string
		args        []string
		wantCmd     string
		wantArgs    []string
	}{
		{
			name:        "git-get with no args",
			programName: "git-get",
			args:        []string{"git-get"},
			wantCmd:     "get",
			wantArgs:    []string{},
		},
		{
			name:        "git-get with repo arg",
			programName: "git-get",
			args:        []string{"git-get", "user/repo"},
			wantCmd:     "get",
			wantArgs:    []string{"user/repo"},
		},
		{
			name:        "git-get with get subcommand",
			programName: "git-get",
			args:        []string{"git-get", "get", "user/repo"},
			wantCmd:     "get",
			wantArgs:    []string{"user/repo"},
		},
		{
			name:        "git-get with list subcommand",
			programName: "git-get",
			args:        []string{"git-get", "list"},
			wantCmd:     "list",
			wantArgs:    []string{},
		},
		{
			name:        "git-list with no args",
			programName: "git-list",
			args:        []string{"git-list"},
			wantCmd:     "list",
			wantArgs:    []string{},
		},
		{
			name:        "git-list with args",
			programName: "git-list",
			args:        []string{"git-list", "--fetch"},
			wantCmd:     "list",
			wantArgs:    []string{"--fetch"},
		},
		{
			name:        "git-get.exe on Windows",
			programName: "git-get.exe",
			args:        []string{"git-get.exe", "user/repo"},
			wantCmd:     "get",
			wantArgs:    []string{"user/repo"},
		},
		{
			name:        "unknown program name with args",
			programName: "some-other-name",
			args:        []string{"some-other-name", "get", "user/repo"},
			wantCmd:     "get",
			wantArgs:    []string{"user/repo"},
		},
		{
			name:        "unknown program name with no args",
			programName: "some-other-name",
			args:        []string{"some-other-name"},
			wantCmd:     "get",
			wantArgs:    []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original os.Args
			oldArgs := os.Args

			defer func() { os.Args = oldArgs }()

			// Set test args
			os.Args = tt.args

			gotCmd, gotArgs := determineCommand()

			if gotCmd != tt.wantCmd {
				t.Errorf("determineCommand() command = %v, want %v", gotCmd, tt.wantCmd)
			}

			if !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Errorf("determineCommand() args = %v, want %v", gotArgs, tt.wantArgs)
			}
		})
	}
}

//nolint:paralleltest // Tests modifies global state (os.Args) and cannot run in parallel
func TestHandleGitGetInvocation(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantCmd  string
		wantArgs []string
	}{
		{
			name:     "no args",
			args:     []string{"git-get"},
			wantCmd:  "get",
			wantArgs: []string{},
		},
		{
			name:     "with repo arg",
			args:     []string{"git-get", "user/repo"},
			wantCmd:  "get",
			wantArgs: []string{"user/repo"},
		},
		{
			name:     "with get subcommand",
			args:     []string{"git-get", "get", "user/repo"},
			wantCmd:  "get",
			wantArgs: []string{"user/repo"},
		},
		{
			name:     "with list subcommand",
			args:     []string{"git-get", "list", "--fetch"},
			wantCmd:  "list",
			wantArgs: []string{"--fetch"},
		},
		{
			name:     "with invalid subcommand",
			args:     []string{"git-get", "invalid", "user/repo"},
			wantCmd:  "get",
			wantArgs: []string{"invalid", "user/repo"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original os.Args
			oldArgs := os.Args

			defer func() { os.Args = oldArgs }()

			// Set test args
			os.Args = tt.args

			gotCmd, gotArgs := handleGitGetInvocation()

			if gotCmd != tt.wantCmd {
				t.Errorf("handleGitGetInvocation() command = %v, want %v", gotCmd, tt.wantCmd)
			}

			if !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Errorf("handleGitGetInvocation() args = %v, want %v", gotArgs, tt.wantArgs)
			}
		})
	}
}

//nolint:paralleltest // Tests modifies global state (os.Args) and cannot run in parallel
func TestHandleGitListInvocation(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantCmd  string
		wantArgs []string
	}{
		{
			name:     "no args",
			args:     []string{"git-list"},
			wantCmd:  "list",
			wantArgs: []string{},
		},
		{
			name:     "with flags",
			args:     []string{"git-list", "--fetch", "--out", "flat"},
			wantCmd:  "list",
			wantArgs: []string{"--fetch", "--out", "flat"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original os.Args
			oldArgs := os.Args

			defer func() { os.Args = oldArgs }()

			// Set test args
			os.Args = tt.args

			gotCmd, gotArgs := handleGitListInvocation()

			if gotCmd != tt.wantCmd {
				t.Errorf("handleGitListInvocation() command = %v, want %v", gotCmd, tt.wantCmd)
			}

			if !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Errorf("handleGitListInvocation() args = %v, want %v", gotArgs, tt.wantArgs)
			}
		})
	}
}

//nolint:paralleltest // Tests modifies global state (os.Args) and cannot run in parallel
func TestHandleDefaultInvocation(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantCmd  string
		wantArgs []string
	}{
		{
			name:     "no args",
			args:     []string{"some-program"},
			wantCmd:  "get",
			wantArgs: []string{},
		},
		{
			name:     "with command arg",
			args:     []string{"some-program", "list"},
			wantCmd:  "list",
			wantArgs: []string{},
		},
		{
			name:     "with command and args",
			args:     []string{"some-program", "get", "user/repo", "--branch", "main"},
			wantCmd:  "get",
			wantArgs: []string{"user/repo", "--branch", "main"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original os.Args
			oldArgs := os.Args

			defer func() { os.Args = oldArgs }()

			// Set test args
			os.Args = tt.args

			gotCmd, gotArgs := handleDefaultInvocation()

			if gotCmd != tt.wantCmd {
				t.Errorf("handleDefaultInvocation() command = %v, want %v", gotCmd, tt.wantCmd)
			}

			if !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Errorf("handleDefaultInvocation() args = %v, want %v", gotArgs, tt.wantArgs)
			}
		})
	}
}
