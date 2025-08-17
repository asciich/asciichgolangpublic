package tarutils

import (
	"archive/tar"
	"bytes"
	"io"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/archiveutils/tarutils/tarparameteroptions"
)

func TestFileToTarReader(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		// 1. Create a temporary file with test content.
		testContent := "This is the content for the tar archive."
		fileName := "testfile.txt"

		tmpDir, err := os.MkdirTemp("", "tar_test_dir-*")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		localFilePath := filepath.Join(tmpDir, fileName)

		err = os.WriteFile(localFilePath, []byte(testContent), 0644)
		require.NoError(t, err)

		tarBuffer, err := FileToTarReader(localFilePath, &tarparameteroptions.FileToTarOptions{})
		require.NoError(t, err)
		require.NotNil(t, tarBuffer)

		tr := tar.NewReader(tarBuffer)
		header, err := tr.Next()
		require.NoError(t, err)

		expectedTarFileName := path.Base(localFilePath)
		require.EqualValues(t, expectedTarFileName, header.Name)

		readContentBytes := new(bytes.Buffer)
		_, err = io.Copy(readContentBytes, tr)
		require.NoError(t, err)

		readContent := readContentBytes.String()
		require.EqualValues(t, testContent, readContent)

		// Ensure no more entries in the tar:
		_, err = tr.Next()
		require.EqualValues(t, io.EOF, err)
	})

	t.Run("Nonexistent file", func(t *testing.T) {
		nonExistentPath := "non_existent_file.txt"
		_, err := FileToTarReader(nonExistentPath, &tarparameteroptions.FileToTarOptions{})
		require.Error(t, err)
	})
}
