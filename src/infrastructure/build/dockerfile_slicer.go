package build

import (
	core "github.com/zeflq/dockpoint/src/core/build"
	coreErr "github.com/zeflq/dockpoint/src/core/errors"
	"github.com/zeflq/dockpoint/src/domain"
)

type DockerfileSlicerImpl struct{}

func NewDockerfileSlicer() core.DockerfileSlicer {
	return &DockerfileSlicerImpl{}
}

func (s *DockerfileSlicerImpl) Slice(lines []string, savepoint domain.Savepoint) ([]string, error) {
	if savepoint.StartLine < 0 || savepoint.StartLine >= len(lines) {
		return nil, coreErr.ErrSavepointNotFound
	}

	end := savepoint.EndLine
	if end == -1 || end >= len(lines) {
		end = len(lines) - 1
	}

	if end < savepoint.StartLine {
		return nil, coreErr.ErrSavepointNotFound
	}

	return lines[savepoint.StartLine : end+1], nil
}
