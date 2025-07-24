package parameteroptions

import (
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type ChownOptions struct {
	UserName  string
	GroupName string
	UseSudo   bool
	Verbose   bool
}

func NewChownOptions() (c *ChownOptions) {
	return new(ChownOptions)
}

func (c *ChownOptions) GetGroupName() (GroupName string, err error) {
	if c.GroupName == "" {
		return "", tracederrors.TracedErrorf("GroupName not set")
	}

	return c.GroupName, nil
}

func (c *ChownOptions) GetUseSudo() (useSudo bool) {

	return c.UseSudo
}

func (c *ChownOptions) GetUserAndOptionallyGroupForChownCommand() (userAndGroup string, err error) {
	userName, err := c.GetUserName()
	if err != nil {
		return "", err
	}

	if c.GroupName == "" {
		return userName, nil
	}

	groupName, err := c.GetGroupName()
	if err != nil {
		return "", err
	}

	return userName + ":" + groupName, err
}

func (c *ChownOptions) GetUserName() (userName string, err error) {
	if c.UserName == "" {
		return "", tracederrors.TracedErrorf("UserName not set")
	}

	return c.UserName, nil
}

func (c *ChownOptions) GetVerbose() (verbose bool) {

	return c.Verbose
}

func (c *ChownOptions) IsGroupNameSet() (isSet bool) {
	return c.GroupName != ""
}

func (c *ChownOptions) SetGroupName(GroupName string) (err error) {
	if GroupName == "" {
		return tracederrors.TracedErrorf("GroupName is empty string")
	}

	c.GroupName = GroupName

	return nil
}

func (c *ChownOptions) SetUseSudo(useSudo bool) {
	c.UseSudo = useSudo
}

func (c *ChownOptions) SetUserName(userName string) (err error) {
	if userName == "" {
		return tracederrors.TracedErrorf("userName is empty string")
	}

	c.UserName = userName

	return nil
}

func (c *ChownOptions) SetVerbose(verbose bool) {
	c.Verbose = verbose
}
