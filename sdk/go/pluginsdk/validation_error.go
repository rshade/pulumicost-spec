package pluginsdk

import "fmt"

// ValidationError represents a structured validation error with discrete fields
// for programmatic error inspection. This enables callers to parse and categorize
// validation failures without string parsing.
//
// Usage:
//
//	var valErr *ValidationError
//	if errors.As(err, &valErr) {
//	    fmt.Printf("Field %s failed: %s\n", valErr.FieldName, valErr.Constraint)
//	}
type ValidationError struct {
	// FieldName is the name of the field that failed validation.
	FieldName string

	// Constraint describes the validation rule that was violated.
	Constraint string

	// ActualValue is a string representation of the actual value found.
	ActualValue string

	// ExpectedValue is a string representation of what was expected.
	ExpectedValue string
}

// NewValidationError creates a new ValidationError with the specified fields.
// This constructor ensures consistent field population and makes it easier
// to evolve the type in the future without breaking existing code.
func NewValidationError(fieldName, constraint, actualValue, expectedValue string) *ValidationError {
	return &ValidationError{
		FieldName:     fieldName,
		Constraint:    constraint,
		ActualValue:   actualValue,
		ExpectedValue: expectedValue,
	}
}

// Error implements the error interface.
// Format: "{FieldName}: {Constraint} (actual: {ActualValue}, expected: {ExpectedValue})".
func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s (actual: %s, expected: %s)",
		e.FieldName, e.Constraint, e.ActualValue, e.ExpectedValue)
}
