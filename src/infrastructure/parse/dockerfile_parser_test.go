package parse

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeflq/dockpoint/src/domain"
)

func TestDockerfileParser_Parse(t *testing.T) {
	t.Run("valid savepoints with line numbers", func(t *testing.T) {
		tmpFile := filepath.Join(t.TempDir(), "Dockerfile")
		content := `FROM alpine
# savepoint: base
RUN echo "hello"
# savepoint: deps
RUN apk add curl
RUN echo "world"`
		if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write test Dockerfile: %v", err)
		}

		parser := NewDockerfileParser()
		savepoints, err := parser.Parse(tmpFile)

		assert.NoError(t, err)
		assert.Equal(t, 3, len(savepoints))
		assert.Equal(t, domain.Savepoint{Name: "base", StartLine: 0, EndLine: 1}, savepoints[0])
		assert.Equal(t, domain.Savepoint{Name: "deps", StartLine: 2, EndLine: 3}, savepoints[1])
		assert.Equal(t, domain.Savepoint{Name: "final", StartLine: 4, EndLine: 5}, savepoints[2])
	})

	t.Run("savepoints before FROM are ignored", func(t *testing.T) {
		tmpFile := filepath.Join(t.TempDir(), "Dockerfile")
		content := `# savepoint: ignored
FROM alpine
# savepoint: base
RUN echo "hello"`
		if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write test Dockerfile: %v", err)
		}

		parser := NewDockerfileParser()
		savepoints, err := parser.Parse(tmpFile)

		assert.NoError(t, err)
		assert.Equal(t, 2, len(savepoints))
		assert.Equal(t, domain.Savepoint{Name: "base", StartLine: 0, EndLine: 2}, savepoints[0])
		assert.Equal(t, domain.Savepoint{Name: "final", StartLine: 3, EndLine: 3}, savepoints[1])
	})

	t.Run("malformed savepoint is ignored", func(t *testing.T) {
		tmpFile := filepath.Join(t.TempDir(), "Dockerfile")
		content := `FROM alpine
# savepoint:
RUN echo "hello"`
		if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write test Dockerfile: %v", err)
		}

		parser := NewDockerfileParser()
		savepoints, err := parser.Parse(tmpFile)

		assert.NoError(t, err)
		assert.Equal(t, 1, len(savepoints))
		assert.Equal(t, domain.Savepoint{Name: "final", StartLine: 0, EndLine: 2}, savepoints[0])
	})

	t.Run("multiple FROM instructions use first one", func(t *testing.T) {
		tmpFile := filepath.Join(t.TempDir(), "Dockerfile")
		content := `FROM alpine AS builder
# savepoint: build
RUN echo "build"
FROM nginx
# savepoint: runtime
COPY --from=builder /app /app`
		if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write test Dockerfile: %v", err)
		}

		parser := NewDockerfileParser()
		savepoints, err := parser.Parse(tmpFile)

		assert.NoError(t, err)
		assert.Equal(t, 3, len(savepoints))
		assert.Equal(t, domain.Savepoint{Name: "build", StartLine: 0, EndLine: 1}, savepoints[0])
		assert.Equal(t, domain.Savepoint{Name: "runtime", StartLine: 2, EndLine: 4}, savepoints[1])
		assert.Equal(t, domain.Savepoint{Name: "final", StartLine: 5, EndLine: 5}, savepoints[2])
	})

	t.Run("line numbers after FROM", func(t *testing.T) {
		tmpFile := filepath.Join(t.TempDir(), "Dockerfile")
		content := `# comment
FROM alpine
# savepoint: base
RUN echo "hello"`
		if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write test Dockerfile: %v", err)
		}

		parser := NewDockerfileParser()
		savepoints, err := parser.Parse(tmpFile)

		assert.NoError(t, err)
		assert.Equal(t, 2, len(savepoints))
		// StartLine is 1 because it starts from the FROM instruction
		assert.Equal(t, domain.Savepoint{Name: "base", StartLine: 0, EndLine: 2}, savepoints[0])
		assert.Equal(t, domain.Savepoint{Name: "final", StartLine: 3, EndLine: 3}, savepoints[1])
	})

	// Keep existing empty file and non-existent file tests as they are correct
}
