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
			name: "valid savepoints with base",
			savepoints: []domain.Savepoint{
				{Name: "base", StartLine: 0, EndLine: 1},
				{Name: "deps", StartLine: 2, EndLine: 3},
			},
			wantError: false,
		},
		{
			name: "duplicate savepoint names",
			savepoints: []domain.Savepoint{
				{Name: "base", StartLine: 0, EndLine: 1},
				{Name: "base", StartLine: 2, EndLine: 3},
			},
			wantError: true,
			errorMsg:  "Duplicate savepoint name detected: base",
		},
		{
			name: "missing base savepoint",
			savepoints: []domain.Savepoint{
				{Name: "deps", StartLine: 0, EndLine: 1},
				{Name: "build", StartLine: 2, EndLine: 3},
			},
			wantError: true,
			errorMsg:  "No FROM instruction detected before savepoints",
		},
		{
			name:       "empty savepoints list",
			savepoints: []domain.Savepoint{},
			wantError:  true,
			errorMsg:   "No FROM instruction detected before savepoints",
		},
		{
			name: "base-prefixed savepoint is valid",
			savepoints: []domain.Savepoint{
				{Name: "base-deps", StartLine: 0, EndLine: 1},
				{Name: "build", StartLine: 2, EndLine: 3},
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