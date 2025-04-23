package parse

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeflq/dockpoint/src/domain"
)

func TestDockerfileParser_Parse(t *testing.T) {
	t.Run("valid savepoints", func(t *testing.T) {
		tmpFile := filepath.Join(t.TempDir(), "Dockerfile")
		content := `
# savepoint: base
FROM alpine
RUN echo "hello"
# savepoint: deps
RUN apk add curl
`
		if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write test Dockerfile: %v", err)
		}

		parser := NewDockerfileParser()
		savepoints, err := parser.Parse(tmpFile)

		assert.NoError(t, err)
		assert.Equal(t, 2, len(savepoints))
		assert.Equal(t, domain.Savepoint{Name: "base", StartLine: 1, EndLine: 3}, savepoints[0])
		assert.Equal(t, domain.Savepoint{Name: "deps", StartLine: 4, EndLine: -1}, savepoints[1])
	})

	t.Run("empty file", func(t *testing.T) {
		tmpFile := filepath.Join(t.TempDir(), "Dockerfile")
		if err := os.WriteFile(tmpFile, []byte(""), 0644); err != nil {
			t.Fatalf("Failed to write test Dockerfile: %v", err)
		}

		parser := NewDockerfileParser()
		savepoints, err := parser.Parse(tmpFile)

		assert.NoError(t, err)
		assert.Empty(t, savepoints)
	})

	t.Run("file with no savepoints", func(t *testing.T) {
		tmpFile := filepath.Join(t.TempDir(), "Dockerfile")
		content := `
FROM alpine
RUN echo "hello"
RUN apk add curl
`
		if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write test Dockerfile: %v", err)
		}

		parser := NewDockerfileParser()
		savepoints, err := parser.Parse(tmpFile)

		assert.NoError(t, err)
		assert.Empty(t, savepoints)
	})

	t.Run("invalid savepoint format", func(t *testing.T) {
		tmpFile := filepath.Join(t.TempDir(), "Dockerfile")
		content := `
# savepoint base
FROM alpine
`
		if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write test Dockerfile: %v", err)
		}

		parser := NewDockerfileParser()
		savepoints, err := parser.Parse(tmpFile)

		assert.NoError(t, err)
		assert.Empty(t, savepoints)
	})

	t.Run("non-existent file", func(t *testing.T) {
		parser := NewDockerfileParser()
		_, err := parser.Parse("non-existent-dockerfile")

		assert.Error(t, err)
	})

	t.Run("multiple savepoints with instructions", func(t *testing.T) {
		tmpFile := filepath.Join(t.TempDir(), "Dockerfile")
		content := `
# savepoint: base
FROM alpine
RUN echo "hello"
# savepoint: deps
RUN apk add curl
RUN apk add git
# savepoint: final
RUN echo "done"
`
		if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write test Dockerfile: %v", err)
		}

		parser := NewDockerfileParser()
		savepoints, err := parser.Parse(tmpFile)

		assert.NoError(t, err)
		assert.Equal(t, 3, len(savepoints))
		assert.Equal(t, domain.Savepoint{Name: "base", StartLine: 1, EndLine: 3}, savepoints[0])
		assert.Equal(t, domain.Savepoint{Name: "deps", StartLine: 4, EndLine: 6}, savepoints[1])
		assert.Equal(t, domain.Savepoint{Name: "final", StartLine: 7, EndLine: -1}, savepoints[2])
	})
}
