package build

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/zeflq/dockpoint/src/domain"
)

type DockerfileHasherImpl struct{}

func NewDockerfileHasher() domain.DockerfileHasher {
	return &DockerfileHasherImpl{}
}

func (h *DockerfileHasherImpl) Hash(content string) (string, error) {
	hash := sha256.Sum256([]byte(content))
	shortHash := hex.EncodeToString(hash[:])[:8]
	return shortHash, nil
}
