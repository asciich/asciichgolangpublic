package nativefilesoo

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func (f *File) CopyToFile(ctx context.Context, destFile filesinterfaces.File) (err error) {
	return tracederrors.TracedErrorNotImplemented()
}

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

func (f *File) GetSizeBytes() (fileSize int64, err error) {
	return -1, tracederrors.TracedErrorNotImplemented()
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
func (f *File) Truncate(ctx context.Context, newSizeBytes int64) (err error) {
	return tracederrors.TracedErrorNotImplemented()
}
