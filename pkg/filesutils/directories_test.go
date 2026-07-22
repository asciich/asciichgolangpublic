package filesutils_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfileoo"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefilesoo"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func getDirectoryToTest(implementationName string, testPath string) (directory filesinterfaces.Directory) {
	if implementationName == "localDirectory" {
		ctx := contextutils.ContextVerbose()
		dir, err := files.GetLocalDirectoryByPath(ctx, testPath)
		if err != nil {
			logging.LogGoErrorFatal(err)
		}

		return dir
	}

	if implementationName == "localCommandExecutorDirectory" {
		return mustutils.Must(files.GetLocalCommandExecutorDirectoryByPath(testPath))
	}

	if implementationName == "nativedirectoryoo" {
		return mustutils.Must(nativefilesoo.NewDirectoryByPath(testPath))
	}

	if implementationName == "commandexecutorfileoo" {
		commandExecutor := commandexecutorexecoo.Exec()
		return mustutils.Must(commandexecutorfileoo.NewDirectory(commandExecutor, testPath))
	}

	panic(fmt.Sprintf("unknown implementationName='%s'", implementationName))
}

func TestGetLocalPath(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"localDirectory"},
		{"localCommandExecutorDirectory"},
		{"nativedirectoryoo"},
	}

	for _, tt := range tests {
		t.Run("getlocalpath_"+tt.implementationName, func(t *testing.T) {
			const testPath = "/testfile"

			sourceFile := getDirectoryToTest(tt.implementationName, testPath)

			localPath, err := sourceFile.GetLocalPath()
			require.NoError(t, err)

			require.EqualValues(t, "/testfile", localPath)
		})
	}
}

func TestDirectory_GetParentDirectory(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"localDirectory"},                // This is a legacy implementaion we should get rid off.
		{"localCommandExecutorDirectory"}, // This is a legacy implementaion we should get rid off.
		{"nativedirectoryoo"},
		{"commandexecutorfileoo"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				tempDirPath, err := tempfiles.CreateTempDir(ctx)
				require.NoError(t, err)

				dir := getDirectoryToTest(tt.implementationName, tempDirPath)
				defer dir.Delete(ctx, &filesoptions.DeleteOptions{})

				subDir, err := dir.CreateSubDirectory(ctx, "subdir", &filesoptions.CreateOptions{})
				require.NoError(t, err)

				dirPath, err := dir.GetPath()
				require.NoError(t, err)

				subDirPath, err := subDir.GetPath()
				require.NoError(t, err)

				require.NotEqualValues(t, dirPath, subDirPath)

				parentDir, err := subDir.GetParentDirectory(ctx)
				require.NoError(t, err)

				parentDirPath, err := parentDir.GetPath()
				require.NoError(t, err)

				require.EqualValues(t, dirPath, parentDirPath)
			},
		)
	}
}

func TestDirectory_CreateSubDirectory_NoPathTraversal(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"localDirectory"},                // This is a legacy implementaion we should get rid off.
		{"localCommandExecutorDirectory"}, // This is a legacy implementaion we should get rid off.
		{"nativedirectoryoo"},
		{"commandexecutorfileoo"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				tempDirPath, err := tempfiles.CreateTempDir(ctx)
				require.NoError(t, err)

				dir := getDirectoryToTest(tt.implementationName, tempDirPath)
				defer dir.Delete(ctx, &filesoptions.DeleteOptions{})

				// Attempt to create a subdirectory with a leading slash.
				// This must NOT result in a path rooted at "/" (path traversal).
				subDir, err := dir.CreateSubDirectory(ctx, "/subdir", &filesoptions.CreateOptions{})
				require.NoError(t, err)

				dirPath, err := dir.GetPath()
				require.NoError(t, err)

				subDirPath, err := subDir.GetPath()
				require.NoError(t, err)

				// The subdirectory path must be different from the parent
				require.NotEqualValues(t, dirPath, subDirPath)

				// The subdirectory path must start with the parent directory path,
				// ensuring it is nested within it and not rooted at "/"
				require.True(
					t,
					strings.HasPrefix(subDirPath, dirPath),
					"subdirectory path '%s' must be nested under parent directory '%s' (potential path traversal)",
					subDirPath,
					dirPath,
				)

				// Explicitly ensure the subdirectory is NOT rooted at the filesystem root
				require.False(
					t,
					strings.HasPrefix(subDirPath, "/subdir"),
					"subdirectory path '%s' must not be rooted at '/' (path traversal detected)",
					subDirPath,
				)
			},
		)
	}
}
