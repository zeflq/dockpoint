package errors

import "errors"

var (
	// ErrSavepointNotFound is returned when the requested savepoint doesn't exist.
	ErrSavepointNotFound = errors.New("savepoint not found")

	// ErrRepoMissing is returned when the .dockpointrc.json doesn't specify a repo.
	ErrRepoMissing = errors.New("repo not defined")

	// ErrDockerUnavailable is returned when docker CLI is not available on the system.
	ErrDockerUnavailable = errors.New("docker CLI not available")
)
