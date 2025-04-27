package validator

import (
	"fmt"
	"strings"

	"github.com/zeflq/dockpoint/src/domain"
)

type SavepointValidatorImpl struct{}

func NewSavepointValidator() domain.SavepointValidator {
	return &SavepointValidatorImpl{}
}

func (v *SavepointValidatorImpl) Validate(savepoints []domain.Savepoint) error {
	seen := make(map[string]bool)
	firstFromSeen := false

	for _, sp := range savepoints {
		if seen[sp.Name] {
			return fmt.Errorf("❌ Duplicate savepoint name detected: %s", sp.Name)
		}
		seen[sp.Name] = true

		if strings.HasPrefix(strings.ToLower(sp.Name), "base") {
			firstFromSeen = true
		}
	}

	if !firstFromSeen {
		return fmt.Errorf("❌ No FROM instruction detected before savepoints")
	}

	return nil
}
