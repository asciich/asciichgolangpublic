package asciichgolangpublic

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTemporaryFile(t *testing.T) {
	tests := []struct {
		content string
	}{
		{""},
		{"a"},
		{"hello world"},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose = true

				file := TemporaryFiles().MustCreateFromString(tt.content, verbose)
				defer file.Delete(verbose)

				assert.True(file.MustExists())
				assert.EqualValues(tt.content, file.MustReadAsString())
			},
		)
	}
}

func TestCreateEmptyTemporaryFile(t *testing.T) {
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
				assert.EqualValues("", file.MustReadAsString())
				assert.True(strings.HasPrefix(file.MustGetLocalPath(), "/tmp/"))
			},
		)
	}
}

func TestCreateEmptyTemporaryFileAndGetPath(t *testing.T) {
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

				filePath := TemporaryFiles().MustCreateEmptyTemporaryFileAndGetPath(verbose)
				file := MustNewLocalFileByPath(filePath)
				defer file.Delete(verbose)

				assert.True(file.MustExists())
				assert.EqualValues("", file.MustReadAsString())
				assert.True(strings.HasPrefix(filePath, "/tmp/"))
				assert.True(strings.HasPrefix(file.MustGetPath(), "/tmp/"))
			},
		)
	}
}
