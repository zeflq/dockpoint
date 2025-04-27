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
		fullTarget  string
		wantSkipped bool
		wantTag     string
		wantError   bool
		errorMsg    string
	}{
		{
			name: "successful dry run",
			dockerfile: []string{
				"# savepoint: base",
				"FROM node:20",
				"RUN npm install",
			},
			savepoint:   "base",
			dryRun:      true,
			fullTarget:  "docker.io/user/app:latest",
			wantSkipped: true,
			wantTag:     "docker.io/user/app:base",
		},
		{
			name: "successful build",
			dockerfile: []string{
				"# savepoint: base",
				"FROM node:20",
				"RUN npm install",
			},
			savepoint:   "base",
			dryRun:      false,
			fullTarget:  "docker.io/user/app:latest",
			wantSkipped: false,
			wantTag:     "docker.io/user/app:base",
		},
		{
			name: "non-existent savepoint",
			dockerfile: []string{
				"# savepoint: base",
				"FROM node:20",
			},
			savepoint:  "nonexistent",
			dryRun:     true,
			fullTarget: "docker.io/user/app:latest",
			wantError:  true,
		},
		{
			name:       "empty dockerfile",
			dockerfile: []string{},
			savepoint:  "base",
			dryRun:     true,
			fullTarget: "docker.io/user/app:latest",
			wantError:  true,
		},
		{
			name:      "missing target",
			dockerfile: []string{"FROM node:20"},
			savepoint: "base",
			wantError: true,
			errorMsg:  "Missing required -t flag",
		},
		{
			name:       "invalid target format",
			dockerfile: []string{"FROM node:20"},
			savepoint:  "base",
			fullTarget: "invalid",
			wantError:  true,
			errorMsg:   "missing ':tag'",  // Updated to match actual error
		},
		// Add new test for missing slash
		{
			name:       "missing slash in repository",
			dockerfile: []string{"FROM node:20"},
			savepoint:  "base",
			fullTarget: "invalid:tag",
			wantError:  true,
			errorMsg:   "must contain at least one '/'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmp := t.TempDir()
			_ = os.Chdir(tmp)

			err := os.WriteFile("Dockerfile", []byte(strings.Join(tt.dockerfile, "\n")), 0644)
			if err != nil {
				t.Fatalf("Failed to create test Dockerfile: %v", err)
			}

			usecase := NewBuildSavepointUseCase(
				&mockParser{dockerfile: tt.dockerfile},
				&mockSlicer{content: tt.dockerfile},
				&mockWriter{},
				&mockBuilder{},
				&mockChecker{},
				&mockPusher{},
				&mockValidator{},
			)

			req := BuildSavepointRequest{
				Savepoint:  tt.savepoint,
				DryRun:     tt.dryRun,
				FilePath:   "Dockerfile",
				FullTarget: tt.fullTarget,
			}

			results, err := usecase.Execute(context.Background(), req)

			if tt.wantError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
				return
			}

			assert.NoError(t, err)
			assert.NotEmpty(t, results.Results)
			result := results.Results[0]
			assert.Equal(t, tt.wantSkipped, result.Skipped)
			assert.Equal(t, tt.wantTag, result.Tag)
			if tt.dryRun {
				assert.Contains(t, result.DockerfileOut, "FROM node:20")
			}
		})
	}
}
