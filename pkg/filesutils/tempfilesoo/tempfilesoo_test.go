package tempfilesoo_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfilesoo"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
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
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose = true

				file, err := tempfilesoo.CreateFromString(tt.content, verbose)
				require.NoError(t, err)
				defer file.Delete(verbose)

				require.True(t, file.MustExists(verbose))
				require.EqualValues(t, tt.content, file.MustReadAsString())
			},
		)
	}
}

func TestCreateEmptyTemporaryFile(t *testing.T) {
	const verbose bool = true

	file, err := tempfilesoo.CreateEmptyTemporaryFile(verbose)
	require.NoError(t, err)
	defer file.Delete(verbose)

	require.True(t, file.MustExists(verbose))
	require.EqualValues(t, "", file.MustReadAsString())
	require.True(t, strings.HasPrefix(file.MustGetLocalPath(), "/tmp/"))
}

func TestCreateEmptyTemporaryFileAndGetPath(t *testing.T) {
	const verbose bool = true

	filePath, err := tempfilesoo.CreateEmptyTemporaryFileAndGetPath(verbose)
	require.NoError(t, err)

	file := files.MustNewLocalFileByPath(filePath)
	defer file.Delete(verbose)

	require.True(t, file.MustExists(verbose))
	require.EqualValues(t, "", file.MustReadAsString())
	require.True(t, strings.HasPrefix(filePath, "/tmp/"))
	require.True(t, strings.HasPrefix(file.MustGetPath(), "/tmp/"))
}
