package build

import (
	"testing"

	"github.com/stretchr/testify/assert"
	coreErr "github.com/zeflq/dockpoint/src/core/errors"
	"github.com/zeflq/dockpoint/src/domain"
)

func TestDockerfileSlicer_Slice(t *testing.T) {
	lines := []string{
		"# savepoint: base",
		"FROM node:20",
		"",
		"# savepoint: deps",
		"COPY package.json ./",
		"RUN npm ci",
	}

	slicer := NewDockerfileSlicer()

	t.Run("valid slice", func(t *testing.T) {
		savepoint := domain.Savepoint{Name: "deps", StartLine: 3, EndLine: 5}
		result, err := slicer.Slice(lines, savepoint)
		assert.NoError(t, err)
		assert.Equal(t, []string{
			"# savepoint: deps",
			"COPY package.json ./",
			"RUN npm ci",
		}, result)
	})

	t.Run("end -1 slices to EOF", func(t *testing.T) {
		savepoint := domain.Savepoint{Name: "deps", StartLine: 3, EndLine: -1}
		result, err := slicer.Slice(lines, savepoint)
		assert.NoError(t, err)
		assert.Equal(t, 3, len(result))
	})

	t.Run("start out of bounds", func(t *testing.T) {
		savepoint := domain.Savepoint{Name: "oops", StartLine: 10, EndLine: -1}
		_, err := slicer.Slice(lines, savepoint)
		assert.ErrorIs(t, err, coreErr.ErrSavepointNotFound)
	})

	t.Run("end < start", func(t *testing.T) {
		savepoint := domain.Savepoint{Name: "bad", StartLine: 3, EndLine: 2}
		_, err := slicer.Slice(lines, savepoint)
		assert.ErrorIs(t, err, coreErr.ErrSavepointNotFound)
	})
}
