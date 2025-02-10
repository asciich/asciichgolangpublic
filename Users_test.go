package asciichgolangpublic

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
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
				require := require.New(t)

				require.True(
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
				require := require.New(t)

				fileInHome := Users().MustGetFileInHomeDirectory(tt.filePath...)

				filePath := fileInHome.MustGetLocalPath()

				require.True(
					strings.HasPrefix(
						filePath,
						"/home/",
					),
				)

				expectedPrefix := "/" + strings.Join(tt.filePath, "/")

				require.True(
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
				require := require.New(t)

				fileInHome := Users().MustGetDirectoryInHomeDirectory(tt.filePath...)

				filePath := fileInHome.MustGetLocalPath()

				require.True(
					strings.HasPrefix(
						filePath,
						"/home/",
					),
				)

				expectedPrefix := "/" + strings.Join(tt.filePath, "/")

				require.True(
					strings.HasSuffix(
						filePath,
						expectedPrefix,
					),
				)
			},
		)
	}
}
