package kubeadmutils

// JoinNodeOptions contains options for joining a worker node with
// "kubeadm join".
type JoinNodeOptions struct {
	// JoinCommand is the full "kubeadm join ..." command to run.
	//
	// In contrast to JoinControlPlaneOptions there is intentionally no
	// ControlPlaneEndpoint / KubeVip configuration here: a worker node neither
	// serves the VIP nor uploads control-plane certificates. The API endpoint it
	// talks to is already encoded in the join command.
	JoinCommand string

	// UseSudo determines if sudo should be used.
	// Default: true
	UseSudo bool
}
