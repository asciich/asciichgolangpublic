package kubernetesutils

import (
	"github.com/asciich/asciichgolangpublic/datatypes/mapsutils"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type CreateConfigMapOptions struct {
	ConfigMapData map[string]string
}

func (c *CreateConfigMapOptions) GetConfigMapData() (map[string]string, error) {
	if c.ConfigMapData == nil {
		return nil, tracederrors.TracedError("ConfigMapData not set")
	}

	return mapsutils.DeepCopyStringsMap(c.ConfigMapData), nil
}
