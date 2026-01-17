package linuxuseroptions

import "github.com/asciich/asciichgolangpublic/pkg/tracederrors"

type DeleteOptions struct {
	UserName string

	// If true a force delete is performed.
	Force bool

	UseSudo bool
}

func (d *DeleteOptions) GetUSerName() (string, error) {
	if d.UserName == "" {
		return "", tracederrors.TracedError("UserName not set")
	}

	return d.UserName, nil
}
