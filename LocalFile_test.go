package asciichgolangpublic

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/asciich/asciichgolangpublic/datatypes/pointers"
)

func TestLocalFileImplementsFileInterface(t *testing.T) {
	var file File = MustNewLocalFileByPath("/example/path")
	assert.EqualValues(t, "/example/path", file.MustGetLocalPath())
}

func TestLocalFileIsPathSetOnEmptyFile(t *testing.T) {
	assert.EqualValues(t, false, NewLocalFile().IsPathSet())
}

func TestLocalFileSetAndGetPath(t *testing.T) {
	assert := assert.New(t)

	var localFile = LocalFile{}

	err := localFile.SetPath("testpath")
	assert.EqualValues(nil, err)

	receivedPath, err := localFile.GetPath()
	assert.EqualValues(nil, err)
	assert.True(strings.HasSuffix(receivedPath, "/testpath"))
}


// TODO: Move to File_test.go and test for all File implementations.
func TestLocalFileGetUriAsString(t *testing.T) {

	tests := []struct {
		path        string
		expectedUri string
	}{
		{"/etc/hello.txt", "file:///etc/hello.txt"},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				var file File = MustNewLocalFileByPath(tt.path)

				uri := file.MustGetUriAsString()

				assert.EqualValues(tt.expectedUri, uri)
			},
		)
	}
}

// TODO: Move to File_test.go and test for all File implementations.
func TestLocalFileReadAndWriteAsBytes(t *testing.T) {
	tests := []struct {
		content []byte
	}{
		{[]byte("hello world")},
		{[]byte("hello world\n")},
		{[]byte("hello\nworld\n")},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = false

				var file File = TemporaryFiles().MustCreateEmptyTemporaryFile(verbose)

				assert.EqualValues([]byte{}, file.MustReadAsBytes())

				for i := 0; i < 2; i++ {
					file.MustWriteBytes(tt.content, verbose)

					assert.EqualValues(tt.content, file.MustReadAsBytes())
				}
			},
		)
	}
}

// TODO: Move to File_test.go and test for all File implementations.
func TestLocalFileReadAndWriteAsInt64(t *testing.T) {
	tests := []struct {
		content int64
	}{
		{1},
		{2},
		{3},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = false

				var file File = TemporaryFiles().MustCreateEmptyTemporaryFile(verbose)

				for i := 0; i < 2; i++ {
					file.MustWriteInt64(tt.content, verbose)

					assert.EqualValues(tt.content, file.MustReadAsInt64())
				}
			},
		)
	}
}

// TODO: Move to File_test.go and test for all File implementations.
func TestLocalFileReadAndWriteAsString(t *testing.T) {
	tests := []struct {
		content string
	}{
		{"hello world"},
		{"hello world\n"},
		{"hello\nworld\n"},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = false

				var file File = TemporaryFiles().MustCreateEmptyTemporaryFile(verbose)

				assert.EqualValues("", file.MustReadAsString())

				for i := 0; i < 2; i++ {
					file.MustWriteString(tt.content, verbose)

					assert.EqualValues(tt.content, file.MustReadAsString())
				}
			},
		)
	}
}

// TODO: Move to File_test.go and test for all File implementations.
func TestLocalFileGetBaseName(t *testing.T) {
	tests := []struct {
		path             string
		expectedBaseName string
	}{
		{"hello", "hello"},
		{"this/hello", "hello"},
		{"/this/hello", "hello"},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				var file File = MustGetLocalFileByPath(tt.path)

				assert.EqualValues(tt.expectedBaseName, file.MustGetBaseName())
			},
		)
	}
}

// TODO: Move to File_test.go and test for all File implementations.
func TestLocalFileGetSha256Sum(t *testing.T) {
	tests := []struct {
		input            string
		expectedChecksum string
	}{
		{"", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		{"hello world", "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				temporaryFile := TemporaryFiles().MustCreateFromString(tt.input, verbose)
				defer temporaryFile.Delete(verbose)

				assert.EqualValues(
					tt.expectedChecksum,
					temporaryFile.MustGetSha256Sum(),
				)
			},
		)
	}
}

// TODO: Move to File_test.go and test for all File implementations.
func TestLocalFileIsMatchingSha256Sum(t *testing.T) {
	tests := []struct {
		input              string
		sha256sum          string
		expectedIsMatching bool
	}{
		{"", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855", true},
		{"hello world", "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9", true},
		{"", "aaaaaae3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855", false},
		{"hello world", "aaaaaab94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9", false},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				temporaryFile := TemporaryFiles().MustCreateFromString(tt.input, verbose)
				defer temporaryFile.Delete(verbose)

				assert.EqualValues(
					tt.expectedIsMatching,
					temporaryFile.MustIsMatchingSha256Sum(tt.sha256sum),
				)
			},
		)
	}
}

// TODO: Move to File_test.go and test for all File implementations.
func TestLocalFileGetParentDirectory(t *testing.T) {
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

				temporaryDir := TemporaryDirectories().MustCreateEmptyTemporaryDirectory(verbose)
				defer temporaryDir.Delete(verbose)

				temporaryFile := temporaryDir.MustCreateFileInDirectory(verbose, "test.txt")
				parentDir := temporaryFile.MustGetParentDirectory()

				assert.EqualValues(
					temporaryDir.MustGetLocalPath(),
					parentDir.MustGetLocalPath(),
				)
			},
		)
	}
}

// TODO: Move to File_test.go and test for all File implementations.
func TestLocalFileIsContentEqualByComparingSha256Sum(t *testing.T) {
	tests := []struct {
		contentFile1    string
		contentFile2    string
		expectedIsEqual bool
	}{
		{"", "", true},
		{"testcase", "testcase", true},
		{"testcase1", "testcase", false},
		{"testcase1", "testcase2", false},
		{"testcase", "testcase3", false},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				tempFile1 := TemporaryFiles().MustCreateFromString(tt.contentFile1, verbose)
				defer tempFile1.Delete(verbose)

				tempFile2 := TemporaryFiles().MustCreateFromString(tt.contentFile2, verbose)
				defer tempFile2.Delete(verbose)

				assert.EqualValues(tt.expectedIsEqual, tempFile1.MustIsContentEqualByComparingSha256Sum(tempFile2, verbose))
				assert.EqualValues(tt.expectedIsEqual, tempFile2.MustIsContentEqualByComparingSha256Sum(tempFile1, verbose))
			},
		)
	}
}

// TODO: Move to File_test.go and test for all File implementations.
func TestLocalFileGetLocalPathIsAbsolute(t *testing.T) {
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

				localFile := MustGetLocalFileByPath(tt.pathToTest)

				localPath := localFile.MustGetLocalPath()

				assert.True(Paths().IsAbsolutePath(localPath))
			},
		)
	}
}

// TODO: Move to File_test.go and test for all File implementations.
func TestFileGetTextBlocksGolangWithCommentAboveFunction(t *testing.T) {
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

				testData := "package main\n"
				testData += "\n"
				testData += "//this comment\n"
				testData += "func belongsToMe() (err error) {\n"
				testData += "\n" // ensure empty lines do not split the function block
				testData += "\treturn nil\n"
				testData += "}\n"

				testFile := TemporaryFiles().MustCreateFromString(testData, verbose)
				blocks := testFile.MustGetTextBlocks(verbose)

				assert.Len(blocks, 2)
			},
		)
	}
}

// TODO: Move to File_test.go and test for all File implementations.
func TestFileGetTextBlocksYamlWithoutLeadingThreeMinuses(t *testing.T) {
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

				testData := "a: b\n"
				testData += "\n"
				testData += "c: d\n"

				testFile := TemporaryFiles().MustCreateFromString(testData, verbose)
				blocks := testFile.MustGetTextBlocks(verbose)

				assert.Len(blocks, 2)
			},
		)
	}
}

// TODO: Move to File_test.go and test for all File implementations.
func TestFileGetTextBlocksYamlWithLeadingThreeMinuses(t *testing.T) {
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

				testData := "---\n"
				testData += "a: b\n"
				testData += "\n"
				testData += "c: d\n"

				testFile := TemporaryFiles().MustCreateFromString(testData, verbose)
				blocks := testFile.MustGetTextBlocks(verbose)

				assert.Len(blocks, 3)
			},
		)
	}
}

// TODO: Move to File_test.go and test for all File implementations.
func TestLocalFileGetDeepCopy(t *testing.T) {
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

				var testFile File = TemporaryFiles().MustCreateEmptyTemporaryFile(verbose)
				localTestFile := MustGetLocalFileByFile(testFile)

				copy := testFile.GetDeepCopy()
				assert.EqualValues(
					testFile.MustGetLocalPath(),
					copy.MustGetLocalPath(),
				)
				localCopy := MustGetLocalFileByFile(copy)

				assert.False(pointers.MustPointersEqual(
					localTestFile.MustGetParentFileForBaseClassAsLocalFile(),
					localCopy.MustGetParentFileForBaseClassAsLocalFile(),
				))

				assert.True(testFile.MustExists(verbose))
				assert.True(copy.MustExists(verbose))
			},
		)
	}
}

// TODO: Move to File_test.go and test for all File implementations.
func TestFileReplaceLineAfterLine(t *testing.T) {
	tests := []struct {
		input                     string
		lineToSearch              string
		replaceLineAfterFoundWith string
		expectedContent           string
		expectedChanged           bool
	}{
		{"a\nb\nc\n", "a", "d", "a\nd\nc\n", true},
		{"a\nb\nc\n", "x", "d", "a\nb\nc\n", false},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				testFile := TemporaryFiles().MustCreateFromString(tt.input, verbose)
				defer testFile.Delete(verbose)

				changeSummary := testFile.MustReplaceLineAfterLine(tt.lineToSearch, tt.replaceLineAfterFoundWith, verbose)

				content := testFile.MustReadAsString()
				assert.EqualValues(tt.expectedContent, content)
				assert.EqualValues(tt.expectedChanged, changeSummary.IsChanged())
			},
		)
	}
}

// TODO: Move to File_test.go and test for all File implementations.
// Test if GetPath always returns an absolute value which stays the same even if the current working directory is changed.
func TestLocalFileGetPathReturnsAbsoluteValue(t *testing.T) {
	tests := []struct {
		path string
	}{
		{"test.txt"},
		{"./test.txt"},
		{"../test.txt"},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				startPath, err := os.Getwd()
				if err != nil {
					LogFatalWithTrace(err)
				}

				var path1 string
				var path2 string

				var waitGroup sync.WaitGroup

				testFunction := func() {
					defer os.Chdir(startPath)
					defer waitGroup.Done()

					file := MustGetLocalFileByPath(tt.path)
					path1 = file.MustGetPath()
					os.Chdir("..")
					path2 = file.MustGetPath()
				}

				waitGroup.Add(1)
				go testFunction()
				waitGroup.Wait()

				assert.True(Paths().IsAbsolutePath(path1))
				assert.True(Paths().IsAbsolutePath(path2))

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

// TODO: Move to File_test.go and test for all File implementations.
func TestLocalFileSortBlocksInFile(t *testing.T) {
	type TestCase struct {
		testDataDir string
	}

	tests := []TestCase{}

	testDataDirectory := MustGetLocalGitRepositoryByPath(".").MustGetSubDirectory("testdata", "File", "SortBlocksInFile")
	for _, testDirectory := range testDataDirectory.MustListSubDirectories(&ListDirectoryOptions{Recursive: false}) {
		tests = append(tests, TestCase{testDirectory.MustGetLocalPath()})
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose = true

				testDataDir := MustGetLocalDirectoryByPath(tt.testDataDir)
				testFile := testDataDir.MustCopyFileToTemporaryFile(verbose, "input")
				expectedFile := testDataDir.MustGetFileInDirectory("expectedOutput")
				testFile.MustSortBlocksInFile(verbose)

				sortedChecksum := testFile.MustGetSha256Sum()
				expectedChecksum := expectedFile.MustGetSha256Sum()

				if os.Getenv("UPDATE_EXPECTED") == "1" {
					testFile.MustCopyToFile(expectedFile, verbose)
				}

				assert.EqualValues(expectedChecksum, sortedChecksum)
			},
		)
	}
}

// TODO: Move to File_test.go and test for all File implementations.
func TestLocalFileGetLastCharAsString(t *testing.T) {
	tests := []struct {
		content  string
		lastChar string
	}{
		{" ", " "},
		{"a", "a"},
		{" \n", "\n"},
		{" \nb", "b"},
		{" \nb\n", "\n"},
		{"a\n", "\n"},
		{"\n", "\n"},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose = true

				testFile := TemporaryFiles().MustCreateFromString(tt.content, verbose)
				defer testFile.Delete(verbose)

				lastChar := testFile.MustReadLastCharAsString()

				assert.EqualValues(tt.lastChar, lastChar)
			},
		)
	}
}

// TODO: Move to File_test.go and test for all File implementations.
func TestLocalFileGetAsFloat64(t *testing.T) {
	tests := []struct {
		content       string
		expectedFloat float64
	}{
		{"0", 0.0},
		{"0.", 0.0},
		{"0.0", 0.0},
		{"0.1", 0.1},
		{"0.10", 0.1},
		{"-3", -3.0},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose = true

				testFile := TemporaryFiles().MustCreateFromString(tt.content, verbose)
				defer testFile.Delete(verbose)

				readFloat := testFile.MustReadAsFloat64()

				assert.EqualValues(tt.expectedFloat, readFloat)
			},
		)
	}
}

// TODO: Move to File_test.go and test for all File implementations.
func TestFileGetAsInt64(t *testing.T) {
	tests := []struct {
		content     string
		expectedInt int64
	}{
		{"0", 0},
		{"1", 1},
		{"1\n", 1},
		{"10\n", 10},
		{" 10\n", 10},
		{" -110\n", -110},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose = true

				testFile := TemporaryFiles().MustCreateFromString(tt.content, verbose)
				defer testFile.Delete(verbose)

				readInt64 := testFile.MustReadAsInt64()

				assert.EqualValues(tt.expectedInt, readInt64)
			},
		)
	}
}

// TODO: Move to File_test.go and test for all File implementations.
func TestFileGetAsInt(t *testing.T) {
	tests := []struct {
		content     string
		expectedInt int
	}{
		{"0", 0},
		{"1", 1},
		{"1\n", 1},
		{"10\n", 10},
		{" 10\n", 10},
		{" -110\n", -110},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose = true

				testFile := TemporaryFiles().MustCreateFromString(tt.content, verbose)
				defer testFile.Delete(verbose)

				readInt64 := testFile.MustReadAsInt()

				assert.EqualValues(tt.expectedInt, readInt64)
			},
		)
	}
}

// TODO: Move to File_test.go and test for all File implementations.
func TestFileGetParentDirectoryPath(t *testing.T) {
	tests := []struct {
		inputPath          string
		expectedParentPath string
	}{
		{"/abc", "/"},
		{"/abc/d", "/abc"},
		{"/abc/d.go", "/abc"},
		{"/abc/d.txt", "/abc"},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				testFile := MustGetLocalFileByPath(tt.inputPath)
				parentPath := testFile.MustGetParentDirectoryPath()
				assert.EqualValues(tt.expectedParentPath, parentPath)
			},
		)
	}
}

// TODO: Move to File_test.go and test for all File implementations.
func TestFileIsPgpEncrypted_Case1_unencrypted(t *testing.T) {
	tests := []struct {
		unencrypted string
	}{
		{""},
		{"\n"},
		{"---"},
		{"testcase"},
		{"testcase\n"},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				temporaryFile := TemporaryFiles().MustCreateFromString(tt.unencrypted, verbose)
				defer temporaryFile.Delete(verbose)

				assert.False(temporaryFile.MustIsPgpEncrypted(verbose))
			},
		)
	}
}

// TODO: Move to File_test.go and test for all File implementations.
func TestFileIsPgpEncrypted_Case2_encryptedBinary(t *testing.T) {
	tests := []struct {
		testcase string
	}{
		{"testcase\n"},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				temporaryFile := TemporaryFiles().MustCreateEmptyTemporaryFile(verbose)
				defer temporaryFile.Delete(verbose)

				createCommand := []string{
					"bash",
					"-c",
					fmt.Sprintf(
						"exec 3<<<$(echo hallo) ; echo test | gpg --batch --symmetric --passphrase-fd=3 > '%s'",
						temporaryFile.MustGetLocalPath(),
					),
				}
				Bash().MustRunCommand(
					&RunCommandOptions{
						Command: createCommand,
						Verbose: verbose,
					},
				)

				assert.True(temporaryFile.MustIsPgpEncrypted(verbose))
			},
		)
	}
}

// TODO: Move to File_test.go and test for all File implementations.
func TestFileIsPgpEncrypted_Case3_encryptedAsciiArmor(t *testing.T) {
	tests := []struct {
		testcase string
	}{
		{"testcase\n"},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				temporaryFile := TemporaryFiles().MustCreateEmptyTemporaryFile(verbose)
				defer temporaryFile.Delete(verbose)

				createCommand := []string{
					"bash",
					"-c",
					fmt.Sprintf(
						"exec 3<<<$(echo hallo) ; echo test | gpg --batch --symmetric --passphrase-fd=3 -a > '%s'",
						temporaryFile.MustGetLocalPath(),
					),
				}
				Bash().MustRunCommand(
					&RunCommandOptions{
						Command: createCommand,
						Verbose: verbose,
					},
				)

				assert.True(temporaryFile.MustIsPgpEncrypted(verbose))
			},
		)
	}
}

// TODO: Move to File_test.go and test for all File implementations.
func TestFileGetMimeTypeOfEmptyFile(t *testing.T) {
	tests := []struct {
		testcase string
	}{
		{"testcase\n"},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				temporaryFile := TemporaryFiles().MustCreateEmptyTemporaryFile(verbose)
				expectedMimeType := "inode/x-empty"

				mimeType := temporaryFile.MustGetMimeType(verbose)

				assert.EqualValues(expectedMimeType, mimeType)
			},
		)
	}
}

// TODO: Move to File_test.go and test for all File implementations.
func TestFileGetCreationDateByFileName(t *testing.T) {
	tests := []struct {
		filename string
		expected time.Time
	}{
		{"20231121_140424", time.Date(2023, 11, 21, 14, 04, 24, 0, time.UTC)},
		{"20231121-140424", time.Date(2023, 11, 21, 14, 04, 24, 0, time.UTC)},
		{"20231121-140424thisisignored", time.Date(2023, 11, 21, 14, 04, 24, 0, time.UTC)},
		{"20231121_140424.jpg", time.Date(2023, 11, 21, 14, 04, 24, 0, time.UTC)},
		{"20231121-140424.jpg", time.Date(2023, 11, 21, 14, 04, 24, 0, time.UTC)},
		{"20231121-140424thisisignored.jpg", time.Date(2023, 11, 21, 14, 04, 24, 0, time.UTC)},
		{"signal-2023-04-05-19-47-40-414-1.jpg", time.Date(2023, 04, 05, 19, 47, 40, 0, time.UTC)},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				temporaryDir := TemporaryDirectories().MustCreateEmptyTemporaryDirectory(verbose)
				defer temporaryDir.Delete(verbose)

				file := temporaryDir.MustWriteStringToFileInDirectory("content", verbose, tt.filename)
				readDate := file.MustGetCreationDateByFileName(verbose)

				assert.EqualValues(
					tt.expected,
					*readDate,
				)
			},
		)
	}
}

// TODO: Move to File_test.go and test for all File implementations.
func TestFileHasYYYYmmdd_HHMMSSPrefix(t *testing.T) {
	tests := []struct {
		filename          string
		expectedHasPrefix bool
	}{
		{"20231121_140424", true},
		{"20231121_140424.jpg", true},
		{"20231121_140424_test.jpg", true},
		{"a20231121_140424_test.jpg", false},
		{"a.jpg", false},
		{"a", false},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				temporaryDir := TemporaryDirectories().MustCreateEmptyTemporaryDirectory(verbose)
				defer temporaryDir.Delete(verbose)

				file := temporaryDir.MustWriteStringToFileInDirectory("content", verbose, tt.filename)
				assert.EqualValues(tt.expectedHasPrefix, file.MustIsYYYYmmdd_HHMMSSPrefix())
			},
		)
	}
}

// TODO: Move to File_test.go and test for all File implementations.
func TestFileGetSizeBytes(t *testing.T) {
	tests := []struct {
		content      []byte
		expectedSize int64
	}{
		{[]byte{}, 0},
		{[]byte("helloWorld"), 10},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				testFile := TemporaryFiles().MustCreateFromBytes(tt.content, verbose)
				sizeBytes := testFile.MustGetSizeBytes()

				assert.EqualValues(tt.expectedSize, sizeBytes)
			},
		)
	}
}

// TODO: Move to File_test.go and test for all File implementations.
func TestFileEnsureEndsWithLineBreakOnEmptyFile(t *testing.T) {
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

				emptyFile := TemporaryFiles().MustCreateEmptyTemporaryFile(verbose)
				emptyFile.MustEnsureEndsWithLineBreak(verbose)

				assert.EqualValues("\n", emptyFile.MustReadLastCharAsString())
			},
		)
	}
}

// TODO: Move to File_test.go and test for all File implementations.
func TestFileEnsureEndsWithLineBreakOnNonExitistingFile(t *testing.T) {
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

				nonExistingFile := TemporaryFiles().MustCreateEmptyTemporaryFile(verbose)
				nonExistingFile.MustDelete(verbose)

				assert.False(nonExistingFile.MustExists(verbose))

				nonExistingFile.MustEnsureEndsWithLineBreak(verbose)

				assert.EqualValues("\n", nonExistingFile.MustReadLastCharAsString())
			},
		)
	}
}

// TODO: Move to File_test.go and test for all File implementations.
func TestFileTrimSpacesAtBeginningOfFile(t *testing.T) {
	tests := []struct {
		input           string
		expectedContent string
	}{
		{"testcase", "testcase"},
		{" testcase", "testcase"},
		{"\ntestcase", "testcase"},
		{"\ttestcase", "testcase"},
		{"  testcase", "testcase"},
		{"\n testcase", "testcase"},
		{"\t testcase", "testcase"},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				testFile := TemporaryFiles().MustCreateFromString(tt.input, verbose)
				defer testFile.Delete(verbose)

				testFile.MustTrimSpacesAtBeginningOfFile(verbose)

				content := testFile.MustReadAsString()
				assert.EqualValues(tt.expectedContent, content)
			},
		)
	}
}

// TODO: Move to File_test.go and test for all File implementations.
func TestFileReplaceBetweenMarkers(t *testing.T) {
	type TestCase struct {
		testDataDir string
	}

	tests := []TestCase{}

	testDataDirectory := MustGetLocalGitRepositoryByPath(".").MustGetSubDirectory("testdata", "File", "ReplaceBetweenMarkers")
	for _, testDirectory := range testDataDirectory.MustListSubDirectories(&ListDirectoryOptions{Recursive: false}) {
		tests = append(tests, TestCase{testDirectory.MustGetLocalPath()})
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				sourceDir := MustGetLocalDirectoryByPath(filepath.Join(tt.testDataDir, "input"))
				expectedDir := MustGetLocalDirectoryByPath(filepath.Join(tt.testDataDir, "expectedOutput"))

				tempDir := TemporaryDirectories().MustCreateEmptyTemporaryDirectory(verbose)

				sourceDir.MustCopyContentToLocalDirectory(
					tempDir,
					verbose,
				)

				tempDir.MustReplaceBetweenMarkers(verbose)

				listOptions := NewListFileOptions()
				listOptions.NonRecursive = false
				listOptions.ReturnRelativePaths = true
				expectedFilePaths := expectedDir.MustListFilePaths(listOptions)
				currentFilePaths := tempDir.MustListFilePaths(listOptions)

				assert.EqualValues(expectedFilePaths, currentFilePaths)

				for _, toCheck := range currentFilePaths {
					currentFile := tempDir.MustGetFileInDirectory(toCheck)
					expectedFile := expectedDir.MustGetFileInDirectory(toCheck)

					currentChecksum := currentFile.MustGetSha256Sum()
					expectedChecksum := expectedFile.MustGetSha256Sum()

					if currentChecksum != expectedChecksum {

						if os.Getenv("UPDATE_EXPECTED") == "1" {
							currentFile.CopyToFile(expectedFile, verbose)
						}
					}

					assert.EqualValuesf(
						expectedChecksum,
						currentChecksum,
						"File '%s' missmatch: '%s' '%s'",
						toCheck,
						currentFile.MustGetLocalPath(),
						expectedFile.MustGetLocalPath(),
					)
				}
			},
		)
	}
}

// TODO: Move to File_test.go and test for all File implementations.
func TestLocalFileGetNumberOfNonEmptyLines(t *testing.T) {
	tests := []struct {
		content               string
		expectedNonEmptyLines int
	}{
		{"", 0},
		{"testcase", 1},
		{"testcase\n", 1},
		{"testcase\n\n", 1},
		{"testcase\n\na", 2},
	}
	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				localFile := TemporaryFiles().MustCreateFromString(tt.content, verbose)
				assert.EqualValues(
					tt.expectedNonEmptyLines,
					localFile.MustGetNumberOfNonEmptyLines())
			},
		)
	}
}
