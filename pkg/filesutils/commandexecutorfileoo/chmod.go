package commandexecutorfileoo

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfile"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
)

func (f *File) Chmod(ctx context.Context, options *filesoptions.ChmodOptions) (err error) {
	path, err := f.GetPath()
	if err != nil {
		return err
	}

	commandExecutor, err := f.GetCommandExecutor()
	if err != nil {
		return err
	}

	return commandexecutorfile.Chmod(ctx, commandExecutor, path, options)
}

func (f *File) GetAccessPermissions() (permission int, err error) {
	path, err := f.GetPath()
	if err != nil {
		return 0, err
	}

	commandExecutor, err := f.GetCommandExecutor()
	if err != nil {
		return 0, err
	}

	return commandexecutorfile.GetAccessPermissions(commandExecutor, path)
}

func (f *File) GetAccessPermissionsString() (permissionString string, err error) {
	path, err := f.GetPath()
	if err != nil {
		return "", err
	}

	commandExecutor, err := f.GetCommandExecutor()
	if err != nil {
		return "", err
	}

	return commandexecutorfile.GetAccessPermissionsString(commandExecutor, path)
}
