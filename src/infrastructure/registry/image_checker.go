package registry

import (
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
	cmd := exec.Command("docker", "inspect", tag)
	output, err := cmd.CombinedOutput()

	if err != nil {
		if strings.Contains(string(output), "No such image") ||
			strings.Contains(string(output), "not found") {
			return false, nil
		}
		if strings.Contains(err.Error(), "executable file not found") {
			return false, coreErr.ErrDockerUnavailable
		}
		return false, err
	}

	return true, nil
}
