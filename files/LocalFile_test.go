package files

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/datatypes/pointersutils"
	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/pathsutils"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func TestLocalFileImplementsFileInterface(t *testing.T) {
	var file File = MustNewLocalFileByPath("/example/path")
	require.EqualValues(t, "/example/path", file.MustGetLocalPath())
}

func TestLocalFileIsPathSetOnEmptyFile(t *testing.T) {
	require.EqualValues(t, false, NewLocalFile().IsPathSet())
}

func TestLocalFileSetAndGetPath(t *testing.T) {
	require := require.New(t)

	var localFile = LocalFile{}

	err := localFile.SetPath("testpath")
	require.EqualValues(nil, err)

	receivedPath, err := localFile.GetPath()
	require.EqualValues(nil, err)
	require.True(strings.HasSuffix(receivedPath, "/testpath"))
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				var file File = MustNewLocalFileByPath(tt.path)

				uri := file.MustGetUriAsString()

				require.EqualValues(tt.expectedUri, uri)
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = false

				var file File = getFileToTest("localFile")

				require.EqualValues([]byte{}, file.MustReadAsBytes())

				for i := 0; i < 2; i++ {
					file.MustWriteBytes(tt.content, verbose)

					require.EqualValues(tt.content, file.MustReadAsBytes())
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = false

				var file File = getFileToTest("localFile")

				for i := 0; i < 2; i++ {
					file.MustWriteInt64(tt.content, verbose)

					require.EqualValues(tt.content, file.MustReadAsInt64())
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = false

				var file File = getFileToTest("localFile")

				require.EqualValues("", file.MustReadAsString())

				for i := 0; i < 2; i++ {
					file.MustWriteString(tt.content, verbose)

					require.EqualValues(tt.content, file.MustReadAsString())
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				var file File = MustGetLocalFileByPath(tt.path)

				require.EqualValues(tt.expectedBaseName, file.MustGetBaseName())
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				temporaryFile := getFileToTest("localFile")
				temporaryFile.MustWriteString(tt.input, verbose)
				defer temporaryFile.Delete(verbose)

				require.EqualValues(
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				temporaryFile := getFileToTest("localFile")
				temporaryFile.MustWriteString(tt.input, verbose)
				defer temporaryFile.Delete(verbose)

				require.EqualValues(
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				temporaryDir := getDirectoryToTest("localDirectory")
				defer temporaryDir.Delete(verbose)

				temporaryFile := temporaryDir.MustCreateFileInDirectory(verbose, "test.txt")
				parentDir := temporaryFile.MustGetParentDirectory()

				require.EqualValues(
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				tempFile1 := getFileToTest("localFile")
				tempFile1.MustWriteString(tt.contentFile1, verbose)
				defer tempFile1.Delete(verbose)

				tempFile2 := getFileToTest("localFile")
				tempFile2.MustWriteString(tt.contentFile2, verbose)
				defer tempFile2.Delete(verbose)

				require.EqualValues(tt.expectedIsEqual, tempFile1.MustIsContentEqualByComparingSha256Sum(tempFile2, verbose))
				require.EqualValues(tt.expectedIsEqual, tempFile2.MustIsContentEqualByComparingSha256Sum(tempFile1, verbose))
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				localFile := MustGetLocalFileByPath(tt.pathToTest)

				localPath := localFile.MustGetLocalPath()

				require.True(pathsutils.IsAbsolutePath(localPath))
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				testData := "package main\n"
				testData += "\n"
				testData += "//this comment\n"
				testData += "func belongsToMe() (err error) {\n"
				testData += "\n" // ensure empty lines do not split the function block
				testData += "\treturn nil\n"
				testData += "}\n"

				testFile := getFileToTest("localFile")
				testFile.MustWriteString(testData, verbose)
				blocks := testFile.MustGetTextBlocks(verbose)

				require.Len(blocks, 2)
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				testData := "a: b\n"
				testData += "\n"
				testData += "c: d\n"

				testFile := getFileToTest("localFile")
				testFile.WriteString(testData, verbose)
				blocks := testFile.MustGetTextBlocks(verbose)

				require.Len(blocks, 2)
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				testData := "---\n"
				testData += "a: b\n"
				testData += "\n"
				testData += "c: d\n"

				testFile := getFileToTest("localFile")
				testFile.WriteString(testData, verbose)
				blocks := testFile.MustGetTextBlocks(verbose)

				require.Len(blocks, 3)
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				var testFile File = getFileToTest("localFile")
				localTestFile := MustGetLocalFileByFile(testFile)

				copy := testFile.GetDeepCopy()
				require.EqualValues(
					testFile.MustGetLocalPath(),
					copy.MustGetLocalPath(),
				)
				localCopy := MustGetLocalFileByFile(copy)

				require.False(pointersutils.MustPointersEqual(
					localTestFile.MustGetParentFileForBaseClassAsLocalFile(),
					localCopy.MustGetParentFileForBaseClassAsLocalFile(),
				))

				require.True(testFile.MustExists(verbose))
				require.True(copy.MustExists(verbose))
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				testFile := getFileToTest("localFile")
				testFile.WriteString(tt.input, verbose)
				defer testFile.Delete(verbose)

				changeSummary := testFile.MustReplaceLineAfterLine(tt.lineToSearch, tt.replaceLineAfterFoundWith, verbose)

				content := testFile.MustReadAsString()
				require.EqualValues(tt.expectedContent, content)
				require.EqualValues(tt.expectedChanged, changeSummary.IsChanged())
			},
		)
	}
}

// TODO: Move to File_test.go and test for all File implementations.
// Test if GetPath always returns an absolute value which stays the same even if the current working directory is changed.
func TestLocalFile_GetPathReturnsAbsoluteValue(t *testing.T) {
	tests := []struct {
		path string
	}{
		{"test.txt"},
		{"./test.txt"},
		{"../test.txt"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				startPath, err := os.Getwd()
				if err != nil {
					logging.LogFatalWithTrace(err)
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

func getRepoRootDir(ctx context.Context, t *testing.T) (repoRoot Directory) {
	path, err := commandexecutor.Bash().RunOneLinerAndGetStdoutAsString(ctx, "git rev-parse --show-toplevel")
	require.NoError(t, err)
	path = strings.TrimSpace(path)

	return MustGetLocalDirectoryByPath(path)
}

// TODO: Move to File_test.go and test for all File implementations.
func TestLocalFileSortBlocksInFile(t *testing.T) {
	type TestCase struct {
		testDataDir string
	}

	tests := []TestCase{}
	ctx := getCtx()

	testDataDirectory := getRepoRootDir(ctx, t).MustGetSubDirectory("testdata", "File", "SortBlocksInFile")
	for _, testDirectory := range mustutils.Must(testDataDirectory.ListSubDirectories(&parameteroptions.ListDirectoryOptions{Recursive: false})) {
		tests = append(tests, TestCase{testDirectory.MustGetLocalPath()})
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose = true

				testDataDir := MustGetLocalDirectoryByPath(tt.testDataDir)
				testInput := testDataDir.MustReadFileInDirectoryAsString("input")

				testFile := getFileToTest("localFile")
				testFile.MustWriteString(testInput, verbose)

				expectedFile := testDataDir.MustGetFileInDirectory("expectedOutput")
				testFile.MustSortBlocksInFile(verbose)

				sortedChecksum := testFile.MustGetSha256Sum()
				expectedChecksum := expectedFile.MustGetSha256Sum()

				if os.Getenv("UPDATE_EXPECTED") == "1" {
					testFile.MustCopyToFile(expectedFile, verbose)
				}

				require.EqualValues(expectedChecksum, sortedChecksum)
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose = true

				testFile := getFileToTest("localFile")
				testFile.WriteString(tt.content, verbose)
				defer testFile.Delete(verbose)

				lastChar := testFile.MustReadLastCharAsString()

				require.EqualValues(tt.lastChar, lastChar)
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose = true

				testFile := getFileToTest("localFile")
				testFile.WriteString(tt.content, verbose)
				defer testFile.Delete(verbose)

				readFloat := testFile.MustReadAsFloat64()

				require.EqualValues(tt.expectedFloat, readFloat)
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose = true

				testFile := getFileToTest("localFile")
				testFile.WriteString(tt.content, verbose)
				defer testFile.Delete(verbose)

				readInt64 := testFile.MustReadAsInt64()

				require.EqualValues(tt.expectedInt, readInt64)
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose = true

				testFile := getFileToTest("localFile")
				testFile.WriteString(tt.content, verbose)
				defer testFile.Delete(verbose)

				readInt64 := testFile.MustReadAsInt()

				require.EqualValues(tt.expectedInt, readInt64)
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				testFile := MustGetLocalFileByPath(tt.inputPath)
				parentPath := testFile.MustGetParentDirectoryPath()
				require.EqualValues(tt.expectedParentPath, parentPath)
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				testFile := getFileToTest("localFile")
				testFile.WriteString(tt.unencrypted, verbose)
				defer testFile.Delete(verbose)

				require.False(testFile.MustIsPgpEncrypted(verbose))
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				ctx := getCtx()
				const verbose bool = true

				temporaryFile := getFileToTest("localFile")
				defer temporaryFile.Delete(verbose)

				createCommand := []string{
					"bash",
					"-c",
					fmt.Sprintf(
						"exec 3<<<$(echo hallo) ; echo test | gpg --batch --symmetric --passphrase-fd=3 > '%s'",
						temporaryFile.MustGetLocalPath(),
					),
				}
				commandexecutor.Bash().RunCommand(
					ctx,
					&parameteroptions.RunCommandOptions{
						Command: createCommand,
					},
				)

				require.True(temporaryFile.MustIsPgpEncrypted(verbose))
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()
				const verbose bool = true

				temporaryFile := getFileToTest("localFile")
				defer temporaryFile.Delete(verbose)

				createCommand := []string{
					"bash",
					"-c",
					fmt.Sprintf(
						"exec 3<<<$(echo hallo) ; echo test | gpg --batch --symmetric --passphrase-fd=3 -a > '%s'",
						temporaryFile.MustGetLocalPath(),
					),
				}
				_, err := commandexecutor.Bash().RunCommand(
					ctx,
					&parameteroptions.RunCommandOptions{
						Command: createCommand,
					},
				)
				require.NoError(t, err)

				require.True(t, temporaryFile.MustIsPgpEncrypted(verbose))
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				temporaryFile := getFileToTest("localFile")
				expectedMimeType := "inode/x-empty"

				mimeType := temporaryFile.MustGetMimeType(verbose)

				require.EqualValues(expectedMimeType, mimeType)
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				temporaryDir := getDirectoryToTest("localDirectory")
				defer temporaryDir.Delete(verbose)

				file := temporaryDir.MustWriteStringToFileInDirectory("content", verbose, tt.filename)
				readDate := file.MustGetCreationDateByFileName(verbose)

				require.EqualValues(
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				temporaryDir := getDirectoryToTest("localDirectory")
				defer temporaryDir.Delete(verbose)

				file := temporaryDir.MustWriteStringToFileInDirectory("content", verbose, tt.filename)
				require.EqualValues(tt.expectedHasPrefix, file.MustIsYYYYmmdd_HHMMSSPrefix())
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				testFile := getFileToTest("localFile")
				testFile.MustWriteBytes(tt.content, verbose)
				sizeBytes := testFile.MustGetSizeBytes()

				require.EqualValues(tt.expectedSize, sizeBytes)
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				tempFile, err := os.CreateTemp("", "testfile")
				require.Nil(err)

				emptyFile := MustGetLocalFileByPath(tempFile.Name())
				defer emptyFile.MustDelete(verbose)

				emptyFile.MustEnsureEndsWithLineBreak(verbose)

				require.EqualValues("\n", emptyFile.MustReadLastCharAsString())
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				tempFile, err := os.CreateTemp("", "testfile")
				require.Nil(err)

				nonExistingFile := MustGetLocalFileByPath(tempFile.Name())
				defer nonExistingFile.MustDelete(verbose)
				nonExistingFile.MustDelete(verbose)

				require.False(nonExistingFile.MustExists(verbose))

				nonExistingFile.MustEnsureEndsWithLineBreak(verbose)

				require.EqualValues("\n", nonExistingFile.MustReadLastCharAsString())
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				testFile := getFileToTest("localFile")
				testFile.MustWriteString(tt.input, verbose)

				testFile.MustTrimSpacesAtBeginningOfFile(verbose)

				content := testFile.MustReadAsString()
				require.EqualValues(tt.expectedContent, content)
			},
		)
	}
}

/* TODO remove or move
// TODO: Move to File_test.go and test for all File implementations.
func TestFileReplaceBetweenMarkers(t *testing.T) {
	type TestCase struct {
		testDataDir string
	}

	tests := []TestCase{}

	testDataDirectory := MustGetLocalGitRepositoryByPath(".").MustGetSubDirectory("testdata", "File", "ReplaceBetweenMarkers")
	for _, testDirectory := range testDataDirectory.MustListSubDirectories(&parameteroptions.ListDirectoryOptions{Recursive: false}) {
		tests = append(tests, TestCase{testDirectory.MustGetLocalPath()})
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				sourceDir := MustGetLocalDirectoryByPath(filepath.Join(tt.testDataDir, "input"))
				expectedDir := MustGetLocalDirectoryByPath(filepath.Join(tt.testDataDir, "expectedOutput"))

				tempDir := TemporaryDirectories().MustCreateEmptyTemporaryDirectory(verbose)

				sourceDir.MustCopyContentToLocalDirectory(
					tempDir,
					verbose,
				)

				tempDir.MustReplaceBetweenMarkers(verbose)

				listOptions := parameteroptions.NewListFileOptions()
				listOptions.NonRecursive = false
				listOptions.ReturnRelativePaths = true
				expectedFilePaths := expectedDir.MustListFilePaths(listOptions)
				currentFilePaths := tempDir.MustListFilePaths(listOptions)

				require.EqualValues(expectedFilePaths, currentFilePaths)

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

					require.EqualValuesf(
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
*/

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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				testFile := getFileToTest("localFile")
				testFile.MustWriteString(tt.content, verbose)

				require.EqualValues(
					tt.expectedNonEmptyLines,
					testFile.MustGetNumberOfNonEmptyLines())
			},
		)
	}
}

func Test_SecureDelete(t *testing.T) {
	ctx := getCtx()

	t.Run("delete", func(t *testing.T) {
		const verbose bool = true

		testPath := createTempFileAndGetPath()
		require.True(t, filesutils.IsFile(ctx, testPath))

		localFile, err := GetLocalFileByPath(testPath)
		require.NoError(t, err)
		exists, err := localFile.Exists(verbose)
		require.NoError(t, err)
		require.True(t, exists)

		err = localFile.SecurelyDelete(ctx)

		exists, err = localFile.Exists(verbose)
		require.NoError(t, err)
		require.False(t, exists)
	})
}
