package build_savepoint

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeflq/dockpoint/src/domain"
	"github.com/zeflq/dockpoint/src/infrastructure/build"
)

type emptyParser struct{}

func (p *emptyParser) Parse(string) ([]domain.Savepoint, error) {
	return []domain.Savepoint{}, nil
}

func TestFallbackBehavior(t *testing.T) {
	tests := []struct {
		name           string
		dockerfile     string
		fullTarget     string
		expectedTag    string
		expectedOutput string
		wantError      bool
		errorContains  string
	}{
		{
			name:           "builds full dockerfile if no savepoints",
			dockerfile:     "FROM scratch\nRUN echo hi",
			fullTarget:     "docker.io/user/app:1.0",
			expectedTag:    "docker.io/user/app:1.0",
			expectedOutput: "FROM scratch\nRUN echo hi",
		},
		{
			name:          "fails without target",
			dockerfile:    "FROM scratch",
			wantError:     true,
			errorContains: "Missing required -t flag",
		},
		{
			name:           "uses custom target tag",
			dockerfile:     "FROM alpine",
			fullTarget:     "custom.registry/app:v2",
			expectedTag:    "custom.registry/app:v2",
			expectedOutput: "FROM alpine",
		},
		{
			name:          "fails with invalid target format",
			dockerfile:    "FROM alpine",
			fullTarget:    "invalid-format",
			wantError:     true,
			errorContains: "missing ':tag'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup temp directory
			tmp := t.TempDir()
			fullFile := filepath.Join(tmp, "Dockerfile")
			err := os.WriteFile(fullFile, []byte(tt.dockerfile), 0644)
			assert.NoError(t, err)

			// Save and change directory
			original, err := os.Getwd()
			assert.NoError(t, err)
			err = os.Chdir(tmp)
			assert.NoError(t, err)
			defer os.Chdir(original)

			uc := NewBuildSavepointUseCase(
				&emptyParser{},
				build.NewDockerfileSlicer(),
				&mockWriter{},
				&mockBuilder{},
				&mockChecker{},
				&mockPusher{},
			)

			result, err := uc.Execute(context.Background(), BuildSavepointRequest{
				FilePath:   fullFile,
				DryRun:     true,
				FullTarget: tt.fullTarget,
			})

			if tt.wantError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
				return
			}

			assert.NoError(t, err)
			assert.True(t, result.Skipped)
			assert.Equal(t, tt.expectedTag, result.Tag)
			assert.Equal(t, tt.expectedOutput, result.DockerfileOut)
		})
	}
}