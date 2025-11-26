package io

import (
	"fmt"

	"github.com/hydrocode-de/tool-spec-go/validate"
)

func WriteValidationError(validationError validate.ValidationError, verbose bool) string {
	if verbose {
		return fmt.Sprintf("%s: %s\n  Field: %s\n  Expected: %s\n  Actual: %s",
			validationError.Type,
			validationError.Message,
			validationError.Field,
			validationError.Expected,
			validationError.Actual,
		)
	}

	return fmt.Sprintf("%s: %s (%s)",
		validationError.Type,
		validationError.Message,
		validationError.Field,
	)
}
