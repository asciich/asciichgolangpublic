package kubernetesparameteroptions

import (
	"github.com/asciich/asciichgolangpublic/pkg/datatypes/mapsutils"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type CreateSecretOptions struct {
	SecretData map[string][]byte
}

func (c *CreateSecretOptions) GetSecretData() (map[string][]byte, error) {
	if c.SecretData == nil {
		return nil, tracederrors.TracedError("SecretData not set")
	}

	return mapsutils.DeepCopyBytesMap(c.SecretData), nil
}
