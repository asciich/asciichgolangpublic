package files_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorbashoo"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/datatypes/pointersutils"
	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/pathsutils"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func TestLocalFileImplementsFileInterface(t *testing.T) {
	var file filesinterfaces.File = files.MustNewLocalFileByPath("/example/path")
	localPath, err := file.GetLocalPath()
	require.NoError(t, err)
	require.EqualValues(t, "/example/path", localPath)
}

func TestLocalFileIsPathSetOnEmptyFile(t *testing.T) {
	require.EqualValues(t, false, files.NewLocalFile().IsPathSet())
}

func TestLocalFileSetAndGetPath(t *testing.T) {
	require := require.New(t)

	var localFile = files.LocalFile{}

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
				var file filesinterfaces.File = files.MustNewLocalFileByPath(tt.path)

				uri, err := file.GetUriAsString()
				require.NoError(t, err)

				require.EqualValues(t, tt.expectedUri, uri)
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
				const verbose bool = false

				var file filesinterfaces.File = getFileToTest("localFile")

				content, err := file.ReadAsBytes()
				require.NoError(t, err)
				require.EqualValues(t, []byte{}, content)

				for i := 0; i < 2; i++ {
					err = file.WriteBytes(tt.content, verbose)
					require.NoError(t, err)

					content, err := file.ReadAsBytes()
					require.NoError(t, err)
					require.EqualValues(t, tt.content, content)
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
				const verbose bool = false

				var file filesinterfaces.File = getFileToTest("localFile")

				for i := 0; i < 2; i++ {
					err := file.WriteInt64(tt.content, verbose)
					require.NoError(t, err)

					require.EqualValues(t, tt.content, file.MustReadAsInt64())
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
				const verbose bool = false

				var file filesinterfaces.File = getFileToTest("localFile")

				require.EqualValues(t, "", file.MustReadAsString())

				for i := 0; i < 2; i++ {
					err := file.WriteString(tt.content, verbose)
					require.NoError(t, err)

					require.EqualValues(t, tt.content, file.MustReadAsString())
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
				var file filesinterfaces.File
				var err error

				file, err = files.GetLocalFileByPath(tt.path)
				require.NoError(t, err)

				baseName, err := file.GetBaseName()
				require.NoError(t, err)
				require.EqualValues(t, tt.expectedBaseName, baseName)
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
				const verbose bool = true

				temporaryFile := getFileToTest("localFile")
				err := temporaryFile.WriteString(tt.input, verbose)
				require.NoError(t, err)
				defer temporaryFile.Delete(verbose)

				sha256Sum, err := temporaryFile.GetSha256Sum()
				require.NoError(t, err)
				require.EqualValues(t, tt.expectedChecksum, sha256Sum)
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
				const verbose bool = true

				temporaryFile := getFileToTest("localFile")
				err := temporaryFile.WriteString(tt.input, verbose)
				require.NoError(t, err)
				defer temporaryFile.Delete(verbose)

				require.EqualValues(t, tt.expectedIsMatching, temporaryFile.MustIsMatchingSha256Sum(tt.sha256sum))
			},
		)
	}
}

// TODO: Move to File_test.go and test for all File implementations.
func TestLocalFileGetParentDirectory(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		ctx := getCtx()

		const verbose bool = true

		temporaryDir := getDirectoryToTest("localDirectory")
		defer temporaryDir.Delete(verbose)

		temporaryFile, err := temporaryDir.CreateFileInDirectory(ctx, "test.txt", &filesoptions.CreateOptions{})
		require.NoError(t, err)
		parentDir, err := temporaryFile.GetParentDirectory()
		require.NoError(t, err)

		tmpPath, err := temporaryDir.GetLocalPath()
		require.NoError(t, err)

		parentPath, err := parentDir.GetLocalPath()
		require.NoError(t, err)

		require.EqualValues(t, tmpPath, parentPath)
	})
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
				const verbose bool = true

				tempFile1 := getFileToTest("localFile")
				err := tempFile1.WriteString(tt.contentFile1, verbose)
				require.NoError(t, err)
				defer tempFile1.Delete(verbose)

				tempFile2 := getFileToTest("localFile")
				err = tempFile2.WriteString(tt.contentFile2, verbose)
				require.NoError(t, err)
				defer tempFile2.Delete(verbose)

				require.EqualValues(t, tt.expectedIsEqual, mustutils.Must(tempFile1.IsContentEqualByComparingSha256Sum(tempFile2, verbose)))
				require.EqualValues(t, tt.expectedIsEqual, mustutils.Must(tempFile2.IsContentEqualByComparingSha256Sum(tempFile1, verbose)))
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
				localFile, err := files.GetLocalFileByPath(tt.pathToTest)
				require.NoError(t, err)

				localPath := localFile.MustGetLocalPath()

				require.True(t, pathsutils.IsAbsolutePath(localPath))
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
				const verbose bool = true

				testData := "package main\n"
				testData += "\n"
				testData += "//this comment\n"
				testData += "func belongsToMe() (err error) {\n"
				testData += "\n" // ensure empty lines do not split the function block
				testData += "\treturn nil\n"
				testData += "}\n"

				testFile := getFileToTest("localFile")
				err := testFile.WriteString(testData, verbose)
				require.NoError(t, err)
				blocks, err := testFile.GetTextBlocks(verbose)
				require.NoError(t, err)

				require.Len(t, blocks, 2)
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
				const verbose bool = true

				testData := "a: b\n"
				testData += "\n"
				testData += "c: d\n"

				testFile := getFileToTest("localFile")
				testFile.WriteString(testData, verbose)
				blocks, err := testFile.GetTextBlocks(verbose)
				require.NoError(t, err)

				require.Len(t, blocks, 2)
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
				const verbose bool = true

				testData := "---\n"
				testData += "a: b\n"
				testData += "\n"
				testData += "c: d\n"

				testFile := getFileToTest("localFile")
				testFile.WriteString(testData, verbose)
				blocks, err := testFile.GetTextBlocks(verbose)
				require.NoError(t, err)

				require.Len(t, blocks, 3)
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
				const verbose bool = true

				var testFile filesinterfaces.File = getFileToTest("localFile")
				localTestFile := files.MustGetLocalFileByFile(testFile)

				copy := testFile.GetDeepCopy()
				require.EqualValues(
					t,
					mustutils.Must(testFile.GetLocalPath()),
					mustutils.Must(copy.GetLocalPath()),
				)
				localCopy := files.MustGetLocalFileByFile(copy)

				require.False(t, pointersutils.MustPointersEqual(
					localTestFile.MustGetParentFileForBaseClassAsLocalFile(),
					localCopy.MustGetParentFileForBaseClassAsLocalFile(),
				))

				require.True(t, mustutils.Must(testFile.Exists(verbose)))
				require.True(t, mustutils.Must(copy.Exists(verbose)))
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

					file, err := files.GetLocalFileByPath(tt.path)
					require.NoError(t, err)
					path1 = file.MustGetPath()
					os.Chdir("..")
					path2 = file.MustGetPath()
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

func getRepoRootDir(ctx context.Context, t *testing.T) (repoRoot filesinterfaces.Directory) {
	path, err := commandexecutorbashoo.Bash().RunOneLinerAndGetStdoutAsString(ctx, "git rev-parse --show-toplevel")
	require.NoError(t, err)
	path = strings.TrimSpace(path)

	return files.MustGetLocalDirectoryByPath(path)
}

// TODO: Move to File_test.go and test for all File implementations.
func TestLocalFileSortBlocksInFile(t *testing.T) {
	type TestCase struct {
		testDataDir string
	}

	tests := []TestCase{}
	ctx := getCtx()

	testDataDirectory, err := getRepoRootDir(ctx, t).GetSubDirectory("testdata", "File", "SortBlocksInFile")
	require.NoError(t, err)
	for _, testDirectory := range mustutils.Must(testDataDirectory.ListSubDirectories(&parameteroptions.ListDirectoryOptions{Recursive: false})) {
		localPath, err := testDirectory.GetLocalPath()
		require.NoError(t, err)
		tests = append(tests, TestCase{localPath})
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose = true

				testDataDir := files.MustGetLocalDirectoryByPath(tt.testDataDir)
				testInput := testDataDir.MustReadFileInDirectoryAsString("input")

				testFile := getFileToTest("localFile")
				err = testFile.WriteString(testInput, verbose)
				require.NoError(t, err)

				expectedFile := testDataDir.MustGetFileInDirectory("expectedOutput")
				err = testFile.SortBlocksInFile(verbose)
				require.NoError(t, err)

				sortedChecksum, err := testFile.GetSha256Sum()
				require.NoError(t, err)

				expectedChecksum, err := expectedFile.GetSha256Sum()
				require.NoError(t, err)

				if os.Getenv("UPDATE_EXPECTED") == "1" {
					err = testFile.CopyToFile(expectedFile, verbose)
					require.NoError(t, err)
				}

				require.EqualValues(t, expectedChecksum, sortedChecksum)
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
				testFile, err := files.GetLocalFileByPath(tt.inputPath)
				require.NoError(t, err)

				parentPath := testFile.MustGetParentDirectoryPath()
				require.EqualValues(t, tt.expectedParentPath, parentPath)
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
				ctx := getCtx()
				const verbose bool = true

				temporaryFile := getFileToTest("localFile")
				defer temporaryFile.Delete(verbose)

				localPath, err := temporaryFile.GetLocalPath()
				require.NoError(t, err)

				createCommand := []string{
					"bash",
					"-c",
					fmt.Sprintf(
						"exec 3<<<$(echo hallo) ; echo test | gpg --batch --symmetric --passphrase-fd=3 > '%s'",
						localPath,
					),
				}
				commandexecutorbashoo.Bash().RunCommand(
					ctx,
					&parameteroptions.RunCommandOptions{
						Command: createCommand,
					},
				)

				require.True(t, temporaryFile.MustIsPgpEncrypted(verbose))
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

				localPath, err := temporaryFile.GetLocalPath()
				require.NoError(t, err)

				createCommand := []string{
					"bash",
					"-c",
					fmt.Sprintf(
						"exec 3<<<$(echo hallo) ; echo test | gpg --batch --symmetric --passphrase-fd=3 -a > '%s'",
						localPath,
					),
				}
				_, err = commandexecutorbashoo.Bash().RunCommand(
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
				const verbose bool = true

				temporaryFile := getFileToTest("localFile")
				expectedMimeType := "inode/x-empty"

				mimeType, err := temporaryFile.GetMimeType(verbose)
				require.NoError(t, err)

				require.EqualValues(t, expectedMimeType, mimeType)
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
				const verbose bool = true

				temporaryDir := getDirectoryToTest("localDirectory")
				defer temporaryDir.Delete(verbose)

				file, err := temporaryDir.WriteStringToFileInDirectory("content", verbose, tt.filename)
				require.NoError(t, err)
				readDate, err := file.GetCreationDateByFileName(verbose)
				require.NoError(t, err)

				require.EqualValues(t, tt.expected, *readDate)
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
				const verbose bool = true

				temporaryDir := getDirectoryToTest("localDirectory")
				defer temporaryDir.Delete(verbose)

				file, err := temporaryDir.WriteStringToFileInDirectory("content", verbose, tt.filename)
				require.NoError(t, err)
				require.EqualValues(t, tt.expectedHasPrefix, file.MustIsYYYYmmdd_HHMMSSPrefix())
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
				const verbose bool = true

				testFile := getFileToTest("localFile")
				err := testFile.WriteBytes(tt.content, verbose)
				require.NoError(t, err)

				sizeBytes, err := testFile.GetSizeBytes()
				require.NoError(t, err)

				require.EqualValues(t, tt.expectedSize, sizeBytes)
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
				const verbose bool = true

				tempFile, err := os.CreateTemp("", "testfile")
				require.NoError(t, err)

				emptyFile, err := files.GetLocalFileByPath(tempFile.Name())
				require.NoError(t, err)
				defer func() { _ = emptyFile.Delete(verbose) }()

				emptyFile.MustEnsureEndsWithLineBreak(verbose)

				require.EqualValues(t, "\n", emptyFile.MustReadLastCharAsString())
			},
		)
	}
}

// TODO: Move to File_test.go and test for all File implementations.
func TestFileEnsureEndsWithLineBreakOnNonExitistingFile(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		const verbose bool = true

		tempFile, err := os.CreateTemp("", "testfile")
		require.Nil(t, err)

		nonExistingFile, err := files.GetLocalFileByPath(tempFile.Name())
		require.NoError(t, err)
		defer func() { _ = nonExistingFile.Delete(verbose) }()
		err = nonExistingFile.Delete(verbose)
		require.NoError(t, err)

		require.False(t, nonExistingFile.MustExists(verbose))

		nonExistingFile.MustEnsureEndsWithLineBreak(verbose)

		require.EqualValues(t, "\n", nonExistingFile.MustReadLastCharAsString())
	})
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
				const verbose bool = true

				testFile := getFileToTest("localFile")
				err := testFile.WriteString(tt.input, verbose)
				require.NoError(t, err)

				err = testFile.TrimSpacesAtBeginningOfFile(verbose)
				require.NoError(t, err)

				content := testFile.MustReadAsString()
				require.EqualValues(t, tt.expectedContent, content)
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
				const verbose bool = true

				testFile := getFileToTest("localFile")
				err := testFile.WriteString(tt.content, verbose)
				require.NoError(t, err)

				nEmptyLines, err := testFile.GetNumberOfNonEmptyLines()
				require.NoError(t, err)

				require.EqualValues(t, tt.expectedNonEmptyLines, nEmptyLines)
			},
		)
	}
}

func Test_SecureDelete(t *testing.T) {
	ctx := getCtx()

	t.Run("delete", func(t *testing.T) {
		const verbose bool = true

		testPath := createTempFileAndGetPath()
		require.True(t, nativefiles.IsFile(ctx, testPath))

		localFile, err := files.GetLocalFileByPath(testPath)
		require.NoError(t, err)
		exists, err := localFile.Exists(verbose)
		require.NoError(t, err)
		require.True(t, exists)

		err = localFile.SecurelyDelete(ctx)
		require.NoError(t, err)

		exists, err = localFile.Exists(verbose)
		require.NoError(t, err)
		require.False(t, exists)
	})
}
