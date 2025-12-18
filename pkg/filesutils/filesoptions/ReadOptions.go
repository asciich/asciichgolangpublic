package filesoptions

type ReadOptions struct {
	// Use sudo to read the file.
	// Useful when access as current user is denied.
	UseSudo bool
}