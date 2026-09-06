package kubeletutils

// InstallKubeletOptions contains options for installing kubelet.
type InstallKubeletOptions struct {
	// InstallPath is the path where kubelet will be installed.
	// Default: /bin/kubelet
	InstallPath string

	// UseSudo determines if sudo should be used for installation.
	// Default: true
	UseSudo bool

	// Version is the kubelet version to install (e.g., "v1.36.2").
	// Default: "v1.36.2"
	Version string
}

// DefaultInstallKubeletOptions returns the default options for installing kubelet.
func DefaultInstallKubeletOptions() *InstallKubeletOptions {
	return &InstallKubeletOptions{
		InstallPath: "/bin/kubelet",
		UseSudo:     true,
		Version:     "v1.36.2",
	}
}
