package main

import (
	"fmt"
	"testing"

	git "github.com/libgit2/git2go/v30"
)

func TestBranches(t *testing.T) {
	repo, err := git.OpenRepository("/home/grdl/workspace/gitlab.com/grdl/git-get")
	checkFatal(t, err)

	branches, err := Branches(repo)
	checkFatal(t, err)

	fmt.Println(len(branches))
}
