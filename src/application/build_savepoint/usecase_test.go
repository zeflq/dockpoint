package build_savepoint

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeflq/dockpoint/src/domain"
)

// ✅ Mocks
// Update mockParser to handle empty Dockerfile
type mockParser struct {
	dockerfile []string
}

func (m *mockParser) Parse(path string) ([]domain.Savepoint, error) {
	if len(m.dockerfile) == 0 {
		return nil, errors.New("empty dockerfile")
	}
	return []domain.Savepoint{{Name: "base", StartLine: 0, EndLine: 1}}, nil
}

type mockSlicer struct {
	content []string
}
func (m *mockSlicer) Slice(lines []string, sp domain.Savepoint) ([]string, error) {
	return []string{"FROM node:20"}, nil
}

type mockWriter struct{}
func (m *mockWriter) Write(lines []string, savepoint string) (string, error) {
	return "/tmp/fake.Dockerfile", nil
}

type mockBuilder struct{}
func (m *mockBuilder) Build(ctx context.Context, dockerfilePath, contextDir, tag string) error {
	return nil
}

type mockChecker struct{}
func (m *mockChecker) TagExists(tag string) (bool, error) {
	return false, nil
}

type mockPusher struct{}
func (m *mockPusher) Push(tag string) error {
	return nil
}

type mockConfig struct{}
func (m *mockConfig) GetRepo() (string, error) {
	return "ghcr.io/test/repo", nil
}

func TestUseCase_DryRun(t *testing.T) {
	tests := []struct {
		name        string
		dockerfile  []string
		savepoint   string
		dryRun      bool
		wantSkipped bool
		wantTag     string
		wantError   bool
	}{
		{
			name: "successful dry run",
			dockerfile: []string{
				"FROM node:20",
				"# savepoint: base",
				"RUN npm install",
			},
			savepoint:   "base",
			dryRun:      true,
			wantSkipped: true,
			wantTag:     "ghcr.io/test/repo:base",
		},
		{
			name: "successful build",
			dockerfile: []string{
				"FROM node:20",
				"# savepoint: base",
				"RUN npm install",
			},
			savepoint:   "base",
			dryRun:      false,
			wantSkipped: false,
			wantTag:     "ghcr.io/test/repo:base",
		},
		{
			name: "non-existent savepoint",
			dockerfile: []string{
				"FROM node:20",
				"# savepoint: base",
				"RUN npm install",
			},
			savepoint:   "nonexistent",
			dryRun:      true,
			wantError:   true,
		},
		{
			name: "empty dockerfile",
			dockerfile: []string{},
			savepoint:  "base",
			dryRun:     true,
			wantError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := os.WriteFile("Dockerfile", []byte(strings.Join(tt.dockerfile, "\n")), 0644)
			if err != nil {
				t.Fatalf("Failed to create test Dockerfile: %v", err)
			}
			defer os.Remove("Dockerfile")

			usecase := NewBuildSavepointUseCase(
				&mockParser{dockerfile: tt.dockerfile}, // Pass the dockerfile content
				&mockSlicer{content: tt.dockerfile},
				&mockWriter{},
				&mockBuilder{},
				&mockChecker{},
				&mockPusher{},
				&mockConfig{},
			)

			req := BuildSavepointRequest{
				Savepoint: tt.savepoint,
				DryRun:    tt.dryRun,
			}

			result, err := usecase.Execute(context.Background(), req)

			if tt.wantError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wantSkipped, result.Skipped)
			assert.Equal(t, tt.wantTag, result.Tag)
			if tt.dryRun {
				assert.Contains(t, result.DockerfileOut, "FROM node:20")
			}
		})
	}
}
