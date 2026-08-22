package filesutils_test

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
	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfilesoo"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/pathsutils"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

var allFileImplementations = []string{
	"nativefilesoo",
	"commandExecutorFileExec",
	"commandExecutorFileBash",
}

func getTestFileToTest(implementationName string) filesinterfaces.File {
	tempFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		panic(fmt.Sprintf("failed to create temp file: %v", err))
	}
	tempFile.Close()

	return getFileToTest(implementationName, tempFile.Name())
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
	var localFile = files.LocalFile{}

	err := localFile.SetPath("testpath")
	require.EqualValues(t, nil, err)

	receivedPath, err := localFile.GetPath()
	require.EqualValues(t, nil, err)
	require.True(t, strings.HasSuffix(receivedPath, "/testpath"))
}

func TestLocalFileGetUriAsString(t *testing.T) {
	tests := []struct {
		implementationName string
		path               string
		expectedUri        string
	}{
		{"localFile", "/etc/hello.txt", "file:///etc/hello.txt"},
		{"nativefilesoo", "/etc/hello.txt", "file:///etc/hello.txt"},
		{"commandExecutorFileExec", "/etc/hello.txt", "file:///etc/hello.txt"},
		{"commandExecutorFileBash", "/etc/hello.txt", "file:///etc/hello.txt"},
	}

	for _, tt := range tests {
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

func TestFileReadAndWriteAsBytes(t *testing.T) {
	tests := []struct {
		implementationName string
		content            []byte
	}{}

	for _, impl := range allFileImplementations {
		for _, content := range [][]byte{
			[]byte("hello world"),
			[]byte("hello world\n"),
			[]byte("hello\nworld\n"),
		} {
			tests = append(tests, struct {
				implementationName string
				content            []byte
			}{impl, content})
		}
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				var file filesinterfaces.File = getTestFileToTest(tt.implementationName)

				content, err := file.ReadAsBytes(ctx)
				require.NoError(t, err)
				require.EqualValues(t, []byte{}, content)

				for i := 0; i < 2; i++ {
					err = file.WriteBytes(ctx, tt.content, &filesoptions.WriteOptions{})
					require.NoError(t, err)

					content, err := file.ReadAsBytes(ctx)
					require.NoError(t, err)
					require.EqualValues(t, tt.content, content)
				}
			},
		)
	}
}

func TestFileReadAndWriteAsInt64(t *testing.T) {
	tests := []struct {
		implementationName string
		content            int64
	}{}

	for _, impl := range allFileImplementations {
		for _, content := range []int64{1, 2, 3} {
			tests = append(tests, struct {
				implementationName string
				content            int64
			}{impl, content})
		}
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				var file filesinterfaces.File = getTestFileToTest(tt.implementationName)

				for i := 0; i < 2; i++ {
					err := file.WriteInt64(ctx, tt.content)
					require.NoError(t, err)

					got, err := file.ReadAsInt64(ctx)
					require.NoError(t, err)
					require.EqualValues(t, tt.content, got)
				}
			},
		)
	}
}

func TestFileReadAndWriteAsString(t *testing.T) {
	tests := []struct {
		implementationName string
		content            string
	}{}

	for _, impl := range allFileImplementations {
		for _, content := range []string{
			"hello world",
			"hello world\n",
			"hello\nworld\n",
		} {
			tests = append(tests, struct {
				implementationName string
				content            string
			}{impl, content})
		}
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				var file filesinterfaces.File = getTestFileToTest(tt.implementationName)

				content, err := file.ReadAsString(ctx)
				require.NoError(t, err)
				require.EqualValues(t, "", content)

				for i := 0; i < 2; i++ {
					err := file.WriteString(ctx, tt.content, &filesoptions.WriteOptions{})
					require.NoError(t, err)

					content, err := file.ReadAsString(ctx)
					require.NoError(t, err)
					require.EqualValues(t, tt.content, content)
				}
			},
		)
	}
}

func TestFileGetBaseName(t *testing.T) {
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

func TestFileGetSha256Sum(t *testing.T) {
	tests := []struct {
		implementationName string
		input              string
		expectedChecksum   string
	}{}

	for _, impl := range allFileImplementations {
		for _, tc := range []struct {
			input            string
			expectedChecksum string
		}{
			{"", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
			{"hello world", "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"},
		} {
			tests = append(tests, struct {
				implementationName string
				input              string
				expectedChecksum   string
			}{impl, tc.input, tc.expectedChecksum})
		}
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				temporaryFile := getTestFileToTest(tt.implementationName)
				err := temporaryFile.WriteString(ctx, tt.input, &filesoptions.WriteOptions{})
				require.NoError(t, err)
				defer temporaryFile.Delete(ctx, &filesoptions.DeleteOptions{})

				sha256Sum, err := temporaryFile.GetSha256Sum(ctx)
				require.NoError(t, err)
				require.EqualValues(t, tt.expectedChecksum, sha256Sum)
			},
		)
	}
}

func TestFileIsMatchingSha256Sum(t *testing.T) {
	tests := []struct {
		implementationName string
		input              string
		sha256sum          string
		expectedIsMatching bool
	}{}

	for _, impl := range allFileImplementations {
		for _, tc := range []struct {
			input              string
			sha256sum          string
			expectedIsMatching bool
		}{
			{"", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855", true},
			{"hello world", "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9", true},
			{"", "aaaaaae3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855", false},
			{"hello world", "aaaaaab94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9", false},
		} {
			tests = append(tests, struct {
				implementationName string
				input              string
				sha256sum          string
				expectedIsMatching bool
			}{impl, tc.input, tc.sha256sum, tc.expectedIsMatching})
		}
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				temporaryFile := getTestFileToTest(tt.implementationName)
				err := temporaryFile.WriteString(ctx, tt.input, &filesoptions.WriteOptions{})
				require.NoError(t, err)
				defer temporaryFile.Delete(ctx, &filesoptions.DeleteOptions{})

				isMatching, err := temporaryFile.IsMatchingSha256Sum(tt.sha256sum)
				require.NoError(t, err)
				require.EqualValues(t, tt.expectedIsMatching, isMatching)
			},
		)
	}
}

func TestFileGetParentDirectory(t *testing.T) {
	for _, impl := range allFileImplementations {
		t.Run(impl, func(t *testing.T) {
			ctx := getCtx()

			tempDir, err := tempfilesoo.CreateEmptyTemporaryDirectoryAndGetPath(ctx)
			require.NoError(t, err)

			temporaryDir := getDirectoryToTest("localDirectory", tempDir)
			defer temporaryDir.Delete(ctx, &filesoptions.DeleteOptions{})

			temporaryFile, err := temporaryDir.CreateFileInDirectory(ctx, "test.txt", &filesoptions.CreateOptions{})
			require.NoError(t, err)
			parentDir, err := temporaryFile.GetParentDirectory(ctx)
			require.NoError(t, err)

			tmpPath, err := temporaryDir.GetLocalPath()
			require.NoError(t, err)

			parentPath, err := parentDir.GetLocalPath()
			require.NoError(t, err)

			require.EqualValues(t, tmpPath, parentPath)
		})
	}
}

func TestFileIsContentEqualByComparingSha256Sum(t *testing.T) {
	tests := []struct {
		implementationName string
		contentFile1       string
		contentFile2       string
		expectedIsEqual    bool
	}{}

	for _, impl := range allFileImplementations {
		for _, tc := range []struct {
			contentFile1    string
			contentFile2    string
			expectedIsEqual bool
		}{
			{"", "", true},
			{"testcase", "testcase", true},
			{"testcase1", "testcase", false},
			{"testcase1", "testcase2", false},
			{"testcase", "testcase3", false},
		} {
			tests = append(tests, struct {
				implementationName string
				contentFile1       string
				contentFile2       string
				expectedIsEqual    bool
			}{impl, tc.contentFile1, tc.contentFile2, tc.expectedIsEqual})
		}
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				tempFile1 := getTestFileToTest(tt.implementationName)
				err := tempFile1.WriteString(ctx, tt.contentFile1, &filesoptions.WriteOptions{})
				require.NoError(t, err)
				defer tempFile1.Delete(ctx, &filesoptions.DeleteOptions{})

				tempFile2 := getTestFileToTest(tt.implementationName)
				err = tempFile2.WriteString(ctx, tt.contentFile2, &filesoptions.WriteOptions{})
				require.NoError(t, err)
				defer tempFile2.Delete(ctx, &filesoptions.DeleteOptions{})

				require.EqualValues(t, tt.expectedIsEqual, mustutils.Must(tempFile1.IsContentEqualByComparingSha256Sum(ctx, tempFile2)))
				require.EqualValues(t, tt.expectedIsEqual, mustutils.Must(tempFile2.IsContentEqualByComparingSha256Sum(ctx, tempFile1)))
			},
		)
	}
}

func TestFileGetLocalPathIsAbsolute(t *testing.T) {
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
				localFile, err := files.GetLocalFileByPath(tt.pathToTest)
				require.NoError(t, err)

				localPath, err := localFile.GetLocalPath()
				require.NoError(t, err)

				require.True(t, pathsutils.IsAbsolutePath(localPath))
			},
		)
	}
}

func TestFileGetTextBlocksGolangWithCommentAboveFunction(t *testing.T) {
	for _, impl := range allFileImplementations {
		t.Run(impl, func(t *testing.T) {
			ctx := getCtx()

			testData := "package main\n"
			testData += "\n"
			testData += "//this comment\n"
			testData += "func belongsToMe() (err error) {\n"
			testData += "\n"
			testData += "\treturn nil\n"
			testData += "}\n"

			testFile := getTestFileToTest(impl)
			err := testFile.WriteString(ctx, testData, &filesoptions.WriteOptions{})
			require.NoError(t, err)
			blocks, err := testFile.GetTextBlocks(ctx)
			require.NoError(t, err)

			require.Len(t, blocks, 2)
		})
	}
}

func TestFileGetTextBlocksYamlWithoutLeadingThreeMinuses(t *testing.T) {
	for _, impl := range allFileImplementations {
		t.Run(impl, func(t *testing.T) {
			ctx := getCtx()

			testData := "a: b\n"
			testData += "\n"
			testData += "c: d\n"

			testFile := getTestFileToTest(impl)
			testFile.WriteString(ctx, testData, &filesoptions.WriteOptions{})
			blocks, err := testFile.GetTextBlocks(ctx)
			require.NoError(t, err)

			require.Len(t, blocks, 2)
		})
	}
}

func TestFileGetTextBlocksYamlWithLeadingThreeMinuses(t *testing.T) {
	for _, impl := range allFileImplementations {
		t.Run(impl, func(t *testing.T) {
			ctx := getCtx()

			testData := "---\n"
			testData += "a: b\n"
			testData += "\n"
			testData += "c: d\n"

			testFile := getTestFileToTest(impl)
			testFile.WriteString(ctx, testData, &filesoptions.WriteOptions{})
			blocks, err := testFile.GetTextBlocks(ctx)
			require.NoError(t, err)

			require.Len(t, blocks, 3)
		})
	}
}

func TestFileReplaceLineAfterLine(t *testing.T) {
	tests := []struct {
		implementationName        string
		input                     string
		lineToSearch              string
		replaceLineAfterFoundWith string
		expectedContent           string
		expectedChanged           bool
	}{}

	for _, impl := range allFileImplementations {
		for _, tc := range []struct {
			input                     string
			lineToSearch              string
			replaceLineAfterFoundWith string
			expectedContent           string
			expectedChanged           bool
		}{
			{"a\nb\nc\n", "a", "d", "a\nd\nc\n", true},
			{"a\nb\nc\n", "x", "d", "a\nb\nc\n", false},
		} {
			tests = append(tests, struct {
				implementationName        string
				input                     string
				lineToSearch              string
				replaceLineAfterFoundWith string
				expectedContent           string
				expectedChanged           bool
			}{impl, tc.input, tc.lineToSearch, tc.replaceLineAfterFoundWith, tc.expectedContent, tc.expectedChanged})
		}
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				testFile := getTestFileToTest(tt.implementationName)
				testFile.WriteString(ctx, tt.input, &filesoptions.WriteOptions{})
				defer testFile.Delete(ctx, &filesoptions.DeleteOptions{})

				changeSummary, err := testFile.ReplaceLineAfterLine(ctx, tt.lineToSearch, tt.replaceLineAfterFoundWith)
				require.NoError(t, err)

				content, err := testFile.ReadAsString(ctx)
				require.NoError(t, err)
				require.EqualValues(t, tt.expectedContent, content)
				require.EqualValues(t, tt.expectedChanged, changeSummary.IsChanged())
			},
		)
	}
}

func TestFile_GetPathReturnsAbsoluteValue(t *testing.T) {
	tests := []struct {
		path string
	}{
		{"test.txt"},
		{"./test.txt"},
		{"../test.txt"},
	}

	for _, tt := range tests {
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
					path1, err = file.GetPath()
					require.NoError(t, err)
					os.Chdir("..")
					path2, err = file.GetPath()
					require.NoError(t, err)
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

	repoRoot, err = files.GetLocalDirectoryByPath(ctx, path)
	require.NoError(t, err)

	return repoRoot
}

func TestFileSortBlocksInFile(t *testing.T) {
	type TestCase struct {
		implementationName string
		testDataDir        string
	}

	tests := []TestCase{}
	ctx := getCtx()

	testDataDirectory, err := getRepoRootDir(ctx, t).GetSubDirectory(ctx, "testdata", "File", "SortBlocksInFile")
	require.NoError(t, err)

	for _, impl := range allFileImplementations {
		for _, testDirectory := range mustutils.Must(testDataDirectory.ListSubDirectories(ctx, &parameteroptions.ListDirectoryOptions{Recursive: false})) {
			localPath, err := testDirectory.GetLocalPath()
			require.NoError(t, err)
			tests = append(tests, TestCase{impl, localPath})
		}
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				testDataDir, err := files.GetLocalDirectoryByPath(ctx, tt.testDataDir)
				require.NoError(t, err)

				testInput, err := testDataDir.ReadFileInDirectoryAsString(ctx, "input")
				require.NoError(t, err)

				testFile := getTestFileToTest(tt.implementationName)
				err = testFile.WriteString(ctx, testInput, &filesoptions.WriteOptions{})
				require.NoError(t, err)

				expectedFile, err := testDataDir.GetFileInDirectory("expectedOutput")
				require.NoError(t, err)

				err = testFile.SortBlocksInFile(ctx)
				require.NoError(t, err)

				sortedChecksum, err := testFile.GetSha256Sum(ctx)
				require.NoError(t, err)

				expectedChecksum, err := expectedFile.GetSha256Sum(ctx)
				require.NoError(t, err)

				if os.Getenv("UPDATE_EXPECTED") == "1" {
					err = testFile.CopyToFile(ctx, expectedFile, &filesoptions.CopyOptions{})
					require.NoError(t, err)
				}

				require.EqualValues(t, expectedChecksum, sortedChecksum)
			},
		)
	}
}

func TestFileGetLastCharAsString(t *testing.T) {
	tests := []struct {
		implementationName string
		content            string
		lastChar           string
	}{}

	for _, impl := range allFileImplementations {
		for _, tc := range []struct {
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
		} {
			tests = append(tests, struct {
				implementationName string
				content            string
				lastChar           string
			}{impl, tc.content, tc.lastChar})
		}
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				testFile := getTestFileToTest(tt.implementationName)
				testFile.WriteString(ctx, tt.content, &filesoptions.WriteOptions{})
				defer testFile.Delete(ctx, &filesoptions.DeleteOptions{})

				lastChar, err := testFile.ReadLastCharAsString(ctx)
				require.NoError(t, err)

				require.EqualValues(t, tt.lastChar, lastChar)
			},
		)
	}
}

func TestFileGetAsFloat64(t *testing.T) {
	tests := []struct {
		implementationName string
		content            string
		expectedFloat      float64
	}{}

	for _, impl := range allFileImplementations {
		for _, tc := range []struct {
			content       string
			expectedFloat float64
		}{
			{"0", 0.0},
			{"0.", 0.0},
			{"0.0", 0.0},
			{"0.1", 0.1},
			{"0.10", 0.1},
			{"-3", -3.0},
		} {
			tests = append(tests, struct {
				implementationName string
				content            string
				expectedFloat      float64
			}{impl, tc.content, tc.expectedFloat})
		}
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				testFile := getTestFileToTest(tt.implementationName)
				testFile.WriteString(ctx, tt.content, &filesoptions.WriteOptions{})
				defer testFile.Delete(ctx, &filesoptions.DeleteOptions{})

				readFloat, err := testFile.ReadAsFloat64(ctx)
				require.NoError(t, err)

				require.EqualValues(t, tt.expectedFloat, readFloat)
			},
		)
	}
}

func TestFileGetAsInt64(t *testing.T) {
	tests := []struct {
		implementationName string
		content            string
		expectedInt        int64
	}{}

	for _, impl := range allFileImplementations {
		for _, tc := range []struct {
			content     string
			expectedInt int64
		}{
			{"0", 0},
			{"1", 1},
			{"1\n", 1},
			{"10\n", 10},
			{" 10\n", 10},
			{" -110\n", -110},
		} {
			tests = append(tests, struct {
				implementationName string
				content            string
				expectedInt        int64
			}{impl, tc.content, tc.expectedInt})
		}
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				testFile := getTestFileToTest(tt.implementationName)
				testFile.WriteString(ctx, tt.content, &filesoptions.WriteOptions{})
				defer testFile.Delete(ctx, &filesoptions.DeleteOptions{})

				readInt64, err := testFile.ReadAsInt64(ctx)
				require.NoError(t, err)

				require.EqualValues(t, tt.expectedInt, readInt64)
			},
		)
	}
}

func TestFileGetAsInt(t *testing.T) {
	tests := []struct {
		implementationName string
		content            string
		expectedInt        int
	}{}

	for _, impl := range allFileImplementations {
		for _, tc := range []struct {
			content     string
			expectedInt int
		}{
			{"0", 0},
			{"1", 1},
			{"1\n", 1},
			{"10\n", 10},
			{" 10\n", 10},
			{" -110\n", -110},
		} {
			tests = append(tests, struct {
				implementationName string
				content            string
				expectedInt        int
			}{impl, tc.content, tc.expectedInt})
		}
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				testFile := getTestFileToTest(tt.implementationName)
				testFile.WriteString(ctx, tt.content, &filesoptions.WriteOptions{})
				defer testFile.Delete(ctx, &filesoptions.DeleteOptions{})

				readInt, err := testFile.ReadAsInt(ctx)
				require.NoError(t, err)

				require.EqualValues(t, tt.expectedInt, readInt)
			},
		)
	}
}

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
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()
				testFile, err := files.GetLocalFileByPath(tt.inputPath)
				require.NoError(t, err)

				parentPath, err := testFile.GetParentDirectoryPath(ctx)
				require.NoError(t, err)
				require.EqualValues(t, tt.expectedParentPath, parentPath)
			},
		)
	}
}

func TestFileIsPgpEncrypted_Case1_unencrypted(t *testing.T) {
	tests := []struct {
		implementationName string
		unencrypted        string
	}{}

	for _, impl := range allFileImplementations {
		for _, content := range []string{"", "\n", "---", "testcase", "testcase\n"} {
			tests = append(tests, struct {
				implementationName string
				unencrypted        string
			}{impl, content})
		}
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				testFile := getTestFileToTest(tt.implementationName)
				testFile.WriteString(ctx, tt.unencrypted, &filesoptions.WriteOptions{})
				defer testFile.Delete(ctx, &filesoptions.DeleteOptions{})

				isPgpEncrypted, err := testFile.IsPgpEncrypted(ctx)
				require.NoError(t, err)
				require.False(t, isPgpEncrypted)
			},
		)
	}
}

func TestFileIsPgpEncrypted_Case2_encryptedBinary(t *testing.T) {
	for _, impl := range allFileImplementations {
		t.Run(impl, func(t *testing.T) {
			ctx := getCtx()

			temporaryFile := getTestFileToTest(impl)
			defer temporaryFile.Delete(ctx, &filesoptions.DeleteOptions{})

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

			isPgpEncrypted, err := temporaryFile.IsPgpEncrypted(ctx)
			require.NoError(t, err)
			require.True(t, isPgpEncrypted)
		})
	}
}

func TestFileIsPgpEncrypted_Case3_encryptedAsciiArmor(t *testing.T) {
	for _, impl := range allFileImplementations {
		t.Run(impl, func(t *testing.T) {
			ctx := getCtx()

			temporaryFile := getTestFileToTest(impl)
			defer temporaryFile.Delete(ctx, &filesoptions.DeleteOptions{})

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

			isPgpEncrypted, err := temporaryFile.IsPgpEncrypted(ctx)
			require.True(t, isPgpEncrypted)
		})
	}
}

func TestFileGetMimeTypeOfEmptyFile(t *testing.T) {
	for _, impl := range allFileImplementations {
		t.Run(impl, func(t *testing.T) {
			ctx := getCtx()
			temporaryFile := getTestFileToTest(impl)
			expectedMimeType := "inode/x-empty"

			mimeType, err := temporaryFile.GetMimeType(ctx)
			require.NoError(t, err)

			require.EqualValues(t, expectedMimeType, mimeType)
		})
	}
}

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
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				tempPath, err := tempfilesoo.CreateEmptyTemporaryDirectoryAndGetPath(ctx)
				require.NoError(t, err)

				temporaryDir := getDirectoryToTest("localDirectory", tempPath)
				defer temporaryDir.Delete(ctx, &filesoptions.DeleteOptions{})

				file, err := temporaryDir.WriteStringToFile(ctx, tt.filename, "content", &filesoptions.WriteOptions{})
				require.NoError(t, err)
				readDate, err := file.GetCreationDateByFileName(ctx)
				require.NoError(t, err)

				require.EqualValues(t, tt.expected, *readDate)
			},
		)
	}
}

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
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				tempPath, err := tempfilesoo.CreateEmptyTemporaryDirectoryAndGetPath(ctx)
				require.NoError(t, err)

				temporaryDir := getDirectoryToTest("localDirectory", tempPath)
				defer temporaryDir.Delete(ctx, &filesoptions.DeleteOptions{})

				file, err := temporaryDir.WriteStringToFile(ctx, tt.filename, "content", &filesoptions.WriteOptions{})
				require.NoError(t, err)
				got, err := file.IsYYYYmmdd_HHMMSSPrefix()
				require.NoError(t, err)
				require.EqualValues(t, tt.expectedHasPrefix, got)
			},
		)
	}
}

func TestFileGetSizeBytes(t *testing.T) {
	tests := []struct {
		implementationName string
		content            []byte
		expectedSize       int64
	}{}

	for _, impl := range allFileImplementations {
		for _, tc := range []struct {
			content      []byte
			expectedSize int64
		}{
			{[]byte{}, 0},
			{[]byte("helloWorld"), 10},
		} {
			tests = append(tests, struct {
				implementationName string
				content            []byte
				expectedSize       int64
			}{impl, tc.content, tc.expectedSize})
		}
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				testFile := getTestFileToTest(tt.implementationName)
				err := testFile.WriteBytes(ctx, tt.content, &filesoptions.WriteOptions{})
				require.NoError(t, err)

				sizeBytes, err := testFile.GetSizeBytes(ctx)
				require.NoError(t, err)

				require.EqualValues(t, tt.expectedSize, sizeBytes)
			},
		)
	}
}

func TestFileEnsureEndsWithLineBreakOnEmptyFile(t *testing.T) {
	for _, impl := range allFileImplementations {
		t.Run(impl, func(t *testing.T) {
			ctx := getCtx()

			testFile := getTestFileToTest(impl)
			defer testFile.Delete(ctx, &filesoptions.DeleteOptions{})

			err := testFile.EnsureEndsWithLineBreak(ctx)
			require.NoError(t, err)

			got, err := testFile.ReadLastCharAsString(ctx)
			require.NoError(t, err)
			require.EqualValues(t, "\n", got)
		})
	}
}

func TestFileEnsureEndsWithLineBreakOnNonExistingFile(t *testing.T) {
	for _, impl := range allFileImplementations {
		t.Run(impl, func(t *testing.T) {
			ctx := getCtx()

			tempFile, err := os.CreateTemp("", "testfile")
			require.NoError(t, err)

			nonExistingFile, err := files.GetLocalFileByPath(tempFile.Name())
			require.NoError(t, err)
			defer func() { _ = nonExistingFile.Delete(ctx, &filesoptions.DeleteOptions{}) }()
			err = nonExistingFile.Delete(ctx, &filesoptions.DeleteOptions{})
			require.NoError(t, err)

			exists, err := nonExistingFile.Exists(ctx)
			require.NoError(t, err)
			require.False(t, exists)

			err = nonExistingFile.EnsureEndsWithLineBreak(ctx)
			require.NoError(t, err)

			got, err := nonExistingFile.ReadLastCharAsString(ctx)
			require.NoError(t, err)

			require.EqualValues(t, "\n", got)
		})
	}
}

func TestFileTrimSpacesAtBeginningOfFile(t *testing.T) {
	tests := []struct {
		implementationName string
		input              string
		expectedContent    string
	}{}

	for _, impl := range allFileImplementations {
		for _, tc := range []struct {
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
		} {
			tests = append(tests, struct {
				implementationName string
				input              string
				expectedContent    string
			}{impl, tc.input, tc.expectedContent})
		}
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				testFile := getTestFileToTest(tt.implementationName)
				err := testFile.WriteString(ctx, tt.input, &filesoptions.WriteOptions{})
				require.NoError(t, err)

				err = testFile.TrimSpacesAtBeginningOfFile(ctx)
				require.NoError(t, err)

				content, err := testFile.ReadAsString(ctx)
				require.NoError(t, err)
				require.EqualValues(t, tt.expectedContent, content)
			},
		)
	}
}

func TestFileGetNumberOfNonEmptyLines(t *testing.T) {
	tests := []struct {
		implementationName    string
		content               string
		expectedNonEmptyLines int
	}{}

	for _, impl := range allFileImplementations {
		for _, tc := range []struct {
			content               string
			expectedNonEmptyLines int
		}{
			{"", 0},
			{"testcase", 1},
			{"testcase\n", 1},
			{"testcase\n\n", 1},
			{"testcase\n\na", 2},
		} {
			tests = append(tests, struct {
				implementationName    string
				content               string
				expectedNonEmptyLines int
			}{impl, tc.content, tc.expectedNonEmptyLines})
		}
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				testFile := getTestFileToTest(tt.implementationName)
				err := testFile.WriteString(ctx, tt.content, &filesoptions.WriteOptions{})
				require.NoError(t, err)

				nEmptyLines, err := testFile.GetNumberOfNonEmptyLines(ctx)
				require.NoError(t, err)

				require.EqualValues(t, tt.expectedNonEmptyLines, nEmptyLines)
			},
		)
	}
}

func Test_SecureDelete(t *testing.T) {
	for _, impl := range allFileImplementations {
		t.Run(impl, func(t *testing.T) {
			ctx := getCtx()

			testPath, err := tempfilesoo.CreateEmptyTemporaryFileAndGetPath(ctx)
			require.True(t, nativefiles.IsFile(ctx, testPath))

			localFile, err := files.GetLocalFileByPath(testPath)
			require.NoError(t, err)
			exists, err := localFile.Exists(ctx)
			require.NoError(t, err)
			require.True(t, exists)

			err = localFile.SecurelyDelete(ctx)
			require.NoError(t, err)

			exists, err = localFile.Exists(ctx)
			require.NoError(t, err)
			require.False(t, exists)
		})
	}
}
