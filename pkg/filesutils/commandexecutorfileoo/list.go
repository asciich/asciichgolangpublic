package commandexecutorfileoo

import (
	"context"
	"sort"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/pathsutils"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func (d *Directory) ListFiles(ctx context.Context, listFileOptions *parameteroptions.ListFileOptions) (files []filesinterfaces.File, err error) {
	if listFileOptions == nil {
		return nil, tracederrors.TracedErrorNil("listFileOptions")
	}

	optionsToUse := listFileOptions.GetDeepCopy()

	optionsToUse.ReturnRelativePaths = true

	paths, err := d.ListFilePaths(ctx, optionsToUse)
	if err != nil {
		return nil, err
	}

	files = []filesinterfaces.File{}
	for _, path := range paths {
		toAdd, err := d.GetFileInDirectory(path)
		if err != nil {
			return nil, err
		}

		files = append(files, toAdd)
	}

	return files, nil
}

func (d *Directory) ListFilePaths(ctx context.Context, listFileOptions *parameteroptions.ListFileOptions) (filePaths []string, err error) {
	if listFileOptions == nil {
		return nil, tracederrors.TracedErrorNil("listFileOptions")
	}

	commandExecutor, err := d.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	dirPath, err := d.GetPath()
	if err != nil {
		return nil, err
	}

	commandToUse := []string{"find", dirPath, "-type", "f"}
	if listFileOptions.NonRecursive {
		commandToUse = []string{"find", dirPath, "-type", "f", "-maxdepth", "1"}
	}

	foundPaths, err := commandExecutor.RunCommandAndGetStdoutAsLines(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: commandToUse,
		},
	)
	if err != nil {
		return nil, err
	}

	filePaths, err = pathsutils.FilterPaths(foundPaths, listFileOptions)
	if err != nil {
		return nil, err
	}

	if listFileOptions.ReturnRelativePaths {
		filePaths, err = pathsutils.GetRelativePathsTo(filePaths, dirPath)
		if err != nil {
			return nil, err
		}
	}

	sort.Strings(filePaths)

	return filePaths, nil
}
