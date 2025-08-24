# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.6.0] - 2025-08-24

### Added
- Added support for [Scoop](https://scoop.sh/) installation for Windows.

### Changed
- [#5](https://github.com/grdl/git-get/issues/5) git-get is now built as a single binary with `git-list` symlinked. It automatically detect which command is invoked.
- Updated Go version from `1.16` to `1.24`.
- Updated all Go dependencies and GitHub Action workflows to recent versions.
- Removed deprecated Go modules.
### Fixed

- [#14](https://github.com/grdl/git-get/issues/14) Fixed git-list crashing when running on an empty repository.


## [0.5.0] - 2021-06-03
### Changed
- [#15](https://github.com/grdl/git-get/pull/15) Bump Go version to 1.16. 


## [0.4.0] - 2020-09-02
### Added
- `--scheme` flag for `git get` to set the default scheme to use when scheme is missing from the URL.

### Changed
- Default scheme is now `ssh` instead of `https`.


## [0.3.0] - 2020-07-31
### Added
- More meaningful error messages.
- Show list of errors which ocurred when trying to load repository status.

### Fixed
- Remove empty directories after a failed `git get` command.


## [0.2.0] - 2020-07-08
### Changed
- `git list` won't traverse nested git repositories anymore. This significantly improves performance when listing repos with vendored dependencies (eg, node_modules).


## [0.1.0] - 2020-07-07
### Added
- `--skip-host` flag to skip creating a host directory when cloning 


## [0.0.7] - 2020-07-02
### Fixed
- Missing fetch call on `git list`.


## [0.0.6] - 2020-07-01
### Added
- `.deb` and `.rpm` releases.

### Fixed
- Tree view indentation.
- Missing stdout of git commands.
- Incorrect gitconfig file loading.


## [0.0.5] - 2020-06-30
### Changed
- Remove dependency on [go-git](https://github.com/go-git/go-git) and major refactor to fix performance issues on big repos.

### Fixed
- Correctly expand `--root` pointing to a path containing home variable (eg, `~/my-repos`).
- Correctly process paths on Windows.


## [0.0.4] - 2020-06-19
### Added
- `--dump` flag that allows to clone multiple repos listed in a dump file.
- New `dump` output option for `git list` to generate a dump file.
- Readme with documentation.
- Description of CLI flags and usage when running `--help`.

### Changed
- Split `git-get` and `git-list` into separate binaries.
- Refactor code structure by bringing the `pkg` dir back.


## [0.0.3] - 2020-06-11
### Added
- Homebrew release configuration in goreleaser.
- Different ways to print `git list` output: flat, simple tree and smart tree.
- `--brach` flag that specifies which branch to check out after cloning.
- `--fetch` flag that tells `git list` to fetch from remotes before printing repos status.
- Count number of commits a branch is ahead or behind the upstream.
- SSH key authentication.
- Detect if branch has a detached HEAD.

### Changed
- Refactor configuration provider using [viper](https://github.com/spf13/viper).
- Keep `master` branch on top of sorted branches names.

### Fixed
- Fix panic when trying to walk directories we don't have permissions to access.


## [0.0.1] - 2020-06-01
### Added
- Initial release using [goreleaser](https://github.com/goreleaser/goreleaser).


[0.6.0]: https://github.com/grdl/git-get/compare/v0.5.0...v0.6.0
[0.5.0]: https://github.com/grdl/git-get/compare/v0.4.0...v0.5.0
[0.4.0]: https://github.com/grdl/git-get/compare/v0.3.0...v0.4.0
[0.3.0]: https://github.com/grdl/git-get/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/grdl/git-get/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/grdl/git-get/compare/v0.0.7...v0.1.0
[0.0.7]: https://github.com/grdl/git-get/compare/v0.0.6...v0.0.7
[0.0.6]: https://github.com/grdl/git-get/compare/v0.0.5...v0.0.6
[0.0.5]: https://github.com/grdl/git-get/compare/v0.0.4...v0.0.5
[0.0.4]: https://github.com/grdl/git-get/compare/v0.0.3...v0.0.4
[0.0.3]: https://github.com/grdl/git-get/compare/v0.0.1...v0.0.3
[0.0.1]: https://github.com/grdl/git-get/releases/tag/v0.0.1
