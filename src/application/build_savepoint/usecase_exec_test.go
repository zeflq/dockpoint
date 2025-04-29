package build_savepoint

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecBuildPushPath(t *testing.T) {
	tests := []struct {
		name       string
		dockerfile string
		request    BuildSavepointRequest
		wantTag    string
		wantError  bool
		errorMsg   string
	}{
		{
			name:       "successful build and push",
			dockerfile: "FROM alpine\n# savepoint: base",
			request: BuildSavepointRequest{
				Savepoint:  "base",
				Force:      true,
				Push:       true,
				DryRun:     false,
				FullTarget: "docker.io/user/app:latest",
			},
			wantTag: "docker.io/user/app:latest", // Fixed: should match FullTarget
		},
		{
			name:       "missing target fails",
			dockerfile: "# savepoint: base\nFROM alpine",
			request: BuildSavepointRequest{
				Savepoint: "base",
				Force:     true,
				Push:      true,
			},
			wantError: true,
			errorMsg:  "Missing required -t flag",
		},
		{
			name:       "invalid target format",
			dockerfile: "# savepoint: base\nFROM alpine",
			request: BuildSavepointRequest{
				Savepoint:  "base",
				Force:      true,
				Push:      true,
				FullTarget: "invalid-format",
			},
			wantError: true,
			errorMsg:  "missing ':tag'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmp := t.TempDir()
			dockerfile := filepath.Join(tmp, "Dockerfile")
			err := os.WriteFile(dockerfile, []byte(tt.dockerfile), 0644)
			assert.NoError(t, err)

			orig, err := os.Getwd()
			assert.NoError(t, err)
			err = os.Chdir(tmp)
			assert.NoError(t, err)
			defer os.Chdir(orig)

			writer := &spyWriter{}
			builder := &spyBuilder{}
			pusher := &spyPusher{}

			uc := NewBuildSavepointUseCase(
				&parserReturningOne{},
				&simpleSlicer{},
				writer,
				builder,
				&neverTagExists{},
				pusher,
				&mockValidator{},
				&mockHasher{},
				&mockTagBuilder{},
			)

			tt.request.FilePath = dockerfile

			results, err := uc.Execute(context.Background(), tt.request)

			if tt.wantError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
				return
			}

			assert.NoError(t, err)
			result := results.Results[0]
			assert.False(t, result.Skipped)
			assert.Equal(t, tt.wantTag, result.Tag)
			assert.True(t, writer.Called)
			assert.True(t, builder.WasCalled)
			assert.True(t, pusher.WasCalled)
		})
	}
}