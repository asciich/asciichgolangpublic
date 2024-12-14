package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommandExecutorDirectoryRead_GetFileInDirectory(t *testing.T) {
	tests := []struct {
		testContent string
	}{
		{"testcase"},
		{"testcase\n"},
		{"multyLine\nContent"},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				temporaryFile := TemporaryFiles().MustCreateFromString(tt.testContent, verbose)
				defer temporaryFile.Delete(verbose)

				parentDirectoryPath := temporaryFile.MustGetParentDirectoryPath()

				dir := MustGetLocalCommandExecutorDirectoryByPath(parentDirectoryPath)
				assert.NotNil(
					dir.MustGetCommandExecutor(),
				)

				commandExecutorFile := dir.MustGetFileInDirectory(temporaryFile.MustGetBaseName())

				assert.EqualValues(
					tt.testContent,
					commandExecutorFile.MustReadAsString(),
				)
			},
		)
	}
}
