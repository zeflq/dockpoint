package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeflq/dockpoint/src/domain"
)

func TestSavepointValidator_Validate(t *testing.T) {
	tests := []struct {
		name       string
		savepoints []domain.Savepoint
		wantError  bool
		errorMsg   string
	}{
		{
			name: "valid savepoints",
			savepoints: []domain.Savepoint{
				{Name: "base"},
				{Name: "deps"},
				{Name: "build"},
			},
			wantError: false,
		},
		{
			name: "valid savepoint with hyphens and underscores",
			savepoints: []domain.Savepoint{
				{Name: "base-deps"},
				{Name: "build_stage"},
			},
			wantError: false,
		},
		{
			name: "empty savepoint name",
			savepoints: []domain.Savepoint{
				{Name: ""},
			},
			wantError: true,
			errorMsg:  "Empty savepoint name detected",
		},
		{
			name: "invalid savepoint name characters",
			savepoints: []domain.Savepoint{
				{Name: "base@stage"},
			},
			wantError: true,
			errorMsg:  "Invalid savepoint name",
		},
		{
			name: "duplicate savepoint names",
			savepoints: []domain.Savepoint{
				{Name: "base"},
				{Name: "base"},
			},
			wantError: true,
			errorMsg:  "Duplicate savepoint name detected",
		},
		{
			name:       "empty savepoints list",
			savepoints: []domain.Savepoint{},
			wantError:  false,
		},
		{
			name: "savepoint with numbers",
			savepoints: []domain.Savepoint{
				{Name: "stage1"},
				{Name: "stage2"},
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := NewSavepointValidator()
			err := validator.Validate(tt.savepoints)

			if tt.wantError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				return
			}

			assert.NoError(t, err)
		})
	}
}