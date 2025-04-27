package build_savepoint

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeflq/dockpoint/src/domain"
)

// ✅ Spy mocks

type spyWriter struct {
	Called bool
	Lines  []string
}

func (s *spyWriter) Write(lines []string, savepoint string) (string, error) {
	s.Called = true
	s.Lines = lines
	return "/tmp/fake.Dockerfile", nil
}

type spyBuilder struct {
	WasCalled bool
}

func (b *spyBuilder) Build(ctx context.Context, dockerfilePath, contextDir, tag string) error {
	b.WasCalled = true
	return nil
}

type spyPusher struct {
	WasCalled bool
}

func (p *spyPusher) Push(tag string) error {
	p.WasCalled = true
	return nil
}

type alwaysFalseChecker struct{}

func (c *alwaysFalseChecker) TagExists(tag string) (bool, error) {
	return false, nil
}

type parserReturningOne struct{}

func (p *parserReturningOne) Parse(path string) ([]domain.Savepoint, error) {
	return []domain.Savepoint{
		{Name: "base", StartLine: 0, EndLine: 1},
	}, nil
}

type simpleSlicer struct{}

func (s *simpleSlicer) Slice([]string, domain.Savepoint) ([]string, error) {
	return []string{"FROM alpine"}, nil
}

// Add validator mock
type mockValidator struct{}

func (m *mockValidator) Validate(savepoints []domain.Savepoint) error {
	return nil
}

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
			dockerfile: "# savepoint: base\nFROM alpine",
			request: BuildSavepointRequest{
				Savepoint:  "base",
				Force:      true,
				Push:       true,
				DryRun:     false,
				FullTarget: "docker.io/user/app:latest",
			},
			wantTag: "docker.io/user/app:base",
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
			// Setup
			tmp := t.TempDir()
			dockerfile := filepath.Join(tmp, "Dockerfile")
			err := os.WriteFile(dockerfile, []byte(tt.dockerfile), 0644)
			assert.NoError(t, err)

			// Save original dir
			orig, err := os.Getwd()
			assert.NoError(t, err)
			err = os.Chdir(tmp)
			assert.NoError(t, err)
			defer os.Chdir(orig)

			// Setup spies
			writer := &spyWriter{}
			builder := &spyBuilder{}
			pusher := &spyPusher{}

			uc := NewBuildSavepointUseCase(
				&parserReturningOne{},
				&simpleSlicer{},
				writer,
				builder,
				&alwaysFalseChecker{},
				pusher,
				&mockValidator{}, // Fixed: using correct validator mock
			)

			// Set filepath in request
			tt.request.FilePath = dockerfile

			// Execute
			results, err := uc.Execute(context.Background(), tt.request)

			// Assert
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