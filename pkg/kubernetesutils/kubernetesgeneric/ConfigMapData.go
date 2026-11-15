package kubernetesgeneric

import (
	"github.com/asciich/asciichgolangpublic/pkg/fileformats/yamlutils"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// Converts the data of a kubernetes ConfigMap to a single yaml document.
// Every entry of the ConfigMap data becomes a top level yaml field named
// like the entry key. The entry value is parsed as yaml and embedded:
//
//	ConfigMap data:
//	    my-configmap: |
//	        hello: world
//
//	Resulting yaml string:
//	    my-configmap:
//	        hello: world
//
// This way values stored in the ConfigMap can be read by a yaml path
// like `.my-configmap.hello` afterwards. Data entries which are no valid
// yaml are embedded as plain string.
func ConfigMapDataToYamlString(configMapName string, data map[string]string) (yamlString string, err error) {
	if configMapName == "" {
		return "", tracederrors.TracedErrorEmptyString("configMapName")
	}

	if len(data) <= 0 {
		return "", tracederrors.TracedErrorf("ConfigMap '%s' has no data to convert to yaml.", configMapName)
	}

	entries := map[string]interface{}{}

	for key, value := range data {
		valueAsYaml, err := yamlutils.LoadGeneric(value)
		if err != nil || valueAsYaml == nil {
			// Value is no valid yaml. Keep the raw string value:
			entries[key] = value
			continue
		}

		entries[key] = valueAsYaml
	}

	yamlString, err = yamlutils.DataToYamlString(entries)
	if err != nil {
		return "", tracederrors.TracedErrorf("Failed to convert data of ConfigMap '%s' to yaml: %w", configMapName, err)
	}

	return yamlString, nil
}
