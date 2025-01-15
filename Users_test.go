package asciichgolangpublic

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/asciich/asciichgolangpublic/testutils"
)

func TestUserGetHomeDirectory(t *testing.T) {
	tests := []struct {
		testcase string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				assert.True(
					strings.HasPrefix(
						Users().MustGetHomeDirectoryAsString(),
						"/home/",
					),
				)
			},
		)
	}
}

func TestGetFileInHomeDirectory(t *testing.T) {
	tests := []struct {
		filePath []string
	}{
		{[]string{"test"}},
		{[]string{"test", "case"}},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				fileInHome := Users().MustGetFileInHomeDirectory(tt.filePath...)

				filePath := fileInHome.MustGetLocalPath()

				assert.True(
					strings.HasPrefix(
						filePath,
						"/home/",
					),
				)

				expectedPrefix := "/" + strings.Join(tt.filePath, "/")

				assert.True(
					strings.HasSuffix(
						filePath,
						expectedPrefix,
					),
				)
			},
		)
	}
}

func TestGetDirectoryInHomeDirectory(t *testing.T) {
	tests := []struct {
		filePath []string
	}{
		{[]string{"test"}},
		{[]string{"test", "case"}},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				fileInHome := Users().MustGetDirectoryInHomeDirectory(tt.filePath...)

				filePath := fileInHome.MustGetLocalPath()

				assert.True(
					strings.HasPrefix(
						filePath,
						"/home/",
					),
				)

				expectedPrefix := "/" + strings.Join(tt.filePath, "/")

				assert.True(
					strings.HasSuffix(
						filePath,
						expectedPrefix,
					),
				)
			},
		)
	}
}
