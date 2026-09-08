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

	// UseSudo determines if sudo should be used.
	// Default: true
	UseSudo bool
}

// DefaultInitControlPlaneOptions returns the default options for initializing a
// kubernetes control-plane node.
func DefaultInitControlPlaneOptions() *InitControlPlaneOptions {
	return &InitControlPlaneOptions{
		PodNetworkCidr: "10.244.0.0/16",
		UseSudo:        true,
	}
}
