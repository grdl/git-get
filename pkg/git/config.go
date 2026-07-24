// Package git implements functionalities to read and manipulate git repositories
package git

import (
	"github.com/grdl/git-get/pkg/run"
)

// ConfigGlobal represents a global gitconfig file.
type ConfigGlobal struct{}

// Get reads a value from gitconfig using git's native scope resolution.
// Returns empty string when key is missing in all scopes.
func (c *ConfigGlobal) Get(key string) string {
	// Use git config without scope flags to let git handle precedence naturally.
	// This ensures XDG config dir (~/.config/git/config) is checked when ~/.gitconfig exists.
	out, err := run.Git("config", "--get", key).AndCaptureLine()
	if err != nil {
		return ""
	}

	return out
}
