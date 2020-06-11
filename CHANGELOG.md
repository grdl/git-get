# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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


[0.0.3]: https://github.com/grdl/git-get/compare/v0.0.1...v0.0.3
[0.0.1]: https://github.com/grdl/git-get/releases/tag/v0.0.1