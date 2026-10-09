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
	"github.com/asciich/asciichgolangpublic/pkg/netutils/iputils"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

const (
	// kubeVipManifestPath is where the KubeVip static-pod manifest is written.
	// The kubelet picks up static pods from /etc/kubernetes/manifests.
	kubeVipManifestPath = "/etc/kubernetes/manifests/kube-vip.yaml"

	// kubeVipManifestDir is the static-pod manifest directory kubeadm/kubelet watches.
	kubeVipManifestDir = "/etc/kubernetes/manifests"

	// kubeVipManifestTmpPath is used to write the manifest before moving it into
	// place. The kubelet ignores hidden files (starting with '.') in the
	// static-pod directory, so a half written temporary file is never picked up
	// as a static pod.
	kubeVipManifestTmpPath = kubeVipManifestDir + "/.kube-vip.yaml.tmp"

	// defaultKubeVipVersion is used when no explicit version is given.
	defaultKubeVipVersion = "v0.8.7"

	// kubeVipImage is the container image used to render the static-pod manifest.
	kubeVipImage = "ghcr.io/kube-vip/kube-vip"

	// containerdNamespace is the containerd namespace used by the kubelet/CRI.
	// The kube-vip image must live here and "ctr" must operate in it, otherwise
	// the image pulled/used does not match what the kubelet can see.
	containerdNamespace = "k8s.io"
)

// DeployKubeVipManifestOptions contains options for deploying/ensuring the
// KubeVip static-pod manifest.
type DeployKubeVipManifestOptions struct {
	// Vip is the virtual IP KubeVip advertises (an IP, not "host:port").
	Vip string

	// Interface is the network interface KubeVip binds the VIP to (e.g. "eth0").
	// Optional: the first active NIC is auto-detected when empty.
	Interface string

	// Version is the KubeVip image version (e.g. "v0.8.7").
	// Optional: defaultKubeVipVersion is used when empty.
	Version string

	// EnableServices lets KubeVip announce LoadBalancer Service IPs
	// ("svc_enable"). Keep this false when another load balancer (e.g. MetalLB)
	// handles Services. Otherwise both announce the same IPs and KubeVip claims
	// them on control planes without local endpoints, which breaks Services
	// using "externalTrafficPolicy: Local".
	// Default: false
	EnableServices bool

	// SkipLeftoverIpCleanup disables removing leftover "/32" addresses from the
	// KubeVip interface. When "svc_enable" is switched off, KubeVip does not
	// remove the LoadBalancer IPs it added before, so the kernel keeps
	// answering ARP for them.
	// Only used by EnsureKubeVipManifest...: by default every "/32" address on
	// the interface except Vip is removed when EnableServices is false.
	// Default: false (leftover IPs are removed).
	SkipLeftoverIpCleanup bool

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

// sudoShellPrefix returns "sudo " when useSudo is true. Used for commands
// passed to "sh -c".
func sudoShellPrefix(useSudo bool) string {
	if useSudo {
		return "sudo "
	}
	return ""
}

// validateKubeVipManifestOptions validates the input parameters shared by the
// deploy and ensure functions.
func validateKubeVipManifestOptions(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, options *DeployKubeVipManifestOptions) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	if options.Vip == "" {
		return tracederrors.TracedErrorEmptyString("options.Vip")
	}

	err := iputils.CheckValidIP(ctx, options.Vip)
	if err != nil {
		return err
	}

	return nil
}

// checkRootPrivilegesOnTarget ensures the commands on the target are executed
// with root privileges. When useSudo is set, the privileges are provided by
// sudo. Otherwise "id -u" is evaluated on the target machine, since the
// privileges of the local user are irrelevant for remote execution.
func checkRootPrivilegesOnTarget(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, useSudo bool, operation string, hostDescription string) error {
	if useSudo {
		return nil
	}

	stdout, err := commandExecutor.RunCommandAndGetStdoutAsString(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{"id", "-u"},
		},
	)
	if err != nil {
		return err
	}

	uid := strings.TrimSpace(stdout)
	if uid == "" {
		return tracederrors.TracedErrorf("'id -u' returned an empty user id on '%s'.", hostDescription)
	}

	if uid != "0" {
		return tracederrors.TracedErrorf("%s on '%s' requires root privileges. Run as root or set UseSudo.", operation, hostDescription)
	}

	return nil
}

// resolveKubeVipSettings returns the interface and version used to render the
// manifest. The interface is auto-detected when not explicitly set.
func resolveKubeVipSettings(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, options *DeployKubeVipManifestOptions, hostDescription string) (string, string, error) {
	nicInterface := options.Interface
	if nicInterface == "" {
		var err error
		nicInterface, err = commandexecutornicutils.GetFirstActiveNicName(ctx, commandExecutor)
		if err != nil {
			return "", "", err
		}
		logging.LogInfoByCtxf(ctx, "options.Interface not set. Auto-detected interface '%s' on '%s'.", nicInterface, hostDescription)
	}

	version := options.Version
	if version == "" {
		version = defaultKubeVipVersion
	}

	return nicInterface, version, nil
}

// renderKubeVipManifest pulls the KubeVip image and renders the static-pod
// manifest on the target. Deploy and Ensure both use this function so the
// rendered manifest is identical for the same options.
func renderKubeVipManifest(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, options *DeployKubeVipManifestOptions, nicInterface string, version string, hostDescription string) (string, error) {
	image := kubeVipImage + ":" + version
	sudoPrefix := sudoShellPrefix(options.UseSudo)

	// Step 1: Pull the image into the SAME containerd namespace the kubelet uses.
	//
	// NOTE: A previous implementation used the default containerd namespace,
	// which is not the namespace the CRI/kubelet operates in.
	pullCmd := sudoPrefix + "ctr -n " + containerdNamespace + " image pull " + image

	_, err := commandExecutor.RunCommand(
		commandexecutorgeneric.WithLiveOutputOnStdoutIfVerbose(ctx),
		&parameteroptions.RunCommandOptions{
			Command: []string{"sh", "-c", pullCmd},
		},
	)
	if err != nil {
		return "", err
	}

	// Step 2: Render the manifest.
	//
	// IMPORTANT: On containerd 2.x, "ctr run" parses flags that appear AFTER the
	// container command (e.g. "--interface", "--controlplane", "--arp") as its
	// OWN flags. It then fails, prints its usage text to stdout and exits. An
	// old implementation piped that stdout directly into "tee", which wrote the
	// usage text into kube-vip.yaml while the failed exit code was swallowed by
	// the pipe. The kubelet then ignored the invalid manifest and the VIP was
	// never created -> "no route to host".
	//
	// To avoid flag intermixing entirely we configure kube-vip through
	// environment variables (which "ctr" passes through untouched via --env)
	// and call "/kube-vip manifest pod" WITHOUT any trailing flags.
	svcEnable := "false"
	if options.EnableServices {
		svcEnable = "true"
	}

	generateCmd := sudoPrefix + "ctr -n " + containerdNamespace + " run --rm --net-host " +
		"--env vip_interface=" + nicInterface + " " +
		"--env address=" + options.Vip + " " +
		"--env cp_enable=true " +
		"--env svc_enable=" + svcEnable + " " +
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
		return "", err
	}

	manifest, err := generateOutput.GetStdoutAsString()
	if err != nil {
		return "", err
	}

	// Step 3: Validate the rendered manifest BEFORE it is used. A valid
	// kube-vip static-pod manifest is YAML starting with "apiVersion:".
	// This guards against ever writing garbage (e.g. "ctr" usage output)
	// into the static-pod directory again.
	trimmedManifest := strings.TrimSpace(manifest)
	if !strings.HasPrefix(trimmedManifest, "apiVersion:") {
		return "", tracederrors.TracedErrorf(
			"Rendering the KubeVip manifest on '%s' did not produce a valid Kubernetes manifest. "+
				"Expected output to start with 'apiVersion:' but got:\n%s",
			hostDescription,
			trimmedManifest,
		)
	}

	ret := trimmedManifest + "\n"

	return ret, nil
}

// readKubeVipManifest returns the content of the currently deployed manifest.
// The returned bool is false if no manifest exists.
func readKubeVipManifest(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, useSudo bool) (string, bool, error) {
	exists, err := commandexecutorfile.Exists(ctx, commandExecutor, kubeVipManifestPath)
	if err != nil {
		return "", false, err
	}

	if !exists {
		return "", false, nil
	}

	ret, err := commandExecutor.RunCommandAndGetStdoutAsString(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: sudoCommand(useSudo, "cat", kubeVipManifestPath),
		},
	)
	if err != nil {
		return "", false, err
	}

	return ret, true, nil
}

// writeKubeVipManifest writes the manifest atomically: the content is written
// to a hidden temporary file, read back and verified, and then moved into
// place. This way the kubelet never reads a half written or corrupted manifest.
func writeKubeVipManifest(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, useSudo bool, manifest string, hostDescription string) error {
	if manifest == "" {
		return tracederrors.TracedErrorEmptyString("manifest")
	}

	sudoPrefix := sudoShellPrefix(useSudo)

	// Idempotent cleanup: after a successful "mv" the temporary file is already
	// gone and "rm -f" is a no-op.
	defer func() {
		_, err := commandExecutor.RunCommand(
			ctx,
			&parameteroptions.RunCommandOptions{
				Command: sudoCommand(useSudo, "rm", "-f", kubeVipManifestTmpPath),
			},
		)
		if err != nil {
			logging.LogInfoByCtxf(ctx, "Unable to remove temporary file '%s' on '%s', it can be removed manually: %s", kubeVipManifestTmpPath, hostDescription, err)
		}
	}()

	// base64 is used to avoid any shell quoting/escaping issues with the YAML
	// content. base64 only contains [A-Za-z0-9+/=], so single quotes are safe.
	encodedManifest := base64.StdEncoding.EncodeToString([]byte(manifest))
	writeCmd := sudoPrefix + "mkdir -p " + kubeVipManifestDir + " && " +
		"printf '%s' '" + encodedManifest + "' | base64 -d | " +
		sudoPrefix + "tee " + kubeVipManifestTmpPath + " >/dev/null"

	_, err := commandExecutor.RunCommand(
		commandexecutorgeneric.WithLiveOutputOnStdoutIfVerbose(ctx),
		&parameteroptions.RunCommandOptions{
			Command: []string{"sh", "-c", writeCmd},
		},
	)
	if err != nil {
		return err
	}

	// The pipe above only reports the exit code of "tee". Read the temporary
	// file back so a failed "base64 -d" can never end up as manifest.
	written, err := commandExecutor.RunCommandAndGetStdoutAsString(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: sudoCommand(useSudo, "cat", kubeVipManifestTmpPath),
		},
	)
	if err != nil {
		return err
	}

	if strings.TrimSpace(written) != strings.TrimSpace(manifest) {
		return tracederrors.TracedErrorf("Verification of the written KubeVip manifest '%s' on '%s' failed: content differs from the rendered manifest.", kubeVipManifestTmpPath, hostDescription)
	}

	_, err = commandExecutor.RunCommand(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: sudoCommand(useSudo, "chmod", "0600", kubeVipManifestTmpPath),
		},
	)
	if err != nil {
		return err
	}

	_, err = commandExecutor.RunCommand(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: sudoCommand(useSudo, "mv", "-f", kubeVipManifestTmpPath, kubeVipManifestPath),
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// removeLeftoverKubeVipServiceIps removes every "/32" IPv4 address except the
// control-plane VIP from the given interface. KubeVip with "svc_enable" adds
// LoadBalancer IPs as "/32" to the interface but does not remove them when
// Service handling is switched off. As long as they exist the kernel answers
// ARP for them and the real load balancer (e.g. MetalLB) cannot take over.
func removeLeftoverKubeVipServiceIps(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, useSudo bool, nicInterface string, vip string, hostDescription string) error {
	if nicInterface == "" {
		return tracederrors.TracedErrorEmptyString("nicInterface")
	}

	if vip == "" {
		return tracederrors.TracedErrorEmptyString("vip")
	}

	logging.LogInfoByCtxf(ctx, "Remove leftover KubeVip service IPs from interface '%s' on '%s' started.", nicInterface, hostDescription)

	stdout, err := commandExecutor.RunCommandAndGetStdoutAsString(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{"ip", "-o", "-4", "addr", "show", "dev", nicInterface},
		},
	)
	if err != nil {
		return err
	}

	removed := 0

	// Line format: "2: enp1s0    inet 192.168.10.233/32 scope global enp1s0 ..."
	for _, line := range strings.Split(stdout, "\n") {
		fields := strings.Fields(line)

		for i := 0; i+1 < len(fields); i++ {
			if fields[i] != "inet" {
				continue
			}

			ip, prefixLength, found := strings.Cut(fields[i+1], "/")
			if !found || prefixLength != "32" {
				break
			}

			if ip == vip {
				break
			}

			isValid, err := iputils.IsValidIP(ctx, ip)
			if err != nil {
				return err
			}

			if !isValid {
				return tracederrors.TracedErrorf("Unexpected address '%s' found on interface '%s' on '%s'.", fields[i+1], nicInterface, hostDescription)
			}

			_, err = commandExecutor.RunCommand(
				ctx,
				&parameteroptions.RunCommandOptions{
					Command: sudoCommand(useSudo, "ip", "addr", "del", ip+"/32", "dev", nicInterface),
				},
			)
			if err != nil {
				return err
			}

			logging.LogChangedByCtxf(ctx, "Removed leftover KubeVip service IP '%s/32' from interface '%s' on '%s'.", ip, nicInterface, hostDescription)
			removed++

			break
		}
	}

	if removed == 0 {
		logging.LogInfoByCtxf(ctx, "No leftover KubeVip service IPs on interface '%s' on '%s'.", nicInterface, hostDescription)
	}

	logging.LogInfoByCtxf(ctx, "Remove leftover KubeVip service IPs from interface '%s' on '%s' finished. Removed '%d' addresses.", nicInterface, hostDescription, removed)

	return nil
}

// DeployKubeVipManifest deploys the KubeVip static-pod manifest on the local machine.
func DeployKubeVipManifest(ctx context.Context, options *DeployKubeVipManifestOptions) error {
	return DeployKubeVipManifestUsingCommandExecutor(ctx, commandexecutorexecoo.Exec(), options)
}

// DeployKubeVipManifestUsingCommandExecutor renders the KubeVip static-pod
// manifest and writes it to /etc/kubernetes/manifests/kube-vip.yaml using the
// given command executor. The operation is skipped when the manifest already
// exists (idempotent). Use EnsureKubeVipManifestUsingCommandExecutor to update
// an already existing manifest.
//
// Requires root privileges on the target (or UseSudo).
func DeployKubeVipManifestUsingCommandExecutor(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, options *DeployKubeVipManifestOptions) error {
	err := validateKubeVipManifestOptions(ctx, commandExecutor, options)
	if err != nil {
		return err
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return err
	}

	err = checkRootPrivilegesOnTarget(ctx, commandExecutor, options.UseSudo, "Deploy KubeVip static-pod manifest", hostDescription)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Deploy KubeVip static-pod manifest on '%s' started.", hostDescription)

	alreadyDeployed, err := commandexecutorfile.Exists(ctx, commandExecutor, kubeVipManifestPath)
	if err != nil {
		return err
	}

	if alreadyDeployed {
		logging.LogInfoByCtxf(ctx, "KubeVip manifest '%s' already exists on '%s'. Skip deployment.", kubeVipManifestPath, hostDescription)
		logging.LogInfoByCtxf(ctx, "Deploy KubeVip static-pod manifest on '%s' finished.", hostDescription)
		return nil
	}

	nicInterface, version, err := resolveKubeVipSettings(ctx, commandExecutor, options, hostDescription)
	if err != nil {
		return err
	}

	manifest, err := renderKubeVipManifest(ctx, commandExecutor, options, nicInterface, version, hostDescription)
	if err != nil {
		return err
	}

	err = writeKubeVipManifest(ctx, commandExecutor, options.UseSudo, manifest, hostDescription)
	if err != nil {
		return err
	}

	logging.LogChangedByCtxf(ctx, "Deployed KubeVip static-pod manifest '%s' (vip='%s', interface='%s', version='%s', services='%t') on '%s'.", kubeVipManifestPath, options.Vip, nicInterface, version, options.EnableServices, hostDescription)
	logging.LogInfoByCtxf(ctx, "Deploy KubeVip static-pod manifest on '%s' finished.", hostDescription)

	return nil
}

// EnsureKubeVipManifest ensures the KubeVip static-pod manifest on the local
// machine matches the given options. Returns true if the manifest was changed.
func EnsureKubeVipManifest(ctx context.Context, options *DeployKubeVipManifestOptions) (bool, error) {
	return EnsureKubeVipManifestUsingCommandExecutor(ctx, commandexecutorexecoo.Exec(), options)
}

// EnsureKubeVipManifestUsingCommandExecutor ensures the KubeVip static-pod
// manifest matches the given options. In contrast to
// DeployKubeVipManifestUsingCommandExecutor an existing manifest is compared
// with the rendered one and replaced if it differs (idempotent).
//
// Returns true if the manifest was changed. In this case the kubelet restarts
// KubeVip, so callers updating multiple control planes should wait until the
// API is reachable through the VIP again before continuing with the next node.
//
// If EnableServices is false and SkipLeftoverIpCleanup is not set, leftover
// "/32" addresses (except Vip) are removed from the KubeVip interface.
//
// Requires root privileges on the target (or UseSudo).
func EnsureKubeVipManifestUsingCommandExecutor(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, options *DeployKubeVipManifestOptions) (bool, error) {
	err := validateKubeVipManifestOptions(ctx, commandExecutor, options)
	if err != nil {
		return false, err
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return false, err
	}

	err = checkRootPrivilegesOnTarget(ctx, commandExecutor, options.UseSudo, "Ensure KubeVip static-pod manifest", hostDescription)
	if err != nil {
		return false, err
	}

	logging.LogInfoByCtxf(ctx, "Ensure KubeVip static-pod manifest on '%s' started.", hostDescription)

	nicInterface, version, err := resolveKubeVipSettings(ctx, commandExecutor, options, hostDescription)
	if err != nil {
		return false, err
	}

	desiredManifest, err := renderKubeVipManifest(ctx, commandExecutor, options, nicInterface, version, hostDescription)
	if err != nil {
		return false, err
	}

	currentManifest, exists, err := readKubeVipManifest(ctx, commandExecutor, options.UseSudo)
	if err != nil {
		return false, err
	}

	// Compare trimmed content since executors may strip the trailing newline of stdout.
	upToDate := exists && strings.TrimSpace(currentManifest) == strings.TrimSpace(desiredManifest)

	changed := false
	if upToDate {
		logging.LogInfoByCtxf(ctx, "KubeVip manifest '%s' on '%s' is already up to date.", kubeVipManifestPath, hostDescription)
	} else {
		err = writeKubeVipManifest(ctx, commandExecutor, options.UseSudo, desiredManifest, hostDescription)
		if err != nil {
			return false, err
		}

		changed = true
		logging.LogChangedByCtxf(ctx, "Updated KubeVip static-pod manifest '%s' (vip='%s', interface='%s', version='%s', services='%t') on '%s'.", kubeVipManifestPath, options.Vip, nicInterface, version, options.EnableServices, hostDescription)
	}

	if !options.EnableServices && !options.SkipLeftoverIpCleanup {
		err = removeLeftoverKubeVipServiceIps(ctx, commandExecutor, options.UseSudo, nicInterface, options.Vip, hostDescription)
		if err != nil {
			return false, err
		}
	}

	logging.LogInfoByCtxf(ctx, "Ensure KubeVip static-pod manifest on '%s' finished.", hostDescription)

	return changed, nil
}
