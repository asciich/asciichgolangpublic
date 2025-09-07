package files_test

import (
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/pathsutils"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func TestLocalDirectoryExists(t *testing.T) {

	tests := []struct {
		testcase string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				var directory filesinterfaces.Directory = getDirectoryToTest("localDirectory")
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
			},
		)
	}
}

func TestLocalDirectoryGetFileInDirectory(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		homeDir, err := files.GetLocalDirectoryByPath("/home/")
		require.NoError(t, err)

		file, err := homeDir.GetFileInDirectory("testfile")
		require.NoError(t, err)

		localPath, err := file.GetLocalPath()
		require.NoError(t, err)

		require.EqualValues(t, "/home/testfile", localPath)
	})

	t.Run("with sub path", func(t *testing.T) {
		homeDir, err := files.GetLocalDirectoryByPath("/home/")
		require.NoError(t, err)

		file, err := homeDir.GetFileInDirectory("subdir", "another_file")
		require.NoError(t, err)

		localPath, err := file.GetLocalPath()
		require.NoError(t, err)

		require.EqualValues(t, "/home/subdir/another_file", localPath)
	})
}

func TestLocalDirectoryGetFilePathInDirectory(t *testing.T) {
	t.Run("testfile in home", func(t *testing.T) {
		homeDir, err := files.GetLocalDirectoryByPath("/home/")
		require.NoError(t, err)

		path, err := homeDir.GetFilePathInDirectory("testfile")
		require.NoError(t, err)
		require.EqualValues(t, "/home/testfile", path)

	})

	t.Run("testfile subdir in in home", func(t *testing.T) {
		homeDir, err := files.GetLocalDirectoryByPath("/home/")
		require.NoError(t, err)

		path, err := homeDir.GetFilePathInDirectory("subdir", "another_file")
		require.NoError(t, err)

		require.EqualValues(t, "/home/subdir/another_file", path)
	})
}

func TestLocalDirectoryGetSubDirectory(t *testing.T) {
	t.Run("single file", func(t *testing.T) {
		homeDir, err := files.GetLocalDirectoryByPath("/home/")
		require.NoError(t, err)

		subDir, err := homeDir.GetSubDirectory("testfile")
		require.NoError(t, err)

		localPath, err := subDir.GetLocalPath()
		require.NoError(t, err)

		require.EqualValues(t, "/home/testfile", localPath)

	})

	t.Run("subdir and file", func(t *testing.T) {
		homeDir, err := files.GetLocalDirectoryByPath("/home/")
		require.NoError(t, err)

		subDir, err := homeDir.GetSubDirectory("subdir", "another_file")
		require.NoError(t, err)

		localPath, err := subDir.GetLocalPath()
		require.NoError(t, err)
		require.EqualValues(t, "/home/subdir/another_file", localPath)
	})
}

func TestLocalDirectoryParentForBaseClassSet(t *testing.T) {

	tests := []struct {
		testcase string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				dir := files.NewLocalDirectory()
				parent, err := dir.GetParentDirectoryForBaseClass()
				require.NoError(t, err)
				require.NotNil(t, parent)
			},
		)
	}
}

func TestLocalDirectoryCreateFileInDirectoryFromString(t *testing.T) {

	tests := []struct {
		filename []string
		content  string
	}{
		{[]string{"testcase"}, "content"},
		{[]string{"testcase", "test.txt"}, "content"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true
				ctx := getCtx()

				tempDirPath, err := os.MkdirTemp("", "testDir")
				require.NoError(t, err)

				dir, err := files.GetLocalDirectoryByPath(tempDirPath)
				require.NoError(t, err)
				defer dir.Delete(ctx, &filesoptions.DeleteOptions{})

				createdFile, err := dir.CreateFileInDirectoryFromString(tt.content, verbose, tt.filename...)
				require.NoError(t, err)

				localPath, err := dir.GetLocalPath()
				require.NoError(t, err)
				pathElements := []string{localPath}
				pathElements = append(pathElements, tt.filename...)
				expectedFileName := filepath.Join(pathElements...)

				localPath, err = createdFile.GetLocalPath()
				require.NoError(t, err)
				require.EqualValues(t, expectedFileName, localPath)
				require.EqualValues(t, tt.content, createdFile.MustReadAsString())
			},
		)
	}
}

func TestLocalDirectoryGetLocalPathIsAbsolute(t *testing.T) {
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
				localDir, err := files.GetLocalDirectoryByPath(tt.pathToTest)
				require.NoError(t, err)

				localPath, err := localDir.GetLocalPath()
				require.NoError(t, err)

				require.True(t, pathsutils.IsAbsolutePath(localPath))
			},
		)
	}
}

func TestLocalDirectoryWriteStringToFile(t *testing.T) {
	tests := []struct {
		fileName string
		content  string
	}{
		{"a.txt", "testcase"},
		{"b.txt", "testcase\nmultiline"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				tempDirPath, err := os.MkdirTemp("", "testDir")
				require.NoError(t, err)

				testDirectory, err := files.GetLocalDirectoryByPath(tempDirPath)
				require.NoError(t, err)
				defer testDirectory.Delete(ctx, &filesoptions.DeleteOptions{})

				require.False(t, mustutils.Must(testDirectory.FileInDirectoryExists(ctx, tt.fileName)))

				testFile, err := testDirectory.WriteStringToFile(ctx, tt.fileName, tt.content, &filesoptions.WriteOptions{})
				require.NoError(t, err)

				require.True(t, mustutils.Must(testDirectory.FileInDirectoryExists(ctx, tt.fileName)))
				require.EqualValues(t, tt.content, testFile.MustReadAsString())
			},
		)
	}
}

func TestDirectoryListFilesInDirectory(t *testing.T) {
	tests := []struct {
		fileNames     []string
		listOptions   parameteroptions.ListFileOptions
		expectedPaths []string
	}{
		{[]string{"a.go", "b.go"}, parameteroptions.ListFileOptions{ReturnRelativePaths: true}, []string{"a.go", "b.go"}},
		{[]string{"a.go", "a/b.go"}, parameteroptions.ListFileOptions{ReturnRelativePaths: true}, []string{"a.go", "a/b.go"}},
		{[]string{"a.go", "a/b.go", "b.go"}, parameteroptions.ListFileOptions{ReturnRelativePaths: true, ExcludeBasenamePattern: []string{"a.*"}}, []string{"a/b.go", "b.go"}},
		{[]string{"a.go", "a/b.go", "b.go"}, parameteroptions.ListFileOptions{ReturnRelativePaths: true, ExcludeBasenamePattern: []string{"a.*"}}, []string{"a/b.go", "b.go"}},
		{[]string{"a.go", "a/b.go", "b.go"}, parameteroptions.ListFileOptions{ReturnRelativePaths: true, ExcludeBasenamePattern: []string{"b.*"}}, []string{"a.go"}},
		{[]string{"b.go", "a.go"}, parameteroptions.ListFileOptions{ReturnRelativePaths: true}, []string{"a.go", "b.go"}},
		{[]string{"b.go", "a.go", "go.mod", "go.sum"}, parameteroptions.ListFileOptions{ReturnRelativePaths: true}, []string{"a.go", "b.go", "go.mod", "go.sum"}},
		{[]string{"b.go", "a.go", "go.mod", "go.sum"}, parameteroptions.ListFileOptions{ReturnRelativePaths: true, MatchBasenamePattern: []string{".*.go"}}, []string{"a.go", "b.go"}},
		{[]string{"b.go", "a.go", "go.mod", "go.sum"}, parameteroptions.ListFileOptions{ReturnRelativePaths: true, ExcludeBasenamePattern: []string{".*.go"}}, []string{"go.mod", "go.sum"}},
		{[]string{"b.go", "a.go", "go.go", "go.mod", "go.sum"}, parameteroptions.ListFileOptions{ReturnRelativePaths: true, MatchBasenamePattern: []string{"go.*"}, ExcludeBasenamePattern: []string{".*.go", ".*.mod"}}, []string{"go.sum"}},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				tempDirPath, err := os.MkdirTemp("", "tempToTest")
				require.NoError(t, err)

				temporaryDirectory, err := files.GetLocalDirectoryByPath(tempDirPath)
				_, err = temporaryDirectory.CreateFilesInDirectory(ctx, tt.fileNames, &filesoptions.CreateOptions{})
				require.NoError(t, err)

				listedFiles, err := temporaryDirectory.ListFilePaths(ctx, &tt.listOptions)
				require.NoError(t, err)
				require.EqualValues(t, tt.expectedPaths, listedFiles)
			},
		)
	}
}

func TestLocalDirectoryCreate(t *testing.T) {
	tests := []struct {
		subDirPath []string
	}{
		{[]string{"a"}},
		{[]string{"a", "b"}},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				tempDir := getDirectoryToTest("localDirectory")
				subDir, err := tempDir.GetSubDirectory(tt.subDirPath...)
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

// Test if GetPath always returns an absolute value which stays the same even if the current working directory is changed.
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

					directory, err := files.GetLocalDirectoryByPath(tt.path)
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
	tests := []struct {
		testcase string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				path, err := getDirectoryToTest("localDirectory").GetPath()
				require.NoError(t, err)

				tempDir, err := files.GetLocalDirectoryByPath(path)
				require.NoError(t, err)
				defer tempDir.Delete(ctx, &filesoptions.DeleteOptions{})

				isEmptyDir, err := tempDir.IsEmptyDirectory(ctx)
				require.NoError(t, err)

				require.True(t, isEmptyDir)
			},
		)
	}
}

func TestDirectory_CheckExists(t *testing.T) {
	tests := []struct {
		testcase string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				temporaryDirectory := getDirectoryToTest("localDirectory")
				defer temporaryDirectory.Delete(ctx, &filesoptions.DeleteOptions{})

				require.Nil(t, temporaryDirectory.CheckExists(ctx))

				err := temporaryDirectory.Delete(ctx, &filesoptions.DeleteOptions{})
				require.NoError(t, err)

				require.NotNil(t, temporaryDirectory.CheckExists(ctx))
			},
		)
	}
}
