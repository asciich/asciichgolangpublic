package kubeadmutils

import (
	"context"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfile"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

const (
	// adminConfPath is the kubeconfig kubeadm writes on the control-plane after "kubeadm init".
	adminConfPath = "/etc/kubernetes/admin.conf"

	// controlPlaneTaintKey is the modern control-plane taint key
	// (node-role.kubernetes.io/master is deprecated).
	controlPlaneTaintKey = "node-role.kubernetes.io/control-plane"
)

// prepareInitControlPlaneOptions applies defaults to the given options.
func prepareInitControlPlaneOptions(options *InitControlPlaneOptions) *InitControlPlaneOptions {
	if options == nil {
		options = DefaultInitControlPlaneOptions()
	}
	if options.PodNetworkCidr == "" {
		options.PodNetworkCidr = "10.244.0.0/16"
	}
	// UseSudo defaults to true if not set
	return options
}

// sudoCommand prefixes the given command with "sudo" when useSudo is true.
func sudoCommand(useSudo bool, command ...string) []string {
	if useSudo {
		return append([]string{"sudo"}, command...)
	}
	return command
}

// IsControlPlaneInitialized returns true if this node has already been
// initialized as a control-plane (i.e. /etc/kubernetes/admin.conf exists).
func IsControlPlaneInitialized(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor) (bool, error) {
	if commandExecutor == nil {
		return false, tracederrors.TracedErrorNil("commandExecutor")
	}

	return commandexecutorfile.Exists(ctx, commandExecutor, adminConfPath)
}

// InitControlPlane initializes the control-plane on the local machine.
func InitControlPlane(ctx context.Context, options *InitControlPlaneOptions) error {
	return InitControlPlaneUsingCommandExecutor(ctx, commandexecutorexecoo.Exec(), options)
}

// InitControlPlaneUsingCommandExecutor runs "kubeadm init" using the given
// command executor. The operation is skipped when the node is already
// initialized (idempotent).
func InitControlPlaneUsingCommandExecutor(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, options *InitControlPlaneOptions) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	options = prepareInitControlPlaneOptions(options)

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Init kubernetes control-plane on '%s' started.", hostDescription)

	alreadyInitialized, err := IsControlPlaneInitialized(ctx, commandExecutor)
	if err != nil {
		return err
	}

	if alreadyInitialized {
		logging.LogInfoByCtxf(ctx, "Control-plane on '%s' is already initialized ('%s' exists). Skip 'kubeadm init'.", hostDescription, adminConfPath)
		logging.LogInfoByCtxf(ctx, "Init kubernetes control-plane on '%s' finished.", hostDescription)
		return nil
	}

	command := []string{"kubeadm", "init"}
	if options.PodNetworkCidr != "" {
		command = append(command, "--pod-network-cidr="+options.PodNetworkCidr)
	}
	if options.ApiServerAdvertiseAddress != "" {
		command = append(command, "--apiserver-advertise-address="+options.ApiServerAdvertiseAddress)
	}
	if options.KubernetesVersion != "" {
		command = append(command, "--kubernetes-version="+options.KubernetesVersion)
	}

	_, err = commandExecutor.RunCommand(
		commandexecutorgeneric.WithLiveOutputOnStdoutIfVerbose(ctx),
		&parameteroptions.RunCommandOptions{
			Command: sudoCommand(options.UseSudo, command...),
		},
	)
	if err != nil {
		return err
	}

	logging.LogChangedByCtxf(ctx, "Initialized kubernetes control-plane on '%s'.", hostDescription)
	logging.LogInfoByCtxf(ctx, "Init kubernetes control-plane on '%s' finished.", hostDescription)

	return nil
}

// SetupAdminKubeconfig copies /etc/kubernetes/admin.conf to the invoking user's
// ~/.kube/config on the local machine.
func SetupAdminKubeconfig(ctx context.Context) error {
	return SetupAdminKubeconfigUsingCommandExecutor(ctx, commandexecutorexecoo.Exec())
}

// SetupAdminKubeconfigUsingCommandExecutor copies /etc/kubernetes/admin.conf to
// the invoking user's ~/.kube/config using the given command executor, so
// kubectl can talk to the cluster.
func SetupAdminKubeconfigUsingCommandExecutor(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Setup admin kubeconfig on '%s' started.", hostDescription)

	// Create ~/.kube, copy admin.conf and take ownership so kubectl works for the
	// invoking (non-root) user. Uses a shell so $HOME / id are resolved on the target.
	setupCmd := "mkdir -p \"$HOME/.kube\" && " +
		"sudo cp -f " + adminConfPath + " \"$HOME/.kube/config\" && " +
		"sudo chown \"$(id -u):$(id -g)\" \"$HOME/.kube/config\""

	_, err = commandExecutor.RunCommand(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{"sh", "-c", setupCmd},
		},
	)
	if err != nil {
		return err
	}

	logging.LogChangedByCtxf(ctx, "Set up admin kubeconfig in '$HOME/.kube/config' on '%s'.", hostDescription)
	logging.LogInfoByCtxf(ctx, "Setup admin kubeconfig on '%s' finished.", hostDescription)

	return nil
}

// UntaintControlPlane removes the control-plane NoSchedule taint on the local machine.
func UntaintControlPlane(ctx context.Context) error {
	return UntaintControlPlaneUsingCommandExecutor(ctx, commandexecutorexecoo.Exec())
}

// UntaintControlPlaneUsingCommandExecutor removes the control-plane NoSchedule
// taint from all nodes using the given command executor, so workloads can be
// scheduled on the control-plane (useful for single-node clusters).
func UntaintControlPlaneUsingCommandExecutor(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Untaint kubernetes control-plane on '%s' started.", hostDescription)

	// "kubectl taint ... <key>-" is idempotent enough: it returns non-zero only
	// if the taint is not found. We tolerate that by checking the output.
	stdout, err := commandExecutor.RunCommandAndGetStdoutAsString(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{
				"sh", "-c",
				"kubectl taint nodes --all " + controlPlaneTaintKey + "- 2>&1 || true",
			},
		},
	)
	if err != nil {
		return err
	}

	if strings.Contains(stdout, "not found") {
		logging.LogInfoByCtxf(ctx, "Control-plane taint already removed on '%s'. Skip untainting.", hostDescription)
	} else {
		logging.LogChangedByCtxf(ctx, "Removed control-plane taint on '%s'.", hostDescription)
	}

	logging.LogInfoByCtxf(ctx, "Untaint kubernetes control-plane on '%s' finished.", hostDescription)

	return nil
}

// ApplyManifest applies a kubernetes manifest (file path or URL) on the local machine.
func ApplyManifest(ctx context.Context, manifestUrlOrPath string) error {
	return ApplyManifestUsingCommandExecutor(ctx, commandexecutorexecoo.Exec(), manifestUrlOrPath)
}

// ApplyManifestUsingCommandExecutor runs "kubectl apply -f <manifest>" using the
// given command executor.
func ApplyManifestUsingCommandExecutor(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, manifestUrlOrPath string) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	if manifestUrlOrPath == "" {
		return tracederrors.TracedErrorEmptyString("manifestUrlOrPath")
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Apply manifest '%s' on '%s' started.", manifestUrlOrPath, hostDescription)

	_, err = commandExecutor.RunCommand(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{"kubectl", "apply", "-f", manifestUrlOrPath},
		},
	)
	if err != nil {
		return err
	}

	logging.LogChangedByCtxf(ctx, "Applied manifest '%s' on '%s'.", manifestUrlOrPath, hostDescription)
	logging.LogInfoByCtxf(ctx, "Apply manifest '%s' on '%s' finished.", manifestUrlOrPath, hostDescription)

	return nil
}
