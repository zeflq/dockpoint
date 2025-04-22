package build

import "github.com/zeflq/dockpoint/src/domain"

// DockerfileSlicer slices a Dockerfile up to a given savepoint.
type DockerfileSlicer interface {
	Slice(lines []string, savepoint domain.Savepoint) ([]string, error)
}
