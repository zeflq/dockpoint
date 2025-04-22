package parse

import "github.com/zeflq/dockpoint/src/domain"

type DockerfileParser interface {
	Parse(path string) ([]domain.Savepoint, error)
}
