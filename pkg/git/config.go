package git

import "os/exec"

// Gitconfig provides methods to get value from gitconfig files.
type Gitconfig interface {
	GetCfg(key string) string
}

// globalCfg represents a global gitconfig file.
type globalCfg struct{}

func (c *globalCfg) GetCfg(key string) string {
	cmd := exec.Command("git", "config", "--global", key)
	out, err := cmd.Output()

	// In case of error return an empty string, the missing value will fall back to a default.
	if err != nil {
		return ""
	}

	lines := lines(out)
	return lines[0]
}
