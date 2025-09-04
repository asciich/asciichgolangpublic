package nativefiles

import (
	"context"
	"os"
	"path/filepath"

	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/datatypes/slicesutils"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/pathsutils"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func ListFiles(ctx context.Context, path string, listOptions *parameteroptions.ListFileOptions) ([]string, error) {
	if path == "" {
		return nil, tracederrors.TracedErrorEmptyString("path")
	}

	if listOptions == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	listOptions = listOptions.GetDeepCopy()
	listOptions.OnlyFiles = true

	filePathList := []string{}
	err := filepath.Walk(
		path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			isSymlink, err := IsSymlinkToDirectory(contextutils.WithSilent(ctx), path)
			if err != nil {
				return err
			}

			if isSymlink {
				resolved, err := ResolveSymlink(ctx, path)
				if err != nil {
					return err
				}

				toExtend, err := ListFiles(contextutils.WithSilent(ctx), resolved, &parameteroptions.ListFileOptions{ReturnRelativePaths: true})
				if err != nil {
					return err
				}

				for _, toAdd := range toExtend {
					filePathList = append(filePathList, filepath.Join(path, toAdd))
				}
			} else {
				filePathList = append(filePathList, path)
			}
			return nil
		},
	)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Unable to filepath.Walk: '%w'", err)
	}

	filePathList = slicesutils.RemoveEmptyStrings(filePathList)

	filePathList, err = pathsutils.FilterPaths(filePathList, listOptions)
	if err != nil {
		return nil, err
	}

	if listOptions.ReturnRelativePaths {
		filePathList, err = pathsutils.GetRelativePathsTo(filePathList, path)
		if err != nil {
			return nil, err
		}
	}

	filePathList = slicesutils.SortStringSliceAndRemoveEmpty(filePathList)

	if len(filePathList) <= 0 {
		if !listOptions.AllowEmptyListIfNoFileIsFound {
			return nil, tracederrors.TracedErrorf("No files in '%s' found", path)
		}
	}

	return filePathList, nil
}
