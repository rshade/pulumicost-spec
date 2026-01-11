// Copyright 2026 PulumiCost/FinFocus Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
//	    fmt.Printf("Field %s failed: %s
", valErr.FieldName, valErr.Constraint)
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
