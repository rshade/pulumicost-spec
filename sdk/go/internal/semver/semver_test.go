package semver_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/rshade/finfocus-spec/sdk/go/internal/semver"
)

func TestIsValid(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    bool
	}{
		// Valid versions
		{"basic version", "v1.0.0", true},
		{"zero patch", "v1.0.0", true},
		{"multi-digit", "v10.20.30", true},
		{"large numbers", "v123.456.789", true},
		{"zero major", "v0.1.0", true},
		{"zero minor", "v1.0.1", true},
		{"all zeros", "v0.0.0", true},
		{"mixed digits", "v0.4.11", true},

		// Invalid versions
		{"empty string", "", false},
		{"no v prefix", "1.0.0", false},
		{"missing patch", "v1.0", false},
		{"missing minor and patch", "v1", false},
		{"extra part", "v1.2.3.4", false},
		{"leading zero major", "v01.2.3", false},
		{"leading zero minor", "v1.02.3", false},
		{"leading zero patch", "v1.2.03", false},
		{"pre-release suffix", "v1.0.0-alpha", false},
		{"build metadata", "v1.0.0+build", false},
		{"negative number", "v-1.0.0", false},
		{"text in version", "v1.a.0", false},
		{"spaces", "v1.0.0 ", false},
		{"uppercase V", "V1.0.0", false},
		{"just v", "v", false},
		{"partial version", "v1.", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := semver.IsValid(tt.version); got != tt.want {
				t.Errorf("IsValid(%q) = %v, want %v", tt.version, got, tt.want)
			}
		})
	}
}

func TestRegex(t *testing.T) {
	// Verify regex is exported and usable
	if semver.Regex == nil {
		t.Error("Regex is nil")
	}

	// Test direct regex usage
	if !semver.Regex.MatchString("v1.0.0") {
		t.Error("Regex should match v1.0.0")
	}
}

func BenchmarkIsValid(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = semver.IsValid("v0.4.11")
	}
}

func BenchmarkRegexMatch(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = semver.Regex.MatchString("v0.4.11")
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name        string
		version     string
		wantErr     bool
		errContains string // substring that error should contain
	}{
		// Valid versions - no error
		{"valid basic", "v1.0.0", false, ""},
		{"valid multi-digit", "v10.20.30", false, ""},
		{"valid zeros", "v0.0.0", false, ""},

		// Empty version - specific error
		{"empty string", "", true, "version is empty"},

		// Missing v prefix - specific error (sentinel error message)
		{"no v prefix", "1.0.0", true, "must start with 'v'"},
		{"uppercase V", "V1.0.0", true, "must start with 'v'"},

		// Invalid format - general format error (sentinel error message)
		{"missing patch", "v1.0", true, "does not match semantic versioning format"},
		{"missing minor and patch", "v1", true, "does not match semantic versioning format"},
		{"extra part", "v1.2.3.4", true, "does not match semantic versioning format"},
		{"leading zero major", "v01.2.3", true, "does not match semantic versioning format"},
		{"leading zero minor", "v1.02.3", true, "does not match semantic versioning format"},
		{"leading zero patch", "v1.2.03", true, "does not match semantic versioning format"},
		{"pre-release suffix", "v1.0.0-alpha", true, "does not match semantic versioning format"},
		{"build metadata", "v1.0.0+build", true, "does not match semantic versioning format"},
		{"text in version", "v1.a.0", true, "does not match semantic versioning format"},
		{"spaces", "v1.0.0 ", true, "does not match semantic versioning format"},
		{"just v", "v", true, "does not match semantic versioning format"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := semver.Validate(tt.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate(%q) error = %v, wantErr %v", tt.version, err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errContains != "" {
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Validate(%q) error = %q, want error containing %q",
						tt.version, err.Error(), tt.errContains)
				}
			}
		})
	}
}

func TestValidate_ErrEmptyVersion(t *testing.T) {
	err := semver.Validate("")
	if !errors.Is(err, semver.ErrEmptyVersion) {
		t.Errorf("Validate(\"\") = %v, want ErrEmptyVersion", err)
	}
}

func TestValidate_ErrMissingVPrefix(t *testing.T) {
	err := semver.Validate("1.0.0")
	if !errors.Is(err, semver.ErrMissingVPrefix) {
		t.Errorf("Validate(\"1.0.0\") = %v, want ErrMissingVPrefix", err)
	}
	// Also verify the error message contains the version
	if !strings.Contains(err.Error(), "1.0.0") {
		t.Errorf("error message should contain the version, got: %s", err.Error())
	}
}

func TestValidate_ErrInvalidFormat(t *testing.T) {
	err := semver.Validate("v1.0")
	if !errors.Is(err, semver.ErrInvalidFormat) {
		t.Errorf("Validate(\"v1.0\") = %v, want ErrInvalidFormat", err)
	}
	// Also verify the error message contains the version
	if !strings.Contains(err.Error(), "v1.0") {
		t.Errorf("error message should contain the version, got: %s", err.Error())
	}
}

func TestValidate_SentinelErrorsAreProgrammatic(t *testing.T) {
	// Verify that callers can use errors.Is for programmatic error handling
	testCases := []struct {
		name    string
		version string
		want    error
	}{
		{"empty version", "", semver.ErrEmptyVersion},
		{"missing v prefix", "1.0.0", semver.ErrMissingVPrefix},
		{"uppercase V prefix", "V1.0.0", semver.ErrMissingVPrefix},
		{"invalid format missing patch", "v1.0", semver.ErrInvalidFormat},
		{"invalid format leading zero", "v01.0.0", semver.ErrInvalidFormat},
		{"invalid format prerelease", "v1.0.0-alpha", semver.ErrInvalidFormat},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := semver.Validate(tc.version)
			if !errors.Is(err, tc.want) {
				t.Errorf("Validate(%q): errors.Is(err, %v) = false, got err = %v", tc.version, tc.want, err)
			}
		})
	}
}

func TestValidate_ErrorMessages(t *testing.T) {
	// Test specific error message formats
	tests := []struct {
		version string
		want    string
	}{
		{"", "version is empty"},
		{"1.0.0", `version must start with 'v': "1.0.0"`},
		{"v1.0", `version does not match semantic versioning format vMAJOR.MINOR.PATCH: "v1.0"`},
		{"v01.0.0", `version does not match semantic versioning format vMAJOR.MINOR.PATCH: "v01.0.0"`},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("version=%q", tt.version), func(t *testing.T) {
			err := semver.Validate(tt.version)
			if err == nil {
				t.Fatalf("Validate(%q) = nil, want error", tt.version)
			}
			if got := err.Error(); got != tt.want {
				t.Errorf("Validate(%q) error = %q, want %q", tt.version, got, tt.want)
			}
		})
	}
}

func BenchmarkValidate(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = semver.Validate("v0.4.11")
	}
}

func BenchmarkValidate_Invalid(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = semver.Validate("1.0.0")
	}
}
