package common

import (
	"io"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3" // for multi-doc support
)

// GenericYAML represents generic YAML document
type GenericYAML map[string]any

// ParseYAML parses a YAML string into a slice of GenericYAML documents
func ParseYAML(s string) ([]GenericYAML, error) {
	decoder := yaml.NewDecoder(strings.NewReader(s))
	var docs []GenericYAML
	for {
		var y GenericYAML
		err := decoder.Decode(&y)
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return nil, errors.Wrap(err, "Failed to decode YAML")
		}
		docs = append(docs, y)
	}
	return docs, nil
}

// EncodeYAML encodes a slice of GenericYAML documents into a multi-doc YAML string
func EncodeYAML(docs []GenericYAML) (string, error) {
	var sb strings.Builder
	encoder := yaml.NewEncoder(&sb)
	encoder.SetIndent(2) // match test cases
	for _, doc := range docs {
		if err := encoder.Encode(doc); err != nil {
			return "", errors.Wrap(err, "Failed to encode document")
		}
	}
	if err := encoder.Close(); err != nil {
		return "", errors.Wrap(err, "Failed to close encoder")
	}
	return "---\n" + sb.String(), nil // always start with a document separator
}
