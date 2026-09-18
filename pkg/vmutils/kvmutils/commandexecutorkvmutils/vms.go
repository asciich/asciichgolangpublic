package commandexecutorkvmutils

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfilesoo"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"github.com/asciich/asciichgolangpublic/pkg/vmutils/kvmutils/kvmutilsgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/vmutils/kvmutils/kvmutilsinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/vmutils/kvmutils/kvmutilsoptions"
)

func (k *CommandExecutrKvmHypervisor) CreateVm(ctx context.Context, createOptions *kvmutilsoptions.KvmCreateVmOptions) (createdVm kvmutilsinterfaces.VM, err error) {
	if createOptions == nil {
		return nil, tracederrors.TracedError("createOptions is nil")
	}

	vmName, err := createOptions.GetVmName()
	if err != nil {
		return nil, err
	}

	exists, err := k.VmByNameExists(ctx, vmName)
	if err != nil {
		return nil, err
	}

	if exists {
		logging.LogInfoByCtxf(ctx, "VM '%s' already exists", vmName)

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

	if !diskImageExists {
		return nil, tracederrors.TracedErrorf("Disk image '%s' does not exist to create VM.", diskImagePath)
	}

	vmXml, err := tempfilesoo.CreateEmptyTemporaryFile(ctx)
	if err != nil {
		return nil, err
	}
	defer vmXml.Delete(ctx, &filesoptions.DeleteOptions{})

	err = kvmutilsgeneric.WriteXmlForVm(ctx, createOptions, vmXml)
	if err != nil {
		return nil, err
	}

	vmXmlPath, err := vmXml.GetLocalPath()
	if err != nil {
		return nil, err
	}

	// 'virsh define' registers a PERSISTENT domain (survives host reboots).
	// This replaces 'virsh create' which would create a transient domain that
	// is lost on reboot.
	defineOutput, err := k.RunKvmCommandAndGetStdout(
		ctx,
		[]string{"define", vmXmlPath},
	)
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Output of VM '%s' definition:\n%s", vmName, defineOutput)

	// 'virsh autostart' enables automatic start of the VM when the host boots.
	// Persistence alone only keeps the definition; autostart also powers the
	// VM on after a host reboot.
	autostartOutput, err := k.RunKvmCommandAndGetStdout(
		ctx,
		[]string{"autostart", vmName},
	)
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Output of VM '%s' autostart enable:\n%s", vmName, autostartOutput)

	// 'virsh start' boots the now-defined (but still inactive) domain.
	startOutput, err := k.RunKvmCommandAndGetStdout(
		ctx,
		[]string{"start", vmName},
	)
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Output of VM '%s' start:\n%s", vmName, startOutput)

	createdVm, err = k.GetVmByName(ctx, vmName)
	if err != nil {
		return nil, err
	}

	logging.LogChangedByCtxf(ctx, "Vm '%s' created.", vmName)

	return createdVm, nil
}
