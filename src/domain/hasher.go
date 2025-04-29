package domain

type DockerfileHasher interface {
	Hash(content string) (string, error)
}
