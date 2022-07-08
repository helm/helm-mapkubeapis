package common

import (
	"log"
	"strings"

	"github.com/helm/helm-mapkubeapis/pkg/mapping"
	"github.com/pkg/errors"
	"golang.org/x/mod/semver"
	"sigs.k8s.io/yaml"
)

// Manifest is used to grab the minimum information from a resources.
// The Rest
type Manifest struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`

	// RestOfManifest will hold all other data in the resources
	// yaml document. For example, everything in spec, metadata,
	// or any other key in the document.
	RestOfManifest map[string]interface{} `json:"-"`
}

// UnmarshallManifest will parse the apiVersion and kind from a yaml document,
// then parse the rest of the document. Finally, we'll put it into a Manifest
// that we can easily use to check if resources need updates.
func UnmarshallManifest(manifest string) (*Manifest, error) {
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

// Marshall will transform a manifest into a string with a prepended
// yaml doc seperator ---
func (m *Manifest) Marshall() (string, error) {
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

// CheckAllMappings takes a yaml document, not the entire manifest, and checks if any of the mappings
// will update it.
func CheckAllMappings(yamlDocument string, mappings []*mapping.Mapping, kubeVersionStr string) (string, error) {
	m, err := UnmarshallManifest(yamlDocument)
	if err != nil {
		return "", err
	}

	for _, mapping := range mappings {
		updated, err := CheckForOldAPI(m, mapping, kubeVersionStr)
		if err != nil {
			log.Println(err)
		}

		if updated {
			break
		}
	}
	return m.Marshall()
}

// CheckForOldAPI will check one yaml document (one kube resources) and see if it needs to be
// updated given a one mapping.
func CheckForOldAPI(manifest *Manifest, mapping *mapping.Mapping, kubeVersionStr string) (bool, error) {
	updated := false
	var apiVersionStr string
	if mapping.DeprecatedInVersion != "" {
		apiVersionStr = mapping.DeprecatedInVersion
	} else {
		apiVersionStr = mapping.RemovedInVersion
	}
	if !semver.IsValid(apiVersionStr) {
		return false, errors.Errorf("Failed to get the deprecated or removed Kubernetes version for API: %s %s", mapping.DeprecatedAPI.APIVersion, mapping.DeprecatedAPI.Kind)
	}

	if semver.Compare(apiVersionStr, kubeVersionStr) > 0 {
		log.Printf("The resource %s/%s does not require mapping as the "+
			"API is not deprecated or removed in Kubernetes '%s'\n",
			mapping.DeprecatedAPI.APIVersion,
			mapping.DeprecatedAPI.Kind,
			kubeVersionStr,
		)
		return false, nil
	}

	if manifest.APIVersion == mapping.DeprecatedAPI.APIVersion && manifest.Kind == mapping.DeprecatedAPI.Kind {
		log.Printf(
			"Found instance of deprecated or removed Kubernetes API:%s kind:%s\nSupported API:%s kind:%s",
			mapping.DeprecatedAPI.APIVersion,
			mapping.DeprecatedAPI.Kind,
			mapping.NewAPI.APIVersion,
			mapping.NewAPI.Kind,
		)
		updated = true
		manifest.APIVersion = mapping.NewAPI.APIVersion
		manifest.Kind = mapping.NewAPI.Kind
	}

	return updated, nil
}

// SeparateDocuments splits one yaml file into several documents, separated by
// ---
func SeparateDocuments(manifest string) []string {
	lines := strings.Split(manifest, "\n")
	documents := []string{}
	document := []string{}
	for _, line := range lines {
		if line == "---" && len(document) > 1 {
			documents = append(documents, strings.Join(document, "\n")+"\n")
			document = []string{}
		}
		document = append(document, line)
	}

	// Add the last doc
	documents = append(documents, strings.Join(document, "\n"))

	return documents
}
