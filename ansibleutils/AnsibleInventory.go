package ansibleutils

import (
	"context"
	"slices"

	"github.com/asciich/asciichgolangpublic/datatypes/slicesutils"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type AnsibleInventory struct {
	name      string
	hostNames []string
}

func NewAnsibleInventoryByName(name string) (inventory *AnsibleInventory) {
	inventory = new(AnsibleInventory)

	inventory.name = name

	return inventory
}

func NewAnsibleInventory() (inventory *AnsibleInventory) {
	return NewAnsibleInventoryByName("in memory ansible inventory")
}

func (a *AnsibleInventory) Name() (name string) {
	return a.name
}

func (a *AnsibleInventory) GetHostNames() (hostNames []string) {
	return slicesutils.GetSortedDeepCopyOfStringsSlice(a.hostNames)
}

func (a *AnsibleInventory) MustAddHostByName(ctx context.Context, hostName string) {
	err := a.AddHostByName(ctx, hostName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (a *AnsibleInventory) HostByNameExists(ctx context.Context, hostName string) (exists bool, err error) {
	if hostName == "" {
		return false, tracederrors.TracedErrorEmptyString("hostName")
	}

	exists = slices.Contains(a.hostNames, hostName)

	if exists {
		logging.LogInfoByCtxf(ctx, "Host '%s' exists in ansible inventory '%s'.", hostName, a.Name())
	} else {
		logging.LogInfoByCtxf(ctx, "Host '%s' does not exist in ansible inventory '%s'.", hostName, a.Name())
	}

	return exists, nil
}

func (a *AnsibleInventory) AddHostByName(ctx context.Context, hostName string) (err error) {
	if hostName == "" {
		return tracederrors.TracedErrorEmptyString("hostName")
	}

	exists, err := a.HostByNameExists(ctx, hostName)
	if err != nil {
		return err
	}

	if exists {
		logging.LogInfoByCtxf(ctx, "Host '%s' already exists in inventory '%s'.", hostName, a.Name())
	} else {
		a.hostNames = append(a.hostNames, hostName)
		logging.LogChangedByCtxf(ctx, "Host '%s' added to ansible inventory '%s'.", hostName, a.Name())
	}

	return nil
}

func (a *AnsibleInventory) GetNumberOfHosts(ctx context.Context) (nHosts int) {
	nHosts = len(a.hostNames)

	logging.LogInfoByCtxf(ctx, "There are '%d' hosts in ansible inventory '%s'.", nHosts, a.Name())

	return nHosts
}
