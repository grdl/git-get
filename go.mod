module git-get

go 1.14

require (
	github.com/libgit2/git2go/v30 v30.0.3
	github.com/pkg/errors v0.9.1
)

replace github.com/libgit2/git2go/v30 => ./static/git2go
