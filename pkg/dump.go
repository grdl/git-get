package pkg

import (
	"bufio"
	"git-get/pkg/cfg"
	"os"
	"strings"

	"github.com/pkg/errors"
)

var (
	errInvalidNumberOfElements = errors.New("More than two space-separated 2 elements on the line")
	errEmptyLine               = errors.New("Empty line")
)

type parsedLine struct {
	rawurl string
	branch string
}

// ParseDumpFile opens a given gitgetfile and parses its content into a slice of CloneOpts.
func parseDumpFile(path string) ([]parsedLine, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed opening dump file %s", path)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var parsedLines []parsedLine
	var line int
	for scanner.Scan() {
		line++
		parsed, err := parseLine(scanner.Text())
		if err != nil && !errors.Is(errEmptyLine, err) {
			return nil, errors.Wrapf(err, "Failed parsing line %d", line)
		}

		parsedLines = append(parsedLines, parsed)
	}

	return parsedLines, nil
}

// parseLine splits a dump file line into space-separated segments.
// First part is the URL to clone. Second, optional, is the branch (or tag) to checkout after cloning
func parseLine(line string) (parsedLine, error) {
	var parsed parsedLine

	if len(strings.TrimSpace(line)) == 0 {
		return parsed, errEmptyLine
	}

	parts := strings.Split(strings.TrimSpace(line), " ")
	if len(parts) > 2 {
		return parsed, errInvalidNumberOfElements
	}

	parsed.rawurl = parts[0]
	parsed.branch = cfg.DefBranch
	if len(parts) == 2 {
		parsed.branch = parts[1]
	}

	return parsed, nil
}
