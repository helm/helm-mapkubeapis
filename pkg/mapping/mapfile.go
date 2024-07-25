/*
Copyright

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package mapping

import (
	"os"

	"sigs.k8s.io/yaml"
)

// LoadMapfile loads a Map.yaml file into a *Metadata.
func LoadMapfile(filename string) (*Metadata, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	y := new(Metadata)
	err = yaml.Unmarshal(b, y)
	return y, err
}
