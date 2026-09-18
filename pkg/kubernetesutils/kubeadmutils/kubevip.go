package kubeadmutils

import (
	"context"
	"encoding/base64"
	"net"
	"strings"

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

	// containerdNamespace is the containerd namespace used by the kubelet/CRI.
	// The kube-vip image must live here and "ctr" must operate in it, otherwise
	// the image pulled/used does not match what the kubelet can see.
	containerdNamespace = "k8s.io"
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

	sudoPrefix := ""
	if options.UseSudo {
		sudoPrefix = "sudo "
	}

	// Step 1: Pull the image into the SAME containerd namespace the kubelet uses.
	//
	// NOTE: The previous implementation used the default containerd namespace,
	// which is not the namespace the CRI/kubelet operates in.
	pullCmd := sudoPrefix + "ctr -n " + containerdNamespace + " image pull " + image

	_, err = commandExecutor.RunCommand(
		commandexecutorgeneric.WithLiveOutputOnStdoutIfVerbose(ctx),
		&parameteroptions.RunCommandOptions{
			Command: []string{"sh", "-c", pullCmd},
		},
	)
	if err != nil {
		return err
	}

	// Step 2: Render the manifest.
	//
	// IMPORTANT: On containerd 2.x, "ctr run" parses flags that appear AFTER the
	// container command (e.g. "--interface", "--controlplane", "--arp") as its
	// OWN flags. It then fails, prints its usage text to stdout and exits. The
	// old code piped that stdout directly into "tee", which happily wrote the
	// usage text into kube-vip.yaml (9 KB of "ctr run" help) while the failed
	// exit code was swallowed by the pipe. The kubelet then ignored the invalid
	// manifest and the VIP was never created -> "no route to host".
	//
	// To avoid flag intermixing entirely we configure kube-vip through
	// environment variables (which "ctr" passes through untouched via --env)
	// and call "/kube-vip manifest pod" WITHOUT any trailing flags.
	generateCmd := sudoPrefix + "ctr -n " + containerdNamespace + " run --rm --net-host " +
		"--env vip_interface=" + nicInterface + " " +
		"--env address=" + options.Vip + " " +
		"--env cp_enable=true " +
		"--env svc_enable=true " +
		"--env vip_arp=true " +
		"--env vip_leaderelection=true " +
		"--env port=6443 " +
		"--env cp_namespace=kube-system " +
		image + " kube-vip-generate /kube-vip manifest pod"

	generateOutput, err := commandExecutor.RunCommand(
		commandexecutorgeneric.WithLiveOutputOnStdoutIfVerbose(ctx),
		&parameteroptions.RunCommandOptions{
			Command: []string{"sh", "-c", generateCmd},
		},
	)
	if err != nil {
		return err
	}

	manifest, err := generateOutput.GetStdoutAsString()
	if err != nil {
		return err
	}

	// Step 3: Validate the rendered manifest BEFORE writing it. A valid
	// kube-vip static-pod manifest is YAML starting with "apiVersion:".
	// This guards against ever writing garbage (e.g. "ctr" usage output)
	// into the static-pod directory again.
	trimmedManifest := strings.TrimSpace(manifest)
	if !strings.HasPrefix(trimmedManifest, "apiVersion:") {
		return tracederrors.TracedErrorf(
			"Rendering the KubeVip manifest on '%s' did not produce a valid Kubernetes manifest. "+
				"Expected output to start with 'apiVersion:' but got:\n%s",
			hostDescription,
			trimmedManifest,
		)
	}

	// Step 4: Write the validated manifest atomically. base64 is used to avoid
	// any shell quoting/escaping issues with the YAML content and to keep the
	// non-zero exit codes intact (no lossy pipe like the old "| tee").
	encodedManifest := base64.StdEncoding.EncodeToString([]byte(trimmedManifest + "\n"))
	writeCmd := sudoPrefix + "mkdir -p " + kubeVipManifestDir + " && " +
		"printf '%s' '" + encodedManifest + "' | base64 -d | " +
		sudoPrefix + "tee " + kubeVipManifestPath + " >/dev/null"

	_, err = commandExecutor.RunCommand(
		commandexecutorgeneric.WithLiveOutputOnStdoutIfVerbose(ctx),
		&parameteroptions.RunCommandOptions{
			Command: []string{"sh", "-c", writeCmd},
		},
	)
	if err != nil {
		return err
	}

	logging.LogChangedByCtxf(ctx, "Deployed KubeVip static-pod manifest '%s' (vip='%s', interface='%s', version='%s') on '%s'.", kubeVipManifestPath, options.Vip, nicInterface, version, hostDescription)
	logging.LogInfoByCtxf(ctx, "Deploy KubeVip static-pod manifest on '%s' finished.", hostDescription)

	return nil
}
