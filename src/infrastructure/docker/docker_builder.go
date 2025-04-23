package docker

import (
	"context"
	"os"
	"os/exec"

	core "github.com/zeflq/dockpoint/src/core/docker"
	coreErr "github.com/zeflq/dockpoint/src/core/errors"
)

type DockerBuilderImpl struct{}

func NewDockerBuilder() core.DockerBuilder {
	return &DockerBuilderImpl{}
}

func (b *DockerBuilderImpl) Build(ctx context.Context, dockerfilePath, contextDir, tag string) error {
	cmd := exec.CommandContext(ctx, "docker", "build",
		"-f", dockerfilePath,
		"-t", tag,
		contextDir,
	)

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