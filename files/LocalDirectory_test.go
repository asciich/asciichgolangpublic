package files

import (
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pathsutils"
	"github.com/asciich/asciichgolangpublic/testutils"
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
				assert := assert.New(t)

				const verbose bool = true

				var directory Directory = getDirectoryToTest("localDirectory")
				defer directory.Delete(verbose)

				assert.True(directory.MustExists(verbose))

				for i := 0; i < 2; i++ {
					directory.MustDelete(verbose)
					assert.False(directory.MustExists(verbose))
				}

				for i := 0; i < 2; i++ {
					directory.MustCreate(verbose)
					assert.True(directory.MustExists(verbose))
				}

				for i := 0; i < 2; i++ {
					directory.MustDelete(verbose)
					assert.False(directory.MustExists(verbose))
				}

			},
		)
	}
}

func TestLocalDirectoryGetFileInDirectory(t *testing.T) {

	tests := []struct {
		testcase string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				homeDir := MustGetLocalDirectoryByPath("/home/")

				assert.EqualValues(
					"/home/testfile",
					homeDir.MustGetFileInDirectory("testfile").MustGetLocalPath(),
				)

				assert.EqualValues(
					"/home/subdir/another_file",
					homeDir.MustGetFileInDirectory("subdir", "another_file").MustGetLocalPath(),
				)
			},
		)
	}
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
				assert := assert.New(t)

				homeDir := MustGetLocalDirectoryByPath("/home/")

				assert.EqualValues(
					"/home/testfile",
					homeDir.MustGetFilePathInDirectory("testfile"),
				)

				assert.EqualValues(
					"/home/subdir/another_file",
					homeDir.MustGetFilePathInDirectory("subdir", "another_file"),
				)
			},
		)
	}
}

func TestLocalDirectoryGetSubDirectory(t *testing.T) {

	tests := []struct {
		testcase string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				homeDir := MustGetLocalDirectoryByPath("/home/")

				assert.EqualValues(
					"/home/testfile",
					homeDir.MustGetSubDirectory("testfile").MustGetLocalPath(),
				)

				assert.EqualValues(
					"/home/subdir/another_file",
					homeDir.MustGetSubDirectory("subdir", "another_file").MustGetLocalPath(),
				)
			},
		)
	}
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
				assert := assert.New(t)

				dir := NewLocalDirectory()
				assert.NotNil(dir.MustGetParentDirectoryForBaseClass())
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
				assert := assert.New(t)

				const verbose bool = true

				tempDirPath, err := os.MkdirTemp("", "testDir")
				require.Nil(t, err)

				dir := MustGetLocalDirectoryByPath(tempDirPath)
				defer dir.Delete(verbose)

				createdFile := dir.MustCreateFileInDirectoryFromString(tt.content, verbose, tt.filename...)

				pathElements := []string{dir.MustGetLocalPath()}
				pathElements = append(pathElements, tt.filename...)
				expectedFileName := filepath.Join(pathElements...)

				assert.EqualValues(expectedFileName, createdFile.MustGetLocalPath())
				assert.EqualValues(tt.content, createdFile.MustReadAsString())
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
				assert := assert.New(t)

				localDir := MustGetLocalDirectoryByPath(tt.pathToTest)

				localPath := localDir.MustGetLocalPath()

				assert.True(pathsutils.IsAbsolutePath(localPath))
			},
		)
	}
}

/* TODO remove or move
func TestLocalDirectoryGetGitRepositories(t *testing.T) {
	tests := []struct {
		testcase string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose = true

				testDirectory := TemporaryDirectories().MustCreateEmptyTemporaryDirectory(verbose)

				test1 := testDirectory.MustCreateSubDirectory("test1", verbose)
				test2 := testDirectory.MustCreateSubDirectory("test2", verbose)

				test2.MustCreateSubDirectory("a", verbose)
				test2.MustCreateSubDirectory("b", verbose)
				c := test2.MustCreateSubDirectory("c", verbose)

				test1GitRepo := MustGetLocalGitReposioryFromDirectory(test1)
				test1GitRepo.MustInit(&CreateRepositoryOptions{
					Verbose: true,
				})

				cGitRepo := MustGetLocalGitReposioryFromDirectory(c)
				cGitRepo.MustInit(&CreateRepositoryOptions{
					Verbose: true,
				})

				gitRepos := testDirectory.MustGetGitRepositoriesAsLocalGitRepositories(verbose)

				assert.Len(gitRepos, 2)
				assert.EqualValues("test1", gitRepos[0].MustGetBaseName(), "test1")
				assert.EqualValues("c", gitRepos[1].MustGetBaseName(), "c")
				assert.EqualValues(testDirectory.MustGetLocalPath(), gitRepos[0].MustGetDirName())
				assert.EqualValues(filepath.Join(testDirectory.MustGetLocalPath(), "test2"), gitRepos[1].MustGetDirName())
			},
		)
	}
}
*/

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
				assert := assert.New(t)

				const verbose = true

				tempDirPath, err := os.MkdirTemp("", "testDir")
				require.Nil(t, err)

				testDirectory := MustGetLocalDirectoryByPath(tempDirPath)
				defer testDirectory.Delete(verbose)

				assert.False(testDirectory.MustFileInDirectoryExists(verbose, tt.fileName))

				testFile := testDirectory.MustWriteStringToFileInDirectory(tt.content, verbose, tt.fileName)

				assert.True(testDirectory.MustFileInDirectoryExists(verbose, tt.fileName))
				assert.EqualValues(
					tt.content,
					testFile.MustReadAsString(),
				)
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
				assert := assert.New(t)

				const verbose = true
				tt.listOptions.Verbose = verbose

				tempDirPath, err := os.MkdirTemp("", "tempToTest")
				require.Nil(t, err)

				temporaryDirectory := MustGetLocalDirectoryByPath(tempDirPath)
				temporaryDirectory.MustCreateFilesInDirectory(tt.fileNames, verbose)
				listedFiles := temporaryDirectory.MustListFilePaths(&tt.listOptions)
				assert.EqualValues(tt.expectedPaths, listedFiles)
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
				assert := assert.New(t)

				const verbose = true

				tempDir := getDirectoryToTest("localDirectory")
				subDir := tempDir.MustGetSubDirectory(tt.subDirPath...)
				assert.False(subDir.MustExists(verbose))
				subDir.MustCreate(verbose)
				assert.True(subDir.MustExists(verbose))
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
				assert := assert.New(t)

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

					directory := MustGetLocalDirectoryByPath(tt.path)
					path1 = directory.MustGetLocalPath()
					os.Chdir("..")
					path2 = directory.MustGetLocalPath()
				}

				waitGroup.Add(1)
				go testFunction()
				waitGroup.Wait()

				assert.True(pathsutils.IsAbsolutePath(path1))
				assert.True(pathsutils.IsAbsolutePath(path2))

				assert.EqualValues(path1, path2)

				currentPath, err := os.Getwd()
				if err != nil {
					t.Fatalf("%v", err)
				}

				assert.EqualValues(startPath, currentPath)
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
				assert := assert.New(t)

				const verbose = true

				tempDir := MustGetLocalDirectoryByPath(
					getDirectoryToTest("localDirectory").MustGetPath(),
				)
				defer tempDir.Delete(verbose)

				assert.True(tempDir.MustIsEmptyDirectory(verbose))
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
				assert := assert.New(t)

				const verbose = true

				temporaryDirectory := getDirectoryToTest("localDirectory")
				defer temporaryDirectory.Delete(verbose)

				assert.Nil(temporaryDirectory.CheckExists(verbose))

				temporaryDirectory.MustDelete(verbose)

				assert.NotNil(temporaryDirectory.CheckExists(verbose))
			},
		)
	}
}
