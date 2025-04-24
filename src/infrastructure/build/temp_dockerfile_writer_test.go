package build

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTempDockerfileWriterImpl_Write(t *testing.T) {
	tests := []struct {
		name      string
		lines     []string
		savepoint string
		wantErr   bool
	}{
		{
			name: "successful write",
			lines: []string{
				"FROM alpine:latest",
				"RUN echo 'hello'",
			},
			savepoint: "base",
			wantErr:   false,
		},
		{
			name:      "empty content",
			lines:     []string{},
			savepoint: "empty",
			wantErr:   false,
			// For empty content, we shouldn't expect a trailing newline
		},
		{
			name: "special characters in savepoint",
			lines: []string{
				"FROM alpine:latest",
			},
			savepoint: "test-123_abc",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &TempDockerfileWriterImpl{}
			
			// Execute Write
			gotPath, err := w.Write(tt.lines, tt.savepoint)
			
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			
			// Verify no error
			assert.NoError(t, err)
			
			// Verify file exists
			_, err = os.Stat(gotPath)
			assert.NoError(t, err)
			
			// Verify file content
			content, err := os.ReadFile(gotPath)
			assert.NoError(t, err)
			
			// Verify content matches input
			expectedContent := strings.Join(tt.lines, "\n")
			if len(tt.lines) > 0 {
				expectedContent += "\n"
			}
			assert.Equal(t, expectedContent, string(content))
			
			// Verify file path
			expectedFilename := filepath.Join(".dockpoint/tmp", ".dockpoint.Dockerfile."+tt.savepoint)
			assert.Equal(t, expectedFilename, gotPath)
			
			// Cleanup
			os.Remove(gotPath)
		})
	}
}