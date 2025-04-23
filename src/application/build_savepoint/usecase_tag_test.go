package build_savepoint

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeflq/dockpoint/src/domain"
)

// Shared mocks
type alwaysTagExists struct{}
func (m *alwaysTagExists) TagExists(tag string) (bool, error) { return true, nil }

type neverTagExists struct{}
func (m *neverTagExists) TagExists(tag string) (bool, error) { return false, nil }

type dummyDeps struct {
	mockSlicer
	mockWriter
	mockBuilder
	mockPusher
	mockConfig
}

// Simple shared parser for "base"
type simpleParser struct{}
func (m *simpleParser) Parse(path string) ([]domain.Savepoint, error) {
	return []domain.Savepoint{{Name: "base", StartLine: 0, EndLine: 1}}, nil
}

func prepareDockerfile(t *testing.T) {
	tmp := t.TempDir()
	os.WriteFile(filepath.Join(tmp, "Dockerfile"), []byte("FROM test"), 0644)
	orig, _ := os.Getwd()
	os.Chdir(tmp)
	t.Cleanup(func() { os.Chdir(orig) })
}

// ✅ Test 1: Tag exists and force = false → skip
func TestTagExistsSkipsBuild(t *testing.T) {
	prepareDockerfile(t)

	usecase := NewBuildSavepointUseCase(
		&simpleParser{},
		&dummyDeps{},
		&dummyDeps{},
		&dummyDeps{},
		&alwaysTagExists{},
		&dummyDeps{},
		&dummyDeps{},
	)

	result, err := usecase.Execute(context.Background(), BuildSavepointRequest{
		Savepoint: "base",
		Force:     false,
	})

	assert.NoError(t, err)
	assert.True(t, result.Skipped)
}

// ✅ Test 2: Tag exists but force = true → build anyway
func TestForceIgnoresTag(t *testing.T) {
	prepareDockerfile(t)

	usecase := NewBuildSavepointUseCase(
		&simpleParser{},
		&dummyDeps{},
		&dummyDeps{},
		&dummyDeps{},
		&alwaysTagExists{},
		&dummyDeps{},
		&dummyDeps{},
	)

	result, err := usecase.Execute(context.Background(), BuildSavepointRequest{
		Savepoint: "base",
		Force:     true,
	})

	assert.NoError(t, err)
	assert.False(t, result.Skipped)
}

// ✅ Test 3: Tag does not exist → should build
func TestTagNotExistsTriggersBuild(t *testing.T) {
	prepareDockerfile(t)

	usecase := NewBuildSavepointUseCase(
		&simpleParser{},
		&dummyDeps{},
		&dummyDeps{},
		&dummyDeps{},
		&neverTagExists{},
		&dummyDeps{},
		&dummyDeps{},
	)

	result, err := usecase.Execute(context.Background(), BuildSavepointRequest{
		Savepoint: "base",
	})

	assert.NoError(t, err)
	assert.False(t, result.Skipped)
}
