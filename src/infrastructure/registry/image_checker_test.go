package registry

import (
	"errors"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	coreErr "github.com/zeflq/dockpoint/src/core/errors"
)

type mockCmd struct {
	output []byte
	err    error
}

type testExecCommand struct {
	*exec.Cmd
	mockOutput []byte
	mockErr    error
}

func (c *testExecCommand) CombinedOutput() ([]byte, error) {
	return c.mockOutput, c.mockErr
}

func TestImageCheckerImpl_TagExists(t *testing.T) {
	tests := []struct {
		name           string
		tag            string
		mockOutput     []byte
		mockErr        error
		expectedExists bool
		expectedErr    error
	}{
		{
			name:           "registry tag exists",
			tag:            "registry.example.com/nginx:latest",
			mockOutput:     []byte("manifest found"),
			mockErr:        nil,
			expectedExists: true,
			expectedErr:    nil,
		},
		{
			name:           "registry tag does not exist - no such manifest",
			tag:            "registry.example.com/nonexistent:tag",
			mockOutput:     []byte("no such manifest"),
			mockErr:        errors.New("exit status 1"),
			expectedExists: false,
			expectedErr:    nil,
		},
		{
			name:           "tag does not exist - not found",
			tag:            "nonexistent:tag",
			mockOutput:     []byte("not found"),
			mockErr:        errors.New("exit status 1"),
			expectedExists: false,
			expectedErr:    nil,
		},
		{
			name:           "docker not available",
			tag:            "nginx:latest",
			mockOutput:     []byte(""),
			mockErr:        errors.New("executable file not found"),
			expectedExists: false,
			expectedErr:    coreErr.ErrDockerUnavailable,
		},
		{
			name:           "unknown error",
			tag:            "nginx:latest",
			mockOutput:     []byte("unknown error"),
			mockErr:        errors.New("unknown error"),
			expectedExists: false,
			expectedErr:    errors.New("unknown error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldExecCommand := execCommand
			defer func() { execCommand = oldExecCommand }()
			
			execCommand = func(name string, args ...string) *exec.Cmd {
				cmd := &testExecCommand{
					Cmd:        exec.Command("mock"),
					mockOutput: tt.mockOutput,
					mockErr:    tt.mockErr,
				}
				return cmd.Cmd
			}

			checker := &ImageCheckerImpl{}
			exists, err := checker.TagExists(tt.tag)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedExists, exists)
		})
	}
}

// Mock for exec.Command
var execCommand = exec.Command