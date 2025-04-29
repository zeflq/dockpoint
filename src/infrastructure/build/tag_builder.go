package build

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/zeflq/dockpoint/src/domain"
)

type TagBuilderImpl struct{}

func NewTagBuilder() domain.TagBuilder {
	return &TagBuilderImpl{}
}

func (tb *TagBuilderImpl) BuildFinalTag(ctx domain.TagContext) (string, error) {
	if ctx.IsLast {
		// Last savepoint: use provided baseTag directly
		return ctx.Repo + ":" + ctx.FinalImageTag, nil
	}

	// Hash the already prepared Dockerfile content
	hash := sha256.Sum256([]byte(ctx.FileContent))
	shortHash := hex.EncodeToString(hash[:])[:8]

	// Build final tag: repo:savepointName-hash
	return ctx.Repo + ":" + ctx.SavepointName + "-" + shortHash, nil
}