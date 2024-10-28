package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommandExecutorFileReadAndWrite(t *testing.T) {
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

				temporaryFile := TemporaryFiles().MustCreateEmptyTemporaryFile(verbose)
				defer temporaryFile.Delete(verbose)

				var fileToTest File = MustGetLocalCommandExecutorFileByPath(temporaryFile.MustGetLocalPath())

				assert.True(fileToTest.MustExists())
				assert.EqualValues(
					"",
					fileToTest.MustReadAsString(),
				)

				fileToTest.WriteString(tt.testContent, verbose)
				assert.True(fileToTest.MustExists())
				assert.EqualValues(
					tt.testContent,
					fileToTest.MustReadAsString(),
				)

				fileToTest.MustDelete(verbose)
				assert.False(fileToTest.MustExists())
			},
		)
	}
}
