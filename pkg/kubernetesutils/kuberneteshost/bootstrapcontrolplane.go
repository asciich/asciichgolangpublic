package kuberneteshost

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/ciliumutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubeadmutils"
	"github.com/asciich/asciichgolangpublic/pkg/runbook"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// BootstrapControlPlaneOptions contains the options used to bootstrap a
// kubernetes control-plane node with kubeadm and Cilium (eBPF) as CNI.
type BootstrapControlPlaneOptions struct {
	// PodNetworkCidr is passed to "kubeadm init --pod-network-cidr" when set.
	// Cilium manages its own IPAM (cluster-pool), so this can be left empty to
	// avoid a mismatch between kubeadm's podCIDR and Cilium's pool. When set,
	// Cilium is configured to honor it (ipam.mode=kubernetes).
	// Default: "" (Cilium owns IPAM via its cluster-pool default).
	PodNetworkCidr string

	// ApiServerAdvertiseAddress is passed to
	// "kubeadm init --apiserver-advertise-address" when set.
	// Optional: kubeadm auto-detects the default route interface if empty.
	ApiServerAdvertiseAddress string

	// KubernetesVersion is passed to "kubeadm init --kubernetes-version" when set.
	// Optional: kubeadm uses its built-in default if empty.
	KubernetesVersion string

	// CiliumVersion is the Cilium version to install (e.g. "1.16.5").
	// Optional: the Cilium CLI picks a compatible default if empty.
	CiliumVersion string

	// UntaintControlPlane removes the control-plane NoSchedule taint so
	// workloads can be scheduled on this node (useful for single-node clusters).
	// Default: false
	UntaintControlPlane bool

	// ControlPlaneEndpoint is passed to
	// "kubeadm init --control-plane-endpoint" when set.
	// This is typically the shared/virtual IP (VIP) provided by KubeVip so the
	// API server is reachable through a stable, highly-available address instead
	// of a single node's IP. Can be an IP or a "host:port" value.
	// Optional: leave empty for single-endpoint (non-HA) setups.
	// Default: "" (kubeadm uses the advertise address directly).
	ControlPlaneEndpoint string

	// EnableKubeVip explicitly activates KubeVip for this control-plane node.
	// When true, a KubeVip static-pod manifest is deployed into
	// "/etc/kubernetes/manifests" before "kubeadm init" runs so the virtual IP
	// (see ControlPlaneEndpoint) is served and reachable while the control-plane
	// bootstraps.
	// Requires ControlPlaneEndpoint and KubeVipInterface to be set.
	// Default: false (KubeVip is not activated).
	EnableKubeVip bool

	// KubeVipInterface is the network interface KubeVip binds the VIP to
	// (e.g. "eth0"). Required when EnableKubeVip is true.
	KubeVipInterface string

	// KubeVipVersion is the KubeVip image version to render the manifest with
	// (e.g. "v0.8.7").
	// Optional: a built-in default is used when empty.
	KubeVipVersion string
}

// DefaultBootstrapControlPlaneOptions returns the default options for
// bootstrapping a kubernetes control-plane node with Cilium.
func DefaultBootstrapControlPlaneOptions() *BootstrapControlPlaneOptions {
	return &BootstrapControlPlaneOptions{
		// Empty on purpose: let Cilium's cluster-pool manage pod IPAM.
		PodNetworkCidr:      "",
		UntaintControlPlane: false,
	}
}

func prepareBootstrapControlPlaneOptions(options *BootstrapControlPlaneOptions) *BootstrapControlPlaneOptions {
	if options == nil {
		options = DefaultBootstrapControlPlaneOptions()
	}
	// PodNetworkCidr intentionally has no default (Cilium owns IPAM).
	return options
}

func NewBootstrapControlPlaneRunbook(commandExecutor commandexecutorinterfaces.CommandExecutor, options *BootstrapControlPlaneOptions) (*runbook.RunBook, error) {
	if commandExecutor == nil {
		return nil, tracederrors.TracedErrorNil("commandExecutor")
	}

	options = prepareBootstrapControlPlaneOptions(options)

	runbook := &runbook.RunBook{
		Name:        "Bootstrap kubernetes control-plane",
		Description: "Initialize a kubernetes control-plane node with kubeadm, set up the admin kubeconfig and install Cilium (eBPF) as CNI.",
		Steps: []runbook.Runnable{
			&runbook.Step{
				Name:        "kubeadm-init",
				Description: "Initialize the control-plane with 'kubeadm init'. Skipped if the node is already initialized.",
				Run: func(ctx context.Context) error {
					return kubeadmutils.InitControlPlaneUsingCommandExecutor(ctx, commandExecutor, &kubeadmutils.InitControlPlaneOptions{
						PodNetworkCidr:            options.PodNetworkCidr,
						ApiServerAdvertiseAddress: options.ApiServerAdvertiseAddress,
						KubernetesVersion:         options.KubernetesVersion,
						ControlPlaneEndpoint:      options.ControlPlaneEndpoint,
						EnableKubeVip:             options.EnableKubeVip,
						KubeVipInterface:          options.KubeVipInterface,
						KubeVipVersion:            options.KubeVipVersion,
					})
				},
			},
			&runbook.Step{
				Name:        "setup-admin-kubeconfig",
				Description: "Copy /etc/kubernetes/admin.conf to the invoking user's ~/.kube/config so kubectl can talk to the cluster.",
				Run: func(ctx context.Context) error {
					return kubeadmutils.SetupAdminKubeconfigUsingCommandExecutor(ctx, commandExecutor)
				},
			},
			&runbook.Step{
				Name:        "install-cilium",
				Description: "Install Cilium (eBPF) as the CNI network plugin so pods can communicate and nodes become Ready.",
				Run: func(ctx context.Context) error {
					return ciliumutils.InstallUsingCommandExecutor(ctx, commandExecutor, &ciliumutils.InstallCiliumOptions{
						Version:        options.CiliumVersion,
						PodNetworkCidr: options.PodNetworkCidr,
						ExtraHelmSet: []string{
							"socketLB.enabled=false",
							"bpf.hostLegacyRouting=true",
						},
					})
				},
			},
			&runbook.Step{
				Name:        "wait-cilium-ready",
				Description: "Wait until Cilium reports healthy so nodes become Ready.",
				Run: func(ctx context.Context) error {
					return ciliumutils.WaitUntilReadyUsingCommandExecutor(ctx, commandExecutor, &ciliumutils.InstallCiliumOptions{})
				},
			},
		},
	}

	return runbook, nil
}

func BootstrapControlPlane(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, options *BootstrapControlPlaneOptions) error {
	runbook, err := NewBootstrapControlPlaneRunbook(commandExecutor, options)
	if err != nil {
		return err
	}

	return runbook.Execute(ctx)
}
