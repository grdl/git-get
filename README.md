# git-get

## Build

How to build with libgit2 statically linked into a single executable without dependencies:

- Install development packages for libssh2 and openssl:
  ```
  sudo apt install libssh2-1-del libssl-dev
  ```

- `git2go` library is added as a submodule (pointing to a correct v30 release). This, in turn, contains `libgit2` submodule.
  To ensure the submodules are cloned run:
  ``` 
  git submodule update --init --recursive
  ```

- build the static `git2go` library:
  ```
  cd static/git2go && make install-static
  ```

- ensure our `git-get` module uses the static `git2go` library instead of the one downloaded by Go modules by having
  the following line in `go.mod`:
  ```
  replace github.com/libgit2/git2go/v30 => ./static/git2go
  ```

- build `git-get` with `--tags static` flag:
  ```
  go build -i --tags static
  ```
