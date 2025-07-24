package kubernetesparameteroptions

import (
	"github.com/asciich/asciichgolangpublic/datatypes/mapsutils"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type CreateConfigMapOptions struct {
	ConfigMapData map[string]string
	Labels        map[string]string
}

func (c *CreateConfigMapOptions) GetLabels() map[string]string {
	if len(c.Labels) <= 0 {
		return map[string]string{}
	}

	return c.Labels
}

func (c *CreateConfigMapOptions) GetConfigMapData() (map[string]string, error) {
	if c.ConfigMapData == nil {
		return nil, tracederrors.TracedError("ConfigMapData not set")
	}

	return mapsutils.DeepCopyStringsMap(c.ConfigMapData), nil
}
