package parameteroptions

import (
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type ChmodOptions struct {
	// Set permissions using string like "u=rwx,g=r,o="
	PermissionsString string

	// Use sudo to perform changemod with root priviledges.
	UseSudo bool

	// Enable verbose output
	Verbose bool
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

func (c *ChmodOptions) GetVerbose() (verbose bool) {

	return c.Verbose
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

func (c *ChmodOptions) SetVerbose(verbose bool) {
	c.Verbose = verbose
}
