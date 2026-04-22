package filesutils

import (
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefilesoo"
)

func NewFileByPath(path string) (filesinterfaces.File, error) {
	return nativefilesoo.NewFileByPath(path)
}
