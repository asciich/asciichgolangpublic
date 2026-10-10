package kuberneteshost

import (
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubeadmutils"
)

// EnsureKubeVipManifestOptions contains the options used to ensure the KubeVip
// static-pod manifest in "/etc/kubernetes/manifests" of an already initialized
// control-plane node matches the desired configuration.
//
// The values must match the ones used during bootstrap/join
// (see BootstrapControlPlaneOptions and JoinControlPlaneOptions). Otherwise the
// rendered manifest differs and KubeVip is restarted on every run.
type EnsureKubeVipManifestOptions struct {
	// Vip is the control-plane virtual IP KubeVip advertises
	// (e.g. "192.168.10.175").
	// Required.
	Vip string

	// Interface is the network interface KubeVip binds the VIP to
	// (e.g. "enp1s0").
	// Optional: must be set to the same value as during bootstrap. Leave empty
	// if the bootstrap did not set it either.
	Interface string

	// Version is the KubeVip image version to render the manifest with
	// (e.g. "v0.8.7").
	// Optional: the same built-in default as during bootstrap is used when empty.
	Version string

	// EnableServices lets KubeVip announce LoadBalancer Service IPs
	// ("svc_enable"). Keep this false when another load balancer (e.g. MetalLB)
	// handles Services. Otherwise both announce the same IPs, and KubeVip
	// claims them on control planes without local endpoints, which breaks
	// Services using "externalTrafficPolicy: Local".
	// Default: false
	EnableServices bool

	// SkipLeftoverIpCleanup disables removing leftover "/32" addresses from
	// Interface. When "svc_enable" is switched off, KubeVip does not remove the
	// LoadBalancer IPs it added before, so the kernel keeps answering ARP for
	// them. By default every "/32" address on Interface except Vip is removed
	// once EnableServices is false.
	// Default: false (leftover IPs are removed).
	SkipLeftoverIpCleanup bool

	// UseSudo runs the commands on the target using "sudo".
	// Default: false
	UseSudo bool
}

// DefaultEnsureKubeVipManifestOptions returns the default options for
// ensuring the KubeVip manifest.
func DefaultEnsureKubeVipManifestOptions() *EnsureKubeVipManifestOptions {
	return &EnsureKubeVipManifestOptions{
		EnableServices:        false,
		SkipLeftoverIpCleanup: false,
	}
}

func prepareEnsureKubeVipManifestOptions(options *EnsureKubeVipManifestOptions) *EnsureKubeVipManifestOptions {
	if options == nil {
		options = DefaultEnsureKubeVipManifestOptions()
	}

	return options
}

// toDeployKubeVipManifestOptions converts the options to the
// kubeadmutils.DeployKubeVipManifestOptions, so the bootstrap and the ensure
// step render the manifest from the same values.
func (o *EnsureKubeVipManifestOptions) toDeployKubeVipManifestOptions() *kubeadmutils.DeployKubeVipManifestOptions {
	return &kubeadmutils.DeployKubeVipManifestOptions{
		Vip:            o.Vip,
		Interface:      o.Interface,
		Version:        o.Version,
		EnableServices: o.EnableServices,
		UseSudo:        o.UseSudo,
	}
}
