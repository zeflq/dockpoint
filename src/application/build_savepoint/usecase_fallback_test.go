package build_savepoint

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeflq/dockpoint/src/domain"
	"github.com/zeflq/dockpoint/src/infrastructure/build"  // Add this import
)

// Fallback when no savepoint provided and none in file
func TestFallbackBuildsFullDockerfileIfNoSavepoints(t *testing.T) {
	tmp := t.TempDir()
	fullFile := filepath.Join(tmp, "Dockerfile")
	_ = os.WriteFile(fullFile, []byte("FROM scratch\nRUN echo hi"), 0644)

	// Change to tmp dir so Dockerfile is discoverable
	original, _ := os.Getwd()
	os.Chdir(tmp)
	defer os.Chdir(original)

	uc := NewBuildSavepointUseCase(
		&emptyParser{},
		build.NewDockerfileSlicer(),  // Use real slicer
		&mockWriter{},
		&mockBuilder{},
		&mockChecker{},
		&mockPusher{},
		&mockConfig{},
	)

	result, err := uc.Execute(context.Background(), BuildSavepointRequest{
		FilePath: "",
		DryRun:   true,
	})

	assert.NoError(t, err)
	assert.True(t, result.Skipped)
	assert.Equal(t, "ghcr.io/test/repo:latest", result.Tag)
	assert.Contains(t, result.DockerfileOut, "FROM scratch")
}
type emptyParser struct{}

func (m *emptyParser) Parse(path string) ([]domain.Savepoint, error) {
	return []domain.Savepoint{}, nil
}