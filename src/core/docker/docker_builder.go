package docker

import "context"

// DockerBuilder is responsible for executing docker build commands.
type DockerBuilder interface {
	Build(ctx context.Context, dockerfilePath, contextDir, tag string) error
}
