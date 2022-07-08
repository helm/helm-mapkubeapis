package common

import (
	"log"
	"strings"

	"github.com/helm/helm-mapkubeapis/pkg/mapping"
	"github.com/pkg/errors"
	"golang.org/x/mod/semver"
	"sigs.k8s.io/yaml"
)

type Manifest struct {
	APIVersion     string                 `json:"apiVersion"`
	Kind           string                 `json:"kind"`
	RestOfManifest map[string]interface{} `json:"-"`
}

func unmarshallManifest(manifest string) (*Manifest, error) {
	m := Manifest{}

	if err := yaml.Unmarshal([]byte(manifest), &m); err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal([]byte(manifest), &m.RestOfManifest); err != nil {
		return nil, err
	}
	delete(m.RestOfManifest, "kind")
	delete(m.RestOfManifest, "apiVersion")

	return &m, nil
}

func (m *Manifest) marshall() (string, error) {
	nm, err := yaml.Marshal(m)
	if err != nil {
		return "", err
	}
	newManifest := string(nm)

	rom, err := yaml.Marshal(m.RestOfManifest)
	if err != nil {
		return "", err
	}
	restOfManifest := string(rom)
	return "---\n" + newManifest + restOfManifest, nil
}

func CheckForOldAPI(manifest string, mapping *mapping.Mapping) (string, error) {
	m, err := unmarshallManifest(manifest)
	if err != nil {
		return "", err
	}

	var apiVersionStr string
	if mapping.DeprecatedInVersion != "" {
		apiVersionStr = mapping.DeprecatedInVersion
	} else {
		apiVersionStr = mapping.RemovedInVersion
	}
	if !semver.IsValid(apiVersionStr) {
		return "", errors.Errorf("Failed to get the deprecated or removed Kubernetes version for API: %s %s", mapping.DeprecatedAPI.APIVersion, mapping.DeprecatedAPI.Kind)
	}

	if m.APIVersion == mapping.DeprecatedAPI.APIVersion && m.Kind == mapping.DeprecatedAPI.Kind {
		log.Printf(
			"Found instance of deprecated or removed Kubernetes API:\n\"%s/%s\"\nSupported API equivalent:\n\"%s/%s\"\n",
			mapping.DeprecatedAPI.APIVersion,
			mapping.DeprecatedAPI.Kind,
			mapping.NewAPI.APIVersion,
			mapping.NewAPI.Kind,
		)
		m.APIVersion = mapping.NewAPI.APIVersion
		m.Kind = mapping.NewAPI.Kind
	}

	return m.marshall()
}

func SeparateDocuments(manifest string) []string {
	lines := strings.Split(manifest, "\n")
	documents := []string{}
	currentDocument := []string{}
	for _, line := range lines {
		if line != "---" {
			currentDocument = append(currentDocument, line)
		} else {
			d := strings.Join(currentDocument, "\n")
			documents = append(documents, d)
			currentDocument = []string{}
		}
	}

	return documents
}
