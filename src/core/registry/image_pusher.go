package registry

// ImagePusher is responsible for pushing an image to the remote registry.
type ImagePusher interface {
	Push(tag string) error
}
