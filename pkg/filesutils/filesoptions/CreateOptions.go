package filesoptions

type CreateOptions struct {
	// If true a priviledge escallation is performed to create the file or directory as root.
	UseSudo bool
}
