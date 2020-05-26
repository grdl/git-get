package new

import "github.com/go-git/go-git/v5"

type Repo struct {
	repo   *git.Repository
	Status *RepoStatus
}
