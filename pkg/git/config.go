package git

import (
	"git-get/pkg/run"
)

// ConfigGlobal represents a global gitconfig file.
type ConfigGlobal struct{}

// Get reads a value from global gitconfig file. Returns empty string when key is missing.
func (c *ConfigGlobal) Get(key string) string {
	out, err := run.Git("config", "--global", key).AndCaptureLine()
	// In case of error return an empty string, the missing value will fall back to a default.
	if err != nil {
		return ""
	}

	return out
}
