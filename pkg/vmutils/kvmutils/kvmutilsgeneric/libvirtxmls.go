package kvmutilsgeneric

import (
	"context"
	_ "embed"
	"encoding/xml"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/macaddresses"
	"github.com/asciich/asciichgolangpublic/pkg/templateutils/gotemplateutils"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"github.com/asciich/asciichgolangpublic/pkg/vmutils/kvmutils/kvmutilsoptions"
	libvirtxml "libvirt.org/libvirt-go-xml"
)

//go:embed data/vm.xml.tmpl
var vmXmlTemplate string

func CreateXmlForVmAsString(createOptions *kvmutilsoptions.KvmCreateVmOptions) (libvirtXml string, err error) {
	if createOptions == nil {
		return "", tracederrors.TracedError("createOptions is nil")
	}

	vmName, err := createOptions.GetVmName()
	if err != nil {
		return "", err
	}

	diskPath, err := createOptions.GetDiskImagePath()
	if err != nil {
		return "", err
	}

	// RAM size and CPU count are optional. When unset the defaults
	// (kvmutilsoptions.DefaultVmRamSizeGb / DefaultCpuCount) are used.
	// libvirt domain XML uses KiB as the default memory unit, so the RAM
	// size is passed in KiB.
	ramSizeKiB, err := createOptions.GetVmRamSizeKiBOrDefaultIfUnset()
	if err != nil {
		return "", err
	}

	cpuCount, err := createOptions.GetCpuCountOrDefaultIfUnset()
	if err != nil {
		return "", err
	}

	// The MAC address is optional. When empty libvirt auto-generates one.
	// Only validate the format when a MAC address was explicitly provided.
	macAddress := createOptions.MacAddress
	if macAddress != "" {
		err = macaddresses.CheckStringIsAMacAddress(macAddress)
		if err != nil {
			return "", err
		}
	}

	// The bridged network interface is optional. When empty the template
	// falls back to the libvirt 'default' network. When set the VM is
	// attached to the given host bridge (e.g. br0).
	bridgedNetworkInterface := createOptions.BridgeInterface

	// Additional block devices are optional. They are passed as raw virtio disks
	// (vdb, vdc, ...) so they can be mounted on block device level inside the VM.
	blockDeviceAttachments, err := GetBlockDeviceAttachments(createOptions.AdditionalBlockDevices)
	if err != nil {
		return "", err
	}

	libvirtXml, err = gotemplateutils.RenderTemplateFromStringAsString(
		vmXmlTemplate,
		map[string]any{
			"VM_NAME":                   vmName,
			"DISK_PATH":                 diskPath,
			"MAC_ADDRESS":               macAddress,
			"BRIDGED_NETWORK_INTERFACE": bridgedNetworkInterface,
			"RAM_SIZE_KIB":              ramSizeKiB,
			"CPU_COUNT":                 cpuCount,
			"ADDITIONAL_BLOCK_DEVICES":  blockDeviceAttachments,
		},
	)
	if err != nil {
		return "", err
	}

	return libvirtXml, nil
}

func GetMacAddressFromXmlString(libvirtXml string) (macAddress string, err error) {
	if libvirtXml == "" {
		return "", tracederrors.TracedError("libvirtXml is empty string")
	}

	domcfg := &libvirtxml.Domain{}
	err = domcfg.Unmarshal(libvirtXml)
	if err != nil {
		return "", tracederrors.TracedError(err)
	}

	networkInterfaces := domcfg.Devices.Interfaces
	nInterfaces := len(networkInterfaces)
	if nInterfaces != 1 {
		return "", tracederrors.TracedErrorf(
			"Only exactly one network interface is supported at the moment but got '%d'",
			nInterfaces,
		)
	}

	nativeMac := networkInterfaces[0].MAC
	if nativeMac == nil {
		return "", tracederrors.TracedError("nativeMac is nil after evaluation")
	}

	macAddress = nativeMac.Address
	if macAddress == "" {
		return "", tracederrors.TracedError("macAddress is empty string after evaluation")
	}

	return macAddress, nil
}

func GetNetworkNameFromXmlString(libvirtXml string) (networkName string, err error) {
	if libvirtXml == "" {
		return "", tracederrors.TracedError("libvirtXml is empty string")
	}

	domcfg := &libvirtxml.Domain{}
	err = domcfg.Unmarshal(libvirtXml)
	if err != nil {
		return "", tracederrors.TracedError(err)
	}

	networkInterfaces := domcfg.Devices.Interfaces
	nInterfaces := len(networkInterfaces)
	if nInterfaces != 1 {
		return "", tracederrors.TracedErrorf(
			"Only exactly one network interface is supported at the moment but got '%d'",
			nInterfaces,
		)
	}

	source := networkInterfaces[0].Source
	if source == nil {
		return "", tracederrors.TracedError("interface source is nil after evaluation")
	}

	if source.Network != nil {
		networkName = source.Network.Network
	} else if source.Bridge != nil {
		// Bridged interfaces (e.g. br0) have no network name; report the bridge.
		networkName = source.Bridge.Bridge
	}

	if networkName == "" {
		return "", tracederrors.TracedError("networkName is empty string after evaluation")
	}

	return networkName, nil
}

func GetBridgeInterfaceFromXmlString(libvirtXml string) (bridgeInterface string, err error) {
	if libvirtXml == "" {
		return "", tracederrors.TracedError("libvirtXml is empty string")
	}

	domcfg := &libvirtxml.Domain{}
	err = domcfg.Unmarshal(libvirtXml)
	if err != nil {
		return "", tracederrors.TracedError(err)
	}

	networkInterfaces := domcfg.Devices.Interfaces
	nInterfaces := len(networkInterfaces)
	if nInterfaces != 1 {
		return "", tracederrors.TracedErrorf(
			"Only exactly one network interface is supported at the moment but got '%d'",
			nInterfaces,
		)
	}

	source := networkInterfaces[0].Source
	if source == nil {
		return "", tracederrors.TracedError("interface source is nil after evaluation")
	}

	if source.Bridge == nil {
		return "", tracederrors.TracedError("interface is not attached to a bridge")
	}

	bridgeInterface = source.Bridge.Bridge
	if bridgeInterface == "" {
		return "", tracederrors.TracedError("bridgeInterface is empty string after evaluation")
	}

	return bridgeInterface, nil
}

func GetVncPortFromXmlString(domainXml string) (vncPort int, err error) {
	if domainXml == "" {
		return -1, tracederrors.TracedErrorEmptyString("domainXml")
	}

	parsed := struct {
		Devices struct {
			Graphics []struct {
				Type string `xml:"type,attr"`
				Port int    `xml:"port,attr"`
			} `xml:"graphics"`
		} `xml:"devices"`
	}{}

	err = xml.Unmarshal([]byte(domainXml), &parsed)
	if err != nil {
		return -1, tracederrors.TracedErrorf("Unable to parse domain XML for VNC port: %w", err)
	}

	for _, graphics := range parsed.Devices.Graphics {
		if graphics.Type != "vnc" {
			continue
		}

		// port='-1' means libvirt auto-allocates and no port is assigned yet
		// (e.g. VM not running or autoport in use without a resolved port).
		if graphics.Port < 0 {
			return -1, tracederrors.TracedError("VNC is configured but no concrete port is assigned yet (port='-1'). Is the VM running?")
		}

		return graphics.Port, nil
	}

	return -1, tracederrors.TracedError("No VNC graphics device found in domain XML.")
}

func WriteXmlForVm(ctx context.Context, createOptions *kvmutilsoptions.KvmCreateVmOptions, outputFile filesinterfaces.File) (err error) {
	if createOptions == nil {
		return tracederrors.TracedError("createOptions is nil")
	}

	if outputFile == nil {
		return tracederrors.TracedError("outputFile is nil")
	}

	xmlString, err := CreateXmlForVmAsString(createOptions)
	if err != nil {
		return err
	}

	err = outputFile.WriteString(ctx, xmlString, &filesoptions.WriteOptions{})
	if err != nil {
		return err
	}

	outputPath, err := outputFile.GetLocalPath()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Created xml for VM to: '%s'", outputPath)

	return nil
}

// GetAdditionalBlockDevicesFromXmlString returns all disks of the domain XML which are backed
// by a host block device (<disk type='block'>) in the order they are defined.
// The file based root disk is not included.
func GetAdditionalBlockDevicesFromXmlString(libvirtXml string) ([]BlockDeviceAttachment, error) {
	if libvirtXml == "" {
		return nil, tracederrors.TracedErrorEmptyString("libvirtXml")
	}

	domcfg := &libvirtxml.Domain{}
	err := domcfg.Unmarshal(libvirtXml)
	if err != nil {
		return nil, tracederrors.TracedError(err)
	}

	ret := []BlockDeviceAttachment{}

	if domcfg.Devices == nil {
		return ret, nil
	}

	for _, disk := range domcfg.Devices.Disks {
		if disk.Source == nil || disk.Source.Block == nil {
			// Not a block device backed disk (e.g. the qcow2 root disk).
			continue
		}

		sourcePath := disk.Source.Block.Dev
		if sourcePath == "" {
			return nil, tracederrors.TracedError("Block device disk without source dev found in domain XML.")
		}

		if disk.Target == nil || disk.Target.Dev == "" {
			return nil, tracederrors.TracedErrorf("Block device disk '%s' has no target dev in domain XML.", sourcePath)
		}

		ret = append(ret, BlockDeviceAttachment{
			SourcePath: sourcePath,
			TargetDev:  disk.Target.Dev,
			Serial:     disk.Serial,
			ReadOnly:   disk.ReadOnly != nil,
		})
	}

	return ret, nil
}
