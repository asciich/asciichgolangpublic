package xzutils

import (
	"context"
	"io"
	"os"

	"github.com/ulikunitz/xz"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// Decompress decompresses a plain xz compressed file to 'outputPath'.
//
// This is for files compressed with xz only (not tar.xz archives).
func Decompress(ctx context.Context, archivePath string, outputPath string) error {
	if archivePath == "" {
		return tracederrors.TracedErrorEmptyString("archivePath")
	}

	if outputPath == "" {
		return tracederrors.TracedErrorEmptyString("outputPath")
	}

	logging.LogInfoByCtxf(ctx, "Decompress xz file '%s' to '%s' started.", archivePath, outputPath)

	archiveFile, err := os.Open(archivePath)
	if err != nil {
		return tracederrors.TracedErrorf("Failed to open xz file '%s': %w", archivePath, err)
	}
	defer archiveFile.Close()

	xzReader, err := xz.NewReader(archiveFile)
	if err != nil {
		return tracederrors.TracedErrorf("Failed to create xz reader for '%s': %w", archivePath, err)
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		return tracederrors.TracedErrorf("Failed to create output file '%s': %w", outputPath, err)
	}
	defer outFile.Close()

	written, err := io.Copy(outFile, xzReader)
	if err != nil {
		os.Remove(outputPath)
		return tracederrors.TracedErrorf("Failed to decompress '%s' to '%s': %w", archivePath, outputPath, err)
	}

	logging.LogInfoByCtxf(ctx, "Decompressed '%d' bytes to '%s'.", written, outputPath)
	logging.LogInfoByCtxf(ctx, "Decompress xz file '%s' to '%s' finished.", archivePath, outputPath)

	return nil
}
