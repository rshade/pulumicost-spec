// Package pluginsdk provides a development SDK for PulumiCost plugins.
package pluginsdk

// DefaultAllowedHeaders contains the CORS allowed headers for Connect/gRPC-Web.
// These headers are used when WebConfig.AllowedHeaders is nil.
// The value is pre-joined for efficient use in HTTP headers.
const DefaultAllowedHeaders = "Accept, Content-Type, Content-Length, Accept-Encoding, " +
	"Authorization, X-CSRF-Token, X-Requested-With, Connect-Protocol-Version, " +
	"Connect-Timeout-Ms, Grpc-Timeout, X-Grpc-Web, X-User-Agent"

// DefaultExposedHeaders contains the CORS exposed headers for Connect/gRPC-Web.
// These headers are used when WebConfig.ExposedHeaders is nil.
// The value is pre-joined for efficient use in HTTP headers.
const DefaultExposedHeaders = "Grpc-Status, Grpc-Message, Grpc-Status-Details-Bin, " +
	"Connect-Content-Encoding, Connect-Content-Type"

// DefaultMaxAge is the default CORS preflight cache duration in seconds.
// 24 hours balances performance (fewer preflight requests) with security.
// Note: Most browsers cap max-age at 24 hours regardless of the value set.
const DefaultMaxAge = 86400

// WebConfig holds configuration for gRPC-Web and CORS support.
//
// Field Semantics:
// Many fields use a three-state pattern for flexible configuration:
//   - nil: Use sensible defaults (backward compatible)
//   - empty slice/zero value: Explicitly disable or set to empty
//   - populated: Use custom values
//
// These settings enable browser-based clients to communicate with plugins
// using the gRPC-Web protocol over HTTP/1.1 or HTTP/2.
//
// Default CORS Headers (when nil):
//
//	AllowedHeaders: Accept, Content-Type, Content-Length, Accept-Encoding,
//	  Authorization, X-CSRF-Token, X-Requested-With, Connect-Protocol-Version,
//	  Connect-Timeout-Ms, Grpc-Timeout, X-Grpc-Web, X-User-Agent
//
//	ExposedHeaders: Grpc-Status, Grpc-Message, Grpc-Status-Details-Bin,
//	  Connect-Content-Encoding, Connect-Content-Type
type WebConfig struct {
	// Enabled enables gRPC-Web support on the server.
	// When true, the server wraps the gRPC server with an HTTP handler
	// that accepts gRPC-Web requests alongside native gRPC.
	Enabled bool

	// AllowedOrigins specifies which origins are permitted for CORS.
	// This is required for browser-based clients making cross-origin requests.
	// Examples: []string{"http://localhost:3000", "https://app.example.com"}
	// If empty and Enabled is true, no origins are allowed (secure default).
	AllowedOrigins []string

	// AllowCredentials indicates whether cross-origin requests can include
	// credentials like cookies or authorization headers.
	// Default is false for security.
	AllowCredentials bool

	// EnableHealthEndpoint exposes a /healthz HTTP endpoint that returns
	// 200 OK for health checks. This is useful for load balancers and
	// orchestration platforms that prefer simple HTTP health checks
	// over gRPC health protocol.
	EnableHealthEndpoint bool

	// AllowedHeaders specifies which request headers are permitted in CORS requests.
	// This controls the Access-Control-Allow-Headers response header.
	//
	// Semantics:
	//   - nil: Use DefaultAllowedHeaders (Connect/gRPC-Web compatible set)
	//   - empty []string{}: Set empty Access-Control-Allow-Headers (simple headers only)
	//   - populated []string: Use exactly these headers, joined by ", "
	AllowedHeaders []string

	// ExposedHeaders specifies which response headers are accessible to JavaScript.
	// This controls the Access-Control-Expose-Headers response header.
	//
	// Semantics:
	//   - nil: Use DefaultExposedHeaders (gRPC status headers)
	//   - empty []string{}: Set empty Access-Control-Expose-Headers
	//   - populated []string: Use exactly these headers, joined by ", "
	ExposedHeaders []string

	// MaxAge specifies how long (in seconds) browsers can cache CORS preflight responses.
	// This controls the Access-Control-Max-Age response header.
	//
	// Semantics:
	//   - nil: Use DefaultMaxAge (86400 seconds = 24 hours)
	//   - non-nil: Use the specified value (0 disables caching)
	//
	// Lower values increase security (faster policy updates) but reduce performance.
	// Higher values improve performance but delay CORS policy changes.
	MaxAge *int
}

// DefaultWebConfig returns the default web configuration with web support disabled.
func DefaultWebConfig() WebConfig {
	return WebConfig{
		Enabled:              false,
		AllowedOrigins:       nil,
		AllowCredentials:     false,
		EnableHealthEndpoint: false,
	}
}

// WithWebEnabled returns a copy of the config with web support enabled.
func (c WebConfig) WithWebEnabled(enabled bool) WebConfig {
	c.Enabled = enabled
	return c
}

// WithAllowedOrigins returns a copy of the config with the specified allowed origins.
// A defensive copy is made to prevent mutation of the original slice.
func (c WebConfig) WithAllowedOrigins(origins []string) WebConfig {
	if origins != nil {
		c.AllowedOrigins = make([]string, len(origins))
		copy(c.AllowedOrigins, origins)
	} else {
		c.AllowedOrigins = nil
	}
	return c
}

// WithAllowCredentials returns a copy of the config with credentials support.
//
// WARNING: Enabling credentials allows cookies and authentication headers to be
// sent cross-origin. Only use with specific AllowedOrigins (never wildcard "*").
// The server will reject configurations that combine AllowCredentials with wildcard.
func (c WebConfig) WithAllowCredentials(allow bool) WebConfig {
	c.AllowCredentials = allow
	return c
}

// WithHealthEndpoint returns a copy of the config with health endpoint enabled.
func (c WebConfig) WithHealthEndpoint(enabled bool) WebConfig {
	c.EnableHealthEndpoint = enabled
	return c
}

// WithAllowedHeaders returns a copy of the config with the specified allowed headers.
// A defensive copy is made to prevent mutation of the original slice.
//
// See AllowedHeaders field documentation for nil vs empty slice semantics.
func (c WebConfig) WithAllowedHeaders(headers []string) WebConfig {
	if headers != nil {
		c.AllowedHeaders = make([]string, len(headers))
		copy(c.AllowedHeaders, headers)
	} else {
		c.AllowedHeaders = nil
	}
	return c
}

// WithExposedHeaders returns a copy of the config with the specified exposed headers.
// A defensive copy is made to prevent mutation of the original slice.
//
// See ExposedHeaders field documentation for nil vs empty slice semantics.
func (c WebConfig) WithExposedHeaders(headers []string) WebConfig {
	if headers != nil {
		c.ExposedHeaders = make([]string, len(headers))
		copy(c.ExposedHeaders, headers)
	} else {
		c.ExposedHeaders = nil
	}
	return c
}

// WithMaxAge returns a copy of the config with the specified max-age value.
// The value is in seconds and controls how long browsers cache CORS preflight responses.
//
// Common values:
//   - 0: Disable caching (useful for development)
//   - 3600: 1 hour (balanced for production)
//   - 86400: 24 hours (default, maximum for most browsers)
//
// Note: Most browsers cap max-age at 24 hours regardless of the value set.
// Negative values are treated as 0 by browsers.
func (c WebConfig) WithMaxAge(seconds int) WebConfig {
	c.MaxAge = &seconds
	return c
}
