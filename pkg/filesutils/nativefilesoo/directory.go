package nativefilesoo

import (
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type Directory struct {
	filesgeneric.DirectoryBase
	path string
}

func NewDirectoryByPath(path string) (filesinterfaces.Directory, error) {
	if path == "" {
		return nil, tracederrors.TracedErrorEmptyString("path")
	}

	ret := &Directory{
		path: path,
	}

	err := ret.SetParentDirectoryForBaseClass(ret)
	if err != nil {
		panic(err)
	}

	return ret, nil
}