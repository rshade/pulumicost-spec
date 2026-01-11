package pluginsdk_test

import (
	"testing"

	"github.com/rshade/finfocus-spec/sdk/go/pluginsdk"
)

func TestSpecVersionConstant(t *testing.T) {
	// Verify the SpecVersion constant is a valid semantic version
	if err := pluginsdk.ValidateSpecVersion(pluginsdk.SpecVersion); err != nil {
		t.Errorf("SpecVersion constant %q is not a valid semantic version: %v", pluginsdk.SpecVersion, err)
	}
}

func TestValidateSpecVersion(t *testing.T) {
	tests := []struct {
		name    string
		version string
		wantErr bool
	}{
		// Valid versions
		{name: "valid v0.0.0", version: "v0.0.0", wantErr: false},
		{name: "valid v0.4.11", version: "v0.4.11", wantErr: false},
		{name: "valid v1.0.0", version: "v1.0.0", wantErr: false},
		{name: "valid v1.2.3", version: "v1.2.3", wantErr: false},
		{name: "valid v10.20.30", version: "v10.20.30", wantErr: false},
		{name: "valid v999.999.999", version: "v999.999.999", wantErr: false},

		// Invalid versions - empty
		{name: "empty string", version: "", wantErr: true},

		// Invalid versions - missing v prefix
		{name: "no v prefix", version: "0.4.11", wantErr: true},
		{name: "no v prefix 1.0.0", version: "1.0.0", wantErr: true},

		// Invalid versions - wrong format
		{name: "missing patch", version: "v1.2", wantErr: true},
		{name: "missing minor and patch", version: "v1", wantErr: true},
		{name: "too many parts", version: "v1.2.3.4", wantErr: true},
		{name: "prerelease suffix", version: "v1.0.0-alpha", wantErr: true},
		{name: "build metadata", version: "v1.0.0+build", wantErr: true},

		// Invalid versions - leading zeros
		{name: "leading zero in major", version: "v01.2.3", wantErr: true},
		{name: "leading zero in minor", version: "v1.02.3", wantErr: true},
		{name: "leading zero in patch", version: "v1.2.03", wantErr: true},

		// Invalid versions - non-numeric
		{name: "non-numeric major", version: "va.2.3", wantErr: true},
		{name: "non-numeric minor", version: "v1.b.3", wantErr: true},
		{name: "non-numeric patch", version: "v1.2.c", wantErr: true},

		// Invalid versions - negative numbers
		{name: "negative major", version: "v-1.2.3", wantErr: true},
		{name: "negative minor", version: "v1.-2.3", wantErr: true},
		{name: "negative patch", version: "v1.2.-3", wantErr: true},

		// Invalid versions - special characters
		{name: "spaces", version: "v 1.2.3", wantErr: true},
		{name: "dots only", version: "v...", wantErr: true},
		{name: "uppercase V", version: "V1.2.3", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pluginsdk.ValidateSpecVersion(tt.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSpecVersion(%q) error = %v, wantErr %v", tt.version, err, tt.wantErr)
			}
		})
	}
}

func TestIsValidSpecVersion(t *testing.T) {
	tests := []struct {
		version string
		want    bool
	}{
		{"v0.4.11", true},
		{"v1.0.0", true},
		{"", false},
		{"1.0.0", false},
		{"v1.2", false},
	}

	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			if got := pluginsdk.IsValidSpecVersion(tt.version); got != tt.want {
				t.Errorf("IsValidSpecVersion(%q) = %v, want %v", tt.version, got, tt.want)
			}
		})
	}
}

func BenchmarkValidateSpecVersion(b *testing.B) {
	version := "v0.4.11"
	b.ResetTimer()
	for range b.N {
		_ = pluginsdk.ValidateSpecVersion(version)
	}
}

func BenchmarkIsValidSpecVersion(b *testing.B) {
	version := "v0.4.11"
	b.ResetTimer()
	for range b.N {
		_ = pluginsdk.IsValidSpecVersion(version)
	}
}
