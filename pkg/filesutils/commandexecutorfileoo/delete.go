package commandexecutorfileoo

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfile"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
)

func (f *File) Delete(ctx context.Context, options *filesoptions.DeleteOptions) error {
	commandExecutor, path, err := f.GetCommandExecutorAndFilePath()
	if err != nil {
		return err
	}

	return commandexecutorfile.Delete(ctx, commandExecutor, path, options)
}


func (d *Directory) Delete(ctx context.Context, options *filesoptions.DeleteOptions) (err error) {
	commandExecutor, dirPath, err := d.GetCommandExecutorAndDirectoryPath()
	if err != nil {
		return err
	}

	return commandexecutorfile.DeleteDirectory(ctx, commandExecutor, dirPath, options)
}