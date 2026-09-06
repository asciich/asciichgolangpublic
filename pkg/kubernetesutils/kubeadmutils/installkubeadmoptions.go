package kubeadmutils

// InstallKubeadmOptions contains options for installing kubeadm.
type InstallKubeadmOptions struct {
	// InstallPath is the path where kubeadm will be installed.
	// Default: /bin/kubeadm
	InstallPath string

	// UseSudo determines if sudo should be used for installation.
	// Default: true
	UseSudo bool

	// Version is the kubeadm version to install (e.g., "v1.36.2").
	// Default: "v1.36.2"
	Version string
}

// DefaultInstallKubeadmOptions returns the default options for installing kubeadm.
func DefaultInstallKubeadmOptions() *InstallKubeadmOptions {
	return &InstallKubeadmOptions{
		InstallPath: "/bin/kubeadm",
		UseSudo:     true,
		Version:     "v1.36.2",
	}
}
