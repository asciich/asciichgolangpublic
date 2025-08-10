package virtualenvutils

import "github.com/asciich/asciichgolangpublic/pkg/tracederrors"

type CreateVirtualenvOptions struct {
	// Path of the virtualenv to create:
	Path string

	// Slice of packages to install:
	Packages []string
}

func (c *CreateVirtualenvOptions) GetPath() (string, error) {
	if c.Path == "" {
		return "", tracederrors.TracedError("Path not set")
	}

	return c.Path, nil
}