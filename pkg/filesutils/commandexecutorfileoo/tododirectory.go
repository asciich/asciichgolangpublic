package commandexecutorfileoo

import (
	"context"
	"path/filepath"
	"slices"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfile"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func (d *Directory) Chmod(ctx context.Context, chmodOptions *filesoptions.ChmodOptions) (err error) {
	commandExecutor, err := d.GetCommandExecutor()
	if err != nil {
		return err
	}

	path, err := d.GetPath()
	if err != nil {
		return err
	}

	return commandexecutorfile.Chmod(ctx, commandExecutor, path, chmodOptions)
}

func (d *Directory) CopyContentToDirectory(ctx context.Context, destinationDir filesinterfaces.Directory) (err error) {
	// List all files in source directory
	files, err := d.ListFiles(ctx, &parameteroptions.ListFileOptions{})
	if err != nil {
		return err
	}

	// Copy each file to destination
	for _, file := range files {
		srcFile := file
		// Get just the base name of the file
		baseName, err := srcFile.GetBaseName()
		if err != nil {
			return err
		}

		destFile, err := destinationDir.GetFileInDirectory(baseName)
		if err != nil {
			return err
		}

		err = srcFile.CopyToFile(ctx, destFile, &filesoptions.CopyOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *Directory) GetBaseName() (baseName string, err error) {
	path, err := d.GetPath()
	if err != nil {
		return "", err
	}

	baseName = filepath.Base(path)

	if slices.Contains([]string{"", " ", ".", "/", "\\"}, baseName) {
		return "", tracederrors.TracedErrorf("Evaluated invalid baseName '%s' out of path '%s'.", baseName, path)
	}

	return baseName, nil
}

func (d *Directory) GetDirName() (dirName string, err error) {
	path, err := d.GetPath()
	if err != nil {
		return "", err
	}

	dirName = filepath.Dir(path)

	return dirName, nil
}

func (d *Directory) GetHostDescription() (hostDescription string, err error) {
	commandExecutor, err := d.GetCommandExecutor()
	if err != nil {
		return "", err
	}

	return commandExecutor.GetHostDescription()
}

