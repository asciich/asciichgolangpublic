package filesoptions

import (
	"os"

	"github.com/asciich/asciichgolangpublic/pkg/osutils/unixfilepermissionsutils"
)

type WriteOptions struct {
	// If true a priviledge escallation is performed to write to the file as root.
	UseSudo bool

	Perm *os.FileMode
}

func (w *WriteOptions) GetPermOrDefault() os.FileMode {
	if w.Perm == nil {
		return os.FileMode(0644)
	}

	return *w.Perm
}

func (w *WriteOptions) GetPermissionsStringOrDefault() (string, error) {
	perm := w.GetPermOrDefault()

	return unixfilepermissionsutils.GetPermissionString(int(perm))
}