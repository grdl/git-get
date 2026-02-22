package pkg

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSelectBestRoot(t *testing.T) {
	t.Parallel()

	tempDir := filepath.Join(os.TempDir(), "git-get-test")
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		os.RemoveAll(tempDir)
	})

	root1 := filepath.Join(tempDir, "root1")
	root2 := filepath.Join(tempDir, "root2")

	_ = os.MkdirAll(filepath.Join(root1, "github.com", "user1"), 0755)
	_ = os.MkdirAll(filepath.Join(root2, "gitlab.com", "user2"), 0755)

	roots := []string{root1, root2}

	tests := []struct {
		name       string
		pathSuffix string
		want       string
	}{
		{
			name:       "Match first level in root1",
			pathSuffix: "github.com/user1/repo1",
			want:       root1,
		},
		{
			name:       "Match first level in root2",
			pathSuffix: "gitlab.com/user2/repo2",
			want:       root2,
		},
		{
			name:       "No match, default to first root",
			pathSuffix: "bitbucket.org/user3/repo3",
			want:       root1,
		},
		{
			name:       "Match host only in root1",
			pathSuffix: "github.com/user2/repo3",
			want:       root1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := selectBestRoot(roots, tt.pathSuffix)
			if got != tt.want {
				t.Errorf("selectBestRoot() = %v, want %v", got, tt.want)
			}
		})
	}
}
