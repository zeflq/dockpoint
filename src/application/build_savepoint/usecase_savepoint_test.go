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

// Mock parser with various savepoints
type parserWith struct {
	savepoints []domain.Savepoint
}

func (m *parserWith) Parse(path string) ([]domain.Savepoint, error) {
	return m.savepoints, nil
}

// ✅ 1. Savepoint doesn't exist
func TestSavepointNotFound(t *testing.T) {
	prepDockerfile(t)

	parser := &parserWith{savepoints: []domain.Savepoint{
		{Name: "deps", StartLine: 0, EndLine: 1},
	}}

	uc := NewBuildSavepointUseCase(parser, &mockSlicer{}, &mockWriter{}, &mockBuilder{}, &mockChecker{}, &mockPusher{}, &mockConfig{})

	_, err := uc.Execute(context.Background(), BuildSavepointRequest{
		Savepoint: "base",
	})
	assert.ErrorIs(t, err, errors.ErrSavepointNotFound)
}

// ✅ 2. Valid savepoint — basic test
func TestValidSavepoint(t *testing.T) {
	prepDockerfile(t)

	parser := &parserWith{savepoints: []domain.Savepoint{
		{Name: "base", StartLine: 0, EndLine: 0},
	}}

	uc := NewBuildSavepointUseCase(parser, &mockSlicer{}, &mockWriter{}, &mockBuilder{}, &mockChecker{}, &mockPusher{}, &mockConfig{})

	result, err := uc.Execute(context.Background(), BuildSavepointRequest{
		Savepoint: "base",
		DryRun:    true,
	})
	assert.NoError(t, err)
	assert.Equal(t, "ghcr.io/test/repo:base", result.Tag)
}

// ✅ 3. Savepoint with end < start should fail at slicer level
type slicerWithError struct{}
func (s *slicerWithError) Slice([]string, domain.Savepoint) ([]string, error) {
	return nil, errors.ErrSavepointNotFound
}

func TestSavepointInvalidRange(t *testing.T) {
	prepDockerfile(t)

	parser := &parserWith{savepoints: []domain.Savepoint{
		{Name: "base", StartLine: 3, EndLine: 1}, // invalid range
	}}

	uc := NewBuildSavepointUseCase(parser, &slicerWithError{}, &mockWriter{}, &mockBuilder{}, &mockChecker{}, &mockPusher{}, &mockConfig{})

	_, err := uc.Execute(context.Background(), BuildSavepointRequest{
		Savepoint: "base",
	})
	assert.ErrorIs(t, err, errors.ErrSavepointNotFound)
}

func TestEmptySavepointStillReturnsEmptySlicedDockerfile(t *testing.T) {
	prepDockerfile(t)

	parser := &parserWith{savepoints: []domain.Savepoint{
		{Name: "empty", StartLine: 0, EndLine: -1},
	}}

	slicer := &mockSlicer{} // Returns non-error slice even if empty range

	uc := NewBuildSavepointUseCase(parser, slicer, &mockWriter{}, &mockBuilder{}, &mockChecker{}, &mockPusher{}, &mockConfig{})

	result, err := uc.Execute(context.Background(), BuildSavepointRequest{
		Savepoint: "empty",
		DryRun:    true,
	})

	assert.NoError(t, err)
	assert.True(t, result.Skipped)
	assert.Equal(t, "ghcr.io/test/repo:empty", result.Tag)
}
