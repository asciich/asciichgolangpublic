package filesutils_test

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
)

func Test_OpenAsReadCloser(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    []byte
	}{
		{
			name:    "empty file",
			content: "",
			want:    []byte{},
		},
		{
			name:    "hello world",
			content: "hello world",
			want:    []byte("hello world"),
		},
	}

	implementationNames := []string{
		"localFile",
		"localCommandExecutorFile",
		"commandExecutorFileExec",
		"commandExecutorFileBash",
		"nativefilesoo",
	}

	for _, tt := range tests {
		for _, implementation := range implementationNames {
			t.Run(tt.name, func(t *testing.T) {
				ctx := getCtx()

				file := getTemporaryFileToTest(implementation)
				defer file.Delete(ctx, &filesoptions.DeleteOptions{})

				err := file.WriteString(ctx, tt.content, &filesoptions.WriteOptions{})
				require.NoError(t, err)

				readCloser, err := file.OpenAsReadCloser(ctx)
				require.NoError(t, err)
				defer readCloser.Close()

				got, err := io.ReadAll(readCloser)
				require.NoError(t, err)
				require.EqualValues(t, tt.want, got)
			})
		}
	}
}
