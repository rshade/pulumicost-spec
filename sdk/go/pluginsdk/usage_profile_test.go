package pluginsdk_test

import (
	"testing"

	"github.com/rshade/finfocus-spec/sdk/go/pluginsdk"
	pbc "github.com/rshade/finfocus-spec/sdk/go/proto/finfocus/v1"
)

func TestAllUsageProfiles(t *testing.T) {
	profiles := pluginsdk.AllUsageProfiles()

	// Should contain exactly 4 values
	expected := 4
	if len(profiles) != expected {
		t.Errorf("AllUsageProfiles() returned %d profiles, want %d", len(profiles), expected)
	}

	// Verify all expected profiles are present
	expectedProfiles := map[pbc.UsageProfile]bool{
		pbc.UsageProfile_USAGE_PROFILE_UNSPECIFIED: false,
		pbc.UsageProfile_USAGE_PROFILE_PROD:        false,
		pbc.UsageProfile_USAGE_PROFILE_DEV:         false,
		pbc.UsageProfile_USAGE_PROFILE_BURST:       false,
	}

	for _, p := range profiles {
		if _, ok := expectedProfiles[p]; !ok {
			t.Errorf("Unexpected profile in AllUsageProfiles(): %v", p)
		}
		expectedProfiles[p] = true
	}

	for profile, found := range expectedProfiles {
		if !found {
			t.Errorf("AllUsageProfiles() missing profile: %v", profile)
		}
	}
}

func TestIsValidUsageProfile(t *testing.T) {
	tests := []struct {
		name    string
		profile pbc.UsageProfile
		want    bool
	}{
		{"UNSPECIFIED is valid", pbc.UsageProfile_USAGE_PROFILE_UNSPECIFIED, true},
		{"PROD is valid", pbc.UsageProfile_USAGE_PROFILE_PROD, true},
		{"DEV is valid", pbc.UsageProfile_USAGE_PROFILE_DEV, true},
		{"BURST is valid", pbc.UsageProfile_USAGE_PROFILE_BURST, true},
		{"unknown value 100 is invalid", pbc.UsageProfile(100), false},
		{"unknown value 999 is invalid", pbc.UsageProfile(999), false},
		{"unknown value -1 is invalid", pbc.UsageProfile(-1), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pluginsdk.IsValidUsageProfile(tt.profile); got != tt.want {
				t.Errorf("IsValidUsageProfile(%v) = %v, want %v", tt.profile, got, tt.want)
			}
		})
	}
}

func TestParseUsageProfile(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      pbc.UsageProfile
		wantError bool
	}{
		// Valid inputs - lowercase
		{"empty string", "", pbc.UsageProfile_USAGE_PROFILE_UNSPECIFIED, false},
		{"unspecified", "unspecified", pbc.UsageProfile_USAGE_PROFILE_UNSPECIFIED, false},
		{"prod", "prod", pbc.UsageProfile_USAGE_PROFILE_PROD, false},
		{"production", "production", pbc.UsageProfile_USAGE_PROFILE_PROD, false},
		{"dev", "dev", pbc.UsageProfile_USAGE_PROFILE_DEV, false},
		{"development", "development", pbc.UsageProfile_USAGE_PROFILE_DEV, false},
		{"burst", "burst", pbc.UsageProfile_USAGE_PROFILE_BURST, false},

		// Valid inputs - mixed case
		{"PROD uppercase", "PROD", pbc.UsageProfile_USAGE_PROFILE_PROD, false},
		{"Dev mixed case", "Dev", pbc.UsageProfile_USAGE_PROFILE_DEV, false},
		{"BURST uppercase", "BURST", pbc.UsageProfile_USAGE_PROFILE_BURST, false},
		{"Production mixed", "Production", pbc.UsageProfile_USAGE_PROFILE_PROD, false},
		{"Development mixed", "Development", pbc.UsageProfile_USAGE_PROFILE_DEV, false},

		// Invalid inputs
		{"invalid profile", "invalid", pbc.UsageProfile_USAGE_PROFILE_UNSPECIFIED, true},
		{"unknown profile", "staging", pbc.UsageProfile_USAGE_PROFILE_UNSPECIFIED, true},
		{"partial match", "pr", pbc.UsageProfile_USAGE_PROFILE_UNSPECIFIED, true},
		{"whitespace", " prod ", pbc.UsageProfile_USAGE_PROFILE_UNSPECIFIED, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := pluginsdk.ParseUsageProfile(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("ParseUsageProfile(%q) error = %v, wantError %v", tt.input, err, tt.wantError)
				return
			}
			if got != tt.want {
				t.Errorf("ParseUsageProfile(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestUsageProfileString(t *testing.T) {
	tests := []struct {
		name    string
		profile pbc.UsageProfile
		want    string
	}{
		{"UNSPECIFIED", pbc.UsageProfile_USAGE_PROFILE_UNSPECIFIED, "unspecified"},
		{"PROD", pbc.UsageProfile_USAGE_PROFILE_PROD, "prod"},
		{"DEV", pbc.UsageProfile_USAGE_PROFILE_DEV, "dev"},
		{"BURST", pbc.UsageProfile_USAGE_PROFILE_BURST, "burst"},
		{"unknown value 100", pbc.UsageProfile(100), "unknown(100)"},
		{"unknown value 999", pbc.UsageProfile(999), "unknown(999)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pluginsdk.UsageProfileString(tt.profile); got != tt.want {
				t.Errorf("UsageProfileString(%v) = %q, want %q", tt.profile, got, tt.want)
			}
		})
	}
}

func TestNormalizeUsageProfile(t *testing.T) {
	tests := []struct {
		name    string
		profile pbc.UsageProfile
		want    pbc.UsageProfile
	}{
		{
			"UNSPECIFIED unchanged",
			pbc.UsageProfile_USAGE_PROFILE_UNSPECIFIED,
			pbc.UsageProfile_USAGE_PROFILE_UNSPECIFIED,
		},
		{
			"PROD unchanged",
			pbc.UsageProfile_USAGE_PROFILE_PROD,
			pbc.UsageProfile_USAGE_PROFILE_PROD,
		},
		{
			"DEV unchanged",
			pbc.UsageProfile_USAGE_PROFILE_DEV,
			pbc.UsageProfile_USAGE_PROFILE_DEV,
		},
		{
			"BURST unchanged",
			pbc.UsageProfile_USAGE_PROFILE_BURST,
			pbc.UsageProfile_USAGE_PROFILE_BURST,
		},
		{
			"unknown 100 normalized",
			pbc.UsageProfile(100),
			pbc.UsageProfile_USAGE_PROFILE_UNSPECIFIED,
		},
		{
			"unknown 999 normalized",
			pbc.UsageProfile(999),
			pbc.UsageProfile_USAGE_PROFILE_UNSPECIFIED,
		},
		{
			"unknown -1 normalized",
			pbc.UsageProfile(-1),
			pbc.UsageProfile_USAGE_PROFILE_UNSPECIFIED,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pluginsdk.NormalizeUsageProfile(tt.profile); got != tt.want {
				t.Errorf("NormalizeUsageProfile(%v) = %v, want %v", tt.profile, got, tt.want)
			}
		})
	}
}

func TestDefaultMonthlyHours(t *testing.T) {
	tests := []struct {
		name    string
		profile pbc.UsageProfile
		want    float64
	}{
		{"PROD", pbc.UsageProfile_USAGE_PROFILE_PROD, 730},
		{"DEV", pbc.UsageProfile_USAGE_PROFILE_DEV, 160},
		{"BURST", pbc.UsageProfile_USAGE_PROFILE_BURST, 200},
		{"UNSPECIFIED defaults to PROD", pbc.UsageProfile_USAGE_PROFILE_UNSPECIFIED, 730},
		{"unknown defaults to PROD", pbc.UsageProfile(999), 730},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pluginsdk.DefaultMonthlyHours(tt.profile); got != tt.want {
				t.Errorf("DefaultMonthlyHours(%v) = %v, want %v", tt.profile, got, tt.want)
			}
		})
	}
}

// Benchmarks - Target: <15 ns/op, 0 allocs/op

func BenchmarkIsValidUsageProfile_Valid(b *testing.B) {
	profile := pbc.UsageProfile_USAGE_PROFILE_PROD
	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		_ = pluginsdk.IsValidUsageProfile(profile)
	}
}

func BenchmarkIsValidUsageProfile_Invalid(b *testing.B) {
	profile := pbc.UsageProfile(999)
	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		_ = pluginsdk.IsValidUsageProfile(profile)
	}
}

func BenchmarkParseUsageProfile(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		_, _ = pluginsdk.ParseUsageProfile("prod")
	}
}

func BenchmarkUsageProfileString(b *testing.B) {
	profile := pbc.UsageProfile_USAGE_PROFILE_PROD
	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		_ = pluginsdk.UsageProfileString(profile)
	}
}

func BenchmarkNormalizeUsageProfile_Valid(b *testing.B) {
	profile := pbc.UsageProfile_USAGE_PROFILE_DEV
	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		_ = pluginsdk.NormalizeUsageProfile(profile)
	}
}

func BenchmarkDefaultMonthlyHours(b *testing.B) {
	profile := pbc.UsageProfile_USAGE_PROFILE_PROD
	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		_ = pluginsdk.DefaultMonthlyHours(profile)
	}
}

func BenchmarkAllUsageProfiles(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		_ = pluginsdk.AllUsageProfiles()
	}
}
