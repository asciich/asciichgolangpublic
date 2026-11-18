package nativefiles

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// ListSubDirectories returns the paths of the sub directories of the given path.
//
// If options.Recursive is set, all nested sub directories are returned,
// otherwise only the direct children.
// The given path itself is never part of the result.
func ListSubDirectories(ctx context.Context, path string, options *parameteroptions.ListDirectoryOptions) ([]string, error) {
	if path == "" {
		return nil, tracederrors.TracedError("path is empty")
	}

	if options == nil {
		options = &parameteroptions.ListDirectoryOptions{}
	}

	const hostDescription = "localhost"

	logging.LogInfoByCtxf(ctx, "List subdirectories of '%s' on '%s' started.", path, hostDescription)

	subDirectories := make([]string, 0)

	if options.Recursive {
		err := filepath.WalkDir(path, func(entryPath string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			// Allow callers to abort long running walks.
			if ctxErr := ctx.Err(); ctxErr != nil {
				return ctxErr
			}

			if !entry.IsDir() {
				return nil
			}

			// Skip the root directory itself.
			if entryPath == path {
				return nil
			}

			subDirectories = append(subDirectories, entryPath)
			return nil
		})
		if err != nil {
			return nil, tracederrors.TracedErrorf(
				"Unable to list sub directories of '%s' recursively: %w",
				path,
				err,
			)
		}

		return subDirectories, nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, tracederrors.TracedErrorf(
			"Unable to list sub directories of '%s': %w",
			path,
			err,
		)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		subDirectories = append(subDirectories, filepath.Join(path, entry.Name()))
	}

	logging.LogInfoByCtxf(ctx, "List subdirectories of '%s' on '%s' finished. Found %d sub directories.", path, hostDescription, len(subDirectories))

	return subDirectories, nil
}
