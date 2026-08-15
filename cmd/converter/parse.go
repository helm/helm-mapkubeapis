package main

import "strings"

const (
	apiVersionLabel     = "apiVersion:"
	kindLabel           = "kind:"
	apiVersionLabelSize = len(apiVersionLabel)
	kindLabelSize       = len(kindLabel)
)

// parseAPIString parses the API string into version and kind components
func parseAPIString(apiString string) (version, kind string) {
	idx := strings.Index(apiString, apiVersionLabel)
	if idx != -1 {
		temps := apiString[idx+apiVersionLabelSize:]
		idx = strings.Index(temps, "\n")
		if idx != -1 {
			temps = temps[:idx]
			version = strings.TrimSpace(temps)
		}
	}
	idx = strings.Index(apiString, kindLabel)
	if idx != -1 {
		temps := apiString[idx+kindLabelSize:]
		idx = strings.Index(temps, "\n")
		if idx != -1 {
			temps = temps[:idx]
			kind = strings.TrimSpace(temps)
		}
	}
	return version, kind
}
