package filesutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func TestFile_CopyToFile(t *testing.T) {
	tests := []struct {
		srcImplementationName string
		content               string
	}{
		// legacy implementations:
		{"localFile", "test content\nwith a new line\n"},
		{"localCommandExecutorFile", "test content\nwith a new line\n"},
		// new implementations:
		{"commandExecutorFileExec", "test content\nwith a new line\n"},
		{"commandExecutorFileBash", "test content\nwith a new line\n"},
		{"nativefilesoo", "test content\nwith a new line\n"},
	}
	dests := []struct {
		destImplementationName string
		options                *filesoptions.CopyOptions
	}{
		// nil options
		// ============
		// legacy implementations:
		{"localFile", nil},
		{"localCommandExecutorFile", nil},
		// new implementations:
		{"commandExecutorFileExec", nil},
		{"commandExecutorFileBash", nil},
		{"nativefilesoo", nil},

		// empty options
		// ============
		// legacy implementations:
		{"localFile", &filesoptions.CopyOptions{}},
		{"localCommandExecutorFile", &filesoptions.CopyOptions{}},
		// new implementations:
		{"commandExecutorFileExec", &filesoptions.CopyOptions{}},
		{"commandExecutorFileBash", &filesoptions.CopyOptions{}},
		{"nativefilesoo", &filesoptions.CopyOptions{}},

		// Use sudo
		// ============
		// legacy implementations:
		{"localFile", &filesoptions.CopyOptions{UseSudo: true}},
		{"localCommandExecutorFile", &filesoptions.CopyOptions{UseSudo: true}},
		// new implementations:
		{"commandExecutorFileExec", &filesoptions.CopyOptions{UseSudo: true}},
		{"commandExecutorFileBash", &filesoptions.CopyOptions{UseSudo: true}},
		{"nativefilesoo", &filesoptions.CopyOptions{UseSudo: true}},
	}

	for _, tt := range tests {
		for _, d := range dests {
			t.Run(
				testutils.MustFormatAsTestname(tt)+"_"+testutils.MustFormatAsTestname(d),
				func(t *testing.T) {
					ctx := getCtx()

					srcFile := getTemporaryFileToTest(tt.srcImplementationName)
					err := srcFile.WriteString(ctx, tt.content, &filesoptions.WriteOptions{})
					require.NoError(t, err)
					defer srcFile.Delete(ctx, &filesoptions.DeleteOptions{})

					destFile := getTemporaryFileToTest(d.destImplementationName)
					defer destFile.Delete(ctx, &filesoptions.DeleteOptions{})
					destFile.Delete(ctx, &filesoptions.DeleteOptions{})

					require.True(t, mustutils.Must(srcFile.Exists(ctx)))
					require.False(t, mustutils.Must(destFile.Exists(ctx)))

					err = srcFile.CopyToFile(ctx, destFile, d.options)
					require.NoError(t, err)

					require.True(t, mustutils.Must(srcFile.Exists(ctx)))
					require.True(t, mustutils.Must(destFile.Exists(ctx)))

					content, err := srcFile.ReadAsString(ctx)
					require.NoError(t, err)
					require.EqualValues(t, tt.content, content)

					content, err = destFile.ReadAsString(ctx)
					require.NoError(t, err)
					require.EqualValues(t, tt.content, content)
				},
			)
		}
	}
}
