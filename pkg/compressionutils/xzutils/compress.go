package xzutils

import (
	"context"
	"io"
	"os"

	"github.com/ulikunitz/xz"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// Compress compresses a file using xz and writes the result to 'outputPath'.
//
// This creates a plain xz compressed file (not a tar.xz archive).
func Compress(ctx context.Context, inputPath string, outputPath string) error {
	if inputPath == "" {
		return tracederrors.TracedErrorEmptyString("inputPath")
	}

	if outputPath == "" {
		return tracederrors.TracedErrorEmptyString("outputPath")
	}

	logging.LogInfoByCtxf(ctx, "Compress file '%s' to xz file '%s' started.", inputPath, outputPath)

	inputFile, err := os.Open(inputPath)
	if err != nil {
		return tracederrors.TracedErrorf("Failed to open input file '%s': %w", inputPath, err)
	}
	defer inputFile.Close()

	outFile, err := os.Create(outputPath)
	if err != nil {
		return tracederrors.TracedErrorf("Failed to create output file '%s': %w", outputPath, err)
	}
	defer outFile.Close()

	xzWriter, err := xz.NewWriter(outFile)
	if err != nil {
		os.Remove(outputPath)
		return tracederrors.TracedErrorf("Failed to create xz writer for '%s': %w", outputPath, err)
	}

	written, err := io.Copy(xzWriter, inputFile)
	if err != nil {
		xzWriter.Close()
		os.Remove(outputPath)
		return tracederrors.TracedErrorf("Failed to compress '%s' to '%s': %w", inputPath, outputPath, err)
	}

	if err := xzWriter.Close(); err != nil {
		os.Remove(outputPath)
		return tracederrors.TracedErrorf("Failed to finalize xz compression for '%s': %w", outputPath, err)
	}

	logging.LogInfoByCtxf(ctx, "Compressed '%d' bytes from '%s' to '%s'.", written, inputPath, outputPath)
	logging.LogInfoByCtxf(ctx, "Compress file '%s' to xz file '%s' finished.", inputPath, outputPath)

	return nil
}
