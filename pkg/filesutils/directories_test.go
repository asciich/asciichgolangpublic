package filesutils_test

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorbashoo"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfileoo"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefilesoo"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfilesoo"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
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

	if implementationName == "nativedirectoryoo" || implementationName == "nativefilesoo" {
		return mustutils.Must(nativefilesoo.NewDirectoryByPath(testPath))
	}

	if implementationName == "commandexecutorfileoo" || implementationName == "commandExecutorFileExec" {
		commandExecutor := commandexecutorexecoo.Exec()
		return mustutils.Must(commandexecutorfileoo.NewDirectory(commandExecutor, testPath))
	}

	if implementationName == "commandExecutorFileBash" {
		commandExecutor := commandexecutorbashoo.Bash()
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
				require.True(t, strings.HasPrefix(subDirPath, dirPath), "subdirectory path '%s' must be nested under parent directory '%s' (potential path traversal)", subDirPath, dirPath)

				// Explicitly ensure the subdirectory is NOT rooted at the filesystem root
				require.False(t, strings.HasPrefix(subDirPath, "/subdir"), "subdirectory path '%s' must not be rooted at '/' (path traversal detected)", subDirPath)
			},
		)
	}
}

func TestDirectory_ReadFileInDirectoryAsString(t *testing.T) {
	for _, impl := range allDirectoryImplementations {
		t.Run(impl, func(t *testing.T) {
			ctx := getCtx()

			tempPath, err := tempfilesoo.CreateEmptyTemporaryDirectoryAndGetPath(ctx)
			require.NoError(t, err)

			dir := getDirectoryToTest(impl, tempPath)
			defer dir.Delete(ctx, &filesoptions.DeleteOptions{})

			_, err = dir.WriteStringToFile(ctx, "test.txt", "hello_world", &filesoptions.WriteOptions{})
			require.NoError(t, err)

			content, err := dir.ReadFileInDirectoryAsString(ctx, "test.txt")
			require.NoError(t, err)

			require.EqualValues(t, "hello_world", content)
		})
	}
}

func TestDirectory_ReadFileInDirectoryAsInt64(t *testing.T) {
	tests := []struct {
		implementationName string
		content            string
		expectedInt64      int64
	}{}

	for _, impl := range allDirectoryImplementations {
		for _, tc := range []struct {
			content       string
			expectedInt64 int64
		}{
			{"1234", 1234},
			{"1234\n", 1234},
			{"1234 ", 1234},
			{" 1234", 1234},
			{"\n1234\n", 1234},
			{"\n1234", 1234},
		} {
			tests = append(tests, struct {
				implementationName string
				content            string
				expectedInt64      int64
			}{impl, tc.content, tc.expectedInt64})
		}
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				tempPath, err := tempfilesoo.CreateEmptyTemporaryDirectoryAndGetPath(ctx)
				require.NoError(t, err)

				dir := getDirectoryToTest(tt.implementationName, tempPath)
				defer dir.Delete(ctx, &filesoptions.DeleteOptions{})

				_, err = dir.WriteStringToFile(ctx, "test.txt", tt.content, &filesoptions.WriteOptions{})
				require.NoError(t, err)

				content, err := dir.ReadFileInDirectoryAsInt64(ctx, "test.txt")
				require.NoError(t, err)

				require.EqualValues(t, tt.expectedInt64, content)
			},
		)
	}
}

func TestDirectory_ReadFirstLineOfFileInDirectoryAsString(t *testing.T) {
	for _, impl := range allDirectoryImplementations {
		t.Run(impl, func(t *testing.T) {
			ctx := getCtx()

			tempPath, err := tempfilesoo.CreateEmptyTemporaryDirectoryAndGetPath(ctx)
			require.NoError(t, err)

			dir := getDirectoryToTest(impl, tempPath)
			defer dir.Delete(ctx, &filesoptions.DeleteOptions{})

			_, err = dir.WriteStringToFile(ctx, "test.txt", "1234\nabc\n", &filesoptions.WriteOptions{})
			require.NoError(t, err)

			content, err := dir.ReadFirstLineOfFileInDirectoryAsString(ctx, "test.txt")
			require.NoError(t, err)

			require.EqualValues(t, "1234", content)
		})
	}
}

func TestDirectory_ListSubDirectories_RelativePaths(t *testing.T) {
	for _, impl := range allDirectoryImplementations {
		t.Run(impl, func(t *testing.T) {
			ctx := getCtx()

			tempPath, err := tempfilesoo.CreateEmptyTemporaryDirectoryAndGetPath(ctx)
			require.NoError(t, err)

			testDirectory := getDirectoryToTest(impl, tempPath)
			defer testDirectory.Delete(ctx, &filesoptions.DeleteOptions{})

			_, err = testDirectory.CreateSubDirectory(ctx, "test1", &filesoptions.CreateOptions{})
			require.NoError(t, err)

			test2, err := testDirectory.CreateSubDirectory(ctx, "test2", &filesoptions.CreateOptions{})
			require.NoError(t, err)

			_, err = test2.CreateSubDirectory(ctx, "a", &filesoptions.CreateOptions{})
			require.NoError(t, err)
			_, err = test2.CreateSubDirectory(ctx, "b", &filesoptions.CreateOptions{})
			require.NoError(t, err)
			_, err = test2.CreateSubDirectory(ctx, "c", &filesoptions.CreateOptions{})
			require.NoError(t, err)

			subDirectoryList, err := testDirectory.ListSubDirectoryPaths(
				ctx,
				&parameteroptions.ListDirectoryOptions{
					Recursive:           false,
					ReturnRelativePaths: true,
				},
			)
			require.NoError(t, err)

			require.Len(t, subDirectoryList, 2)
			require.EqualValues(t, "test1", subDirectoryList[0])
			require.EqualValues(t, "test2", subDirectoryList[1])

			subDirectoryList, err = testDirectory.ListSubDirectoryPaths(
				ctx,
				&parameteroptions.ListDirectoryOptions{
					Recursive:           true,
					ReturnRelativePaths: true,
				},
			)
			require.NoError(t, err)

			require.Len(t, subDirectoryList, 5)
			require.EqualValues(t, "test1", subDirectoryList[0])
			require.EqualValues(t, "test2", subDirectoryList[1])
			require.EqualValues(t, "test2/a", subDirectoryList[2])
			require.EqualValues(t, "test2/b", subDirectoryList[3])
			require.EqualValues(t, "test2/c", subDirectoryList[4])
		})
	}
}

func TestDirectory_ListSubDirectories2(t *testing.T) {
	for _, impl := range allDirectoryImplementations {
		t.Run(impl, func(t *testing.T) {
			ctx := getCtx()

			tempPath, err := tempfilesoo.CreateEmptyTemporaryDirectoryAndGetPath(ctx)
			require.NoError(t, err)

			testDirectory := getDirectoryToTest(impl, tempPath)
			defer testDirectory.Delete(ctx, &filesoptions.DeleteOptions{})

			_, err = testDirectory.CreateSubDirectory(ctx, "test1", &filesoptions.CreateOptions{})
			require.NoError(t, err)
			test2, err := testDirectory.CreateSubDirectory(ctx, "test2", &filesoptions.CreateOptions{})
			require.NoError(t, err)

			_, err = test2.CreateSubDirectory(ctx, "a", &filesoptions.CreateOptions{})
			require.NoError(t, err)
			_, err = test2.CreateSubDirectory(ctx, "b", &filesoptions.CreateOptions{})
			require.NoError(t, err)
			_, err = test2.CreateSubDirectory(ctx, "c", &filesoptions.CreateOptions{})
			require.NoError(t, err)

			subDirectoryList, err := testDirectory.ListSubDirectories(
				ctx,
				&parameteroptions.ListDirectoryOptions{
					Recursive: false,
				},
			)
			require.NoError(t, err)

			require.Len(t, subDirectoryList, 2)

			baseName, err := subDirectoryList[0].GetBaseName()
			require.NoError(t, err)
			require.EqualValues(t, "test1", baseName)

			baseName, err = subDirectoryList[1].GetBaseName()
			require.NoError(t, err)
			require.EqualValues(t, "test2", baseName)

			testDirLocalPath, err := testDirectory.GetLocalPath()
			require.NoError(t, err)

			dirName, err := subDirectoryList[0].GetDirName()
			require.NoError(t, err)
			require.EqualValues(t, testDirLocalPath, dirName)

			dirName, err = subDirectoryList[1].GetDirName()
			require.NoError(t, err)
			require.EqualValues(t, testDirLocalPath, dirName)

			subDirectoryList, err = testDirectory.ListSubDirectories(
				ctx,
				&parameteroptions.ListDirectoryOptions{
					Recursive: true,
				},
			)
			require.NoError(t, err)

			require.Len(t, subDirectoryList, 5)

			baseName, err = subDirectoryList[0].GetBaseName()
			require.NoError(t, err)
			require.EqualValues(t, baseName, "test1")

			baseName, err = subDirectoryList[1].GetBaseName()
			require.NoError(t, err)
			require.EqualValues(t, baseName, "test2")

			baseName, err = subDirectoryList[2].GetBaseName()
			require.NoError(t, err)
			require.EqualValues(t, baseName, "a")

			baseName, err = subDirectoryList[3].GetBaseName()
			require.NoError(t, err)
			require.EqualValues(t, baseName, "b")

			baseName, err = subDirectoryList[4].GetBaseName()
			require.NoError(t, err)
			require.EqualValues(t, baseName, "c")

			testDirPath, err := testDirectory.GetLocalPath()
			require.NoError(t, err)

			dirName, err = subDirectoryList[0].GetDirName()
			require.NoError(t, err)
			require.EqualValues(t, dirName, testDirPath)

			dirName, err = subDirectoryList[1].GetDirName()
			require.NoError(t, err)
			require.EqualValues(t, dirName, testDirPath)

			dirName, err = subDirectoryList[2].GetDirName()
			require.NoError(t, err)
			require.EqualValues(t, dirName, filepath.Join(testDirPath, "test2"))

			dirName, err = subDirectoryList[3].GetDirName()
			require.NoError(t, err)
			require.EqualValues(t, dirName, filepath.Join(testDirPath, "test2"))

			dirName, err = subDirectoryList[4].GetDirName()
			require.NoError(t, err)
			require.EqualValues(t, dirName, filepath.Join(testDirPath, "test2"))
		})
	}
}

func TestDirectoriesCreateLocalDirectoryByPath(t *testing.T) {
	for _, impl := range allDirectoryImplementations {
		t.Run(impl, func(t *testing.T) {
			ctx := getCtx()

			tempPath, err := tempfilesoo.CreateEmptyTemporaryDirectoryAndGetPath(ctx)
			require.NoError(t, err)

			directory := getDirectoryToTest(impl, tempPath)
			defer directory.Delete(ctx, &filesoptions.DeleteOptions{})

			exists, err := directory.Exists(ctx)
			require.NoError(t, err)
			require.True(t, exists)

			for i := 0; i < 2; i++ {
				err = directory.Delete(ctx, &filesoptions.DeleteOptions{})
				require.NoError(t, err)

				exists, err = directory.Exists(ctx)
				require.NoError(t, err)
				require.False(t, exists)
			}

			for i := 0; i < 2; i++ {
				localPath, err := directory.GetLocalPath()
				require.NoError(t, err)

				createdDir, err := files.Directories().CreateLocalDirectoryByPath(ctx, localPath, &filesoptions.CreateOptions{})
				require.NoError(t, err)

				dirExists, err := directory.Exists(ctx)
				require.NoError(t, err)
				require.True(t, dirExists)

				createdExists, err := createdDir.Exists(ctx)
				require.NoError(t, err)
				require.True(t, createdExists)
			}
		})
	}
}

func TestDirectory_IsEmptyDirectory(t *testing.T) {
	for _, impl := range allDirectoryImplementations {
		t.Run(impl, func(t *testing.T) {
			ctx := getCtx()

			tempPath, err := tempfilesoo.CreateEmptyTemporaryDirectoryAndGetPath(ctx)
			require.NoError(t, err)

			dir := getDirectoryToTest(impl, tempPath)
			defer dir.Delete(ctx, &filesoptions.DeleteOptions{})

			// A freshly created temporary directory must be empty.
			isEmpty, err := dir.IsEmptyDirectory(ctx)
			require.NoError(t, err)
			require.True(t, isEmpty)

			// After writing a file the directory must no longer be empty.
			_, err = dir.WriteStringToFile(ctx, "test.txt", "hello", &filesoptions.WriteOptions{})
			require.NoError(t, err)

			isEmpty, err = dir.IsEmptyDirectory(ctx)
			require.NoError(t, err)
			require.False(t, isEmpty)
		})
	}
}

func TestDirectory_IsEmptyDirectory_WithSubDirectory(t *testing.T) {
	for _, impl := range allDirectoryImplementations {
		t.Run(impl, func(t *testing.T) {
			ctx := getCtx()

			tempPath, err := tempfilesoo.CreateEmptyTemporaryDirectoryAndGetPath(ctx)
			require.NoError(t, err)

			dir := getDirectoryToTest(impl, tempPath)
			defer dir.Delete(ctx, &filesoptions.DeleteOptions{})

			// A freshly created temporary directory must be empty.
			isEmpty, err := dir.IsEmptyDirectory(ctx)
			require.NoError(t, err)
			require.True(t, isEmpty)

			// After creating a subdirectory the directory must no longer be empty.
			_, err = dir.CreateSubDirectory(ctx, "subdir", &filesoptions.CreateOptions{})
			require.NoError(t, err)

			isEmpty, err = dir.IsEmptyDirectory(ctx)
			require.NoError(t, err)
			require.False(t, isEmpty)
		})
	}
}

func TestDirectory_CreateFilesInDirectory(t *testing.T) {
	for _, impl := range allDirectoryImplementations {
		t.Run(impl, func(t *testing.T) {
			ctx := getCtx()

			tempPath, err := tempfilesoo.CreateEmptyTemporaryDirectoryAndGetPath(ctx)
			require.NoError(t, err)

			dir := getDirectoryToTest(impl, tempPath)
			defer dir.Delete(ctx, &filesoptions.DeleteOptions{})

			// Directory must be empty before creating files.
			isEmpty, err := dir.IsEmptyDirectory(ctx)
			require.NoError(t, err)
			require.True(t, isEmpty)

			filesToCreate := []string{"file1.txt", "file2.txt", "file3.txt"}

			createdFiles, err := dir.CreateFilesInDirectory(ctx, filesToCreate, &filesoptions.CreateOptions{})
			require.NoError(t, err)
			require.Len(t, createdFiles, 3)

			// Directory must no longer be empty after creating files.
			isEmpty, err = dir.IsEmptyDirectory(ctx)
			require.NoError(t, err)
			require.False(t, isEmpty)

			// Verify each created file exists and has the expected base name.
			for i, createdFile := range createdFiles {
				fileExists, err := createdFile.Exists(ctx)
				require.NoError(t, err)
				require.True(t, fileExists)

				baseName, err := createdFile.GetBaseName()
				require.NoError(t, err)
				require.EqualValues(t, filesToCreate[i], baseName)

				// Verify the file is located within the directory.
				dirPath, err := dir.GetLocalPath()
				require.NoError(t, err)

				filePath, err := createdFile.GetLocalPath()
				require.NoError(t, err)
				require.True(t, strings.HasPrefix(filePath, dirPath))
			}
		})
	}
}

func TestDirectory_CreateFilesInDirectory_Idempotent(t *testing.T) {
	for _, impl := range allDirectoryImplementations {
		t.Run(impl, func(t *testing.T) {
			ctx := getCtx()

			tempPath, err := tempfilesoo.CreateEmptyTemporaryDirectoryAndGetPath(ctx)
			require.NoError(t, err)

			dir := getDirectoryToTest(impl, tempPath)
			defer dir.Delete(ctx, &filesoptions.DeleteOptions{})

			filesToCreate := []string{"a.txt", "b.txt"}

			// Creating the same files twice must not fail (idempotent).
			for i := 0; i < 2; i++ {
				createdFiles, err := dir.CreateFilesInDirectory(ctx, filesToCreate, &filesoptions.CreateOptions{})
				require.NoError(t, err)
				require.Len(t, createdFiles, 2)
			}
		})
	}
}
