package registry

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	coreErr "github.com/zeflq/dockpoint/src/core/errors"
)

func TestImageChecker_TagExists(t *testing.T) {
    // Create a fake docker executable in a temporary directory
    tmpDir := t.TempDir()
    dockerPath := filepath.Join(tmpDir, "docker")
    script := `#!/bin/sh
# Fake docker manifest inspect simulator
tag="$3"
case "$tag" in
  exists)
    exit 0
    ;;
  nosuch)
    echo "no such manifest for $tag" >&2
    exit 1
    ;;
  unknown)
    echo "manifest unknown: $tag" >&2
    exit 1
    ;;
  othererror)
    echo "unexpected error" >&2
    exit 1
    ;;
  *)
    exit 0
    ;;
esac
`
    if err := os.WriteFile(dockerPath, []byte(script), 0755); err != nil {
        t.Fatalf("failed to write fake docker binary: %v", err)
    }

    // Use the fake docker script in PATH
    t.Setenv("PATH", tmpDir)

    checker := NewImageChecker()

    t.Run("manifest exists", func(t *testing.T) {
        ok, err := checker.TagExists("exists")
        if err != nil {
            t.Fatalf("expected no error, got %v", err)
        }
        if !ok {
            t.Fatalf("expected exists true, got false")
        }
    })

    t.Run("no such manifest", func(t *testing.T) {
        ok, err := checker.TagExists("nosuch")
        if err != nil {
            t.Fatalf("expected no error for nosuch, got %v", err)
        }
        if ok {
            t.Fatalf("expected exists false for nosuch, got true")
        }
    })

    t.Run("manifest unknown", func(t *testing.T) {
        ok, err := checker.TagExists("unknown")
        if err != nil {
            t.Fatalf("expected no error for unknown, got %v", err)
        }
        if ok {
            t.Fatalf("expected exists false for unknown, got true")
        }
    })

    t.Run("other error", func(t *testing.T) {
        ok, err := checker.TagExists("othererror")
        if err == nil {
            t.Fatal("expected error for othererror, got nil")
        }
        if ok {
            t.Fatalf("expected exists false for othererror, got true")
        }
        if errors.Is(err, coreErr.ErrDockerUnavailable) {
            t.Fatalf("expected generic error, got ErrDockerUnavailable")
        }
    })

    t.Run("docker unavailable", func(t *testing.T) {
        // Unset PATH to simulate docker missing
        t.Setenv("PATH", "")
        ok, err := checker.TagExists("anytag")
        if !errors.Is(err, coreErr.ErrDockerUnavailable) {
            t.Fatalf("expected ErrDockerUnavailable, got %v", err)
        }
        if ok {
            t.Fatalf("expected exists false when docker unavailable, got true")
        }
    })
}
