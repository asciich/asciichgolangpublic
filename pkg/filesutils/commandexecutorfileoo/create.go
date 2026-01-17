package commandexecutorfileoo

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfile"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
)

func (f *File) Create(ctx context.Context, options *filesoptions.CreateOptions) (err error) {
	commandexecutor, path, err := f.GetCommandExecutorAndFilePath()
	if err != nil {
		return err
	}

	return commandexecutorfile.CreateFile(ctx, commandexecutor, path, options)
}


func (d *Directory) Create(ctx context.Context, options *filesoptions.CreateOptions) (err error) {
	commandexecutor, path, err := d.GetCommandExecutorAndDirectoryPath()
	if err != nil {
		return err
	}

	return commandexecutorfile.CreateDirectory(ctx, commandexecutor, path, options)
}
