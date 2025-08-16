package tarutils

import (
	"archive/tar"
	"bytes"
	"io"
	"os"
	"path"

	"github.com/asciich/asciichgolangpublic/pkg/archiveutils/tarutils/tarparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// FileToTarReader creates a tar archive in memory from a given file.
func FileToTarReader(localFilePath string, options *tarparameteroptions.FileToTarOptions) (*bytes.Buffer, error) {
	file, err := os.Open(localFilePath)
	if err != nil {
		return nil, tracederrors.TracedErrorf("failed to open file '%s': %w", localFilePath, err)
	}
	defer file.Close()

	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	defer tw.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, tracederrors.TracedErrorf("failed to get file info for '%s': %w", localFilePath, err)
	}

	header, err := tar.FileInfoHeader(fileInfo, "")
	if err != nil {
		return nil, tracederrors.TracedErrorf("failed to create tar header: %w", err)
	}
	
	// Set the name inside the tar:
	if options.OverrideFileName == "" {
		header.Name = path.Base(localFilePath)
	} else {
		name, err := options.GetOverrideFileName()
		if err != nil {
			return nil, err
		}

		header.Name = name
	}
	
	if err := tw.WriteHeader(header); err != nil {
		return nil, tracederrors.TracedErrorf("failed to write tar header: %w", err)
	}
	if _, err := io.Copy(tw, file); err != nil {
		return nil, tracederrors.TracedErrorf("failed to copy file to tar: %w", err)
	}

	return &buf, nil
}
