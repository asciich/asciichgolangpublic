package kubeadmutils

// InitControlPlaneOptions contains options for initializing a kubernetes
// control-plane node with "kubeadm init".
type InitControlPlaneOptions struct {
	// PodNetworkCidr is passed to "kubeadm init --pod-network-cidr".
	// Must match the CNI plugin installed afterwards.
	PodNetworkCidr string

	// ApiServerAdvertiseAddress is passed to
	// "kubeadm init --apiserver-advertise-address" when set.
	// Optional: kubeadm auto-detects the default route interface if empty.
	ApiServerAdvertiseAddress string

	// KubernetesVersion is passed to "kubeadm init --kubernetes-version" when set.
	// Optional: kubeadm uses its built-in default if empty.
	KubernetesVersion string

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
	// Requires ControlPlaneEndpoint to be set (the VIP KubeVip advertises).
	// Default: false (KubeVip is not activated).
	EnableKubeVip bool

	// KubeVipInterface is the network interface KubeVip binds the VIP to
	// (e.g. "eth0"). Required when EnableKubeVip is true.
	KubeVipInterface string

	// KubeVipVersion is the KubeVip image version to render the manifest with
	// (e.g. "v0.8.7").
	// Optional: a built-in default is used when empty.
	KubeVipVersion string

	// UseSudo determines if sudo should be used.
	// Default: true
	UseSudo bool

	// UploadCerts adds "--upload-certs" to "kubeadm init" so control-plane
	// certificates are uploaded to the "kubeadm-certs" secret, enabling
	// additional control-plane nodes to join.
	// Default: false
	UploadCerts bool
}

// DefaultInitControlPlaneOptions returns the default options for initializing a
// kubernetes control-plane node.
func DefaultInitControlPlaneOptions() *InitControlPlaneOptions {
	return &InitControlPlaneOptions{
		PodNetworkCidr: "10.244.0.0/16",
		UseSudo:        true,
	}
}
