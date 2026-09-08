package ciliumutils

import (
	"context"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// prepareInstallCiliumOptions applies defaults to the given options.
func prepareInstallCiliumOptions(options *InstallCiliumOptions) *InstallCiliumOptions {
	if options == nil {
		options = DefaultInstallCiliumOptions()
	}
	// UseSudo defaults to true if not set
	return options
}

// IsCiliumInstalled returns true if Cilium is already installed in the cluster
// reachable through the given commandExecutor.
func IsCiliumInstalled(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor) (bool, error) {
	if commandExecutor == nil {
		return false, tracederrors.TracedErrorNil("commandExecutor")
	}

	stdout, err := commandExecutor.RunCommandAndGetStdoutAsString(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{
				"sh", "-c",
				"KUBECONFIG=/etc/kubernetes/admin.conf kubectl -n kube-system get daemonset cilium >/dev/null 2>&1 && echo yes || echo no",
			},
		},
	)
	if err != nil {
		return false, err
	}

	isInstalled := strings.TrimSpace(stdout) == "yes"

	if isInstalled {
		logging.LogInfoByCtxf(ctx, "Cilium is installed in kubernetes cluster.")
	} else {
		logging.LogInfoByCtxf(ctx, "Cilium is not installed in kubernetes cluster.")
	}

	return isInstalled, nil
}

// Install installs Cilium (eBPF CNI) into the cluster on the local machine.
func Install(ctx context.Context, options *InstallCiliumOptions) error {
	return InstallUsingCommandExecutor(ctx, commandexecutorexecoo.Exec(), options)
}

// InstallUsingCommandExecutor installs Cilium (eBPF CNI) into the cluster using
// the Cilium CLI via the given command executor. The operation is skipped when
// Cilium is already installed (idempotent).
func InstallUsingCommandExecutor(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, options *InstallCiliumOptions) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	options = prepareInstallCiliumOptions(options)

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Install Cilium into cluster on '%s' started.", hostDescription)

	alreadyInstalled, err := IsCiliumInstalled(ctx, commandExecutor)
	if err != nil {
		return err
	}

	if alreadyInstalled {
		logging.LogInfoByCtxf(ctx, "Cilium is already installed on '%s'. Skip installation.", hostDescription)
		logging.LogInfoByCtxf(ctx, "Install Cilium into cluster on '%s' finished.", hostDescription)
		return nil
	}

	// Ensure the cilium CLI is present (installutils skips if already in place).
	err = InstallCiliumCliUsingCommandExecutor(ctx, commandExecutor, options)
	if err != nil {
		return err
	}

	// Build the "cilium install" command.
	command := []string{"cilium", "install"}
	if options.Version != "" {
		command = append(command, "--version", options.Version)
	}
	if options.PodNetworkCidr != "" {
		// Make Cilium honor kubeadm's pod CIDR.
		command = append(command,
			"--set", "ipam.mode=kubernetes",
			"--set", "ipam.operator.clusterPoolIPv4PodCIDRList="+options.PodNetworkCidr,
		)
	}
	// Additional Helm values, e.g. kernel-compatibility workarounds like
	// disabling socketLB on very new kernels whose eBPF verifier rejects the
	// CGroupSock program.
	for _, kv := range options.ExtraHelmSet {
		command = append(command, "--set", kv)
	}

	// The cilium CLI talks to the cluster via kubeconfig. On a freshly bootstrapped
	// kubeadm control-plane the admin kubeconfig lives at /etc/kubernetes/admin.conf,
	// so we point KUBECONFIG at it explicitly.
	_, err = commandExecutor.RunCommand(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{
				"sh", "-c",
				"KUBECONFIG=/etc/kubernetes/admin.conf " + strings.Join(command, " "),
			},
		},
	)
	if err != nil {
		return err
	}

	logging.LogChangedByCtxf(ctx, "Installed Cilium into cluster on '%s'.", hostDescription)
	logging.LogInfoByCtxf(ctx, "Install Cilium into cluster on '%s' finished.", hostDescription)

	return nil
}

// Uninstall removes Cilium from the cluster on the local machine.
func Uninstall(ctx context.Context) error {
	return UninstallUsingCommandExecutor(ctx, commandexecutorexecoo.Exec())
}

// UninstallUsingCommandExecutor removes Cilium from the cluster using the Cilium
// CLI via the given command executor. Useful to clean up a broken / crash-looping
// installation before reinstalling.
func UninstallUsingCommandExecutor(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Uninstall Cilium on '%s' started.", hostDescription)

	_, err = commandExecutor.RunCommand(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{
				"sh", "-c",
				"KUBECONFIG=/etc/kubernetes/admin.conf cilium uninstall",
			},
		},
	)
	if err != nil {
		return err
	}

	logging.LogChangedByCtxf(ctx, "Uninstalled Cilium on '%s'.", hostDescription)
	logging.LogInfoByCtxf(ctx, "Uninstall Cilium on '%s' finished.", hostDescription)

	return nil
}

// WaitUntilReady waits until Cilium reports a healthy status in the cluster on
// the local machine.
func WaitUntilReady(ctx context.Context, options *InstallCiliumOptions) error {
	return WaitUntilReadyUsingCommandExecutor(ctx, commandexecutorexecoo.Exec(), options)
}

// WaitUntilReadyUsingCommandExecutor blocks until "cilium status --wait" reports
// Cilium as healthy, using the given command executor.
func WaitUntilReadyUsingCommandExecutor(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, options *InstallCiliumOptions) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Wait until Cilium is ready on '%s' started.", hostDescription)

	_, err = commandExecutor.RunCommand(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{
				"sh", "-c",
				"KUBECONFIG=/etc/kubernetes/admin.conf cilium status --wait",
			},
		},
	)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Wait until Cilium is ready on '%s' finished.", hostDescription)

	return nil
}
