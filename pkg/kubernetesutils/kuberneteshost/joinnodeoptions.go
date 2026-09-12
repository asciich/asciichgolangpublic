package kuberneteshost

// JoinNodeOptions contains the options used to join a worker node to an
// existing cluster.
type JoinNodeOptions struct {
	// JoinCommand is the "kubeadm join ..." command obtained from
	// GetNodeJoinCommand.
	//
	// In contrast to JoinControlPlaneOptions there is intentionally no
	// ControlPlaneEndpoint / KubeVip configuration here: a worker node neither
	// serves the VIP nor uploads control-plane certificates. The API endpoint
	// it talks to is already encoded in the join command.
	JoinCommand string
}
