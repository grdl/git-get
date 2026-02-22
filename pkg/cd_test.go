package pkg

import (
	"context"
	"errors"
	"os"
	"testing"
)

var errMockRun = errors.New("mock run error")

//nolint:paralleltest
func TestFindSelectedPathWithQuery(t *testing.T) {
	paths := []string{
		"/root1/github.com/user/repo-foo",
		"/root1/github.com/user/repo-bar",
		"/root2/gitlab.com/user/project-foo",
	}

	tests := []struct {
		name                   string
		query                  string
		want                   string
		wantErr                error
		isInteractiveTestValue bool
	}{
		{
			name:                   "single match, non-interactive",
			query:                  "repo-foo",
			want:                   "/root1/github.com/user/repo-foo",
			wantErr:                nil,
			isInteractiveTestValue: false,
		},
		{
			name:                   "multiple matches, non-interactive (picks first)",
			query:                  "foo",
			want:                   "/root1/github.com/user/repo-foo",
			wantErr:                nil,
			isInteractiveTestValue: false,
		},
		{
			name:                   "no match",
			query:                  "nonexistent",
			want:                   "",
			wantErr:                ErrNoMatchFound,
			isInteractiveTestValue: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Temporarily override pkg.isInteractive for testing purposes.
			oldIsInteractive := isInteractive
			isInteractive = func() bool { return tt.isInteractiveTestValue }

			defer func() { isInteractive = oldIsInteractive }()

			got, err := findSelectedPathWithQuery(tt.query, paths, paths)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("findSelectedPathWithQuery() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if got != tt.want {
				t.Errorf("findSelectedPathWithQuery() got = %v, want %v", got, tt.want)
			}
		})
	}
}

//nolint:paralleltest
func TestFindSelectedPathNoQuery(t *testing.T) {
	paths := []string{
		"/root1/github.com/user/repo1",
		"/root1/github.com/user/repo2",
	}

	tests := []struct {
		name                   string
		isInteractiveTestValue bool
		want                   string
		wantErr                error
	}{
		{
			name:                   "non-interactive, no query (should error)",
			isInteractiveTestValue: false,
			want:                   "",
			wantErr:                ErrMissingQueryTerminal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Temporarily override pkg.isInteractive for testing purposes.
			oldIsInteractive := isInteractive
			isInteractive = func() bool { return tt.isInteractiveTestValue }

			defer func() { isInteractive = oldIsInteractive }()

			got, err := findSelectedPathNoQuery(paths)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("findSelectedPathNoQuery() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if got != tt.want {
				t.Errorf("findSelectedPathNoQuery() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// mockCmd implements Cmd interface for testing.
type mockCmd struct {
	runFunc func() error
	dir     string
	stdin   *os.File
	stdout  *os.File
	stderr  *os.File
}

func (m *mockCmd) Run() error {
	return m.runFunc()
}

func (m *mockCmd) SetDir(dir string) {
	m.dir = dir
}

func (m *mockCmd) SetStdin(stdin *os.File) {
	m.stdin = stdin
}

func (m *mockCmd) SetStdout(stdout *os.File) {
	m.stdout = stdout
}

func (m *mockCmd) SetStderr(stderr *os.File) {
	m.stderr = stderr
}

// mockCommander implements commander interface for testing.
type mockCommander struct {
	mockedCmd    Cmd
	capturedName string
	capturedArgs []string
}

//nolint:ireturn
func (m *mockCommander) CommandContext(_ context.Context, name string, arg ...string) Cmd {
	m.capturedName = name
	m.capturedArgs = arg

	return m.mockedCmd
}

//nolint:cyclop
func TestStartSubshell(t *testing.T) {
	mockedCmd := &mockCmd{}
	mc := &mockCommander{mockedCmd: mockedCmd}

	// Replace the global defaultCommander with our mock.
	oldCommander := defaultCommander
	defaultCommander = mc

	defer func() { defaultCommander = oldCommander }()

	tests := []struct {
		name          string
		repoPath      string
		shellEnv      string
		mockRunError  error
		wantErr       bool
		expectedShell string
	}{
		{
			name:          "successful subshell start (default shell)",
			repoPath:      "",
			shellEnv:      "",
			mockRunError:  nil,
			wantErr:       false,
			expectedShell: "/bin/sh",
		},
		{
			name:          "successful subshell start (custom shell)",
			repoPath:      "",
			shellEnv:      "/bin/bash",
			mockRunError:  nil,
			wantErr:       false,
			expectedShell: "/bin/bash",
		},
		{
			name:          "subshell fails to start",
			repoPath:      "",
			shellEnv:      "/bin/zsh",
			mockRunError:  errMockRun,
			wantErr:       true,
			expectedShell: "/bin/zsh",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoPath := tt.repoPath
			if !tt.wantErr {
				repoPath = t.TempDir()
			}

			mc.capturedName = ""
			mc.capturedArgs = nil
			mockedCmd.dir = ""
			mockedCmd.stdin = nil
			mockedCmd.stdout = nil
			mockedCmd.stderr = nil
			mockedCmd.runFunc = func() error { return tt.mockRunError }

			if tt.shellEnv != "" {
				t.Setenv("SHELL", tt.shellEnv)
			} else {
				os.Unsetenv("SHELL")
			}

			err := startSubshell(repoPath)

			if (err != nil) != tt.wantErr {
				t.Errorf("startSubshell() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if mc.capturedName != tt.expectedShell {
					t.Errorf("captured name = %s, want %s", mc.capturedName, tt.expectedShell)
				}

				if mockedCmd.dir != repoPath {
					t.Errorf("captured Dir = %s, want %s", mockedCmd.dir, repoPath)
				}

				if mockedCmd.stdin != os.Stdin || mockedCmd.stdout != os.Stdout || mockedCmd.stderr != os.Stderr {
					t.Errorf("io not set correctly")
				}
			}
		})
	}
}
