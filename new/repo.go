package new

import (
	"io"
	"os"

	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/pkg/errors"

	"github.com/go-git/go-billy/v5"

	"github.com/go-git/go-git/v5"
)

type Repo struct {
	repo   *git.Repository
	Status *RepoStatus
}

func CloneRepo(url string, path billy.Filesystem, quiet bool) (r *Repo, err error) {
	var output io.Writer
	if !quiet {
		output = os.Stdout
	}

	opts := &git.CloneOptions{
		URL:               url,
		Auth:              nil,
		RemoteName:        git.DefaultRemoteName,
		ReferenceName:     "",
		SingleBranch:      false,
		NoCheckout:        false,
		Depth:             0,
		RecurseSubmodules: git.NoRecurseSubmodules,
		Progress:          output,
		Tags:              git.AllTags,
	}

	dotgit, _ := path.Chroot(git.GitDirName)
	s := filesystem.NewStorage(dotgit, cache.NewObjectLRUDefault())

	repo, err := git.Clone(s, path, opts)

	if err != nil {
		return nil, errors.Wrap(err, "Failed cloning repo")
	}

	r = &Repo{
		repo: repo,
	}
	return r, nil
}
