package print

import (
	"fmt"
	"git-get/git"
	"strings"
	"testing"
)

func TestTree(t *testing.T) {
	var tests = []struct {
		paths []string
		want  string
	}{
		{
			[]string{
				"root/github.com/grdl/repo1",
			}, `
root/
github.com/grdl/repo1
`,
		},
		{
			[]string{
				"root/github.com/grdl/repo1",
				"root/github.com/grdl/repo2",
			}, `
root/
github.com/grdl/
	repo1
	repo2
`,
		},
		{
			[]string{
				"root/gitlab.com/grdl/repo1",
				"root/github.com/grdl/repo1",
			}, `
root/
gitlab.com/grdl/repo1
github.com/grdl/repo1
`,
		},
		{
			[]string{
				"root/gitlab.com/grdl/repo1",
				"root/gitlab.com/grdl/repo2",
				"root/gitlab.com/other/repo1",
				"root/github.com/grdl/repo1",
				"root/github.com/grdl/nested/repo2",
			}, `
root/
gitlab.com/
	grdl/
		repo1
		repo2
	other/repo1
github.com/grdl/
	repo1
	nested/repo2
`,
		},
		{
			[]string{
				"root/gitlab.com/grdl/nested/repo1",
				"root/gitlab.com/grdl/nested/repo2",
				"root/gitlab.com/other/repo1",
			}, `
root/
gitlab.com/
	grdl/nested/
		repo1
		repo2
	other/repo1
`,
		},
		{
			[]string{
				"root/gitlab.com/grdl/double/nested/repo1",
				"root/gitlab.com/grdl/nested/repo2",
				"root/gitlab.com/other/repo1",
			}, `
root/
gitlab.com/
	grdl/
		double/nested/repo1
		nested/repo2
	other/repo1
`,
		},
	}

	for i, test := range tests {
		var repos []*git.Repo
		for _, path := range test.paths {
			repos = append(repos, git.NewRepo(nil, path)) //&Repo{path: path})
		}

		tree := BuildTree("root", repos)
		// Leading and trailing newlines are added to test cases for readability. We also need to add them to the rendering result.
		got := fmt.Sprintf("\n%s\n", RenderSmartTree(tree))

		// Rendered tree uses spaces for indentation but the test cases use tabs.
		if got != strings.ReplaceAll(test.want, "\t", "    ") {
			t.Errorf("Failed test case %d, got: %+v; want: %+v", i, got, test.want)
		}
	}
}
