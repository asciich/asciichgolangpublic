package kubeletutils

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// InstallKubeletAndRunAsService installs the kubelet binary, installs the
// systemd service unit (plus kubeadm drop-in) and finally enables and starts
// the kubelet service on the local machine.
func InstallKubeletAndRunAsService(ctx context.Context, installKubeletOptions *InstallKubeletOptions, installKubeletServiceOptions *InstallKubeletServiceOptions) error {
	return InstallKubeletAndRunAsServiceUsingCommandExecutor(ctx, commandexecutorexecoo.Exec(), installKubeletOptions, installKubeletServiceOptions)
}

// InstallKubeletAndRunAsServiceUsingCommandExecutor installs the kubelet binary,
// installs the systemd service unit (plus kubeadm drop-in) and finally enables
// and starts the kubelet service using the given command executor.
func InstallKubeletAndRunAsServiceUsingCommandExecutor(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, installKubeletOptions *InstallKubeletOptions, installKubeletServiceOptions *InstallKubeletServiceOptions) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	logging.LogInfoByCtxf(ctx, "Install kubelet and run as service started.")

	// 1. Install the kubelet binary.
	err := InstallKubeletUsingCommandExecutor(ctx, commandExecutor, installKubeletOptions)
	if err != nil {
		return err
	}

	// 2. Install the systemd service unit and the kubeadm drop-in.
	err = InstallKubeletServiceUsingCommandExecutor(ctx, commandExecutor, installKubeletServiceOptions)
	if err != nil {
		return err
	}

	// 3. Enable and start the kubelet service.
	err = StartAndEnableKubeletUsingCommandExecutor(ctx, commandExecutor)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Install kubelet and run as service finished.")
	return nil
}
