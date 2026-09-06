package kubeletutils

// InstallKubeletServiceOptions contains options for installing the kubelet
// systemd service unit and the kubeadm drop-in configuration.
type InstallKubeletServiceOptions struct {
	// ReleaseVersion is the kubernetes/release version providing the
	// kubelet.service systemd unit template (e.g., "v0.16.2").
	// Default: "v0.16.2"
	ReleaseVersion string

	// BinPath is the directory where the kubelet binary was installed.
	// Default: /bin
	BinPath string

	// UseSudo determines if sudo should be used for installation.
	// Default: true
	UseSudo bool
}

// DefaultInstallKubeletServiceOptions returns the default options for
// installing the kubelet systemd service.
func DefaultInstallKubeletServiceOptions() *InstallKubeletServiceOptions {
	return &InstallKubeletServiceOptions{
		ReleaseVersion: "v0.16.2",
		BinPath:        "/bin",
		UseSudo:        true,
	}
}
