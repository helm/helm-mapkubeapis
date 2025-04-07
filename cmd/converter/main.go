package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/helm/helm-mapkubeapis/pkg/common"
	"github.com/helm/helm-mapkubeapis/pkg/mapping"
	"sigs.k8s.io/yaml"
)

func readConfig(f string) ([]common.GenericYAML, error) {
	b, err := os.ReadFile(f)
	if err != nil {
		return nil, err
	}
	return common.ParseYAML(string(b))
}

func exitOnError(err error) {
	if err != nil {
		_, _ = os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}

func main() {
	var config string
	if len(os.Args) == 1 || strings.HasPrefix(os.Args[1], "-h") || strings.HasPrefix(os.Args[1], "--h") {
		_, _ = os.Stdout.WriteString("Usage: " + os.Args[0] + " [file]\n")
		os.Exit(0)
	}
	config = os.Args[1]
	input, err := readConfig(config)
	exitOnError(err)
	var mappings []any
	if v, ok := input[0]["mappings"].([]any); ok {
		mappings = v
	}
	if len(mappings) == 0 {
		exitOnError(fmt.Errorf("failed to read mappings"))
	}
	var output mapping.Metadata
	var m common.GenericYAML
	var s string
	var ok bool
	for _, val := range mappings {
		var om mapping.Mapping
		if m, ok = val.(common.GenericYAML); !ok {
			continue // skip invalid mappings
		}
		if s, ok = m["deprecatedAPI"].(string); ok {
			apiVersion, kind := parseAPIString(s)
			if apiVersion == "" || kind == "" {
				continue // skip invalid mappings
			}
			om.DeprecatedAPI.APIVersion = apiVersion
			om.DeprecatedAPI.Kind = kind
		}
		if s, ok = m["newAPI"].(string); ok {
			apiVersion, kind := parseAPIString(s)
			if apiVersion == "" || kind == "" {
				continue // skip invalid mappings
			}
			om.NewAPI.APIVersion = apiVersion
			om.NewAPI.Kind = kind
		}
		if s, ok = m["deprecatedInVersion"].(string); ok {
			om.DeprecatedInVersion = s
		}
		if s, ok = m["removedInVersion"].(string); ok {
			om.RemovedInVersion = s
		}
		output.Mappings = append(output.Mappings, &om)
	}
	var b []byte
	b, err = yaml.Marshal(output)
	exitOnError(err)
	_, _ = os.Stdout.Write(b)
}
