package filesoptions

type CreateOptions struct {
	// If true a priviledge escallation is performed to create the file as root.
	UseSudo bool
}
