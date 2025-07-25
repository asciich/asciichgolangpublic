package files_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
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
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true

				temporaryDir, err := os.MkdirTemp("", "testdir")
				require.Nil(t, err)

				temporaryFile, err := os.CreateTemp(temporaryDir, "testfile")
				require.Nil(t, err)

				_, err = temporaryFile.WriteString(tt.testContent)
				require.Nil(t, err)

				parentDirPath := filepath.Dir(temporaryFile.Name())

				dir, err := files.GetLocalCommandExecutorDirectoryByPath(parentDirPath)
				require.NoError(t, err)

				executor, err := dir.GetCommandExecutor()
				require.NoError(t, err)
				require.NotNil(t, executor)
				defer dir.Delete(verbose)

				commandExecutorFile, err := dir.GetFileInDirectory(filepath.Base(temporaryFile.Name()))
				require.NoError(t, err)

				require.EqualValues(t, tt.testContent, commandExecutorFile.MustReadAsString())
			},
		)
	}
}
