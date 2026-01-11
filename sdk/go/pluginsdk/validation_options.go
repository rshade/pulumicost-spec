package pluginsdk

// ValidationMode controls how validation errors are handled.
type ValidationMode int

const (
	// ValidationModeFailFast stops validation on the first error (default).
	// This is the most performant mode and suitable for real-time validation.
	ValidationModeFailFast ValidationMode = iota

	// ValidationModeAggregate collects all errors before returning.
	// Use this mode for batch data quality workflows where you need
	// a complete picture of all validation issues.
	ValidationModeAggregate
)

// String returns the string representation of the ValidationMode.
func (m ValidationMode) String() string {
	switch m {
	case ValidationModeFailFast:
		return "fail_fast"
	case ValidationModeAggregate:
		return "aggregate"
	default:
		return "unknown"
	}
}

// ValidationOptions configures validation behavior.
type ValidationOptions struct {
	// Mode controls whether validation stops at the first error (FailFast)
	// or collects all errors (Aggregate). Default is FailFast.
	Mode ValidationMode
}
