nanan


# git-get

![build](https://github.com/grdl/git-get/workflows/build/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/grdl/git-get)](https://goreportcard.com/report/github.com/grdl/git-get)

`git-get` is a better way to clone, organize and manage multiple git repositories. 

- [git-get](#git-get)
  - [Description](#description)
  - [Installation](#installation)
    - [macOS](#macos)
    - [Linux](#linux)
    - [Windows](#windows)
  - [Usage](#usage)
    - [git get](#git-get-1)
    - [git list](#git-list)
    - [Dump file](#dump-file)
  - [Configuration](#configuration)
    - [Env variables](#env-variables)
    - [.gitconfig file](#gitconfig-file)
  - [Contributing](#contributing)
  - [Acknowledgments](#acknowledgments)

## Description

`git-get` gives you two new git commands:
- **`git get`** clones repositories into an automatically created directory tree based on repo's URL, owner and name (like golang's [`go get`](https://golang.org/cmd/go/)).
- **`git list`** shows status of all your git repositories.

![Example](./docs/example.svg)

## Installation

Each release contains two binaries: `git-get` and `git-list`. When put on PATH, git automatically recognizes them as custom commands and allows to run them as `git get` or `git list`.

### macOS

Use Homebrew:
```
brew install grdl/tap/git-get
```

### Linux

Download and install `.deb` or `.rpm` file from the [latest release](https://github.com/grdl/git-get/releases/latest).

Or install with [Linuxbrew](https://docs.brew.sh/Homebrew-on-Linux):
```
brew install grdl/tap/git-get
```

### Windows

Grab the `.zip` file from the [latest release](https://github.com/grdl/git-get/releases/latest) and put the binaries on your PATH.


## Usage

### git get
```
git get <REPO> [flags]

Flags:
  -b, --branch              Branch (or tag) to checkout after cloning.
  -d, --dump                Path to a dump file listing repos to clone. Ignored when <REPO> argument is used.
  -h, --help                Print this help and exit.
  -t, --host                Host to use when <REPO> doesn't have a specified host. (default "github.com")
  -r, --root                Path to repos root where repositories are cloned. (default "~/repositories")
  -c, --scheme              Scheme to use when <REPO> doesn't have a specified scheme. (default "ssh")
  -s, --skip-host           Don't create a directory for host.
  -v, --version             Print version and exit.
```

The `<REPO>` argument can be any valid URL supported by git. It also accepts a short `USER/REPO` format. In that case `git-get` will automatically use the configured host (github.com by default).

For example, `git get grdl/git-get` will clone `https://github.com/grdl/git-get`.


### git list
```
Usage:
  git list [flags]

Flags:
  -f, --fetch               First fetch from remotes before listing repositories.
  -h, --help                Print this help and exit.
  -o, --out                 Output format. Allowed values: [dump, flat, tree]. (default "tree")
  -r, --root                Path to repos root where repositories are cloned. (default "~/repositories")
  -v, --version             Print version and exit.
```

`git list` provides different ways to view the list of the repositories and their statuses.

- **tree** (default) - repos printed as a directory tree.

![output_tree](./docs/out_tree.png)

- **flat** - each repo (and each branch) on a new line with full path to the repo.

![output_flat](./docs/out_flat.png)

- **dump** - each repo URL with its current branch on a new line. To be consumed by `git get --dump` command.

![output_dump](./docs/out_dump.png)

### Dump file

`git get` is dotfiles friendly. When run with `--dump` flag, it accepts a file with a list of repositories and clones all of them.

Dump file format is simply:
- Each repo URL on a separate line.
- Each URL can have a space-separated suffix with a branch or tag name to check out after cloning. Without that suffix, repository HEAD is cloned (usually it's `master`).

Example dump file content:
```
https://github.com/grdl/git-get v1.0.0
git@github.com:grdl/another-repository.git
```

You can generate a dump file with all your currently cloned repos by running:
```
git list --out dump > repos.dump
``` 

## Configuration

Each configuration flag listed in the [Usage](#Usage) section can be also specified using environment variables or your global `.gitconfig` file.

The order of precedence for configuration is as follows:
- command line flag (have the highest precedence)
- environment variable
- .gitconfig entry
- default value


### Env variables

Use the `GITGET_` prefix and the uppercase flag name to set the configuration using env variables. For example, to use a different repos root path run:
```
export GITGET_ROOT=/path/to/my/repos
```

### .gitconfig file

You can define a `[gitget]` section inside your global `.gitconfig` file and set the configuration flags there. A recommended pattern is to set `root` and `host` variables there if you don't want to use the defaults. 

If all of your repos come from the same host and you find creating directory for it redundant, you can use the `skip-host` flag to skip creating it.

Here's an example of a working snippet from `.gitconfig` file:
```
[gitget]
    root = /path/to/my/repos
    host = gitlab.com
    skip-host = true
```


## Contributing

Pull requests are welcome. The project is still very much work in progress. Here's some of the missing features planned to be fixed soon:
- improvements to the `git list` output (feedback appreciated)
- info about stashes and submodules
- better recognition of different repo states: conflict, merging, rebasing, cherry picking etc.
- plenty of bugfixes and missing tests


## Acknowledgments

Inspired by:
- golang's [`go get`](https://golang.org/cmd/go/) command
- [x-motemen/ghq](https://github.com/x-motemen/ghq)
- [fboender/multi-git-status](https://github.com/fboender/multi-git-status)
