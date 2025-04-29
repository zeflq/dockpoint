package build_savepoint

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeflq/dockpoint/src/core/errors"
	"github.com/zeflq/dockpoint/src/domain"
)

func prepDockerfile(t *testing.T) {
	tmp := t.TempDir()
	os.WriteFile(filepath.Join(tmp, "Dockerfile"), []byte("FROM node:20"), 0644)
	orig, _ := os.Getwd()
	os.Chdir(tmp)
	t.Cleanup(func() { os.Chdir(orig) })
}

func TestSavepointNotFound(t *testing.T) {
	prepDockerfile(t)

	uc := NewBuildSavepointUseCase(
		&parserWith{savepoints: []domain.Savepoint{{Name: "deps", StartLine: 0, EndLine: 1}}},
		&mockSlicer{},
		&mockWriter{},
		&mockBuilder{},
		&mockChecker{},
		&mockPusher{},
		&mockValidator{},
		&mockHasher{},
		&mockTagBuilder{},
	)

	_, err := uc.Execute(context.Background(), BuildSavepointRequest{
		Savepoint:  "base",
		FilePath:   "Dockerfile",
		FullTarget: "docker.io/user/app:latest",
	})
	assert.ErrorIs(t, err, errors.ErrSavepointNotFound)
}

func TestValidSavepoint(t *testing.T) {
	prepDockerfile(t)

	uc := NewBuildSavepointUseCase(
		&parserWith{savepoints: []domain.Savepoint{{Name: "base", StartLine: 0, EndLine: 0}}},
		&mockSlicer{},
		&mockWriter{},
		&mockBuilder{},
		&mockChecker{},
		&mockPusher{},
		&mockValidator{},
		&mockHasher{},
		&mockTagBuilder{},
	)

	results, err := uc.Execute(context.Background(), BuildSavepointRequest{
		Savepoint:  "base",
		DryRun:     true,
		FilePath:   "Dockerfile",
		FullTarget: "docker.io/user/app:latest",
	})
	assert.NoError(t, err)
	result := results.Results[0]
	assert.Equal(t, "docker.io/user/app:latest", result.Tag)
}