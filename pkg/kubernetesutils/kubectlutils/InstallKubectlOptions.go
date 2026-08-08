package kubectlutils

// InstallKubectlOptions contains options for installing kubectl.
type InstallKubectlOptions struct {
	// InstallPath is the path where kubectl will be installed.
	// Default: /bin/kubectl
	InstallPath string

	// UseSudo determines if sudo should be used for installation.
	// Default: true
	UseSudo bool

	// Version is the kubectl version to install (e.g., "v1.36.2").
	// Default: "v1.36.2"
	Version string
}

// DefaultInstallKubectlOptions returns the default options for installing kubectl.
func DefaultInstallKubectlOptions() *InstallKubectlOptions {
	return &InstallKubectlOptions{
		InstallPath: "/bin/kubectl",
		UseSudo:     true,
		Version:     "v1.36.2",
	}
}
