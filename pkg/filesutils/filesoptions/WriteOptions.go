package filesoptions

type WriteOptions struct {
	// If true a priviledge escallation is performed to write to the file as root.
	UseSudo bool
}