package nativekvmutils

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"github.com/asciich/asciichgolangpublic/pkg/vmutils/kvmutils/kvmutilsgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/vmutils/kvmutils/kvmutilsinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/vmutils/kvmutils/kvmutilsoptions"
)

func (k *NativeKvmHypervisor) CreateVm(ctx context.Context, createOptions *kvmutilsoptions.KvmCreateVmOptions) (createdVm kvmutilsinterfaces.VM, err error) {
	if createOptions == nil {
		return nil, tracederrors.TracedError("createOptions is nil")
	}

	vmName, err := createOptions.GetVmName()
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Create KVM VM '%s' started.", vmName)

	exists, err := k.VmByNameExists(ctx, vmName)
	if err != nil {
		return nil, err
	}

	if exists {
		logging.LogInfoByCtxf(ctx, "VM '%s' already exists.", vmName)

		createdVm, err = k.GetVmByName(ctx, vmName)
		if err != nil {
			return nil, err
		}

		return createdVm, nil
	}

	diskImage, err := createOptions.GetDiskImage()
	if err != nil {
		return nil, err
	}

	diskImagePath, err := diskImage.GetLocalPath()
	if err != nil {
		return nil, err
	}

	diskImageExists, err := diskImage.Exists(ctx)
	if err != nil {
		return nil, err
	}

	if diskImageExists {
		logging.LogInfoByCtxf(ctx, "Going to use existing disk image '%s' to create VM '%s'.", diskImagePath, vmName)
	} else {
		return nil, tracederrors.TracedErrorf("Disk image '%s' does not exist to create VM '%s'.", diskImagePath, vmName)
	}

	domainXml, err := kvmutilsgeneric.CreateXmlForVmAsString(createOptions)
	if err != nil {
		return nil, err
	}

	connection, err := k.getLibvirtConnection(ctx)
	if err != nil {
		return nil, err
	}

	// DomainDefineXML registers a PERSISTENT domain from the given XML, the
	// native equivalent of 'virsh define <xml>'. Unlike DomainCreateXML (which
	// creates a transient domain that is lost on host reboot), a defined domain
	// survives reboots.
	domain, err := connection.DomainDefineXML(domainXml)
	if err != nil {
		return nil, tracederrors.TracedErrorf("failed to define VM '%s' from XML: %w", vmName, err)
	}

	// DomainSetAutostart(domain, 1) enables automatic start of the VM when the
	// host boots, the native equivalent of 'virsh autostart <name>'. Persistence
	// alone only keeps the definition; autostart also powers the VM on after a
	// host reboot.
	err = connection.DomainSetAutostart(domain, 1)
	if err != nil {
		return nil, tracederrors.TracedErrorf("failed to enable autostart for VM '%s': %w", vmName, err)
	}

	// DomainCreate boots the now-defined (but still inactive) domain, the
	// native equivalent of 'virsh start <name>'.
	err = connection.DomainCreate(domain)
	if err != nil {
		return nil, tracederrors.TracedErrorf("failed to start defined VM '%s': %w", vmName, err)
	}

	createdVm, err = k.GetVmByName(ctx, vmName)
	if err != nil {
		return nil, err
	}

	logging.LogChangedByCtxf(ctx, "VM '%s' created.", vmName)

	logging.LogInfoByCtxf(ctx, "Create KVM VM '%s' finished.", vmName)

	return createdVm, nil
}
