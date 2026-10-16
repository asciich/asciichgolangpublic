package hostreconciler

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"github.com/asciich/asciichgolangpublic/pkg/userutils"
	"github.com/asciich/asciichgolangpublic/pkg/vmutils/kvmutils/installkvm"
)

// checkKVMPrivileges returns an error if at least one KVM action is configured but the process is not running as root.
// Called before the before_script and any change is applied, so the host is never left in a half applied state.
func checkKVMPrivileges(ctx context.Context, kvmConfig *KVMConfig) error {
	if kvmConfig == nil {
		return tracederrors.TracedErrorNil("kvmConfig")
	}

	if !kvmConfig.hasActions() {
		return nil
	}

	isRoot, err := userutils.IsRunningAsRoot(ctx)
	if err != nil {
		return err
	}

	if !isRoot {
		return tracederrors.TracedErrorf("Reconcile KVM (start: '%t', enable: '%t') requires root privileges.", kvmConfig.Start, kvmConfig.Enable)
	}

	return nil
}

// reconcileKVM applies the desired 'kvm' configuration to the current host.
//
// Requires root privileges if at least one KVM action is configured.
func reconcileKVM(ctx context.Context, kvmConfig *KVMConfig) error {
	if kvmConfig == nil {
		return tracederrors.TracedErrorNil("kvmConfig")
	}

	if !kvmConfig.hasActions() {
		logging.LogInfoByCtxf(ctx, "No KVM actions configured. Reconcile KVM skipped.")
		return nil
	}

	logging.LogInfoByCtxf(ctx, "Reconcile KVM (start: '%t', enable: '%t') started.", kvmConfig.Start, kvmConfig.Enable)

	if kvmConfig.Enable {
		err := installkvm.Enable(ctx)
		if err != nil {
			return err
		}
	}

	if kvmConfig.Start {
		err := installkvm.Start(ctx)
		if err != nil {
			return err
		}
	}

	logging.LogInfoByCtxf(ctx, "Reconcile KVM (start: '%t', enable: '%t') finished.", kvmConfig.Start, kvmConfig.Enable)

	return nil
}
