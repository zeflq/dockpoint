package build

import (
	"fmt"
	"os"
	"path/filepath"

	core "github.com/zeflq/dockpoint/src/core/build"
)

type TempDockerfileWriterImpl struct{}

func NewTempDockerfileWriter() core.TempDockerfileWriter {
	return &TempDockerfileWriterImpl{}
}

func (w *TempDockerfileWriterImpl) Write(lines []string, savepoint string) (string, error) {
	tmpDir := ".dockpoint/tmp"
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return "", err
	}

	filename := fmt.Sprintf(".dockpoint.Dockerfile.%s", savepoint)
	fullPath := filepath.Join(tmpDir, filename)


	content := ""
	for _, line := range lines {
		content += line + "\n"
	}

	err := os.WriteFile(fullPath, []byte(content), 0644)
	if err != nil {
		return "", err
	}

	return fullPath, nil
}
