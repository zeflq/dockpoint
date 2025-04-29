package build_savepoint

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestParseAndValidateFullTarget(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantRepo  string
		wantTag   string
		wantError string
	}{
		{
			name:      "valid full target",
			input:     "docker.io/user/app:latest",
			wantRepo:  "docker.io/user/app",
			wantTag:   "latest",
		},
		{
			name:      "empty target",
			input:     "",
			wantError: "Full target is empty",
		},
		{
			name:      "missing tag",
			input:     "docker.io/user/app",
			wantError: "missing ':tag'",
		},
		{
			name:      "missing repository path",
			input:     "myapp:latest",
			wantError: "must contain at least one '/'",
		},
		{
			name:      "empty tag",
			input:     "docker.io/user/app:",
			wantError: "Must be repo/image:tag format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, tag, err := ParseAndValidateFullTarget(tt.input)
			if tt.wantError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantError)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.wantRepo, repo)
			assert.Equal(t, tt.wantTag, tag)
		})
	}
}

func TestValidateRequest(t *testing.T) {
	uc := &BuildSavepointUseCase{}
	tests := []struct {
		name      string
		request   BuildSavepointRequest
		wantError string
	}{
		{
			name: "valid request with target",
			request: BuildSavepointRequest{
				FullTarget: "docker.io/user/app:latest",
			},
		},
		{
			name: "missing target",
			request: BuildSavepointRequest{
				FilePath: "Dockerfile",
			},
			wantError: "Missing required -t flag",
		},
		{
			name: "empty filepath gets default",
			request: BuildSavepointRequest{
				FullTarget: "docker.io/user/app:latest",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := uc.validateRequest(&tt.request)
			if tt.wantError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantError)
				return
			}
			assert.NoError(t, err)
			if tt.request.FilePath == "" {
				assert.Equal(t, "Dockerfile", tt.request.FilePath)
			}
		})
	}
}