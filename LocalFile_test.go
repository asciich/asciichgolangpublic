package asciichgolangpublic

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestLocalFileExists(t *testing.T) {

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

				file := TemporaryFiles().MustCreateEmptyTemporaryFile(verbose)
				defer file.Delete(verbose)

				assert.True(file.MustExists())

				file.MustDelete(verbose)

				assert.False(file.MustExists())
			},
		)
	}
}

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

				temporaryFile := temporaryDir.MustCreateFileInDirectory("test.txt")
				parentDir := temporaryFile.MustGetParentDirectory()

				assert.EqualValues(
					temporaryDir.MustGetLocalPath(),
					parentDir.MustGetLocalPath(),
				)
			},
		)
	}
}

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

				assert.False(Pointers().MustPointersEqual(
					localTestFile.MustGetParentFileForBaseClassAsLocalFile(),
					localCopy.MustGetParentFileForBaseClassAsLocalFile(),
				))

				assert.True(testFile.MustExists())
				assert.True(copy.MustExists())
			},
		)
	}
}
