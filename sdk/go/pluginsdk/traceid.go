package pluginsdk

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// GenerateTraceID generates a new valid trace ID using cryptographically secure random bytes.
// The generated ID is a 32-character lowercase hexadecimal string that conforms to
// OpenTelemetry trace ID format requirements (not all zeros).
func GenerateTraceID() (string, error) {
	const traceIDByteLength = 16 // 16 bytes = 32 hex characters
	// Generate 16 random bytes (32 hex characters)
	bytes := make([]byte, traceIDByteLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random trace ID: %w", err)
	}

	// Convert to hex string
	traceID := hex.EncodeToString(bytes)

	// Ensure it's not all zeros (extremely unlikely but check anyway)
	if traceID == "00000000000000000000000000000000" {
		// Regenerate if we somehow got all zeros
		if _, err := rand.Read(bytes); err != nil {
			return "", fmt.Errorf("failed to regenerate trace ID: %w", err)
		}
		traceID = hex.EncodeToString(bytes)
	}

	return traceID, nil
}
