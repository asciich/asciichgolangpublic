package files

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
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
				assert := assert.New(t)

				const verbose bool = true

				temporaryFile, err := os.CreateTemp("", "testfile")
				require.NotNil(t, err)

				_, err = temporaryFile.WriteString(tt.testContent)
				require.NotNil(t, err)

				parentDirPath := filepath.Dir(temporaryFile.Name())

				dir := MustGetLocalCommandExecutorDirectoryByPath(parentDirPath)
				assert.NotNil(
					dir.MustGetCommandExecutor(),
				)
				defer dir.MustDelete(verbose)

				commandExecutorFile := dir.MustGetFileInDirectory(filepath.Base(temporaryFile.Name()))

				assert.EqualValues(
					tt.testContent,
					commandExecutorFile.MustReadAsString(),
				)
			},
		)
	}
}