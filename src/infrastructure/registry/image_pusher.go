package registry

import (
	"os"
	"os/exec"

	coreErr "github.com/zeflq/dockpoint/src/core/errors"
	core "github.com/zeflq/dockpoint/src/core/registry"
)

type ImagePusherImpl struct{}

func NewImagePusher() core.ImagePusher {
	return &ImagePusherImpl{}
}

func (p *ImagePusherImpl) Push(tag string) error {
	cmd := exec.Command("docker", "push", tag)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if execErr, ok := err.(*exec.Error); ok && execErr.Err.Error() == "executable file not found in $PATH" {
			return coreErr.ErrDockerUnavailable
		}
		return err
	}

	return nil
}