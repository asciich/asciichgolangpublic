package kubernetesimplementationindependend

import (
	"sort"
	"strings"

	"github.com/asciich/asciichgolangpublic/datatypes/stringsutils"
	"github.com/asciich/asciichgolangpublic/pkg/fileformats/yamlutils"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"gopkg.in/yaml.v3"
)

type ObjectYamlEntry struct {
	Content string
}

func (r *ObjectYamlEntry) Name() (name string) {
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

func (r *ObjectYamlEntry) Validate() (err error) {
	if r.Content == "" {
		return tracederrors.TracedError("Content not set")
	}

	err = yamlutils.Validate(r.Content)
	if err != nil {
		return err
	}

	if r.Name() == "" {
		return tracederrors.TracedError("Kubernetes object YAML without a name are not valid")
	}

	if r.Kind() == "" {
		return tracederrors.TracedError("Kubernetes object YAML without a kind are not valid")
	}

	return nil
}

func (r *ObjectYamlEntry) Kind() (name string) {
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

func (r *ObjectYamlEntry) ApiVersion() string {
	type ToParse struct {
		ApiVersion string `yaml:"apiVersion"`
	}

	toParse := new(ToParse)

	err := yaml.Unmarshal([]byte(r.Content), toParse)
	if err != nil {
		return ""
	}

	return toParse.ApiVersion
}

func (r *ObjectYamlEntry) Namespace() (namespace string) {
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

func UnmarshalObjectYaml(objectYaml string) (objects []*ObjectYamlEntry, err error) {
	splitted := yamlutils.SplitMultiYaml(objectYaml)

	objects = []*ObjectYamlEntry{}

	for _, s := range splitted {
		toAdd := &ObjectYamlEntry{Content: s}

		err = toAdd.Validate()
		if err != nil {
			return nil, err
		}

		objects = append(objects, toAdd)
	}

	return objects, nil
}

func marshalObjectYaml(objects []*ObjectYamlEntry) (marshalled string, err error) {
	if objects == nil {
		return "", tracederrors.TracedErrorNil("objects")
	}

	splitted := []string{}

	for _, r := range objects {
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

// Sort objects in given multi yaml string.
// Objects are sorted by:
//  1. Namespace
//  2. Object name
//  3. Kind
func SortObjectsYaml(objectsYaml string) (sortedObjectsYaml string, err error) {
	parsed, err := UnmarshalObjectYaml(objectsYaml)
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

	sortedObjectsYaml, err = marshalObjectYaml(parsed)
	if err != nil {
		return "", err
	}

	return sortedObjectsYaml, nil
}
