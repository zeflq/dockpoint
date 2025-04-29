package build_savepoint

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeflq/dockpoint/src/domain"
)

// Simple shared parser for "base"
type simpleParser struct{}
func (m *simpleParser) Parse(path string) ([]domain.Savepoint, error) {
    // Changed: adjust line range to match actual content
    return []domain.Savepoint{{Name: "base", StartLine: 0, EndLine: 0}}, nil
}

func prepareDockerfile(t *testing.T) {
    tmp := t.TempDir()
    // Changed: add newline to ensure proper line counting
    err := os.WriteFile(filepath.Join(tmp, "Dockerfile"), []byte("FROM test\n"), 0644)
    assert.NoError(t, err)
    orig, _ := os.Getwd()
    err = os.Chdir(tmp)
    assert.NoError(t, err)
    t.Cleanup(func() { os.Chdir(orig) })
}

// ✅ Test 1: Tag exists and force = false → skip
func TestTagExistsSkipsBuild(t *testing.T) {
	prepareDockerfile(t)

	usecase := NewBuildSavepointUseCase(
		&simpleParser{},
		&dummySlicer{}, // Changed: use dummySlicer instead of mockSlicer
		&spyWriter{},   // Changed: use spyWriter to track calls
		&spyBuilder{},  // Changed: use spyBuilder to track calls
		&alwaysTagExists{},
		&spyPusher{},   // Changed: use spyPusher to track calls
		&mockValidator{},
		&mockHasher{},
		&mockTagBuilder{},
	)

	results, err := usecase.Execute(context.Background(), BuildSavepointRequest{
		Savepoint:  "base",
		Force:      false,
		FilePath:   "Dockerfile",
		FullTarget: "docker.io/user/app:1.0",
		DryRun:     true, // Added: enable dry run mode
	})

	assert.NoError(t, err)
	result := results.Results[0]
	assert.True(t, result.Skipped, "Build should be skipped when tag exists and force is false")
	assert.Equal(t, "docker.io/user/app:1.0", result.Tag)
}

func TestForceIgnoresTag(t *testing.T) {
	prepareDockerfile(t)

	usecase := NewBuildSavepointUseCase(
		&simpleParser{},
		&mockSlicer{},
		&mockWriter{},
		&mockBuilder{},
		&alwaysTagExists{},
		&mockPusher{},
		&mockValidator{},
		&mockHasher{},
		&mockTagBuilder{},
	)

	results, err := usecase.Execute(context.Background(), BuildSavepointRequest{
		Savepoint:  "base",
		Force:      true,
		FilePath:   "Dockerfile",
		FullTarget: "docker.io/user/app:1.0",
	})

	assert.NoError(t, err)
	result := results.Results[0]
	assert.False(t, result.Skipped)
}

func TestTagNotExistsTriggersBuild(t *testing.T) {
	prepareDockerfile(t)

	usecase := NewBuildSavepointUseCase(
		&simpleParser{},
		&mockSlicer{},
		&mockWriter{},
		&mockBuilder{},
		&neverTagExists{},
		&mockPusher{},
		&mockValidator{},
		&mockHasher{},
		&mockTagBuilder{},
	)

	results, err := usecase.Execute(context.Background(), BuildSavepointRequest{
		Savepoint:  "base",
		FilePath:   "Dockerfile",
		FullTarget: "docker.io/user/app:1.0",
	})

	assert.NoError(t, err)
	result := results.Results[0]
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
				&mockSlicer{},
				&mockWriter{},
				&mockBuilder{},
				&neverTagExists{},
				&mockPusher{},
				&mockValidator{},
				&mockHasher{},
				&mockTagBuilder{},
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
