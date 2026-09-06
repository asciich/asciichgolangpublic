package kubeletutils

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/osutils/linuxutils/systemdutils"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// StartKubelet starts the kubelet service on the local machine.
func StartKubelet(ctx context.Context) error {
	return StartKubeletUsingCommandExecutor(ctx, commandexecutorexecoo.Exec())
}

// StartKubeletUsingCommandExecutor starts the kubelet service using the given command executor.
func StartKubeletUsingCommandExecutor(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	logging.LogInfoByCtxf(ctx, "Start kubelet service started.")

	err := systemdutils.StartService(ctx, commandExecutor, kubeletServiceName)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Start kubelet service finished.")
	return nil
}

// EnableKubelet enables the kubelet service (start on boot) on the local machine.
func EnableKubelet(ctx context.Context) error {
	return EnableKubeletUsingCommandExecutor(ctx, commandexecutorexecoo.Exec())
}

// EnableKubeletUsingCommandExecutor enables the kubelet service using the given command executor.
func EnableKubeletUsingCommandExecutor(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	logging.LogInfoByCtxf(ctx, "Enable kubelet service started.")

	err := systemdutils.EnableService(ctx, commandExecutor, kubeletServiceName)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Enable kubelet service finished.")
	return nil
}

// StartAndEnableKubelet enables and starts the kubelet service on the local
// machine in a single call (equivalent to `systemctl enable --now kubelet`).
func StartAndEnableKubelet(ctx context.Context) error {
	return StartAndEnableKubeletUsingCommandExecutor(ctx, commandexecutorexecoo.Exec())
}

// StartAndEnableKubeletUsingCommandExecutor enables and starts the kubelet
// service using the given command executor (equivalent to
// `systemctl enable --now kubelet`).
func StartAndEnableKubeletUsingCommandExecutor(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	logging.LogInfoByCtxf(ctx, "Start and enable kubelet service started.")

	err := systemdutils.EnableAndStartService(ctx, commandExecutor, kubeletServiceName)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Start and enable kubelet service finished.")
	return nil
}
