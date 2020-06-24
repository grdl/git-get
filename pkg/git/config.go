package git

import "os/exec"

// ConfigGlobal represents a global gitconfig file.
type ConfigGlobal struct{}

// Get reads a value from global gitconfig file. Returns empty string when key is missing.
func (c *ConfigGlobal) Get(key string) string {
	cmd := exec.Command("git", "config", "--global", key)
	out, err := cmd.Output()

	// In case of error return an empty string, the missing value will fall back to a default.
	if err != nil {
		return ""
	}

	lines := lines(out)
	return lines[0]
}
