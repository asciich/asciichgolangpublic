package kvmutilsoptions

// AdditionalBlockDevice describes a host block device which is passed to the VM as an
// additional raw virtio disk. Supported are whole disks (e.g. '/dev/sda'), partitions
// (e.g. '/dev/sda1') and device mapper/ LVM devices (e.g. '/dev/mapper/abc').
//
// The target device name inside the VM ('vdb', 'vdc', ...) is assigned automatically
// based on the order in KvmCreateVmOptions.AdditionalBlockDevices. 'vda' is reserved for the root disk.
type AdditionalBlockDevice struct {
	// Absolute path of the block device on the host. Must start with '/dev/'.
	// Stable paths like '/dev/disk/by-id/...' are recommended for whole disks.
	SourcePath string

	// Optional serial (max 20 chars, [A-Za-z0-9_-]).
	// Inside the VM the disk is available as '/dev/disk/by-id/virtio-<Serial>'.
	Serial string

	// Attach the block device read only.
	ReadOnly bool
}
