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

func (d *Directory) CheckIsLocalDirectory() error {
	hostDescription, err := d.GetHostDescription()
	if err != nil {
		return err
	}

	if hostDescription != "localhost" {
		return tracederrors.TracedErrorf("Directory on host '%s' is not on local machine.", hostDescription)
	}

	return nil
}

func (d *Directory) GetLocalPath() (localPath string, err error) {
	err = d.CheckIsLocalDirectory()
	if err != nil {
		return "", err
	}

	return d.GetPath()
}

func (d *Directory) GetHostDescription() (hostDescription string, err error) {
	return "localhost", err
}

func (d *Directory) GetPath() (dirPath string, err error) {
	if d.path == "" {
		return "", tracederrors.TracedError("path not set")
	}

	return d.path, nil
}
