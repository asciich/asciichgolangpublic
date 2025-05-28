package kubernetesutils

import (
	"github.com/asciich/asciichgolangpublic/datatypes/mapsutils"
	"github.com/asciich/asciichgolangpublic/tracederrors"
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
