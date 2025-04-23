package parse

import (
	"bufio"
	"os"
	"strings"

	"github.com/zeflq/dockpoint/src/core/parse"
	"github.com/zeflq/dockpoint/src/domain"
)

type DockerfileParserImpl struct{}

func NewDockerfileParser() parse.DockerfileParser {
	return &DockerfileParserImpl{}
}

func (p *DockerfileParserImpl) Parse(path string) ([]domain.Savepoint, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var savepoints []domain.Savepoint
	var current *domain.Savepoint

	scanner := bufio.NewScanner(file)
	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()

		if strings.HasPrefix(line, "# savepoint:") {
			if current != nil {
				current.EndLine = i - 1
				savepoints = append(savepoints, *current)
			}

			name := strings.TrimSpace(strings.TrimPrefix(line, "# savepoint:"))
			current = &domain.Savepoint{Name: name, StartLine: i}
		}
	}
	if current != nil {
		current.EndLine = -1 // unknown end (maybe till end-of-file)
		savepoints = append(savepoints, *current)
	}

	return savepoints, scanner.Err()
}
