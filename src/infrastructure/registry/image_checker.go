package registry

import (
	"fmt"
	"os/exec"
	"strings"

	coreErr "github.com/zeflq/dockpoint/src/core/errors"
	core "github.com/zeflq/dockpoint/src/core/registry"
)

type ImageCheckerImpl struct{}

func NewImageChecker() core.ImageChecker {
	return &ImageCheckerImpl{}
}

func (i *ImageCheckerImpl) TagExists(tag string) (bool, error) {
	cmd := exec.Command("docker", "manifest", "inspect", tag)
	output, err := cmd.CombinedOutput()

	if err != nil {
		outputStr := string(output)
		// Handle common manifest not found cases
		if strings.Contains(outputStr, "no such manifest") ||
			strings.Contains(outputStr, "manifest unknown") {
			return false, nil
		}
		// Handle Docker daemon not running
		if strings.Contains(err.Error(), "executable file not found") {
			return false, coreErr.ErrDockerUnavailable
		}
		// Log the full error for debugging
		fmt.Printf("🔍 Docker manifest inspect output: %s\n", outputStr)
		return false, err
	}

	return true, nil
}
