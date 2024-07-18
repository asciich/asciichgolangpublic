package asciichgolangpublic

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocalDirectoryExists(t *testing.T) {

	tests := []struct {
		testcase string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				var directory Directory = TemporaryDirectories().MustCreateEmptyTemporaryDirectory(verbose)
				defer directory.Delete(verbose)

				assert.True(directory.MustExists())

				for i := 0; i < 2; i++ {
					directory.MustDelete(verbose)
					assert.False(directory.MustExists())
				}

				for i := 0; i < 2; i++ {
					directory.MustCreate(verbose)
					assert.True(directory.MustExists())
				}

				for i := 0; i < 2; i++ {
					directory.MustDelete(verbose)
					assert.False(directory.MustExists())
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
			MustFormatAsTestname(tt),
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
			MustFormatAsTestname(tt),
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
			MustFormatAsTestname(tt),
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
			MustFormatAsTestname(tt),
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
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				dir := TemporaryDirectories().MustCreateEmptyTemporaryDirectory(verbose)
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
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				localDir := MustGetLocalDirectoryByPath(tt.pathToTest)

				localPath := localDir.MustGetLocalPath()

				assert.True(Paths().IsAbsolutePath(localPath))
			},
		)
	}
}

func TestLocalDirectoryGetSubDirectories(t *testing.T) {
	tests := []struct {
		testcase string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose = true

				testDirectory := TemporaryDirectories().MustCreateEmptyTemporaryDirectory(verbose)

				testDirectory.MustCreateSubDirectory("test1", verbose)
				test2 := testDirectory.MustCreateSubDirectory("test2", verbose)

				test2.MustCreateSubDirectory("a", verbose)
				test2.MustCreateSubDirectory("b", verbose)
				test2.MustCreateSubDirectory("c", verbose)

				subDirectoryList := testDirectory.MustGetSubDirectories(
					&ListDirectoryOptions{
						Recursive: false,
					},
				)

				assert.Len(subDirectoryList, 2)
				assert.EqualValues(subDirectoryList[0].MustGetBaseName(), "test1")
				assert.EqualValues(subDirectoryList[1].MustGetBaseName(), "test2")
				assert.EqualValues(subDirectoryList[0].MustGetDirName(), testDirectory.MustGetLocalPath())
				assert.EqualValues(subDirectoryList[1].MustGetDirName(), testDirectory.MustGetLocalPath())

				subDirectoryList = testDirectory.MustGetSubDirectories(
					&ListDirectoryOptions{
						Recursive: true,
					},
				)

				assert.Len(subDirectoryList, 5)
				assert.EqualValues(subDirectoryList[0].MustGetBaseName(), "test1")
				assert.EqualValues(subDirectoryList[1].MustGetBaseName(), "test2")
				assert.EqualValues(subDirectoryList[2].MustGetBaseName(), "a")
				assert.EqualValues(subDirectoryList[3].MustGetBaseName(), "b")
				assert.EqualValues(subDirectoryList[4].MustGetBaseName(), "c")
				assert.EqualValues(subDirectoryList[0].MustGetDirName(), testDirectory.MustGetLocalPath())
				assert.EqualValues(subDirectoryList[1].MustGetDirName(), testDirectory.MustGetLocalPath())
				assert.EqualValues(subDirectoryList[2].MustGetDirName(), filepath.Join(testDirectory.MustGetLocalPath(), "test2"))
				assert.EqualValues(subDirectoryList[3].MustGetDirName(), filepath.Join(testDirectory.MustGetLocalPath(), "test2"))
				assert.EqualValues(subDirectoryList[4].MustGetDirName(), filepath.Join(testDirectory.MustGetLocalPath(), "test2"))

			},
		)
	}
}

func TestLocalDirectoryGetGitRepositories(t *testing.T) {
	tests := []struct {
		testcase string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
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
				assert.EqualValues(gitRepos[0].MustGetBaseName(), "test1")
				assert.EqualValues(gitRepos[1].MustGetBaseName(), "c")
				assert.EqualValues(gitRepos[0].MustGetDirName(), testDirectory.MustGetLocalPath())
				assert.EqualValues(gitRepos[1].MustGetDirName(), filepath.Join(testDirectory.MustGetLocalPath(), "test2"))

			},
		)
	}
}
