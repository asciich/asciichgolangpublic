package kvmutilsgeneric

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"github.com/asciich/asciichgolangpublic/pkg/vmutils/kvmutils/kvmutilsinterfaces"
)

type Network struct {
	hypervisor kvmutilsinterfaces.Hypervisor

	name       string
	state      string
	autostart  string
	persistent string
}

func NewNetwork() (ret *Network) {
	return new(Network)
}

func (n *Network) SetHypervisor(hypervisor kvmutilsinterfaces.Hypervisor) (err error) {
	if hypervisor == nil {
		return tracederrors.TracedError("hypervisor is nil")
	}

	n.hypervisor = hypervisor

	return nil
}

func (n *Network) GetHypervisor() (hypervisor kvmutilsinterfaces.Hypervisor, err error) {
	if n.hypervisor == nil {
		return nil, tracederrors.TracedError("hypervisor not set")
	}

	return n.hypervisor, nil
}

func (n *Network) SetName(name string) (err error) {
	if name == "" {
		return tracederrors.TracedErrorEmptyString("name")
	}

	n.name = name

	return nil
}

func (n *Network) GetName() (name string, err error) {
	if n.name == "" {
		return "", tracederrors.TracedError("name not set")
	}

	return n.name, nil
}

func (n *Network) SetState(state string) (err error) {
	if state == "" {
		return tracederrors.TracedErrorEmptyString("state")
	}

	n.state = state

	return nil
}

func (n *Network) GetState() (state string, err error) {
	if n.state == "" {
		return "", tracederrors.TracedError("state not set")
	}

	return n.state, nil
}

func (n *Network) IsActive() (isActive bool, err error) {
	state, err := n.GetState()
	if err != nil {
		return false, err
	}

	return state == "active", nil
}

func (n *Network) SetAutostart(autostart string) (err error) {
	if autostart == "" {
		return tracederrors.TracedErrorEmptyString("autostart")
	}

	n.autostart = autostart

	return nil
}

func (n *Network) GetAutostart() (autostart string, err error) {
	if n.autostart == "" {
		return "", tracederrors.TracedError("autostart not set")
	}

	return n.autostart, nil
}

func (n *Network) SetPersistent(persistent string) (err error) {
	if persistent == "" {
		return tracederrors.TracedErrorEmptyString("persistent")
	}

	n.persistent = persistent

	return nil
}

func (n *Network) GetPersistent() (persistent string, err error) {
	if n.persistent == "" {
		return "", tracederrors.TracedError("persistent not set")
	}

	return n.persistent, nil
}

// getParsedXml runs 'virsh net-dumpxml <name>' and parses the relevant parts.
func (n *Network) getParsedXml(ctx context.Context) (parsed kvmutilsinterfaces.KvmNetworkXml, err error) {
	name, err := n.GetName()
	if err != nil {
		return nil, err
	}

	hypervisor, err := n.GetHypervisor()
	if err != nil {
		return nil, err
	}

	return hypervisor.GetParsedNetworkXml(ctx, name)
}

// GetForwardMode returns the libvirt forward mode of the network
// (e.g. "nat", "route", "bridge", "open"). An empty string means the network
// is isolated (no <forward> element).
func (n *Network) GetForwardMode(ctx context.Context) (forwardMode string, err error) {
	parsed, err := n.getParsedXml(ctx)
	if err != nil {
		return "", err
	}

	return parsed.GetForwardMode()
}

// IsNatted returns true if the network uses NAT forwarding.
func (n *Network) IsNatted(ctx context.Context) (isNatted bool, err error) {
	forwardMode, err := n.GetForwardMode(ctx)
	if err != nil {
		return false, err
	}

	return forwardMode == "nat", nil
}

// GetGatewayIp returns the host/gateway IP address of the network
// (e.g. "192.168.122.1").
func (n *Network) GetGatewayIp(ctx context.Context) (gatewayIp string, err error) {
	parsed, err := n.getParsedXml(ctx)
	if err != nil {
		return "", err
	}

	ip, err := parsed.GetIpAddress()
	if err != nil {
		return "", err
	}

	if ip == "" {
		name, _ := n.GetName()
		return "", tracederrors.TracedErrorf("Network '%s' has no IP address configured.", name)
	}

	return ip, nil
}

// GetNetmask returns the netmask of the network (e.g. "255.255.255.0").
func (n *Network) GetNetmask(ctx context.Context) (netmask string, err error) {
	parsed, err := n.getParsedXml(ctx)
	if err != nil {
		return "", err
	}

	netmask, err = parsed.GetIpNetmask()
	if err != nil {
		return "", err
	}

	if netmask == "" {
		name, _ := n.GetName()
		return "", tracederrors.TracedErrorf("Network '%s' has no netmask configured.", name)
	}

	return netmask, nil
}

// GetDhcpRange returns the DHCP range (start and end IP) of the network.
func (n *Network) GetDhcpRange(ctx context.Context) (rangeStart string, rangeEnd string, err error) {
	parsed, err := n.getParsedXml(ctx)
	if err != nil {
		return "", "", err
	}

	rangeStart, err = parsed.GetIpDhcpRangeStart()
	if err != nil {
		return "", "", err
	}

	rangeEnd, err = parsed.GetIpDhcpRangeEnd()
	if err != nil {
		return "", "", err
	}

	if rangeStart == "" || rangeEnd == "" {
		name, _ := n.GetName()
		return "", "", tracederrors.TracedErrorf("Network '%s' has no DHCP range configured.", name)
	}

	return rangeStart, rangeEnd, nil
}

// SetDhcpStartIp updates the start IP of the network's DHCP range while keeping
// the current end IP unchanged.
//
// libvirt does not support modifying a DHCP range in place
// ("dhcp ranges cannot be modified, only added or deleted"), so this deletes
// the current range and adds the new one. The change is applied live and
// persisted (equivalent to two 'virsh net-update <name> delete/add
// ip-dhcp-range ... --live --config' calls).
func (n *Network) SetDhcpStartIp(ctx context.Context, startIp string) (err error) {
	if startIp == "" {
		return tracederrors.TracedErrorEmptyString("startIp")
	}

	name, err := n.GetName()
	if err != nil {
		return err
	}

	hypervisor, err := n.GetHypervisor()
	if err != nil {
		return err
	}

	hostname, err := hypervisor.GetHostName()
	if err != nil {
		return err
	}

	currentStart, currentEnd, err := n.GetDhcpRange(ctx)
	if err != nil {
		return err
	}

	if currentStart == startIp {
		logging.LogInfoByCtxf(
			ctx,
			"DHCP range start of network '%s' on kvm host '%s' is already '%s'.",
			name, hostname, startIp,
		)
		return nil
	}

	// libvirt only allows adding or deleting DHCP ranges, so delete the current
	// range first and then add the new one.
	err = hypervisor.DeleteIpDhcpRangeInNetwork(ctx, name, currentStart, currentEnd)
	if err != nil {
		return err
	}

	err = hypervisor.AddIpDhcpRangeInNetwork(ctx, name, currentStart, currentEnd)
	if err != nil {
		return err
	}

	logging.LogChangedByCtxf(
		ctx,
		"DHCP range start of network '%s' on kvm host '%s' changed from '%s' to '%s' (end stays '%s').",
		name, hostname, currentStart, startIp, currentEnd,
	)

	return nil
}
