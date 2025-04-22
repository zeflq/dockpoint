package registry

// ImageChecker checks whether a given image tag exists in a remote registry.
type ImageChecker interface {
	TagExists(tag string) (bool, error)
}
