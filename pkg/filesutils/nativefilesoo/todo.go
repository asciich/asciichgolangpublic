package nativefilesoo

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func (f *File) AppendBytes(ctx context.Context, toWrite []byte) (err error) {
	return tracederrors.TracedErrorNotImplemented()
}
func (f *File) AppendString(ctx context.Context, toWrite string) (err error) {
	return tracederrors.TracedErrorNotImplemented()
}

func (f *File) Chown(ctx context.Context, options *parameteroptions.ChownOptions) (err error) {
	return tracederrors.TracedErrorNotImplemented()
}

func (f *File) GetBaseName() (baseName string, err error) {
	return "", tracederrors.TracedErrorNotImplemented()
}

func (f *File) GetHostDescription() (hostDescription string, err error) {
	return "localhost", nil
}
func (f *File) GetLocalPath() (localPath string, err error) {
	return f.GetPath()
}
func (f *File) GetLocalPathOrEmptyStringIfUnset() (localPath string, err error) {
	return "", tracederrors.TracedErrorNotImplemented()
}
func (f *File) GetParentDirectory(ctx context.Context) (parentDirectory filesinterfaces.Directory, err error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}

func (f *File) GetUriAsString() (uri string, err error) {
	return "", tracederrors.TracedErrorNotImplemented()
}
func (f *File) MoveToPath(ctx context.Context, destPath string, useSudo bool) (movedFile filesinterfaces.File, err error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}

func (f *File) SecurelyDelete(ctx context.Context) (err error) {
	return tracederrors.TracedErrorNotImplemented()
}

func (f *File) String() (path string) {
	return f.path
}

func (d *Directory) Chmod(ctx context.Context, chmodOptions *filesoptions.ChmodOptions) error {
	return tracederrors.TracedErrorNotImplemented()
}

func (d *Directory) CopyContentToDirectory(ctx context.Context, destinationDir filesinterfaces.Directory) (err error) {
	return tracederrors.TracedErrorNotImplemented()
}

func (d *Directory) Create(ctx context.Context, options *filesoptions.CreateOptions) (err error) {
	return tracederrors.TracedErrorNotImplemented()
}

func (d *Directory) CreateSubDirectory(ctx context.Context, subDirectoryName string, options *filesoptions.CreateOptions) (createdSubDirectory filesinterfaces.Directory, err error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}

func (d *Directory) Delete(ctx context.Context, options *filesoptions.DeleteOptions) (err error) {
	return tracederrors.TracedErrorNotImplemented()
}

func (d *Directory) Exists(ctx context.Context) (exists bool, err error) {
	return false, tracederrors.TracedErrorNotImplemented()
}

func (d *Directory) GetBaseName() (baseName string, err error) {
	return "", tracederrors.TracedErrorNotImplemented()
}
func (d *Directory) GetDirName() (dirName string, err error) {
	return "", tracederrors.TracedErrorNotImplemented()
}
func (d *Directory) GetFileInDirectory(pathToFile ...string) (file filesinterfaces.File, err error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}

func (d *Directory) GetParentDirectory(ctx context.Context) (parentDirectory filesinterfaces.Directory, err error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}

func (d *Directory) GetSubDirectory(ctx context.Context, path ...string) (subDirectory filesinterfaces.Directory, err error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}
func (d *Directory) IsLocalDirectory() (isLocalDirectory bool, err error) {
	return false, tracederrors.TracedErrorNotImplemented()
}
func (d *Directory) ListFiles(ctx context.Context, listFileOptions *parameteroptions.ListFileOptions) (files []filesinterfaces.File, err error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}
func (d *Directory) ListSubDirectories(ctx context.Context, options *parameteroptions.ListDirectoryOptions) (subDirectories []filesinterfaces.Directory, err error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}
