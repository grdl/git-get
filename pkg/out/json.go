package out

import (
	"encoding/json"
)

// RepoJSON represents a repository in JSON format.
type RepoJSON struct {
	Path     string            `json:"path"`
	Remote   string            `json:"remote,omitempty"`
	Current  string            `json:"current,omitempty"`
	Branches map[string]string `json:"branches,omitempty"`
	Worktree string            `json:"worktree,omitempty"`
	Errors   []string          `json:"errors,omitempty"`
}

// JSONPrinter prints a list of repos in JSON format.
type JSONPrinter struct{}

// NewJSONPrinter creates a JSONPrinter.
func NewJSONPrinter() *JSONPrinter {
	return &JSONPrinter{}
}

// Print generates a JSON array of repository statuses.
func (p *JSONPrinter) Print(repos []Printable) string {
	result := make([]RepoJSON, 0, len(repos))

	for _, r := range repos {
		branches := make(map[string]string)
		for _, branch := range r.Branches() {
			branches[branch] = r.BranchStatus(branch)
		}

		if current := r.Current(); current != "" {
			branches[current] = r.BranchStatus(current)
		}

		repo := RepoJSON{
			Path:     r.Path(),
			Remote:   r.Remote(),
			Current:  r.Current(),
			Branches: branches,
			Worktree: r.WorkTreeStatus(),
			Errors:   r.Errors(),
		}

		result = append(result, repo)
	}

	output, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "[]"
	}

	return string(output) + "\n"
}
