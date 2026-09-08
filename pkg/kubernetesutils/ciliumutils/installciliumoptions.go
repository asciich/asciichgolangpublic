package ciliumutils

// InstallCiliumOptions contains options for installing Cilium (eBPF CNI).
type InstallCiliumOptions struct {
	// Version is the Cilium (dataplane/agent) version to install (e.g. "1.20.1").
	// Optional: the Cilium CLI picks a compatible default if empty.
	Version string

	// CliVersion is the cilium CLI version to install (e.g. "v0.20.0").
	// Default: "v0.20.0"
	CliVersion string

	// CliInstallPath is where the cilium CLI binary is installed.
	// Default: /usr/local/bin/cilium
	CliInstallPath string

	// PodNetworkCidr, when set, makes Cilium honor kubeadm's pod CIDR by using
	// kubernetes-based IPAM (ipam.mode=kubernetes). When empty, Cilium uses its
	// own cluster-pool IPAM default.
	PodNetworkCidr string

	// ExtraHelmSet are extra "key=value" pairs passed as --set to "cilium install".
	// Used for kernel-compatibility workarounds. On very new kernels (e.g. 7.x)
	// the eBPF verifier can reject Cilium's socketLB CGroupSock program, so we
	// disable socketLB and fall back to legacy host routing by default.
	ExtraHelmSet []string

	// UseSudo determines if sudo should be used for CLI installation steps.
	// Default: true
	UseSudo bool
}

// DefaultInstallCiliumOptions returns the default options for installing Cilium.
func DefaultInstallCiliumOptions() *InstallCiliumOptions {
	return &InstallCiliumOptions{
		CliVersion:     "v0.20.0",
		CliInstallPath: "/usr/local/bin/cilium",
		UseSudo:        true,
		// Kernel-compatibility workaround (Option 1): the socketLB CGroupSock
		// program fails to load on newer kernels ("R1 is not a scalar"), so
		// disable it and use legacy host routing.
		ExtraHelmSet: []string{
			"socketLB.enabled=false",
			"bpf.hostLegacyRouting=true",
		},
	}
}
