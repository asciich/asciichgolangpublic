package kvmutils

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
	libvirtxml "libvirt.org/libvirt-go-xml"
)

//go:embed data/vm.xml.tmpl
var vm_on_laptopt_xml_tmpl string

type LibvirtXmlsService struct{}

func LibvirtXmls() (libvirtXmls *LibvirtXmlsService) {
	return NewLibvirtXmlsService()
}

func NewLibvirtXmlsService() (libvirtXmls *LibvirtXmlsService) {
	return new(LibvirtXmlsService)
}

func (l *LibvirtXmlsService) CreateXmlForVmAsString(createOptions *KvmCreateVmOptions) (libvirtXml string, err error) {
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

	// The MAC address is optional. When empty libvirt auto-generates one.
	// Only validate the format when a MAC address was explicitly provided.
	macAddress := createOptions.MacAddress
	if macAddress != "" {
		err = macaddresses.CheckStringIsAMacAddress(macAddress)
		if err != nil {
			return "", err
		}
	}

	libvirtXml, err = gotemplateutils.RenderTemplateFromStringAsString(
		vm_on_laptopt_xml_tmpl,
		map[string]any{
			"VM_NAME":     vmName,
			"DISK_PATH":   diskPath,
			"MAC_ADDRESS": macAddress,
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
		return "", tracederrors.TracedError(err.Error())
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
		return "", tracederrors.TracedError(err.Error())
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

func (l *LibvirtXmlsService) WriteXmlForVmOnLatopToFile(ctx context.Context, createOptions *KvmCreateVmOptions, outputFile filesinterfaces.File) (err error) {
	if createOptions == nil {
		return tracederrors.TracedError("createOptions is nil")
	}

	if outputFile == nil {
		return tracederrors.TracedError("outputFile is nil")
	}

	xmlString, err := l.CreateXmlForVmAsString(createOptions)
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

	logging.LogInfoByCtxf(ctx, "Created xml for laptop on VM to: '%s'", outputPath)

	return nil
}
