package build_savepoint

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestUsesCorrectFilePathFromCLI(t *testing.T) {
	tmp := t.TempDir()
	customFile := filepath.Join(tmp, "MyDockerfile")
	_ = os.WriteFile(customFile, []byte("# savepoint: base\nFROM alpine"), 0644)

	parser := &parserSpy{}

	uc := NewBuildSavepointUseCase(
		parser,
		&dummySlicer{},
		&mockWriter{},
		&mockBuilder{},
		&mockChecker{},
		&mockPusher{},
		&mockValidator{},
		&mockHasher{},
		&mockTagBuilder{},
	)

	_, err := uc.Execute(context.Background(), BuildSavepointRequest{
		FilePath:   customFile,
		Savepoint:  "base",
		DryRun:     true,
		FullTarget: "docker.io/user/app:latest",
	})

	assert.NoError(t, err)
	assert.Equal(t, customFile, parser.Captured)
}

func TestDefaultsToDockerfileInCurrentDir(t *testing.T) {
	tmp := t.TempDir()
	_ = os.Chdir(tmp)
	_ = os.WriteFile("Dockerfile", []byte("FROM alpine"), 0644)

	parser := &parserSpy{}
	uc := NewBuildSavepointUseCase(
		parser,
		&dummySlicer{},
		&mockWriter{},
		&mockBuilder{},
		&mockChecker{},
		&mockPusher{},
		&mockValidator{},
		&mockHasher{},
		&mockTagBuilder{},
	)

	_, err := uc.Execute(context.Background(), BuildSavepointRequest{
		Savepoint:  "base",
		FullTarget: "docker.io/user/app:latest",
	})

	assert.NoError(t, err)
	assert.Equal(t, "Dockerfile", parser.Captured)
}

