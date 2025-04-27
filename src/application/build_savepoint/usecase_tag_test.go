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
	)

	result, err := usecase.Execute(context.Background(), BuildSavepointRequest{
		Savepoint:  "base",
		Force:      false,
		FilePath:   "Dockerfile",
		FullTarget: "docker.io/user/app:1.0",
	})

	assert.NoError(t, err)
	assert.True(t, result.Skipped)
}

func TestForceIgnoresTag(t *testing.T) {
	prepareDockerfile(t)

	usecase := NewBuildSavepointUseCase(
		&simpleParser{},
		&dummyDeps{},
		&dummyDeps{},
		&dummyDeps{},
		&alwaysTagExists{},
		&dummyDeps{},
	)

	result, err := usecase.Execute(context.Background(), BuildSavepointRequest{
		Savepoint:  "base",
		Force:      true,
		FilePath:   "Dockerfile",
		FullTarget: "docker.io/user/app:1.0",
	})

	assert.NoError(t, err)
	assert.False(t, result.Skipped)
}

func TestTagNotExistsTriggersBuild(t *testing.T) {
	prepareDockerfile(t)

	usecase := NewBuildSavepointUseCase(
		&simpleParser{},
		&dummyDeps{},
		&dummyDeps{},
		&dummyDeps{},
		&neverTagExists{},
		&dummyDeps{},
	)

	result, err := usecase.Execute(context.Background(), BuildSavepointRequest{
		Savepoint:  "base",
		FilePath:   "Dockerfile",
		FullTarget: "docker.io/user/app:1.0",
	})

	assert.NoError(t, err)
	assert.False(t, result.Skipped)
}

func TestInvalidImageReference(t *testing.T) {
	tests := []struct {
		name      string
		target    string
		wantError string
	}{
		{
			name:      "missing target",
			target:    "",
			wantError: "❌ Missing required -t flag",
		},
		{
			name:      "missing tag",
			target:    "docker.io/user/app",
			wantError: "missing ':tag'",
		},
		{
			name:      "missing repository",
			target:    ":1.0",
			wantError: "Must be repo/image:tag format",
		},
		{
			name:      "no slash in repo",
			target:    "myapp:1.0",
			wantError: "must contain at least one '/'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prepareDockerfile(t)

			usecase := NewBuildSavepointUseCase(
				&simpleParser{},
				&dummyDeps{},
				&dummyDeps{},
				&dummyDeps{},
				&neverTagExists{},
				&dummyDeps{},
			)

			_, err := usecase.Execute(context.Background(), BuildSavepointRequest{
				Savepoint:  "base",
				FilePath:   "Dockerfile",
				FullTarget: tt.target,
			})

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantError)
		})
	}
}
