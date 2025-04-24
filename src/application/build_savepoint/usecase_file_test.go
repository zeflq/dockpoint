package build_savepoint

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeflq/dockpoint/src/domain"
)

// Captures path sent to the parser
type parserSpy struct {
	Captured string
}

func (p *parserSpy) Parse(path string) ([]domain.Savepoint, error) {
	p.Captured = path
	return []domain.Savepoint{{Name: "base", StartLine: 0, EndLine: 1}}, nil
}

type dummySlicer struct{}
func (s *dummySlicer) Slice(lines []string, sp domain.Savepoint) ([]string, error) {
	return lines, nil
}

func TestUsesCorrectFilePathFromCLI(t *testing.T) {
	tmp := t.TempDir()
	customFile := filepath.Join(tmp, "MyDockerfile")
	_ = os.WriteFile(customFile, []byte("# savepoint: base\nFROM alpine"), 0644)

	parser := &parserSpy{}
	slicer := &dummySlicer{}

	uc := NewBuildSavepointUseCase(parser, slicer, &mockWriter{}, &mockBuilder{}, &mockChecker{}, &mockPusher{}, &mockConfig{})

	_, err := uc.Execute(context.Background(), BuildSavepointRequest{
		FilePath:  customFile,
		Savepoint: "base",
		DryRun:    true,
	})

	assert.NoError(t, err)
	assert.Equal(t, customFile, parser.Captured)
}

