package nativekvmutils

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// IsVmPersistent reports whether the VM is persistent (survives host reboots
// and 'virsh define'). A transient VM (created via DomainCreateXML without a
// config definition) returns false.
func (k *NativeKvmHypervisor) IsVmPersistent(ctx context.Context, vmName string) (isPersistent bool, err error) {
	if vmName == "" {
		return false, tracederrors.TracedErrorEmptyString("vmName")
	}

	connection, err := k.getLibvirtConnection(ctx)
	if err != nil {
		return false, err
	}

	domain, err := connection.DomainLookupByName(vmName)
	if err != nil {
		return false, tracederrors.TracedErrorf("failed to look up domain '%s': %w", vmName, err)
	}

	persistentInt, err := connection.DomainIsPersistent(domain)
	if err != nil {
		return false, tracederrors.TracedErrorf("failed to check if domain '%s' is persistent: %w", vmName, err)
	}

	isPersistent = persistentInt == 1

	logging.LogInfoByCtxf(ctx, "VM '%s' persistent: %v.", vmName, isPersistent)

	return isPersistent, nil
}
