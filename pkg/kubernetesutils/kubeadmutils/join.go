package kubeadmutils

import (
	"context"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/commandexecutornicutils"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// JoinControlPlaneOptions contains options for joining an additional
// control-plane node with "kubeadm join".
type JoinControlPlaneOptions struct {
	// JoinCommand is the full "kubeadm join ... --control-plane --certificate-key ..."
	// command to run.
	JoinCommand string

	// ControlPlaneEndpoint is the shared/virtual IP (VIP) of the cluster.
	ControlPlaneEndpoint string

	// EnableKubeVip deploys the KubeVip static-pod manifest before joining.
	EnableKubeVip bool

	// KubeVipInterface is the interface KubeVip binds the VIP to.
	// Optional: auto-detected when empty.
	KubeVipInterface string

	// KubeVipVersion is the KubeVip image version.
	// Optional: a built-in default is used when empty.
	KubeVipVersion string

	// UseSudo determines if sudo should be used.
	// Default: true
	UseSudo bool
}

// GetControlPlaneJoinCommand returns a control-plane join command from the local machine.
func GetControlPlaneJoinCommand(ctx context.Context) (string, error) {
	return GetControlPlaneJoinCommandUsingCommandExecutor(ctx, commandexecutorexecoo.Exec())
}

// GetControlPlaneJoinCommandUsingCommandExecutor assembles a "kubeadm join"
// command for joining an additional control-plane node using the given command
// executor. Control-plane certificates are re-uploaded so the returned command
// stays valid on an already running cluster.
func GetControlPlaneJoinCommandUsingCommandExecutor(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor) (string, error) {
	if commandExecutor == nil {
		return "", tracederrors.TracedErrorNil("commandExecutor")
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return "", err
	}

	logging.LogInfoByCtxf(ctx, "Get control-plane join command on '%s' started.", hostDescription)

	// Base (worker) join command including token and CA cert hash.
	joinCommand, err := commandExecutor.RunCommandAndGetStdoutAsString(
		commandexecutorgeneric.WithLiveOutputOnStdoutIfVerbose(ctx),
		&parameteroptions.RunCommandOptions{
			Command: []string{"sudo", "kubeadm", "token", "create", "--print-join-command"},
		},
	)
	if err != nil {
		return "", err
	}
	joinCommand = strings.TrimSpace(joinCommand)

	// (Re-)upload the control-plane certificates and grab the certificate key.
	uploadOutput, err := commandExecutor.RunCommandAndGetStdoutAsString(
		commandexecutorgeneric.WithLiveOutputOnStdoutIfVerbose(ctx),
		&parameteroptions.RunCommandOptions{
			Command: []string{"sudo", "kubeadm", "init", "phase", "upload-certs", "--upload-certs"},
		},
	)
	if err != nil {
		return "", err
	}

	// The certificate key is the last non-empty line of the upload output.
	certificateKey := ""
	for _, line := range strings.Split(uploadOutput, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		certificateKey = line
	}

	if certificateKey == "" {
		return "", tracederrors.TracedErrorf("Unable to determine the control-plane certificate key on '%s'.", hostDescription)
	}

	ret := joinCommand + " --control-plane --certificate-key " + certificateKey

	logging.LogInfoByCtxf(ctx, "Get control-plane join command on '%s' finished.", hostDescription)

	return ret, nil
}

// JoinControlPlane joins an additional control-plane node on the local machine.
func JoinControlPlane(ctx context.Context, options *JoinControlPlaneOptions) error {
	return JoinControlPlaneUsingCommandExecutor(ctx, commandexecutorexecoo.Exec(), options)
}

// JoinControlPlaneUsingCommandExecutor runs the given "kubeadm join" command to
// add an additional control-plane node using the given command executor. The
// operation is skipped when the node is already initialized (idempotent).
func JoinControlPlaneUsingCommandExecutor(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, options *JoinControlPlaneOptions) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	if options.JoinCommand == "" {
		return tracederrors.TracedErrorEmptyString("options.JoinCommand")
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Join kubernetes control-plane on '%s' started.", hostDescription)

	alreadyInitialized, err := IsControlPlaneInitialized(ctx, commandExecutor)
	if err != nil {
		return err
	}

	if alreadyInitialized {
		logging.LogInfoByCtxf(ctx, "Control-plane on '%s' is already initialized ('%s' exists). Skip join.", hostDescription, adminConfPath)
		logging.LogInfoByCtxf(ctx, "Join kubernetes control-plane on '%s' finished.", hostDescription)
		return nil
	}

	// KubeVip must serve the VIP on the joining node too so it can float here.
	if options.EnableKubeVip {
		if options.ControlPlaneEndpoint == "" {
			return tracederrors.TracedError("EnableKubeVip is set but ControlPlaneEndpoint is empty: the VIP KubeVip advertises must be provided")
		}

		nicInterface := options.KubeVipInterface
		if nicInterface == "" {
			nicInterface, err = commandexecutornicutils.GetFirstActiveNicName(ctx, commandExecutor)
			if err != nil {
				return err
			}
		}

		err = DeployKubeVipManifestUsingCommandExecutor(ctx, commandExecutor, &DeployKubeVipManifestOptions{
			Vip:       vipAddressFromControlPlaneEndpoint(options.ControlPlaneEndpoint),
			Interface: nicInterface,
			Version:   options.KubeVipVersion,
			UseSudo:   options.UseSudo,
		})
		if err != nil {
			return err
		}
	}

	_, err = commandExecutor.RunCommand(
		commandexecutorgeneric.WithLiveOutputOnStdoutIfVerbose(ctx),
		&parameteroptions.RunCommandOptions{
			Command: sudoCommand(options.UseSudo, strings.Fields(options.JoinCommand)...),
		},
	)
	if err != nil {
		return err
	}

	logging.LogChangedByCtxf(ctx, "Joined kubernetes control-plane on '%s'.", hostDescription)
	logging.LogInfoByCtxf(ctx, "Join kubernetes control-plane on '%s' finished.", hostDescription)

	return nil
}

// GetNodeJoinCommand returns a worker-node join command from the local machine.
func GetNodeJoinCommand(ctx context.Context) (string, error) {
	return GetNodeJoinCommandUsingCommandExecutor(ctx, commandexecutorexecoo.Exec())
}

// GetNodeJoinCommandUsingCommandExecutor assembles a "kubeadm join" command for
// joining a worker node using the given command executor.
//
// In contrast to GetControlPlaneJoinCommandUsingCommandExecutor the returned
// command joins the node as a plain worker: no control-plane certificates are
// uploaded and neither the "--control-plane" flag nor a "--certificate-key" is
// appended.
//
// A fresh bootstrap token is created so the returned command stays valid even
// on an already running cluster (the default bootstrap token expires ~24h after
// creation).
func GetNodeJoinCommandUsingCommandExecutor(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor) (string, error) {
	if commandExecutor == nil {
		return "", tracederrors.TracedErrorNil("commandExecutor")
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return "", err
	}

	logging.LogInfoByCtxf(ctx, "Get worker-node join command on '%s' started.", hostDescription)

	// Base (worker) join command including token and CA cert hash.
	joinCommand, err := commandExecutor.RunCommandAndGetStdoutAsString(
		commandexecutorgeneric.WithLiveOutputOnStdoutIfVerbose(ctx),
		&parameteroptions.RunCommandOptions{
			Command: []string{"sudo", "kubeadm", "token", "create", "--print-join-command"},
		},
	)
	if err != nil {
		return "", err
	}
	joinCommand = strings.TrimSpace(joinCommand)

	if joinCommand == "" {
		return "", tracederrors.TracedErrorf("Unable to determine the worker-node join command on '%s'.", hostDescription)
	}

	ret := joinCommand

	logging.LogInfoByCtxf(ctx, "Get worker-node join command on '%s' finished.", hostDescription)

	return ret, nil
}

// kubeletConfPath is written by "kubeadm join"/"kubeadm init" on every node
// (both workers and control-planes) once it has joined a cluster.
const kubeletConfPath = "/etc/kubernetes/kubelet.conf"

// IsNodeJoined returns true if the given host has already joined a kubernetes
// cluster (i.e. the kubelet config kubeletConfPath exists).
//
// This is true for every joined node, including plain workers. Use
// IsControlPlaneInitialized instead when you specifically need to know whether
// a host is a control-plane (adminConfPath).
func IsNodeJoined(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor) (bool, error) {
	if commandExecutor == nil {
		return false, tracederrors.TracedErrorNil("commandExecutor")
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return false, err
	}

	// Do not rely on the exit code to check for existence: run the check in a
	// shell that always exits 0 and report the result via a well defined
	// "yes"/"no" token instead.
	output, err := commandExecutor.RunCommandAndGetStdoutAsString(
		commandexecutorgeneric.WithLiveOutputOnStdoutIfVerbose(ctx),
		&parameteroptions.RunCommandOptions{
			Command: []string{"sh", "-c", "ls " + kubeletConfPath + " &>/dev/null && echo yes || echo no"},
		},
	)
	if err != nil {
		return false, err
	}
	output = strings.TrimSpace(output)

	var ret bool
	switch output {
	case "yes":
		ret = true
	case "no":
		ret = false
	default:
		return false, tracederrors.TracedErrorf(
			"Unable to determine if node is joined on '%s': unexpected output '%s' while checking for '%s'.",
			hostDescription, output, kubeletConfPath,
		)
	}

	if ret {
		logging.LogInfoByCtxf(ctx, "Node on '%s' is already joined ('%s' exists).", hostDescription, kubeletConfPath)
	} else {
		logging.LogInfoByCtxf(ctx, "Node on '%s' is not joined ('%s' does not exist).", hostDescription, kubeletConfPath)
	}

	return ret, nil
}

// JoinNodeUsingCommandExecutor runs the given "kubeadm join" command to add a
// worker node using the given command executor. The operation is skipped when
// the node has already joined (idempotent).
func JoinNodeUsingCommandExecutor(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, options *JoinNodeOptions) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	if options.JoinCommand == "" {
		return tracederrors.TracedErrorEmptyString("options.JoinCommand")
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Join kubernetes worker node on '%s' started.", hostDescription)

	alreadyJoined, err := IsNodeJoined(ctx, commandExecutor)
	if err != nil {
		return err
	}

	if alreadyJoined {
		logging.LogInfoByCtxf(ctx, "Node on '%s' is already joined ('%s' exists). Skip join.", hostDescription, kubeletConfPath)
		logging.LogInfoByCtxf(ctx, "Join kubernetes worker node on '%s' finished.", hostDescription)
		return nil
	}

	_, err = commandExecutor.RunCommand(
		commandexecutorgeneric.WithLiveOutputOnStdoutIfVerbose(ctx),
		&parameteroptions.RunCommandOptions{
			Command: sudoCommand(options.UseSudo, strings.Fields(options.JoinCommand)...),
		},
	)
	if err != nil {
		return err
	}

	logging.LogChangedByCtxf(ctx, "Joined kubernetes worker node on '%s'.", hostDescription)
	logging.LogInfoByCtxf(ctx, "Join kubernetes worker node on '%s' finished.", hostDescription)

	return nil
}
