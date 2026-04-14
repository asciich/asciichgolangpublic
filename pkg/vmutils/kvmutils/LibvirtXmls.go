package kvmutils

import (
	"context"
	_ "embed"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	libvirtxml "libvirt.org/libvirt-go-xml"
)

// TODO enable again //go:embed data/LibvirtXmls/VmOnLaptop/vm_on_laptop.xml.tmpl
// TODO enable again var vm_on_laptopt_xml_tmpl string

type LibvirtXmlsService struct{}

func LibvirtXmls() (libvirtXmls *LibvirtXmlsService) {
	return NewLibvirtXmlsService()
}

func NewLibvirtXmlsService() (libvirtXmls *LibvirtXmlsService) {
	return new(LibvirtXmlsService)
}

func (l *LibvirtXmlsService) CreateXmlForVmOnLatopAsString(createOptions *KvmCreateVmOptions) (libvirtXml string, err error) {
	return "", tracederrors.TracedErrorNotImplemented()
	/* TODO enable again
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

	macAddress, err := createOptions.GetMacAddress()
	if err != nil {
		return "", err
	}

	_, err = MacAddresses().CheckStringIsAMacAddress(macAddress)
	if err != nil {
		return "", err
	}

	libvirtXml, err = GoTemplate().RenderTemplateFromStringAsString(
		vm_on_laptopt_xml_tmpl,
		map[string]interface{}{
			"VM_NAME":     vmName,
			"DISK_PATH":   diskPath,
			"MAC_ADDRESS": macAddress,
		},
	)
	if err != nil {
		return "", err
	}

	return libvirtXml, nil
	*/
}

func (l *LibvirtXmlsService) GetMacAddressFromXmlString(libvirtXml string) (macAddress string, err error) {
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

func (l *LibvirtXmlsService) WriteXmlForVmOnLatopToFile(ctx context.Context, createOptions *KvmCreateVmOptions, outputFile filesinterfaces.File) (err error) {
	if createOptions == nil {
		return tracederrors.TracedError("createOptions is nil")
	}

	if outputFile == nil {
		return tracederrors.TracedError("outputFile is nil")
	}

	xmlString, err := l.CreateXmlForVmOnLatopAsString(createOptions)
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
