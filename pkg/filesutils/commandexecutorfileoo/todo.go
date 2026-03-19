package commandexecutorfileoo

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
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

func (f *File) GetDeepCopy() (deepCopy filesinterfaces.File) {
	copy := &File{}
	err := copy.SetParentFileForBaseClass(copy)
	if err != nil {
		panic(err)
	}

	if f.commandExecutor != nil {
		copy.commandExecutor = f.commandExecutor.GetDeepCopyAsCommandExecutor()
	}

	copy.path = f.path

	return copy
}

func (f *File) GetLocalPathOrEmptyStringIfUnset() (localPath string, err error) {
	return "", tracederrors.TracedErrorNotImplemented()
}

func (f *File) GetParentDirectory(ctx context.Context) (parentDirectory filesinterfaces.Directory, err error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}

func (f *File) GetPath() (path string, err error) {
	if f.path == "" {
		return "", tracederrors.TracedError("path not set")
	}

	return f.path, nil
}

func (f *File) GetSizeBytes() (fileSize int64, err error) {
	return 0, tracederrors.TracedErrorNotImplemented()
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
	logging.LogFatalWithTrace("Not implemented")
	return ""
}

func (f *File) Truncate(ctx context.Context, newSizeBytes int64) (err error) {
	return tracederrors.TracedErrorNotImplemented()
}

func (f *File) GetCommandExecutor() (commandexecutorinterfaces.CommandExecutor, error) {
	if f.commandExecutor == nil {
		return nil, tracederrors.TracedError("commandExecutor not set")
	}

	return f.commandExecutor, nil
}
