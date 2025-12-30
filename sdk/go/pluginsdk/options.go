// Package pluginsdk provides a development SDK for PulumiCost plugins.
package pluginsdk

// WebConfig holds configuration for gRPC-Web and CORS support.
// These settings enable browser-based clients to communicate with plugins
// using the gRPC-Web protocol over HTTP/1.1 or HTTP/2.
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
