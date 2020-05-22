package pkg

import (
	urlpkg "net/url"
	"os"
	"path"

	"github.com/pkg/errors"

	git "github.com/libgit2/git2go/v30"
)

type Repo struct {
	repo   *git.Repository
	Status *RepoStatus
}

func CloneRepo(url *urlpkg.URL, repoRoot string) (path string, err error) {
	repoPath := URLToPath(url)

	path, err = MakeDir(repoRoot, repoPath)
	if err != nil {
		return path, err
	}

	options := &git.CloneOptions{
		CheckoutOpts:         nil,
		FetchOptions:         nil,
		Bare:                 false,
		CheckoutBranch:       "",
		RemoteCreateCallback: nil,
	}

	_, err = git.Clone(url.String(), path, options)
	if err != nil {
		return path, errors.Wrap(err, "Failed cloning repo")
	}
	return path, nil
}

func OpenRepo(path string) (*Repo, error) {
	r, err := git.OpenRepository(path)
	if err != nil {
		return nil, errors.Wrap(err, "Failed opening repo")
	}

	repoStatus, err := loadStatus(r)
	if err != nil {
		return nil, err
	}

	repo := &Repo{
		repo:   r,
		Status: repoStatus,
	}

	return repo, nil
}

func (r *Repo) Reload() error {
	status, err := loadStatus(r.repo)
	if err != nil {
		return err
	}

	r.Status = status
	return nil
}

func (r *Repo) Fetch() error {
	remoteNames, err := r.repo.Remotes.List()
	if err != nil {
		return errors.Wrap(err, "Failed listing remoteNames")
	}

	for _, name := range remoteNames {
		remote, err := r.repo.Remotes.Lookup(name)
		if err != nil {
			return errors.Wrap(err, "Failed looking up remote")
		}

		err = remote.Fetch(nil, nil, "")
		if err != nil {
			return errors.Wrap(err, "Failed fetching remote")
		}
	}

	return nil
}

func MakeDir(repoRoot, repoPath string) (string, error) {
	dir := path.Join(repoRoot, repoPath)
	err := os.MkdirAll(dir, 0775)
	if err != nil {
		return "", errors.Wrap(err, "Failed creating repo directory")
	}

	return dir, nil
}
