package kvmutils

import (
	"context"
	"encoding/xml"
	"fmt"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type KvmNetwork struct {
	hypervisor *KVMHypervisor

	name       string
	state      string
	autostart  string
	persistent string
}

// Parsed representation of the relevant parts of 'virsh net-dumpxml <name>'.
type kvmNetworkXml struct {
	XMLName xml.Name `xml:"network"`
	Name    string   `xml:"name"`
	Forward struct {
		Mode string `xml:"mode,attr"`
	} `xml:"forward"`
	Ip struct {
		Address string `xml:"address,attr"`
		Netmask string `xml:"netmask,attr"`
		Dhcp    struct {
			Range struct {
				Start string `xml:"start,attr"`
				End   string `xml:"end,attr"`
			} `xml:"range"`
		} `xml:"dhcp"`
	} `xml:"ip"`
}

func NewKvmNetwork() (ret *KvmNetwork) {
	return new(KvmNetwork)
}

func (n *KvmNetwork) SetHypervisor(hypervisor *KVMHypervisor) (err error) {
	if hypervisor == nil {
		return tracederrors.TracedError("hypervisor is nil")
	}

	n.hypervisor = hypervisor

	return nil
}

func (n *KvmNetwork) GetHypervisor() (hypervisor *KVMHypervisor, err error) {
	if n.hypervisor == nil {
		return nil, tracederrors.TracedError("hypervisor not set")
	}

	return n.hypervisor, nil
}

func (n *KvmNetwork) SetName(name string) (err error) {
	if name == "" {
		return tracederrors.TracedErrorEmptyString("name")
	}

	n.name = name

	return nil
}

func (n *KvmNetwork) GetName() (name string, err error) {
	if n.name == "" {
		return "", tracederrors.TracedError("name not set")
	}

	return n.name, nil
}

func (n *KvmNetwork) SetState(state string) (err error) {
	if state == "" {
		return tracederrors.TracedErrorEmptyString("state")
	}

	n.state = state

	return nil
}

func (n *KvmNetwork) GetState() (state string, err error) {
	if n.state == "" {
		return "", tracederrors.TracedError("state not set")
	}

	return n.state, nil
}

func (n *KvmNetwork) IsActive() (isActive bool, err error) {
	state, err := n.GetState()
	if err != nil {
		return false, err
	}

	return state == "active", nil
}

func (n *KvmNetwork) SetAutostart(autostart string) (err error) {
	if autostart == "" {
		return tracederrors.TracedErrorEmptyString("autostart")
	}

	n.autostart = autostart

	return nil
}

func (n *KvmNetwork) GetAutostart() (autostart string, err error) {
	if n.autostart == "" {
		return "", tracederrors.TracedError("autostart not set")
	}

	return n.autostart, nil
}

func (n *KvmNetwork) SetPersistent(persistent string) (err error) {
	if persistent == "" {
		return tracederrors.TracedErrorEmptyString("persistent")
	}

	n.persistent = persistent

	return nil
}

func (n *KvmNetwork) GetPersistent() (persistent string, err error) {
	if n.persistent == "" {
		return "", tracederrors.TracedError("persistent not set")
	}

	return n.persistent, nil
}

// getParsedXml runs 'virsh net-dumpxml <name>' and parses the relevant parts.
func (n *KvmNetwork) getParsedXml(ctx context.Context) (parsed *kvmNetworkXml, err error) {
	name, err := n.GetName()
	if err != nil {
		return nil, err
	}

	hypervisor, err := n.GetHypervisor()
	if err != nil {
		return nil, err
	}

	stdout, err := hypervisor.RunKvmCommandAndGetStdout(ctx, []string{"net-dumpxml", name})
	if err != nil {
		return nil, err
	}

	parsed = &kvmNetworkXml{}
	err = xml.Unmarshal([]byte(stdout), parsed)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to parse net-dumpxml output for network '%s': %w", name, err)
	}

	return parsed, nil
}

// GetForwardMode returns the libvirt forward mode of the network
// (e.g. "nat", "route", "bridge", "open"). An empty string means the network
// is isolated (no <forward> element).
func (n *KvmNetwork) GetForwardMode(ctx context.Context) (forwardMode string, err error) {
	parsed, err := n.getParsedXml(ctx)
	if err != nil {
		return "", err
	}

	return parsed.Forward.Mode, nil
}

// IsNatted returns true if the network uses NAT forwarding.
func (n *KvmNetwork) IsNatted(ctx context.Context) (isNatted bool, err error) {
	forwardMode, err := n.GetForwardMode(ctx)
	if err != nil {
		return false, err
	}

	return forwardMode == "nat", nil
}

// GetGatewayIp returns the host/gateway IP address of the network
// (e.g. "192.168.122.1").
func (n *KvmNetwork) GetGatewayIp(ctx context.Context) (gatewayIp string, err error) {
	parsed, err := n.getParsedXml(ctx)
	if err != nil {
		return "", err
	}

	if parsed.Ip.Address == "" {
		name, _ := n.GetName()
		return "", tracederrors.TracedErrorf("Network '%s' has no IP address configured.", name)
	}

	return parsed.Ip.Address, nil
}

// GetNetmask returns the netmask of the network (e.g. "255.255.255.0").
func (n *KvmNetwork) GetNetmask(ctx context.Context) (netmask string, err error) {
	parsed, err := n.getParsedXml(ctx)
	if err != nil {
		return "", err
	}

	if parsed.Ip.Netmask == "" {
		name, _ := n.GetName()
		return "", tracederrors.TracedErrorf("Network '%s' has no netmask configured.", name)
	}

	return parsed.Ip.Netmask, nil
}

// GetDhcpRange returns the DHCP range (start and end IP) of the network.
func (n *KvmNetwork) GetDhcpRange(ctx context.Context) (rangeStart string, rangeEnd string, err error) {
	parsed, err := n.getParsedXml(ctx)
	if err != nil {
		return "", "", err
	}

	rangeStart = parsed.Ip.Dhcp.Range.Start
	rangeEnd = parsed.Ip.Dhcp.Range.End

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
func (n *KvmNetwork) SetDhcpStartIp(ctx context.Context, startIp string) (err error) {
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

	oldRangeXml := fmt.Sprintf("<range start='%s' end='%s'/>", currentStart, currentEnd)
	newRangeXml := fmt.Sprintf("<range start='%s' end='%s'/>", startIp, currentEnd)

	// libvirt only allows adding or deleting DHCP ranges, so delete the current
	// range first and then add the new one.
	_, err = hypervisor.RunKvmCommandAndGetStdout(
		ctx,
		[]string{"net-update", name, "delete", "ip-dhcp-range", oldRangeXml, "--live", "--config"},
	)
	if err != nil {
		return err
	}

	_, err = hypervisor.RunKvmCommandAndGetStdout(
		ctx,
		[]string{"net-update", name, "add", "ip-dhcp-range", newRangeXml, "--live", "--config"},
	)
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
