package asciichgolangpublic

import (
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
	assert.EqualValues("testpath", receivedPath)
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

				file.Delete(verbose)

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
