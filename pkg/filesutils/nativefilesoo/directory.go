package nativefilesoo

import (
	"context"
	"path/filepath"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
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

func (d *Directory) GetParentDirectory(ctx context.Context) (filesinterfaces.Directory, error) {
	path, err := d.GetPath()
	if err != nil {
		return nil, err
	}

	parentPath := filepath.Dir(path)

	if parentPath == "" {
		return nil, tracederrors.TracedError("parentPath is empty string after evaluation.")
	}

	return NewDirectoryByPath(parentPath)
}

func (d *Directory) GetSubDirectory(ctx context.Context, subdirPath ...string) (filesinterfaces.Directory, error) {
	if len(subdirPath) <= 0 {
		return nil, tracederrors.TracedError("path is empty or nil")
	}

	subDirPath, err := d.GetPath()
	if err != nil {
		return nil, err
	}

	for _, s := range subdirPath {
		subDirPath = filepath.Join(subDirPath, s)
	}

	return NewDirectoryByPath(subDirPath)
}

func (d *Directory) CreateSubDirectory(ctx context.Context, subDirectoryName string, options *filesoptions.CreateOptions) (filesinterfaces.Directory, error) {
	subDir, err := d.GetSubDirectory(ctx, subDirectoryName)
	if err != nil {
		return nil, err
	}

	err = subDir.Create(ctx, &filesoptions.CreateOptions{})
	if err != nil {
		return nil, err
	}

	return subDir, nil
}

func (d *Directory) Create(ctx context.Context, options *filesoptions.CreateOptions) (err error) {
	path, err := d.GetPath()
	if err != nil {
		return err
	}

	return nativefiles.CreateDirectory(ctx, path, options)
}

func (d *Directory) Delete(ctx context.Context, options *filesoptions.DeleteOptions) (err error) {
	path, err := d.GetPath()
	if err != nil {
		return err
	}

	return nativefiles.Delete(ctx, path, options)
}

func (d *Directory) Exists(ctx context.Context) (exists bool, err error) {
	path, err := d.GetPath()
	if err != nil {
		return false, err
	}

	exists = nativefiles.Exists(ctx, path)
	return exists, nil
}
