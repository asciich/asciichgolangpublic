package ansibleutils

import (
	"context"
	"fmt"

	"github.com/asciich/asciichgolangpublic/logging"
)

type AnsibleCliOuput struct {
	name      string
	inventory *AnsibleInventory
}

func NewAnsibleCliOutput() (a *AnsibleCliOuput) {
	a = new(AnsibleCliOuput)

	a.name = "in memory ansible cli output"

	return a
}

func (a *AnsibleCliOuput) Name() (name string) {
	return a.name
}

func (a *AnsibleCliOuput) Inventory() (inventory *AnsibleInventory) {
	return a.inventory
}

func (a *AnsibleCliOuput) MustGetNumberOfHosts(ctx context.Context) (nHosts int) {
	nHosts, err := a.GetNumberOfHosts(ctx)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return nHosts
}

func (a *AnsibleCliOuput) GetNumberOfHosts(ctx context.Context) (nHosts int, err error) {
	inventory := a.Inventory()
	if inventory == nil {
		nHosts = 0
		logging.LogInfoByCtxf(ctx, "There are '%d' hosts in ansible cli output '%s'.", nHosts, a.Name())
	} else {
		nHosts, err = inventory.GetNumberOfHosts(ctx)
		if err != nil {
			return 0, err
		}
	}

	return nHosts, err
}

func (a *AnsibleCliOuput) CreateInventory() (inventory *AnsibleInventory) {
	if a.inventory != nil {
		return a.inventory
	}

	a.inventory = NewAnsibleInventoryByName(
		fmt.Sprintf("in memory ansible inventory for '%s'", a.Name()),
	)

	return a.inventory
}
