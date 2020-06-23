package file

import (
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

// ErrDirectoryAccess indicated a direcotry doesn't exists or can't be accessed
var ErrDirectoryAccess = errors.New("directory doesn't exist or can't be accessed")

// TempDir creates a temporary directory for test repos.
func TempDir() (string, error) {
	dir, err := ioutil.TempDir("", "git-get-repo-")
	if err != nil {
		return "", errors.Wrap(err, "failed creating test repo directory")
	}

	return dir, nil
}

// Write writes string content into a file. If file doesn't exists, it will create it.
func Write(path string, content string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return errors.Wrapf(err, "failed opening a file for writing %s", path)
	}

	_, err = file.Write([]byte(content))
	if err != nil {
		errors.Wrapf(err, "Failed writing to a file %s", path)
	}
	return nil
}

// Exists returns true if a directory exists. If it doesn't or the directory can't be accessed it returns an error.
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)

	if err == nil {
		return true, nil
	}

	if err != nil {
		if os.IsNotExist(err) {
			return false, ErrDirectoryAccess
		}
	}

	// Directory exists but can't be accessed
	return true, ErrDirectoryAccess
}
