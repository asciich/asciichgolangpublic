package commandexecutorfileoo

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func (d *Directory) Chmod(ctx context.Context, chmodOptions *filesoptions.ChmodOptions) (err error) {
	return tracederrors.TracedErrorNotImplemented()
}

func (d *Directory) CopyContentToDirectory(ctx context.Context, destinationDir filesinterfaces.Directory) (err error) {
	return tracederrors.TracedErrorNotImplemented()
}

func (d *Directory) GetBaseName() (baseName string, err error) {
	return "", tracederrors.TracedErrorNotImplemented()
}

func (d *Directory) GetDirName() (dirName string, err error) {
	return "", tracederrors.TracedErrorNotImplemented()
}

func (d *Directory) GetHostDescription() (hostDescription string, err error) {
	commandExecutor, err := d.GetCommandExecutor()
	if err != nil {
		return "", err
	}

	return commandExecutor.GetHostDescription()
}

func (d *Directory) GetParentDirectory() (parentDirectory filesinterfaces.Directory, err error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}


func (d *Directory) ListSubDirectories(ctx context.Context, options *parameteroptions.ListDirectoryOptions) (subDirectories []filesinterfaces.Directory, err error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}
