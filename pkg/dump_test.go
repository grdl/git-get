package pkg

import (
	"testing"
)

func TestParsingRefs(t *testing.T) {
	var tests = []struct {
		name       string
		line       string
		wantBranch string
		wantErr    error
	}{
		{
			name:       "url without branch",
			line:       "https://github.com/grdl/git-get",
			wantBranch: "",
			wantErr:    nil,
		},
		{
			name:       "url with branch",
			line:       "https://github.com/grdl/git-get master",
			wantBranch: "master",
			wantErr:    nil,
		},
		{
			name:       "url with multiple branches",
			line:       "https://github.com/grdl/git-get master branch",
			wantBranch: "",
			wantErr:    errInvalidNumberOfElements,
		},
		{
			name:       "url without path",
			line:       "https://github.com",
			wantBranch: "",
			wantErr:    errEmptyURLPath,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := parseLine(test.line)
			if err != nil && test.wantErr == nil {
				t.Fatalf("got error %q", err)
			}

			// TODO: this should check if we actually got the error we expected
			if err != nil && test.wantErr != nil {
				return
			}

			if got.branch != test.wantBranch {
				t.Errorf("expected %q; got %q", test.wantBranch, got.branch)
			}
		})
	}
}
