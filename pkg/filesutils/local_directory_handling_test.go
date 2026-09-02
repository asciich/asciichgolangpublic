package filesutils_test

import (
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfilesoo"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/pathsutils"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

var allDirectoryImplementations = []string{
	"nativefilesoo",
	"commandExecutorFileExec",
	"commandExecutorFileBash",
}

func Test_LocalDirectoryFulfillsDirectoryInterface(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		ctx := getCtx()
		dir, err := files.GetLocalDirectoryByPath(ctx, ".")
		require.NoError(t, err)
		require.NotNil(t, dir)

		var dirInterface filesinterfaces.Directory = dir
		require.NotNil(t, dirInterface)
	})
}

func TestDirectoryExists(t *testing.T) {
	for _, impl := range allDirectoryImplementations {
		t.Run(impl, func(t *testing.T) {
			ctx := getCtx()

			tempPath, err := tempfilesoo.CreateEmptyTemporaryDirectoryAndGetPath(ctx)
			require.NoError(t, err)

			var directory filesinterfaces.Directory = getDirectoryToTest(impl, tempPath)
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
				err = directory.Create(ctx, &filesoptions.CreateOptions{})
				require.NoError(t, err)

				exists, err = directory.Exists(ctx)
				require.NoError(t, err)
				require.True(t, exists)
			}

			for i := 0; i < 2; i++ {
				err = directory.Delete(ctx, &filesoptions.DeleteOptions{})
				require.NoError(t, err)

				exists, err = directory.Exists(ctx)
				require.NoError(t, err)
				require.False(t, exists)
			}
		})
	}
}

func TestDirectoryGetFileInDirectory(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		ctx := getCtx()
		homeDir, err := files.GetLocalDirectoryByPath(ctx, "/home/")
		require.NoError(t, err)

		file, err := homeDir.GetFileInDirectory("testfile")
		require.NoError(t, err)

		localPath, err := file.GetLocalPath()
		require.NoError(t, err)

		require.EqualValues(t, "/home/testfile", localPath)
	})

	t.Run("with sub path", func(t *testing.T) {
		ctx := getCtx()
		homeDir, err := files.GetLocalDirectoryByPath(ctx, "/home/")
		require.NoError(t, err)

		file, err := homeDir.GetFileInDirectory("subdir", "another_file")
		require.NoError(t, err)

		localPath, err := file.GetLocalPath()
		require.NoError(t, err)

		require.EqualValues(t, "/home/subdir/another_file", localPath)
	})
}

func TestDirectoryGetFilePathInDirectory(t *testing.T) {
	t.Run("testfile in home", func(t *testing.T) {
		ctx := getCtx()
		homeDir, err := files.GetLocalDirectoryByPath(ctx, "/home/")
		require.NoError(t, err)

		path, err := homeDir.GetFilePathInDirectory("testfile")
		require.NoError(t, err)
		require.EqualValues(t, "/home/testfile", path)
	})

	t.Run("testfile subdir in home", func(t *testing.T) {
		ctx := getCtx()
		homeDir, err := files.GetLocalDirectoryByPath(ctx, "/home/")
		require.NoError(t, err)

		path, err := homeDir.GetFilePathInDirectory("subdir", "another_file")
		require.NoError(t, err)

		require.EqualValues(t, "/home/subdir/another_file", path)
	})
}

func TestDirectoryGetSubDirectory(t *testing.T) {
	t.Run("single file", func(t *testing.T) {
		ctx := getCtx()
		homeDir, err := files.GetLocalDirectoryByPath(ctx, "/home/")
		require.NoError(t, err)

		subDir, err := homeDir.GetDirectoryByPath(ctx, "testfile")
		require.NoError(t, err)

		localPath, err := subDir.GetLocalPath()
		require.NoError(t, err)

		require.EqualValues(t, "/home/testfile", localPath)
	})

	t.Run("subdir and file", func(t *testing.T) {
		ctx := getCtx()
		homeDir, err := files.GetLocalDirectoryByPath(ctx, "/home/")
		require.NoError(t, err)

		subDir, err := homeDir.GetDirectoryByPath(ctx, "subdir", "another_file")
		require.NoError(t, err)

		localPath, err := subDir.GetLocalPath()
		require.NoError(t, err)
		require.EqualValues(t, "/home/subdir/another_file", localPath)
	})
}

func TestDirectoryParentForBaseClassSet(t *testing.T) {
	for _, impl := range allDirectoryImplementations {
		t.Run(impl, func(t *testing.T) {
			dir := files.NewLocalDirectory()
			parent, err := dir.GetParentDirectoryForBaseClass()
			require.NoError(t, err)
			require.NotNil(t, parent)
		})
	}
}

func TestDirectoryCreateFileInDirectoryFromString(t *testing.T) {
	tests := []struct {
		implementationName string
		filename           []string
		content            string
	}{}

	for _, impl := range allDirectoryImplementations {
		for _, tc := range []struct {
			filename []string
			content  string
		}{
			{[]string{"testcase"}, "content"},
			{[]string{"testcase", "test.txt"}, "content"},
		} {
			tests = append(tests, struct {
				implementationName string
				filename           []string
				content            string
			}{impl, tc.filename, tc.content})
		}
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				tempDirPath, err := tempfilesoo.CreateEmptyTemporaryDirectoryAndGetPath(ctx)
				require.NoError(t, err)

				dir := getDirectoryToTest(tt.implementationName, tempDirPath)
				defer dir.Delete(ctx, &filesoptions.DeleteOptions{})

				createdFile, err := dir.CreateFileInDirectoryFromString(ctx, tt.content, tt.filename...)
				require.NoError(t, err)

				localPath, err := dir.GetLocalPath()
				require.NoError(t, err)
				pathElements := []string{localPath}
				pathElements = append(pathElements, tt.filename...)
				expectedFileName := filepath.Join(pathElements...)

				localPath, err = createdFile.GetLocalPath()
				require.NoError(t, err)
				require.EqualValues(t, expectedFileName, localPath)

				content, err := createdFile.ReadAsString(ctx)
				require.NoError(t, err)
				require.EqualValues(t, tt.content, content)
			},
		)
	}
}

func TestDirectoryGetLocalPathIsAbsolute(t *testing.T) {
	tests := []struct {
		pathToTest string
	}{
		{"/"},
		{"/tmp"},
		{"abc"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()
				localDir, err := files.GetLocalDirectoryByPath(ctx, tt.pathToTest)
				require.NoError(t, err)

				localPath, err := localDir.GetLocalPath()
				require.NoError(t, err)

				require.True(t, pathsutils.IsAbsolutePath(localPath))
			},
		)
	}
}

func TestDirectoryWriteStringToFile(t *testing.T) {
	tests := []struct {
		implementationName string
		fileName           string
		content            string
	}{}

	for _, impl := range allDirectoryImplementations {
		for _, tc := range []struct {
			fileName string
			content  string
		}{
			{"a.txt", "testcase"},
			{"b.txt", "testcase\nmultiline"},
		} {
			tests = append(tests, struct {
				implementationName string
				fileName           string
				content            string
			}{impl, tc.fileName, tc.content})
		}
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				tempDirPath, err := tempfilesoo.CreateEmptyTemporaryDirectoryAndGetPath(ctx)
				require.NoError(t, err)

				testDirectory := getDirectoryToTest(tt.implementationName, tempDirPath)
				defer testDirectory.Delete(ctx, &filesoptions.DeleteOptions{})

				require.False(t, mustutils.Must(testDirectory.FileInDirectoryExists(ctx, tt.fileName)))

				testFile, err := testDirectory.WriteStringToFile(ctx, tt.fileName, tt.content, &filesoptions.WriteOptions{})
				require.NoError(t, err)

				require.True(t, mustutils.Must(testDirectory.FileInDirectoryExists(ctx, tt.fileName)))

				content, err := testFile.ReadAsString(ctx)
				require.NoError(t, err)
				require.EqualValues(t, tt.content, content)
			},
		)
	}
}

func TestDirectoryListFilesInDirectory(t *testing.T) {
	tests := []struct {
		implementationName string
		fileNames          []string
		listOptions        parameteroptions.ListFileOptions
		expectedPaths      []string
	}{}

	for _, impl := range allDirectoryImplementations {
		for _, tc := range []struct {
			fileNames     []string
			listOptions   parameteroptions.ListFileOptions
			expectedPaths []string
		}{
			{[]string{"a.go", "b.go"}, parameteroptions.ListFileOptions{ReturnRelativePaths: true}, []string{"a.go", "b.go"}},
			{[]string{"a.go", "a/b.go"}, parameteroptions.ListFileOptions{ReturnRelativePaths: true}, []string{"a.go", "a/b.go"}},
			{[]string{"a.go", "a/b.go", "b.go"}, parameteroptions.ListFileOptions{ReturnRelativePaths: true, ExcludeBasenamePattern: []string{"a.*"}}, []string{"a/b.go", "b.go"}},
			{[]string{"a.go", "a/b.go", "b.go"}, parameteroptions.ListFileOptions{ReturnRelativePaths: true, ExcludeBasenamePattern: []string{"b.*"}}, []string{"a.go"}},
			{[]string{"b.go", "a.go"}, parameteroptions.ListFileOptions{ReturnRelativePaths: true}, []string{"a.go", "b.go"}},
			{[]string{"b.go", "a.go", "go.mod", "go.sum"}, parameteroptions.ListFileOptions{ReturnRelativePaths: true}, []string{"a.go", "b.go", "go.mod", "go.sum"}},
			{[]string{"b.go", "a.go", "go.mod", "go.sum"}, parameteroptions.ListFileOptions{ReturnRelativePaths: true, MatchBasenamePattern: []string{".*.go"}}, []string{"a.go", "b.go"}},
			{[]string{"b.go", "a.go", "go.mod", "go.sum"}, parameteroptions.ListFileOptions{ReturnRelativePaths: true, ExcludeBasenamePattern: []string{".*.go"}}, []string{"go.mod", "go.sum"}},
			{[]string{"b.go", "a.go", "go.go", "go.mod", "go.sum"}, parameteroptions.ListFileOptions{ReturnRelativePaths: true, MatchBasenamePattern: []string{"go.*"}, ExcludeBasenamePattern: []string{".*.go", ".*.mod"}}, []string{"go.sum"}},
		} {
			tests = append(tests, struct {
				implementationName string
				fileNames          []string
				listOptions        parameteroptions.ListFileOptions
				expectedPaths      []string
			}{impl, tc.fileNames, tc.listOptions, tc.expectedPaths})
		}
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				tempDirPath, err := tempfilesoo.CreateEmptyTemporaryDirectoryAndGetPath(ctx)
				require.NoError(t, err)

				temporaryDirectory := getDirectoryToTest(tt.implementationName, tempDirPath)
				defer temporaryDirectory.Delete(ctx, &filesoptions.DeleteOptions{})

				_, err = temporaryDirectory.CreateFilesInDirectory(ctx, tt.fileNames, &filesoptions.CreateOptions{})
				require.NoError(t, err)

				listedFiles, err := temporaryDirectory.ListFilePaths(ctx, &tt.listOptions)
				require.NoError(t, err)
				require.EqualValues(t, tt.expectedPaths, listedFiles)
			},
		)
	}
}

func TestDirectoryCreate(t *testing.T) {
	tests := []struct {
		implementationName string
		subDirPath         []string
	}{}

	for _, impl := range allDirectoryImplementations {
		for _, tc := range []struct {
			subDirPath []string
		}{
			{[]string{"a"}},
			{[]string{"a", "b"}},
		} {
			tests = append(tests, struct {
				implementationName string
				subDirPath         []string
			}{impl, tc.subDirPath})
		}
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				tempPath, err := tempfilesoo.CreateEmptyTemporaryDirectoryAndGetPath(ctx)
				require.NoError(t, err)

				tempDir := getDirectoryToTest(tt.implementationName, tempPath)
				defer tempDir.Delete(ctx, &filesoptions.DeleteOptions{})

				subDir, err := tempDir.GetDirectoryByPath(ctx, tt.subDirPath...)
				require.NoError(t, err)

				exists, err := subDir.Exists(ctx)
				require.NoError(t, err)
				require.False(t, exists)

				err = subDir.Create(ctx, &filesoptions.CreateOptions{})
				require.NoError(t, err)

				exists, err = subDir.Exists(ctx)
				require.NoError(t, err)
				require.True(t, exists)
			},
		)
	}
}

func TestDirectoryGetPathReturnsAbsoluteValue(t *testing.T) {
	tests := []struct {
		path string
	}{
		{"."},
		{".."},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := contextutils.ContextVerbose()

				startPath, err := os.Getwd()
				if err != nil {
					t.Fatalf("%v", err)
				}

				var path1 string
				var path2 string

				var waitGroup sync.WaitGroup

				testFunction := func() {
					defer os.Chdir(startPath)
					defer waitGroup.Done()

					directory, err := files.GetLocalDirectoryByPath(ctx, tt.path)
					require.NoError(t, err)

					path1, err = directory.GetLocalPath()
					os.Chdir("..")
					path2, err = directory.GetLocalPath()
				}

				waitGroup.Add(1)
				go testFunction()
				waitGroup.Wait()

				require.True(t, pathsutils.IsAbsolutePath(path1))
				require.True(t, pathsutils.IsAbsolutePath(path2))

				require.EqualValues(t, path1, path2)

				currentPath, err := os.Getwd()
				if err != nil {
					t.Fatalf("%v", err)
				}

				require.EqualValues(t, startPath, currentPath)
			},
		)
	}
}

func TestDirectoryIsEmptyDirectory(t *testing.T) {
	for _, impl := range allDirectoryImplementations {
		t.Run(impl, func(t *testing.T) {
			ctx := getCtx()

			tempPath, err := tempfilesoo.CreateEmptyTemporaryDirectoryAndGetPath(ctx)
			require.NoError(t, err)

			tempDir := getDirectoryToTest(impl, tempPath)
			defer tempDir.Delete(ctx, &filesoptions.DeleteOptions{})

			isEmptyDir, err := tempDir.IsEmptyDirectory(ctx)
			require.NoError(t, err)

			require.True(t, isEmptyDir)
		})
	}
}

func TestDirectory_CheckExists(t *testing.T) {
	for _, impl := range allDirectoryImplementations {
		t.Run(impl, func(t *testing.T) {
			ctx := getCtx()

			tempPath, err := tempfilesoo.CreateEmptyTemporaryDirectoryAndGetPath(ctx)
			require.NoError(t, err)

			temporaryDirectory := getDirectoryToTest(impl, tempPath)
			defer temporaryDirectory.Delete(ctx, &filesoptions.DeleteOptions{})

			require.Nil(t, temporaryDirectory.CheckExists(ctx))

			err = temporaryDirectory.Delete(ctx, &filesoptions.DeleteOptions{})
			require.NoError(t, err)

			require.NotNil(t, temporaryDirectory.CheckExists(ctx))
		})
	}
}
