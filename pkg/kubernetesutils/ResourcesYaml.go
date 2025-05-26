package kubernetesutils

import (
	"sort"
	"strings"

	"github.com/asciich/asciichgolangpublic/datatypes/stringsutils"
	"github.com/asciich/asciichgolangpublic/fileformats/yamlutils"
	"github.com/asciich/asciichgolangpublic/tracederrors"
	"gopkg.in/yaml.v3"
)

type ResourceYamlEntry struct {
	Content string
}

func (r *ResourceYamlEntry) Name() (name string) {
	type ToParse struct {
		Metadata struct {
			Name string `yaml:"name"`
		} `yaml:"metadata"`
	}

	toParse := new(ToParse)

	err := yaml.Unmarshal([]byte(r.Content), toParse)
	if err != nil {
		return ""
	}

	return toParse.Metadata.Name
}

func (r *ResourceYamlEntry) Validate() (err error) {
	if r.Content == "" {
		return tracederrors.TracedError("Content not set")
	}

	err = yamlutils.Validate(r.Content)
	if err != nil {
		return err
	}

	if r.Name() == "" {
		return tracederrors.TracedError("Kubernetes resource YAML without a name are not valid")
	}

	if r.Kind() == "" {
		return tracederrors.TracedError("Kubernetes resource YAML without a kind are not valid")
	}

	return nil
}

func (r *ResourceYamlEntry) Kind() (name string) {
	type ToParse struct {
		Kind string `yaml:"kind"`
	}

	toParse := new(ToParse)

	err := yaml.Unmarshal([]byte(r.Content), toParse)
	if err != nil {
		return ""
	}

	return toParse.Kind
}

func (r *ResourceYamlEntry) Namespace() (namespace string) {
	type ToParse struct {
		Metadata struct {
			Namespace string `yaml:"namespace"`
		} `yaml:"metadata"`
	}

	toParse := new(ToParse)

	err := yaml.Unmarshal([]byte(r.Content), toParse)
	if err != nil {
		return ""
	}

	return toParse.Metadata.Namespace
}

func UnmarshalResourceYaml(resourceYaml string) (resources []*ResourceYamlEntry, err error) {
	splitted := yamlutils.SplitMultiYaml(resourceYaml)

	resources = []*ResourceYamlEntry{}

	for _, s := range splitted {
		toAdd := &ResourceYamlEntry{Content: s}

		err = toAdd.Validate()
		if err != nil {
			return nil, err
		}

		resources = append(resources, toAdd)
	}

	return resources, nil
}

func marshalResourceYaml(resources []*ResourceYamlEntry) (marshalled string, err error) {
	if resources == nil {
		return "", tracederrors.TracedErrorNil("resources")
	}

	splitted := []string{}

	for _, r := range resources {
		toAdd := r.Content

		toAdd = strings.TrimSpace(toAdd)

		if toAdd == "" {
			continue
		}

		splitted = append(splitted, toAdd)
	}

	marshalled, err = yamlutils.MergeMultiYaml(splitted)
	if err != nil {
		return "", err
	}

	return stringsutils.EnsureEndsWithExactlyOneLineBreak(marshalled), err
}

// Sort resources in given multi yaml string.
// Resources are sorted by:
//  1. Namespace
//  2. Resource name
//  3. Kind
func SortResourcesYaml(resourcesYaml string) (sortedResourcesYaml string, err error) {
	parsed, err := UnmarshalResourceYaml(resourcesYaml)
	if err != nil {
		return "", err
	}

	sort.Slice(parsed, func(i, j int) bool {
		namespace_i := parsed[i].Namespace()
		namespace_j := parsed[j].Namespace()

		if stringsutils.IsBeforeInAlphabeth(namespace_i, namespace_j) {
			return true
		}

		name_i := parsed[i].Name()
		name_j := parsed[j].Name()
		if stringsutils.IsBeforeInAlphabeth(name_i, name_j) {
			return true
		}

		kind_i := parsed[i].Kind()
		kind_j := parsed[j].Kind()
		return stringsutils.IsBeforeInAlphabeth(kind_i, kind_j)
	})

	sortedResourcesYaml, err = marshalResourceYaml(parsed)
	if err != nil {
		return "", err
	}

	return sortedResourcesYaml, nil
}
