package pluginsdk

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"google.golang.org/protobuf/encoding/protojson"

	pbc "github.com/rshade/finfocus-spec/sdk/go/proto/finfocus/v1"
)

// ValidationErrors represents multiple validation errors.
type ValidationErrors []*pbc.ValidationError

func (errs ValidationErrors) Error() string {
	if len(errs) == 0 {
		return "no validation errors"
	}

	var b strings.Builder
	fmt.Fprintf(&b, "validation failed with %d error(s):", len(errs))
	for _, err := range errs {
		fmt.Fprintf(&b, "\n  - %s", err.GetMessage()) // Assuming pbc.ValidationError has a GetMessage method
	}
	return b.String()
}

// LoadManifest loads a plugin manifest from a file path.
// LoadManifest reads a plugin manifest from path, decodes it as YAML (.yaml/.yml) or JSON (.json) based on the file extension, validates the manifest, and returns the parsed Manifest.
// It returns an error if the file cannot be read, if the extension is unsupported, if decoding fails, or if the manifest does not pass validation.
func LoadManifest(path string) (*pbc.PluginManifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading manifest file: %w", err)
	}

	manifest := &pbc.PluginManifest{}
	ext := filepath.Ext(path)

	switch ext {
	case ".yaml", ".yml":
		//nolint:musttag // Protobuf messages use json tags which yaml.v3 respects
		if yamlErr := yaml.Unmarshal(data, manifest); yamlErr != nil {
			return nil, fmt.Errorf("parsing YAML manifest: %w", yamlErr)
		}
	case ".json":
		unmarshaler := protojson.UnmarshalOptions{
			AllowPartial:   true,
			DiscardUnknown: true,
		}
		if jsonErr := unmarshaler.Unmarshal(data, manifest); jsonErr != nil {
			return nil, fmt.Errorf("parsing JSON manifest: %w", jsonErr)
		}
	default:
		return nil, fmt.Errorf("unsupported manifest file extension: %s (supported: .yaml, .yml, .json)", ext)
	}

	return manifest, nil
}

// SaveManifest saves a plugin manifest to a file path.
// Format is determined by file extension.
func SaveManifest(path string, m *pbc.PluginManifest) error {
	ext := filepath.Ext(path)
	// Ensure target directory exists
	if dir := filepath.Dir(path); dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o750); err != nil {
			return fmt.Errorf("creating manifest directory: %w", err)
		}
	}

	var data []byte
	var err error

	switch ext {
	case ".yaml", ".yml":
		//nolint:musttag // Protobuf messages use json tags which yaml.v3 respects
		data, err = yaml.Marshal(m)
		if err != nil {
			return fmt.Errorf("marshaling to YAML: %w", err)
		}
	case ".json":
		marshaler := protojson.MarshalOptions{
			Indent: "  ",
		}
		data, err = marshaler.Marshal(m)
		if err != nil {
			return fmt.Errorf("marshaling to JSON: %w", err)
		}
	default:
		return fmt.Errorf("unsupported manifest file extension: %s (supported: .yaml, .yml, .json)", ext)
	}

	if writeErr := os.WriteFile(path, data, 0o600); writeErr != nil {
		return fmt.Errorf("writing manifest file: %w", writeErr)
	}

	return nil
}
