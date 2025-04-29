package validator

import (
	"fmt"
	"regexp"

	"github.com/zeflq/dockpoint/src/domain"
)

var validSavepointName = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

type SavepointValidatorImpl struct{}

func NewSavepointValidator() domain.SavepointValidator {
	return &SavepointValidatorImpl{}
}

func (v *SavepointValidatorImpl) Validate(savepoints []domain.Savepoint) error {
	seen := make(map[string]bool)

	for _, sp := range savepoints {
		if sp.Name == "" {
			return fmt.Errorf("❌ Empty savepoint name detected")
		}

		if !validSavepointName.MatchString(sp.Name) {
			return fmt.Errorf("❌ Invalid savepoint name: %s (only a–z, A–Z, 0–9, _, - allowed)", sp.Name)
		}

		if seen[sp.Name] {
			return fmt.Errorf("❌ Duplicate savepoint name detected: %s", sp.Name)
		}
		seen[sp.Name] = true
	}

	return nil
}
