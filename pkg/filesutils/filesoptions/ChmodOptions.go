package filesoptions

import (
	"github.com/asciich/asciichgolangpublic/pkg/osutils/unixfilepermissionsutils"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type ChmodOptions struct {
	// Set permissions using string like "u=rwx,g=r,o="
	PermissionsString string

	// Use sudo to perform changemod with root priviledges.
	UseSudo bool
}

func NewChmodOptions() (c *ChmodOptions) {
	return new(ChmodOptions)
}

func (c *ChmodOptions) GetPermissionsString() (permissionsString string, err error) {
	if c.PermissionsString == "" {
		return "", tracederrors.TracedErrorf("PermissionsString not set")
	}

	return c.PermissionsString, nil
}

func (c *ChmodOptions) GetUseSudo() (useSudo bool) {

	return c.UseSudo
}

func (c *ChmodOptions) SetPermissionsString(permissionsString string) (err error) {
	if permissionsString == "" {
		return tracederrors.TracedErrorf("permissionsString is empty string")
	}

	c.PermissionsString = permissionsString

	return nil
}

func (c *ChmodOptions) SetUseSudo(useSudo bool) {
	c.UseSudo = useSudo
}

func (c *ChmodOptions) GetPermissions() (int, error) {
	permissionsString, err := c.GetPermissionsString()
	if err != nil {
		return 0, err
	}

	return unixfilepermissionsutils.GetPermissionsValue(permissionsString)
}