package pkg

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/grdl/git-get/pkg/cfg"
)

//nolint:paralleltest
func TestListFuzzy(t *testing.T) {
	tempDir := t.TempDir()
	root := filepath.Join(tempDir, "root")

	err := os.MkdirAll(filepath.Join(root, "github.com", "user", "repo1", ".git"), 0755)
	if err != nil {
		t.Fatal(err)
	}

	err = os.MkdirAll(filepath.Join(root, "github.com", "user", "repo2", ".git"), 0755)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name  string
		query string
		want  int // Number of matches expected
	}{
		{
			name:  "match all",
			query: "repo",
			want:  2,
		},
		{
			name:  "match one",
			query: "repo1",
			want:  1,
		},
		{
			name:  "match none",
			query: "nonexistent",
			want:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := &ListCfg{
				Fetch:  false,
				Output: cfg.OutFlat, // Flat output is easiest
				Roots:  []string{root},
				Query:  tt.query,
			}

			// We use a pipe to suppress output during tests.
			oldStdout := os.Stdout

			defer func() { os.Stdout = oldStdout }()

			_, w, _ := os.Pipe()
			os.Stdout = w

			err := List(conf)

			w.Close()

			if err != nil {
				t.Errorf("List() error = %v", err)
			}
		})
	}
}
