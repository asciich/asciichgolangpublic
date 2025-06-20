package kubernetesparameteroptions

import "github.com/asciich/asciichgolangpublic/tracederrors"

type CreateObjectOptions struct {
	YamlString string

	// Do not check nor try to create missing namespaces.
	// Usefull if the user has no rights to list or create namespaces but requiers to ensure namespaces already exists.
	SkipNamespaceCreation bool
}

func (c CreateObjectOptions) GetYamlString() (string, error) {
	if c.YamlString == "" {
		return "", tracederrors.TracedError("YamlString not set")
	}

	return c.YamlString, nil
}
