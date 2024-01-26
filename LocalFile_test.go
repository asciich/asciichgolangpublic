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
	assert.EqualValues(t, false, LocalFile{}.IsPathSet())
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
