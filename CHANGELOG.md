# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).


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


[0.0.6]: https://github.com/grdl/git-get/compare/v0.0.5...v0.0.6
[0.0.5]: https://github.com/grdl/git-get/compare/v0.0.4...v0.0.5
[0.0.4]: https://github.com/grdl/git-get/compare/v0.0.3...v0.0.4
[0.0.3]: https://github.com/grdl/git-get/compare/v0.0.1...v0.0.3
[0.0.1]: https://github.com/grdl/git-get/releases/tag/v0.0.1