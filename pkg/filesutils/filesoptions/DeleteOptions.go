package filesoptions

type DeleteOptions struct {
	// If true a priviledge escallation is performed to delete the file or directory as root.
	UseSudo bool
}