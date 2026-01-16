package commandexecutorfileoo

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfile"
)

func (f *File) GetCommandExecutorAndFilePath() (commandexecutorinterfaces.CommandExecutor, string, error) {
	ce, err := f.GetCommandExecutor()
	if err != nil {
		return nil, "", err
	}

	path, err := f.GetPath()
	if err != nil {
		return nil, "", err
	}

	return ce, path, nil
}

func (f *File) Exists(ctx context.Context) (exists bool, err error) {
	commandExecutor, path, err := f.GetCommandExecutorAndFilePath()
	if err != nil {
		return
	}

	return commandexecutorfile.Exists(ctx, commandExecutor, path)
}

func (d *Directory) Exists(ctx context.Context) (exists bool, err error) {
	commandExecutor, path, err := d.GetCommandExecutorAndDirectoryPath()
	if err != nil {
		return
	}

	return commandexecutorfile.Exists(ctx, commandExecutor, path)
}
