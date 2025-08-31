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
		homeDir := files.MustGetLocalDirectoryByPath("/home/")

		localPath, err := homeDir.MustGetFileInDirectory("testfile").GetLocalPath()
		require.NoError(t, err)

		require.EqualValues(t, "/home/testfile", localPath)
	})

	t.Run("with sub path", func(t *testing.T) {
		homeDir := files.MustGetLocalDirectoryByPath("/home/")

		localPath, err := homeDir.MustGetFileInDirectory("subdir", "another_file").GetLocalPath()
		require.NoError(t, err)

		require.EqualValues(t, "/home/subdir/another_file", localPath)
	})
}

func TestLocalDirectoryGetFilePathInDirectory(t *testing.T) {
	tests := []struct {
		testcase string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				homeDir := files.MustGetLocalDirectoryByPath("/home/")

				require.EqualValues(
					"/home/testfile",
					homeDir.MustGetFilePathInDirectory("testfile"),
				)

				require.EqualValues(
					"/home/subdir/another_file",
					homeDir.MustGetFilePathInDirectory("subdir", "another_file"),
				)
			},
		)
	}
}

func TestLocalDirectoryGetSubDirectory(t *testing.T) {
	t.Run("single file", func(t *testing.T) {
		homeDir := files.MustGetLocalDirectoryByPath("/home/")
		localPath, err := homeDir.MustGetSubDirectory("testfile").GetLocalPath()
		require.NoError(t, err)
		require.EqualValues(t, "/home/testfile", localPath)

	})

	t.Run("subdir and file", func(t *testing.T) {
		homeDir := files.MustGetLocalDirectoryByPath("/home/")
		localPath, err := homeDir.MustGetSubDirectory("subdir", "another_file").GetLocalPath()
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
				require := require.New(t)

				dir := files.NewLocalDirectory()
				require.NotNil(dir.MustGetParentDirectoryForBaseClass())
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
				require.Nil(t, err)

				dir := files.MustGetLocalDirectoryByPath(tempDirPath)
				defer dir.Delete(ctx, &filesoptions.DeleteOptions{})

				createdFile, err := dir.CreateFileInDirectoryFromString(tt.content, verbose, tt.filename...)
				require.NoError(t, err)

				pathElements := []string{dir.MustGetLocalPath()}
				pathElements = append(pathElements, tt.filename...)
				expectedFileName := filepath.Join(pathElements...)

				localPath, err := createdFile.GetLocalPath()
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
				require := require.New(t)

				localDir := files.MustGetLocalDirectoryByPath(tt.pathToTest)

				localPath := localDir.MustGetLocalPath()

				require.True(pathsutils.IsAbsolutePath(localPath))
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
				const verbose = true

				tempDirPath, err := os.MkdirTemp("", "testDir")
				require.Nil(t, err)

				testDirectory := files.MustGetLocalDirectoryByPath(tempDirPath)
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

				temporaryDirectory := files.MustGetLocalDirectoryByPath(tempDirPath)
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
				require := require.New(t)

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

					directory := files.MustGetLocalDirectoryByPath(tt.path)
					path1 = directory.MustGetLocalPath()
					os.Chdir("..")
					path2 = directory.MustGetLocalPath()
				}

				waitGroup.Add(1)
				go testFunction()
				waitGroup.Wait()

				require.True(pathsutils.IsAbsolutePath(path1))
				require.True(pathsutils.IsAbsolutePath(path2))

				require.EqualValues(path1, path2)

				currentPath, err := os.Getwd()
				if err != nil {
					t.Fatalf("%v", err)
				}

				require.EqualValues(startPath, currentPath)
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
				const verbose = true
				ctx := getCtx()

				path, err := getDirectoryToTest("localDirectory").GetPath()
				require.NoError(t, err)

				tempDir := files.MustGetLocalDirectoryByPath(path)
				defer tempDir.Delete(ctx, &filesoptions.DeleteOptions{})

				require.True(t, tempDir.MustIsEmptyDirectory(verbose))
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
