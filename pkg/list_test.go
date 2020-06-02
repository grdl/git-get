package pkg

import (
	"os"
	"testing"
)

func TestList(t *testing.T) {
	_ = os.Setenv(EnvReposRoot, "/home/gru/workspace")

	paths, err := FindRepos()
	checkFatal(t, err)

	repos, _ := OpenAll(paths)

	PrintRepos(repos)
}
