package kubernetestestsshserver

// StartTestSshServerOptions defines the options for starting a test SSH server in a Kubernetes cluster.
type StartTestSshServerOptions struct {
	// KubernetesNamespace is the namespace where the SSH server pod will be deployed (required).
	KubernetesNamespace string

	// PodName is the name of the SSH server pod (required).
	PodName string

	// SSHUsername is the username for SSH authentication (required).
	SSHUsername string

	// SSHPassword is the password for SSH authentication (optional).
	// If not provided, a random password is generated.
	SSHPassword string

	// SSHPublicKey is the SSH public key for key-based authentication (optional).
	SSHPublicKey string

	// Image is the container image to use for the SSH server (optional).
	// Defaults to a standard SSH server image if not provided.
	Image string

	// SSHPort is the port number for SSH server (optional).
	// Defaults to 22 if not provided.
	SSHPort int
}
