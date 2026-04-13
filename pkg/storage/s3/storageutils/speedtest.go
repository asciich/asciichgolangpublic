package storageutils

import (
	"context"
	"errors"
	"io"
	"math/rand"
	"os"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/datatypes/bytesutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// Performs a simple speedtest by writing random data of size 'size' into the file 'path'.
//
//	Then the data is read again to measure the read speed.
//
// Hint: Use bytesutils.ParseSizeStringAsInt64("1GB") to define the size in a human readable way.
func RunSpeedTest(ctx context.Context, path string, size int64) (*SpeedTestResult, error) {
	if path == "" {
		return nil, tracederrors.TracedErrorEmptyString("path")
	}

	if size < 0 {
		return nil, tracederrors.TracedErrorf("Invalid size = '%d'.", size)
	}

	const chunkSize = 4 * 1024 * 1024

	sizeReadable, err := bytesutils.GetSizeAsHumanReadableString(size)
	if err != nil {
		return nil, err
	}

	chunkSizeReadable, err := bytesutils.GetSizeAsHumanReadableString(chunkSize)
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Storage speed test on file '%s' with write/read size '%s' started. ChunkSize is '%s'.", path, sizeReadable, chunkSizeReadable)

	err = nativefiles.Delete(ctx, path, &filesoptions.DeleteOptions{})
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Write test started...")

	tStart := time.Now()

	file, err := os.Create(path)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Error creating file: %w\n", err)
	}
	defer nativefiles.Delete(ctx, path, &filesoptions.DeleteOptions{})

	buf := make([]byte, chunkSize)
	rng := rand.New(rand.NewSource(42))

	written := 0

	for int64(written) < size {
		rng.Read(buf)

		n, err := file.Write(buf)
		if err != nil {
			err := nativefiles.Delete(ctx, path, &filesoptions.DeleteOptions{})
			if err != nil {
				return nil, err
			}
			return nil, tracederrors.TracedErrorf("Error writing file: %w\n", err)
		}
		written += n
	}

	// Flush OS write cache to get accurate disk write speed
	if err := file.Sync(); err != nil {
		file.Close()
		nativefiles.Delete(ctx, path, &filesoptions.DeleteOptions{})
		return nil, tracederrors.TracedErrorf("Error syncing file: %w\n", err)
	}
	file.Close()

	writeDuration := time.Since(tStart)
	writeSpeed := float64(size) / writeDuration.Seconds()
	logging.LogInfoByCtxf(ctx, "Write complete: %.2f MB/s (took %.2fs)\n\n", writeSpeed/(1024*1024), writeDuration.Seconds())

	file, err = os.Open(path)
	if err != nil {
		nativefiles.Delete(ctx, path, &filesoptions.DeleteOptions{})
		return nil, tracederrors.TracedErrorf("Error opening file: %w\n", err)
	}

	readBuf := make([]byte, chunkSize)
	totalRead := 0

	logging.LogInfoByCtxf(ctx, "Read test started...")

	tStart = time.Now()

	for {
		n, err := file.Read(readBuf)
		totalRead += n
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, tracederrors.TracedErrorf("Failed to read '%s': %w", path, err)
		}
	}
	file.Close()

	readDuration := time.Since(tStart)
	readSpeed := float64(totalRead) / readDuration.Seconds()
	logging.LogInfoByCtxf(ctx, "Read complete:  %.2f MB/s (took %.2fs)\n\n", readSpeed/(1024*1024), readDuration.Seconds())

	result := &SpeedTestResult{
		Size:       size,
		ChunkSize:  chunkSize,
		FileName:   path,
		WriteSpeed: writeSpeed,
		ReadSpeed:  readSpeed,
	}

	message, err := result.GetResultMessage()
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtx(ctx, message)

	logging.LogInfoByCtxf(ctx, "Storage speed test on file '%s' with write/read size '%s' finished. ChunkSize was '%s'.", path, sizeReadable, chunkSizeReadable)

	return result, nil
}
