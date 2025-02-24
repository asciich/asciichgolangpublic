package files

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/testutils"
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

				dir := MustGetLocalCommandExecutorDirectoryByPath(parentDirPath)
				require.NotNil(
					t,
					dir.MustGetCommandExecutor(),
				)
				defer dir.MustDelete(verbose)

				commandExecutorFile := dir.MustGetFileInDirectory(filepath.Base(temporaryFile.Name()))

				require.EqualValues(
					t,
					tt.testContent,
					commandExecutorFile.MustReadAsString(),
				)
			},
		)
	}
}
