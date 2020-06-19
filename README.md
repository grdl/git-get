
# git-get

![build](https://github.com/grdl/git-get/workflows/build/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/grdl/git-get)](https://goreportcard.com/report/github.com/grdl/git-get)

`git-get` is a better way to clone, organize and manage multiple git repositories. 

It gives you two new git commands:
- **`git get`** clones repositories into an organized directory structure (like golang's [`go get`](https://golang.org/cmd/go/)). It's dotfiles friendly - you can clone multiple repositories listed in a file.
- **`git list`** shows status of all your git repositories and their branches.

![Example](./docs/example.svg)

## Installation

Using Homebrew:
```
brew install grdl/tap/git-get
```

Or grab the [latest release](https://github.com/grdl/git-get/releases) and put the binaries on your PATH.

Each release contains two binaries: `git-get` and `git-list`. When put on PATH, git automatically recognizes them as custom commands and allows to call them as `git get` or `git list`.

## Usage


## Configuration


## Features

- statuses
- multiple outputs
- dotfiles friendly



## Contributing

Pull requests are welcome. The project is still very much work in progress. Here's some of the missing features planned to be fixed soon:
- improvements to the `git list` output
- submodules status
- info about stashes
- better recognition of different repo states: conflict, merging, rebasing, cherry picking etc.
- plenty of bugfixes and tests

## Acknowledgments

Inspired by:
- golang's [`go get`](https://golang.org/cmd/go/) command
- [x-motemen/ghq](https://github.com/x-motemen/ghq)
- [fboender/multi-git-status](https://github.com/fboender/multi-git-status)
