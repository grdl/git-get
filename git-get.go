package main

import (
	"fmt"

	git "github.com/libgit2/git2go/v30"
)

func main() {
	options := &git.CloneOptions{
		CheckoutOpts:         nil,
		FetchOptions:         nil,
		Bare:                 false,
		CheckoutBranch:       "",
		RemoteCreateCallback: nil,
	}

	repo, err := git.Clone("https://gitlab.com/grdl/dotfiles.git", "/tmp/dotfiles/", options)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(repo.IsBare())
}
