package parse

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/zeflq/dockpoint/src/core/parse"
	"github.com/zeflq/dockpoint/src/domain"
)

type DockerfileParserImpl struct{}

func NewDockerfileParser() parse.DockerfileParser {
	return &DockerfileParserImpl{}
}

var isValidSavepoint = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

func generateFinalSavepointName() string {
	timestamp := time.Now().Unix()
	return fmt.Sprintf("final-%d", timestamp)
}

func (p *DockerfileParserImpl) Parse(path string) ([]domain.Savepoint, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	fromLine := -1
	for i, line := range lines {
		if strings.HasPrefix(strings.ToUpper(strings.TrimSpace(line)), "FROM ") {
			fromLine = i
			break
		}
	}
	if fromLine == -1 {
		return nil, nil // invalid Dockerfile
	}

	savepoints := []domain.Savepoint{}
	start := 0
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "# savepoint:") {
			if i <= fromLine {
				continue // ignore savepoints before FROM
			}

			name := strings.TrimSpace(strings.TrimPrefix(trimmed, "# savepoint:"))
			if name == "" || !isValidSavepoint.MatchString(name) {
				continue // skip malformed or invalid savepoint
			}

			if start <= i {
				savepoints = append(savepoints, domain.Savepoint{
					Name:      name,
					StartLine: start,
					EndLine:   i, // inclusive
				})
				start = i + 1
			}
		}
	}

	// Append trailing savepoint if lines remain after last savepoint
	if start <= len(lines)-1 {
		savepoints = append(savepoints, domain.Savepoint{
			Name:      generateFinalSavepointName(),
			StartLine: start,
			EndLine:   len(lines) - 1,
		})
	}

	return savepoints, nil
}
