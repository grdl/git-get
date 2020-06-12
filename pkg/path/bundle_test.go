package path

import (
	"testing"
)

func TestParsingRefs(t *testing.T) {
	var tests = []struct {
		line       string
		wantBranch string
		wantErr    error
	}{
		{
			line:       "https://github.com/grdl/git-get",
			wantBranch: "",
			wantErr:    nil,
		},
		{
			line:       "https://github.com/grdl/git-get master",
			wantBranch: "master",
			wantErr:    nil,
		},
		{
			line:       "https://github.com/grdl/git-get refs/tags/v1.0.0",
			wantBranch: "refs/tags/v1.0.0",
			wantErr:    nil,
		},
		{
			line:       "https://github.com/grdl/git-get master branch",
			wantBranch: "",
			wantErr:    ErrInvalidNumberOfElements,
		},
		{
			line:       "https://github.com",
			wantBranch: "",
			wantErr:    ErrEmptyURLPath,
		},
	}

	for i, test := range tests {
		got, err := parseLine(test.line)
		if err != nil && test.wantErr == nil {
			t.Fatalf("Test case %d should not return an error", i)
		}

		if err != nil && test.wantErr != nil {
			continue
		}

		if got.Branch != test.wantBranch {
			t.Errorf("Failed test case %d, got: %s; wantBranch: %s", i, got.Branch, test.wantBranch)
		}
	}

}
