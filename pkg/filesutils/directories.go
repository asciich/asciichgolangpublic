package filesutils

import (
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefilesoo"
)

func NewDirectoryByPath(path string) (filesinterfaces.Directory, error) {
	return nativefilesoo.NewDirectoryByPath(path)
}
