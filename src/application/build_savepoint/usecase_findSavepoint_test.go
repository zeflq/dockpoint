package build_savepoint

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeflq/dockpoint/src/core/errors"
	"github.com/zeflq/dockpoint/src/domain"
)

func TestFindSavepoint(t *testing.T) {
	savepoints := []domain.Savepoint{
		{Name: "base", StartLine: 0, EndLine: 1},
		{Name: "deps", StartLine: 2, EndLine: 3},
	}

	tests := []struct {
		name      string
		search    string
		wantIndex int
		wantError error
	}{
		{
			name:      "find existing savepoint",
			search:    "base",
			wantIndex: 0,
		},
		{
			name:      "find second savepoint",
			search:    "deps",
			wantIndex: 1,
		},
		{
			name:      "savepoint not found",
			search:    "nonexistent",
			wantIndex: -1,
			wantError: errors.ErrSavepointNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sp, idx, err := FindSavepoint(savepoints, tt.search)
			if tt.wantError != nil {
				assert.ErrorIs(t, err, tt.wantError)
				assert.Nil(t, sp)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.wantIndex, idx)
			assert.Equal(t, tt.search, sp.Name)
		})
	}
}