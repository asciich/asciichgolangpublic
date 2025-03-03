package ansibleutils

import (
	"context"
	"encoding/json"
	"slices"

	"github.com/asciich/asciichgolangpublic/datatypes/slicesutils"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type AnsibleInventory struct {
	name   string
	groups []*AnsibleGroup
	hosts  []*AnsibleHost
}

func NewAnsibleInventoryByName(name string) (inventory *AnsibleInventory) {
	inventory = new(AnsibleInventory)

	inventory.name = name

	return inventory
}

func NewAnsibleInventory() (inventory *AnsibleInventory) {
	return NewAnsibleInventoryByName("in memory ansible inventory")
}

func ParseInventoryJson(ctx context.Context, jsonData string) (inventory *AnsibleInventory, err error) {
	if jsonData == "" {
		return nil, tracederrors.TracedErrorEmptyString("jsonData")
	}

	logging.LogInfoByCtx(ctx, "Parse ansible inventory Json started.")

	type AnsibleInventoryJsonEntry struct {
		Hosts []string `json:"hosts"`
	}

	parsed := map[string]AnsibleInventoryJsonEntry{}

	err = json.Unmarshal([]byte(jsonData), &parsed)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Unable to parse jsonData as ansible inventory: %w", err)
	}

	inventory = NewAnsibleInventory()

	for k, v := range parsed {
		if len(v.Hosts) <= 0 {
			continue
		}

		_, err := inventory.CreateGroupByName(ctx, k)
		if err != nil {
			return nil, err
		}

		for _, toAdd := range v.Hosts {
			_, err = inventory.CreateHostByName(ctx, toAdd)
			if err != nil {
				return nil, err
			}
		}
	}

	logging.LogInfoByCtx(ctx, "Parse ansible inventory Json finished.")

	return inventory, nil
}

func MustParseInventoryJson(ctx context.Context, jsonData string) (inventory *AnsibleInventory) {
	inventory, err := ParseInventoryJson(ctx, jsonData)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return inventory
}

func (a *AnsibleInventory) GetHostByName(hostName string) (ansibleHost *AnsibleHost, err error) {
	if hostName == "" {
		return nil, tracederrors.TracedErrorEmptyString("hostName")
	}

	for _, h := range a.hosts {
		name, err := h.GetHostName()
		if err != nil {
			return nil, err
		}

		if name == hostName {
			return h, nil
		}
	}

	return nil, tracederrors.TracedErrorf("Ansible host '%s' not found in inventory '%s'", hostName, a.Name())
}

func (a *AnsibleInventory) MustCreateHostByName(ctx context.Context, hostName string) (addedHost *AnsibleHost) {
	addedHost, err := a.CreateHostByName(ctx, hostName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return addedHost
}

func (a *AnsibleInventory) CreateHostByName(ctx context.Context, hostName string) (addedHost *AnsibleHost, err error) {
	if hostName == "" {
		return nil, tracederrors.TracedErrorEmptyString("hostName")
	}

	exists, err := a.HostByNameExists(ctx, hostName)
	if err != nil {
		return nil, err
	}

	if exists {
		logging.LogInfoByCtxf(ctx, "Host '%s' already present in inventory '%s'. Skip create.", hostName, a.Name())
		return a.GetHostByName(hostName)
	} else {
		toAdd, err := NewAnsibleHostByName(hostName)
		if err != nil {
			return nil, err
		}

		a.hosts = append(a.hosts, toAdd)

		return toAdd, nil
	}
}

func (a *AnsibleInventory) Name() (name string) {
	return a.name
}

func (a *AnsibleInventory) HostByNameExists(ctx context.Context, hostName string) (exists bool, err error) {
	if hostName == "" {
		return false, tracederrors.TracedErrorEmptyString("hostName")
	}

	hostNames, err := a.ListHostNames()
	if err != nil {
		return false, err
	}

	exists = slices.Contains(hostNames, hostName)

	if exists {
		logging.LogInfoByCtxf(ctx, "Host '%s' exists in ansible inventory '%s'.", hostName, a.Name())
	} else {
		logging.LogInfoByCtxf(ctx, "Host '%s' does not exist in ansible inventory '%s'.", hostName, a.Name())
	}

	return exists, nil
}

func (a *AnsibleInventory) ensureAllGroupPresent() (allGroup *AnsibleGroup, err error) {
	for _, g := range a.groups {
		name, err := g.GetGroupName()
		if err != nil {
			return nil, err
		}

		if name == "" {
			return g, nil
		}
	}

	toAdd, err := NewAnsibleGroupByName("all")
	if err != nil {
		return nil, err
	}

	a.groups = append(a.groups, toAdd)

	return toAdd, nil
}

func (a *AnsibleInventory) MustGroupByNameExists(ctx context.Context, groupName string) (groupExists bool) {
	groupExists, err := a.GroupByNameExists(ctx, groupName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return groupExists
}

func (a *AnsibleInventory) GroupByNameExists(ctx context.Context, groupName string) (groupExists bool, err error) {
	for _, group := range a.groups {
		n, err := group.GetGroupName()
		if err != nil {
			return false, err
		}

		if n == groupName {
			groupExists = true
			break
		}
	}

	if groupExists {
		logging.LogInfoByCtxf(ctx, "Ansible group '%s' exists in inventory '%s'", groupName, a.Name())
	} else {
		logging.LogInfoByCtxf(ctx, "Ansible group '%s' does not exist in inventory '%s'", groupName, a.Name())
	}

	return groupExists, nil
}

func (a *AnsibleInventory) MustGetNumberOfHosts(ctx context.Context) (numberOfHosts int) {
	numberOfHosts, err := a.GetNumberOfHosts(ctx)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return numberOfHosts
}

func (a *AnsibleInventory) MustListHostNames() (hostNames []string) {
	hostNames, err := a.ListHostNames()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return hostNames
}

func (a *AnsibleInventory) ListHostNames() (hostNames []string, err error) {
	hostNames = []string{}

	for _, h := range a.hosts {
		toAdd, err := h.GetHostName()
		if err != nil {
			return nil, err
		}

		hostNames = append(hostNames, toAdd)
	}

	hostNames = slicesutils.SortStringSliceAndRemoveDuplicates(hostNames)

	return hostNames, nil
}

func (a *AnsibleInventory) GetNumberOfHosts(ctx context.Context) (numberOfHosts int, err error) {
	list, err := a.ListHostNames()
	if err != nil {
		return 0, err
	}

	numberOfHosts = len(list)

	logging.LogInfoByCtxf(ctx, "There are '%d' hosts in the ansible inventory '%s'.", numberOfHosts, a.Name())

	return numberOfHosts, nil
}

func (a *AnsibleInventory) ListGroupNames() (groupNames []string, err error) {
	groupNames = []string{"all"}

	for _, g := range a.groups {
		toAdd, err := g.GetGroupName()
		if err != nil {
			return nil, err
		}

		groupNames = append(groupNames, toAdd)
	}

	groupNames = slicesutils.SortStringSliceAndRemoveDuplicates(groupNames)

	return groupNames, nil
}

func (a *AnsibleInventory) MustListGroupNames() (groupNames []string) {
	groupNames, err := a.ListGroupNames()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return groupNames
}

func (a *AnsibleInventory) MustCreateGroupByName(ctx context.Context, groupName string) (createdGroup *AnsibleGroup) {
	createdGroup, err := a.CreateGroupByName(ctx, groupName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return createdGroup
}

func (a *AnsibleInventory) GetGroupByName(groupName string) (group *AnsibleGroup, err error) {
	if groupName == "" {
		return nil, tracederrors.TracedErrorEmptyString("groupName")
	}

	for _, g := range a.groups {
		name, err := g.GetGroupName()
		if err != nil {
			return nil, err
		}

		if name == groupName {
			return g, nil
		}
	}

	return nil, tracederrors.TracedErrorf("No group '%s' found in ansible inventory '%s'.", groupName, a.Name())
}

func (a *AnsibleInventory) CreateGroupByName(ctx context.Context, groupName string) (createdGroup *AnsibleGroup, err error) {
	if groupName == "" {
		return nil, tracederrors.TracedErrorEmptyString("groupName")
	}

	exists, err := a.GroupByNameExists(ctx, groupName)
	if err != nil {
		return nil, err
	}

	if exists {
		logging.LogInfoByCtxf(ctx, "Ansible group '%s' already exists in inventory '%s'. Skip create.", groupName, a.Name())

		return a.GetGroupByName(groupName)
	} else {
		toAdd, err := NewAnsibleGroupByName(groupName)
		if err != nil {
			return nil, err
		}

		a.groups = append(a.groups, toAdd)

		logging.LogChangedByCtxf(ctx, "Ansible group '%s' created in inventory '%s'.", groupName, a.Name())

		return toAdd, nil
	}
}
