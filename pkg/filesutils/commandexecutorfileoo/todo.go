package commandexecutorfileoo

import (
	"context"
	"path/filepath"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfile"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func (f *File) AppendBytes(ctx context.Context, toWrite []byte) (err error) {
	commandExecutor, err := f.GetCommandExecutor()
	if err != nil {
		return err
	}

	path, err := f.GetPath()
	if err != nil {
		return err
	}

	return commandexecutorfile.AppendBytes(ctx, commandExecutor, path, toWrite)
}

func (f *File) AppendString(ctx context.Context, toWrite string) (err error) {
	commandExecutor, err := f.GetCommandExecutor()
	if err != nil {
		return err
	}

	path, err := f.GetPath()
	if err != nil {
		return err
	}

	return commandexecutorfile.AppendString(ctx, commandExecutor, path, toWrite)
}

func (f *File) Chown(ctx context.Context, options *parameteroptions.ChownOptions) (err error) {
	commandExecutor, err := f.GetCommandExecutor()
	if err != nil {
		return err
	}

	path, err := f.GetPath()
	if err != nil {
		return err
	}

	return commandexecutorfile.Chown(ctx, commandExecutor, path, options)
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
	if f.path == "" {
		return "", nil
	}
	return f.GetLocalPath()
}

func (f *File) GetParentDirectory(ctx context.Context) (parentDirectory filesinterfaces.Directory, err error) {
	path, err := f.GetPath()
	if err != nil {
		return nil, err
	}

	parentPath := filepath.Dir(path)

	commandExecutor, err := f.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	return NewDirectory(commandExecutor, parentPath)
}

func (f *File) GetUriAsString() (uri string, err error) {
	path, err := f.GetPath()
	if err != nil {
		return "", err
	}

	return path, nil
}

func (f *File) MoveToPath(ctx context.Context, destPath string, useSudo bool) (movedFile filesinterfaces.File, err error) {
	commandExecutor, err := f.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	path, err := f.GetPath()
	if err != nil {
		return nil, err
	}

	options := &filesoptions.MoveOptions{
		UseSudo: useSudo,
	}

	err = commandexecutorfile.Move(ctx, commandExecutor, path, destPath, options)
	if err != nil {
		return nil, err
	}

	return New(commandExecutor, destPath)
}

func (f *File) SecurelyDelete(ctx context.Context) (err error) {
	commandExecutor, err := f.GetCommandExecutor()
	if err != nil {
		return err
	}

	path, err := f.GetPath()
	if err != nil {
		return err
	}

	return commandexecutorfile.SecurelyDelete(ctx, commandExecutor, path)
}

func (f *File) String() (path string) {
	return f.path
}

func (f *File) GetCommandExecutor() (commandExecutor commandexecutorinterfaces.CommandExecutor, err error) {
	if f.commandExecutor == nil {
		return nil, tracederrors.TracedError("commandExecutor not set")
	}

	return f.commandExecutor, nil
}
