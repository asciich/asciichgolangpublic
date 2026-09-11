package kubeadmutils

import (
	"context"
	"net"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfile"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/commandexecutornicutils"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

const (
	// kubeVipManifestPath is where the KubeVip static-pod manifest is written.
	// The kubelet picks up static pods from /etc/kubernetes/manifests.
	kubeVipManifestPath = "/etc/kubernetes/manifests/kube-vip.yaml"

	// kubeVipManifestDir is the static-pod manifest directory kubeadm/kubelet watches.
	kubeVipManifestDir = "/etc/kubernetes/manifests"

	// defaultKubeVipVersion is used when no explicit version is given.
	defaultKubeVipVersion = "v0.8.7"

	// kubeVipImage is the container image used to render the static-pod manifest.
	kubeVipImage = "ghcr.io/kube-vip/kube-vip"
)

// DeployKubeVipManifestOptions contains options for deploying the KubeVip
// static-pod manifest.
type DeployKubeVipManifestOptions struct {
	// Vip is the virtual IP KubeVip advertises (an IP, not "host:port").
	Vip string

	// Interface is the network interface KubeVip binds the VIP to (e.g. "eth0").
	Interface string

	// Version is the KubeVip image version (e.g. "v0.8.7").
	// Optional: defaultKubeVipVersion is used when empty.
	Version string

	// UseSudo determines if sudo should be used.
	UseSudo bool
}

// vipAddressFromControlPlaneEndpoint extracts the bare IP/host from a
// ControlPlaneEndpoint that may be given as "host:port".
func vipAddressFromControlPlaneEndpoint(controlPlaneEndpoint string) string {
	host, _, err := net.SplitHostPort(controlPlaneEndpoint)
	if err != nil {
		// No port present: use the value as-is.
		return controlPlaneEndpoint
	}
	return host
}

// DeployKubeVipManifest deploys the KubeVip static-pod manifest on the local machine.
func DeployKubeVipManifest(ctx context.Context, options *DeployKubeVipManifestOptions) error {
	return DeployKubeVipManifestUsingCommandExecutor(ctx, commandexecutorexecoo.Exec(), options)
}

// DeployKubeVipManifestUsingCommandExecutor renders the KubeVip static-pod
// manifest and writes it to /etc/kubernetes/manifests/kube-vip.yaml using the
// given command executor. The operation is skipped when the manifest already
// exists (idempotent).
func DeployKubeVipManifestUsingCommandExecutor(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, options *DeployKubeVipManifestOptions) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	if options.Vip == "" {
		return tracederrors.TracedErrorEmptyString("options.Vip")
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Deploy KubeVip static-pod manifest on '%s' started.", hostDescription)

	// Auto-detect the interface to bind the VIP to when not explicitly set.
	nicInterface := options.Interface
	if nicInterface == "" {
		nicInterface, err = commandexecutornicutils.GetFirstActiveNicName(ctx, commandExecutor)
		if err != nil {
			return err
		}
		logging.LogInfoByCtxf(ctx, "options.Interface not set. Auto-detected interface '%s' on '%s'.", nicInterface, hostDescription)
	}

	version := options.Version
	if version == "" {
		version = defaultKubeVipVersion
	}

	alreadyDeployed, err := commandexecutorfile.Exists(ctx, commandExecutor, kubeVipManifestPath)
	if err != nil {
		return err
	}

	if alreadyDeployed {
		logging.LogInfoByCtxf(ctx, "KubeVip manifest '%s' already exists on '%s'. Skip deployment.", kubeVipManifestPath, hostDescription)
		logging.LogInfoByCtxf(ctx, "Deploy KubeVip static-pod manifest on '%s' finished.", hostDescription)
		return nil
	}

	image := kubeVipImage + ":" + version

	// Render the manifest with the KubeVip container image and write it to the
	// static-pod directory. A shell is used so pull, run and redirection happen
	// on the target host. containerd's "ctr" is available on kubeadm nodes.
	sudoPrefix := ""
	if options.UseSudo {
		sudoPrefix = "sudo "
	}

	deployCmd := sudoPrefix + "mkdir -p " + kubeVipManifestDir + " && " +
		sudoPrefix + "ctr image pull " + image + " && " +
		sudoPrefix + "ctr run --rm --net-host " + image + " vip " +
		"/kube-vip manifest pod " +
		"--interface " + nicInterface + " " +
		"--address " + options.Vip + " " +
		"--controlplane --services --arp --leaderElection " +
		"| " + sudoPrefix + "tee " + kubeVipManifestPath + " >/dev/null"

	_, err = commandExecutor.RunCommand(
		commandexecutorgeneric.WithLiveOutputOnStdoutIfVerbose(ctx),
		&parameteroptions.RunCommandOptions{
			Command: []string{"sh", "-c", deployCmd},
		},
	)
	if err != nil {
		return err
	}

	logging.LogChangedByCtxf(ctx, "Deployed KubeVip static-pod manifest '%s' (vip='%s', interface='%s', version='%s') on '%s'.", kubeVipManifestPath, options.Vip, nicInterface, version, hostDescription)
	logging.LogInfoByCtxf(ctx, "Deploy KubeVip static-pod manifest on '%s' finished.", hostDescription)

	return nil
}
