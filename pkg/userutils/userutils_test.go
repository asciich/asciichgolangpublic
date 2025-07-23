package userutils_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/userutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/testutils"
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
				homeDir, err := userutils.GetHomeDirectoryPath()
				require.NoError(t, err)
				require.True(t, strings.HasPrefix(homeDir, "/home/"))
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
				fileInHome, err := userutils.GetFileInHomeDirectory(tt.filePath...)
				require.NoError(t, err)

				filePath, err := fileInHome.GetPath()
				require.NoError(t, err)

				require.True(t, strings.HasPrefix(filePath, "/home/"))

				expectedPrefix := "/" + strings.Join(tt.filePath, "/")
				require.True(t, strings.HasSuffix(filePath, expectedPrefix))
			},
		)
	}
}

func TestGetFilePathInHomeDirectory(t *testing.T) {
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
				filePath, err := userutils.GetFilePathInHomeDirectory(tt.filePath...)
				require.NoError(t, err)

				require.True(t, strings.HasPrefix(filePath, "/home/"))

				expectedPrefix := "/" + strings.Join(tt.filePath, "/")
				require.True(t, strings.HasSuffix(filePath, expectedPrefix))
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
				fileInHome, err := userutils.GetDirectoryInHomeDirectory(tt.filePath...)
				require.NoError(t, err)

				filePath := fileInHome.MustGetLocalPath()

				require.True(t, strings.HasPrefix(filePath, "/home/"))

				expectedPrefix := "/" + strings.Join(tt.filePath, "/")
				require.True(t, strings.HasSuffix(filePath, expectedPrefix))
			},
		)
	}
}
