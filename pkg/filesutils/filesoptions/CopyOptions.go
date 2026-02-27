package filesoptions

type CopyOptions struct {
	// UseSudo indicates whether to perform the file operation with
	// root/administrator privileges.
	UseSudo bool

	// ReplaceExisting ensures the destination file is recreated rather
	// than updated in place. This creates a new inode, which allows
	// you to replace/ copy over a file (like a running binary) that is
	// currently mapped in memory by an active process.
	ReplaceExisting bool
}
