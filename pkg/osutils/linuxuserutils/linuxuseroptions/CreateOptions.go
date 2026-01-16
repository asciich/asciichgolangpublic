package linuxuseroptions

import "github.com/asciich/asciichgolangpublic/pkg/tracederrors"

type CreateOptions struct {
	UserName string

	UseSudo bool

	CreateHomeDirectory bool
}

func (c *CreateOptions) GetUserName() (string, error) {
	if c.UserName == "" {
		return "", tracederrors.TracedError("UserName not set")
	}

	return c.UserName, nil
}
