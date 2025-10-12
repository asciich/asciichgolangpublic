package ansibleplaybook

import (
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"gopkg.in/yaml.v3"
)

type Play struct {
	Name  string   `yaml:"name"`
	Hosts Hosts    `yaml:"hosts"`
	Roles []string `yaml:"roles,omitempty"`
}

// Custom type alias so we can implement the UnmarshalYaml
type Hosts []string

// Custom UnmarshalYaml implementation for both strings and list of strings as hosts:
func (h *Hosts) UnmarshalYAML(value *yaml.Node) error {
	switch value.Kind {
	case yaml.ScalarNode:
		var s string
		if err := value.Decode(&s); err != nil {
			return tracederrors.TracedErrorf("failed to decode scalar hosts: %w", err)
		}
		*h = []string{s}
		return nil

	case yaml.SequenceNode:
		var ss []string
		if err := value.Decode(&ss); err != nil {
			return tracederrors.TracedErrorf("failed to decode sequence hosts: %w", err)
		}
		*h = ss
		return nil
	case yaml.MappingNode:
		return tracederrors.TracedErrorf("hosts field cannot be a map/object")
	default:
		if value.Tag == "!!null" {
			*h = Hosts{}
			return nil
		}
		return tracederrors.TracedErrorf("unsupported YAML node kind for hosts: %v", value.Kind)
	}
}
