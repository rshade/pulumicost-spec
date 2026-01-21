package pluginsdk_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rshade/finfocus-spec/sdk/go/pluginsdk"
)

// T015: Test WithAllowedHeaders builder method.
func TestWebConfig_WithAllowedHeaders(t *testing.T) {
	t.Run("nil input sets nil", func(t *testing.T) {
		config := pluginsdk.DefaultWebConfig().WithAllowedHeaders(nil)
		assert.Nil(t, config.AllowedHeaders)
	})

	t.Run("empty slice sets empty slice", func(t *testing.T) {
		config := pluginsdk.DefaultWebConfig().WithAllowedHeaders([]string{})
		assert.NotNil(t, config.AllowedHeaders)
		assert.Empty(t, config.AllowedHeaders)
	})

	t.Run("populated slice sets populated slice", func(t *testing.T) {
		headers := []string{"Content-Type", "Authorization"}
		config := pluginsdk.DefaultWebConfig().WithAllowedHeaders(headers)
		assert.Equal(t, headers, config.AllowedHeaders)
	})
}

// T016: Test WithExposedHeaders builder method.
func TestWebConfig_WithExposedHeaders(t *testing.T) {
	t.Run("nil input sets nil", func(t *testing.T) {
		config := pluginsdk.DefaultWebConfig().WithExposedHeaders(nil)
		assert.Nil(t, config.ExposedHeaders)
	})

	t.Run("empty slice sets empty slice", func(t *testing.T) {
		config := pluginsdk.DefaultWebConfig().WithExposedHeaders([]string{})
		assert.NotNil(t, config.ExposedHeaders)
		assert.Empty(t, config.ExposedHeaders)
	})

	t.Run("populated slice sets populated slice", func(t *testing.T) {
		headers := []string{"X-Request-ID", "Grpc-Status"}
		config := pluginsdk.DefaultWebConfig().WithExposedHeaders(headers)
		assert.Equal(t, headers, config.ExposedHeaders)
	})
}

// T017: Test defensive slice copying in builder methods.
func TestWebConfig_DefensiveSliceCopy(t *testing.T) {
	t.Run("WithAllowedHeaders makes defensive copy", func(t *testing.T) {
		original := []string{"Content-Type", "Authorization"}
		config := pluginsdk.DefaultWebConfig().WithAllowedHeaders(original)

		// Modify original slice
		original[0] = "MODIFIED"

		// Config should NOT be affected
		assert.Equal(t, "Content-Type", config.AllowedHeaders[0])
	})

	t.Run("WithExposedHeaders makes defensive copy", func(t *testing.T) {
		original := []string{"X-Request-ID", "Grpc-Status"}
		config := pluginsdk.DefaultWebConfig().WithExposedHeaders(original)

		// Modify original slice
		original[0] = "MODIFIED"

		// Config should NOT be affected
		assert.Equal(t, "X-Request-ID", config.ExposedHeaders[0])
	})
}

// T018: Test WithMaxAge builder method.
func TestWebConfig_WithMaxAge(t *testing.T) {
	t.Run("default is nil", func(t *testing.T) {
		config := pluginsdk.DefaultWebConfig()
		assert.Nil(t, config.MaxAge)
	})

	t.Run("sets value", func(t *testing.T) {
		config := pluginsdk.DefaultWebConfig().WithMaxAge(3600)
		assert.NotNil(t, config.MaxAge)
		assert.Equal(t, 3600, *config.MaxAge)
	})

	t.Run("zero disables caching", func(t *testing.T) {
		config := pluginsdk.DefaultWebConfig().WithMaxAge(0)
		assert.NotNil(t, config.MaxAge)
		assert.Equal(t, 0, *config.MaxAge)
	})

	t.Run("can override previous value", func(t *testing.T) {
		config := pluginsdk.DefaultWebConfig().WithMaxAge(3600).WithMaxAge(7200)
		assert.NotNil(t, config.MaxAge)
		assert.Equal(t, 7200, *config.MaxAge)
	})
}

// T019: Test builder method chaining.
func TestWebConfig_MethodChaining(t *testing.T) {
	t.Run("all methods can be chained", func(t *testing.T) {
		config := pluginsdk.DefaultWebConfig().
			WithWebEnabled(true).
			WithAllowedOrigins([]string{"http://localhost:3000"}).
			WithAllowCredentials(true).
			WithHealthEndpoint(true).
			WithAllowedHeaders([]string{"Content-Type", "Authorization"}).
			WithExposedHeaders([]string{"X-Request-ID"}).
			WithMaxAge(3600)

		assert.True(t, config.Enabled)
		assert.Equal(t, []string{"http://localhost:3000"}, config.AllowedOrigins)
		assert.True(t, config.AllowCredentials)
		assert.True(t, config.EnableHealthEndpoint)
		assert.Equal(t, []string{"Content-Type", "Authorization"}, config.AllowedHeaders)
		assert.Equal(t, []string{"X-Request-ID"}, config.ExposedHeaders)
		assert.NotNil(t, config.MaxAge)
		assert.Equal(t, 3600, *config.MaxAge)
	})

	t.Run("order of methods does not matter", func(t *testing.T) {
		// Different order, same result
		config := pluginsdk.DefaultWebConfig().
			WithExposedHeaders([]string{"X-Request-ID"}).
			WithAllowedHeaders([]string{"Content-Type"}).
			WithWebEnabled(true)

		assert.True(t, config.Enabled)
		assert.Equal(t, []string{"Content-Type"}, config.AllowedHeaders)
		assert.Equal(t, []string{"X-Request-ID"}, config.ExposedHeaders)
	})
}
