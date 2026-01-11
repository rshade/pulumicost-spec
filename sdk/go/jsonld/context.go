package jsonld

import "net/url"

// Context defines the JSON-LD context configuration for vocabulary mapping.
//
// Context controls how protobuf fields map to JSON-LD properties, including
// Schema.org integration and custom FOCUS namespace definitions.
//
// Thread Safety: Context instances use copy-on-write semantics. All With*
// methods return a new Context instance, leaving the original unchanged.
// This makes Context safe for concurrent use after creation.
//
// Example:
//
//	ctx := jsonld.NewContext().
//	    WithRemoteContext("https://example.com/context.jsonld").
//	    WithCustomMapping("myField", "https://example.com/myField")
//	serializer := jsonld.NewSerializer(jsonld.WithContext(ctx))
type Context struct {
	schemaOrg      bool
	focusNamespace string
	customMappings map[string]interface{}
	remoteContexts []string
}

// NewContext creates a default JSON-LD context.
//
// Default configuration:
//   - Schema.org vocabulary: enabled
//   - FOCUS namespace: https://focus.finops.org/v1#
//   - No custom mappings or remote contexts
func NewContext() *Context {
	return &Context{
		schemaOrg:      true,
		focusNamespace: "https://focus.finops.org/v1#",
		customMappings: make(map[string]interface{}),
		remoteContexts: []string{},
	}
}

// WithCustomMapping adds a custom property mapping to the context.
//
// field is the protobuf field name, iri is the JSON-LD property IRI.
// Custom mappings override default vocabulary definitions.
//
// Thread Safety: This method returns a new Context instance with the
// updated configuration, leaving the original unchanged (copy-on-write).
func (c *Context) WithCustomMapping(field, iri string) *Context {
	// Copy-on-write: create a new instance to avoid race conditions
	copied := &Context{
		schemaOrg:      c.schemaOrg,
		focusNamespace: c.focusNamespace,
		customMappings: make(map[string]interface{}, len(c.customMappings)+1),
		remoteContexts: make([]string, len(c.remoteContexts)),
	}
	// Deep copy the map
	for k, v := range c.customMappings {
		copied.customMappings[k] = v
	}
	// Copy the slice
	copy(copied.remoteContexts, c.remoteContexts)

	// Add the new mapping
	copied.customMappings[field] = iri
	return copied
}

// WithRemoteContext adds a remote context URL to the context.
//
// Remote contexts are loaded from external URLs and merged with inline definitions.
//
// Thread Safety: This method returns a new Context instance with the
// updated configuration, leaving the original unchanged (copy-on-write).
func (c *Context) WithRemoteContext(urlStr string) *Context {
	// Copy-on-write: create a new instance to avoid race conditions
	copied := &Context{
		schemaOrg:      c.schemaOrg,
		focusNamespace: c.focusNamespace,
		customMappings: make(map[string]interface{}, len(c.customMappings)),
		remoteContexts: make([]string, len(c.remoteContexts), len(c.remoteContexts)+1),
	}
	// Deep copy the map
	for k, v := range c.customMappings {
		copied.customMappings[k] = v
	}
	// Copy the slice
	copy(copied.remoteContexts, c.remoteContexts)

	// Add the new remote context
	copied.remoteContexts = append(copied.remoteContexts, urlStr)
	return copied
}

// Validate checks the context configuration for errors.
//
// Returns an error if:
//   - Remote context URLs are invalid (must be valid HTTP(S) URLs)
func (c *Context) Validate() error {
	// Validate remote context URLs
	for _, url := range c.remoteContexts {
		if !isValidURL(url) {
			return &ValidationError{
				Field:      "remoteContexts",
				Message:    "invalid URL: " + url,
				Suggestion: "provide a valid HTTP(S) URL",
			}
		}
	}

	return nil
}

// isValidURL checks if a string is a valid HTTP(S) URL.
// Validates using net/url.Parse to ensure proper URL structure.
func isValidURL(s string) bool {
	u, err := url.Parse(s)
	if err != nil {
		return false
	}
	// Must be absolute URL with http or https scheme and non-empty host
	if !u.IsAbs() {
		return false
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}
	return u.Host != ""
}

// Build generates the @context value for JSON-LD serialization.
//
// Per JSON-LD 1.1, when remote contexts are present, returns an array where
// remote context URLs come first, followed by the inline context object.
// When no remote contexts exist, returns just the inline context object.
//
// Returns:
//   - []interface{} when remote contexts are configured (array form)
//   - map[string]interface{} when no remote contexts (object form)
func (c *Context) Build() interface{} {
	// Build the inline context object
	inline := make(map[string]interface{})

	// Add Schema.org vocabulary (if enabled)
	if c.schemaOrg {
		inline["schema"] = "https://schema.org/"
	}

	// Add FOCUS namespace
	inline["focus"] = c.focusNamespace

	// Add XSD for type coercions
	inline["xsd"] = "http://www.w3.org/2001/XMLSchema#"

	// Add custom mappings
	for field, mapping := range c.customMappings {
		inline[field] = mapping
	}

	// If no remote contexts, return just the inline object
	if len(c.remoteContexts) == 0 {
		return inline
	}

	// Per JSON-LD 1.1: remote contexts + inline object as array
	// Remote URLs come first, inline object comes last
	result := make([]interface{}, 0, len(c.remoteContexts)+1)
	for _, url := range c.remoteContexts {
		result = append(result, url)
	}
	result = append(result, inline)

	return result
}
