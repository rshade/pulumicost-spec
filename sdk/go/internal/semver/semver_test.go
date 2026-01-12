package semver_test

import (
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
