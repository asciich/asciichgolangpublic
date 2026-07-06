package tarutils

import (
	"archive/tar"
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"io"
	"os"
	"sort"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/datatypes/slicesutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func AddFileFromFileContentBytesToTarArchiveBytes(archiveToExtend []byte, fileName string, content []byte) (tarBytes []byte, err error) {
	if archiveToExtend == nil {
		return nil, tracederrors.TracedErrorNil("archiveToExtend")
	}

	if fileName == "" {
		return nil, tracederrors.TracedErrorEmptyString("fileName")
	}

	if content == nil {
		return nil, tracederrors.TracedErrorNil("content")
	}

	sizeBeforeWriting := len(archiveToExtend)

	ioBuffer := bytes.NewBuffer(archiveToExtend)

	err = WriteFileContentBytesIntoWriter(ioBuffer, fileName, content)
	if err != nil {
		return nil, err
	}

	tarBytes = ioBuffer.Bytes()

	sizeAfterWriting := len(tarBytes)

	if sizeAfterWriting <= sizeBeforeWriting {
		return nil, tracederrors.TracedError("Internal error: archive size did not grow after writing.")
	}

	return tarBytes, err
}

func AddFileFromFileContentStringToTarArchiveBytes(archiveToExtend []byte, fileName string, content string) (tarBytes []byte, err error) {
	if archiveToExtend == nil {
		return nil, tracederrors.TracedErrorNil("archiveToExtend")
	}

	if fileName == "" {
		return nil, tracederrors.TracedErrorEmptyString("fileName")
	}

	if content == "" {
		return nil, tracederrors.TracedErrorEmptyString("content")
	}

	tarBytes, err = AddFileFromFileContentBytesToTarArchiveBytes(archiveToExtend, fileName, []byte(content))
	if err != nil {
		return nil, err
	}

	return tarBytes, nil
}

func CreateTarArchiveFromFileContentByteIntoWriter(ioWriter io.Writer, fileName string, content []byte) (err error) {
	if fileName == "" {
		return tracederrors.TracedErrorEmptyString("fileName")
	}

	if content == nil {
		return tracederrors.TracedErrorNil("content")
	}

	if ioWriter == nil {
		return tracederrors.TracedErrorNil("tarWriter")
	}

	err = WriteFileContentBytesIntoWriter(ioWriter, fileName, content)
	if err != nil {
		return err
	}

	return nil
}

func CreateTarArchiveFromFileContentStringAndGetAsBytes(fileName string, content string) (tarBytes []byte, err error) {
	if fileName == "" {
		return nil, tracederrors.TracedErrorEmptyString("fileName")
	}

	if content == "" {
		return nil, tracederrors.TracedErrorEmptyString("content")
	}

	var b bytes.Buffer
	tarWriter := bufio.NewWriter(&b)

	err = CreateTarArchiveFromFileContentStringIntoWriter(fileName, content, tarWriter)
	if err != nil {
		return nil, err
	}

	err = tarWriter.Flush()
	if err != nil {
		return nil, tracederrors.TracedErrorf("flush tarWriter failed: '%w'", err)
	}

	bufferReader := bufio.NewReader(&b)

	bufferLen := b.Len()

	tarBytes = make([]byte, bufferLen)
	nReadBytes, err := io.ReadFull(bufferReader, tarBytes)
	if err != nil {
		return nil, tracederrors.TracedErrorf("ReadFull on tarBytes failed: '%d'", err)
	}

	if bufferLen != nReadBytes {
		return nil, tracederrors.TracedErrorf(
			"Internal error: bufferLen '%d' does not match nReadBytes '%d'",
			bufferLen,
			nReadBytes,
		)
	}

	return tarBytes, nil
}

func CreateTarGzArchiveFromFileContentStringAndGetAsBytes(fileName string, content string) (tarGzBytes []byte, err error) {
	if fileName == "" {
		return nil, tracederrors.TracedErrorEmptyString("fileName")
	}

	if content == "" {
		return nil, tracederrors.TracedErrorEmptyString("content")
	}

	tarBytes, err := CreateTarArchiveFromFileContentStringAndGetAsBytes(fileName, content)
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer
	gzWriter := gzip.NewWriter(&b)

	_, err = gzWriter.Write(tarBytes)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to write tar bytes into gzip writer: '%w'", err)
	}

	err = gzWriter.Close()
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to close gzip writer: '%w'", err)
	}

	return b.Bytes(), nil
}

func CreateTarArchiveFromFileContentStringIntoWriter(fileName string, content string, ioWriter io.Writer) (err error) {
	if fileName == "" {
		return tracederrors.TracedErrorEmptyString("fileName")
	}

	if content == "" {
		return tracederrors.TracedErrorEmptyString("content")
	}

	if ioWriter == nil {
		return tracederrors.TracedErrorNil("tarWriter")
	}

	err = CreateTarArchiveFromFileContentByteIntoWriter(ioWriter, fileName, []byte(content))
	if err != nil {
		return err
	}

	return nil
}

func ListFileNamesFromTarArchiveBytes(archiveBytes []byte) (fileNames []string, err error) {
	if archiveBytes == nil {
		return nil, tracederrors.TracedErrorNil("archiveBytes")
	}

	bytesReader := bytes.NewReader(archiveBytes)

	fileNames = []string{}
	tarReader := tar.NewReader(bytesReader)
	for {
		header, err := tarReader.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return nil, tracederrors.TracedErrorf(
				"Reading next header in tar archive failed: %w",
				err,
			)
		}

		fileNames = append(fileNames, header.Name)
	}

	fileNames = slicesutils.RemoveEmptyStrings(fileNames)

	sort.Strings(fileNames)

	return fileNames, nil
}

func ReadFileFromTarArchiveBytesAsBytes(archiveBytes []byte, fileNameToRead string) (content []byte, err error) {
	if archiveBytes == nil {
		return nil, tracederrors.TracedErrorNil("archiveBytes")
	}

	if fileNameToRead == "" {
		return nil, tracederrors.TracedErrorEmptyString("fileNameToRead")
	}

	if len(archiveBytes) <= 0 {
		return nil, tracederrors.TracedErrorf(
			"Unable to read '%s' from empty tar archive. len(archiveBytes) is 0.",
			fileNameToRead,
		)
	}

	bytesReader := bytes.NewReader(archiveBytes)

	tarReader := tar.NewReader(bytesReader)
	for {
		header, err := tarReader.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return nil, tracederrors.TracedErrorf(
				"Reading next header in tar archive failed: %w",
				err,
			)
		}

		if header.Name == fileNameToRead {
			content = make([]byte, header.Size)
			nBytesRead, err := tarReader.Read(content)
			if err != nil {
				if errors.Is(err, io.EOF) {
					// Reaching end of file can be ok when the last file in the arive is requested.
					// If file was finished to early the next check will catch that error by
					// comparing the expected to the actual read bytes.
				} else {
					return nil, tracederrors.TracedErrorf(
						"Reading file in tar archive failed: %w",
						err,
					)
				}
			}

			if nBytesRead != int(header.Size) {
				return nil, tracederrors.TracedErrorf(
					"Reading file '%s' from tar archive bytes failed. Expected '%d' bytes to read but go '%d'.",
					fileNameToRead,
					header.Size,
					nBytesRead,
				)
			}

			break
		}
	}

	if content == nil {
		return nil, tracederrors.TracedErrorf(
			"Unable to read '%s' from given archive bytes. '%s' was not found.",
			fileNameToRead,
			fileNameToRead,
		)
	}

	return content, nil
}

func ReadFileFromTarArchiveBytesAsString(archiveBytes []byte, fileNameToRead string) (content string, err error) {
	contentBytes, err := ReadFileFromTarArchiveBytesAsBytes(archiveBytes, fileNameToRead)
	if err != nil {
		return "", err
	}

	return string(contentBytes), nil
}

func ReadFileFromTarArchiveAsBytes(ctx context.Context, archivePath string, fileName string) ([]byte, error) {
	if archivePath == "" {
		return nil, tracederrors.TracedErrorEmptyString("arhivePath")
	}

	if fileName == "" {
		return nil, tracederrors.TracedErrorEmptyString("fileName")
	}

	logging.LogInfoByCtxf(ctx, "Read '%s' from tar archive '%s' started.", fileName, archivePath)

	archiveFile, err := os.Open(archivePath)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to open tar archive '%s': %w", archivePath, err)
	}
	defer archiveFile.Close()

	// Detect gzip by magic bytes (0x1f 0x8b) instead of relying on file extension,
	// since temp files may not have a .tar.gz suffix.
	magicBytes := make([]byte, 2)
	_, err = io.ReadFull(archiveFile, magicBytes)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to read magic bytes from '%s': %w", archivePath, err)
	}

	// Seek back to the beginning after reading the magic bytes.
	_, err = archiveFile.Seek(0, io.SeekStart)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to seek back to start of '%s': %w", archivePath, err)
	}

	var tarReader *tar.Reader
	isGzip := magicBytes[0] == 0x1f && magicBytes[1] == 0x8b
	if isGzip {
		gzReader, err := gzip.NewReader(archiveFile)
		if err != nil {
			return nil, tracederrors.TracedErrorf("Failed to create gzip reader for '%s': %w", archivePath, err)
		}
		defer gzReader.Close()

		tarReader = tar.NewReader(gzReader)
	} else {
		tarReader = tar.NewReader(archiveFile)
	}

	var content []byte
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, tracederrors.TracedErrorf("Failed to read tar archive '%s': %w", archivePath, err)
		}

		if header.Name == fileName {
			if header.Typeflag != tar.TypeReg {
				return nil, tracederrors.TracedErrorf(
					"Entry '%s' in archive '%s' is not a regular file",
					fileName, archivePath,
				)
			}

			content, err = io.ReadAll(tarReader)
			if err != nil {
				return nil, tracederrors.TracedErrorf(
					"Failed to read file '%s' from tar archive '%s': %w",
					fileName, archivePath, err,
				)
			}

			break
		}
	}

	if content == nil {
		return nil, tracederrors.TracedErrorf(
			"File '%s' not found in tar archive '%s'",
			fileName, archivePath,
		)
	}

	logging.LogInfoByCtxf(ctx, "Read '%s' from tar archive '%s' finished.", fileName, archivePath)

	return content, nil
}

func ExtractFileFromTarArchive(ctx context.Context, archivePath string, fileName string, destPath string) error {
	if archivePath == "" {
		return tracederrors.TracedErrorEmptyString("archivePath")
	}

	if fileName == "" {
		return tracederrors.TracedErrorEmptyString("fileName")
	}

	if destPath == "" {
		return tracederrors.TracedErrorEmptyString("destPath")
	}

	logging.LogInfoByCtxf(ctx, "Extract '%s' from tar archive '%s' to '%s' started.", fileName, archivePath, destPath)

	content, err := ReadFileFromTarArchiveAsBytes(ctx, archivePath, fileName)
	if err != nil {
		return err
	}

	destFile, err := os.Create(destPath)
	if err != nil {
		return tracederrors.TracedErrorf("Failed to create destination file '%s': %w", destPath, err)
	}
	defer destFile.Close()

	_, err = destFile.Write(content)
	if err != nil {
		return tracederrors.TracedErrorf("Failed to write to destination file '%s': %w", destPath, err)
	}

	logging.LogInfoByCtxf(ctx, "Extract '%s' from tar archive '%s' to '%s' finished.", fileName, archivePath, destPath)

	return nil
}

func WriteFileContentBytesIntoWriter(ioWriter io.Writer, fileName string, content []byte) (err error) {
	if fileName == "" {
		return tracederrors.TracedErrorEmptyString("fileName")
	}

	if content == nil {
		return tracederrors.TracedErrorNil("content")
	}

	if ioWriter == nil {
		return tracederrors.TracedErrorNil("tarWriter")
	}

	fileSize := int64(len(content))

	header := &tar.Header{
		Name:    fileName,
		Size:    fileSize,
		ModTime: time.Now(),
	}

	tarWriter := tar.NewWriter(ioWriter)

	err = tarWriter.WriteHeader(header)
	if err != nil {
		return tracederrors.TracedErrorf(
			"Write tar header failed: '%w'", err,
		)
	}

	contentReader := bytes.NewReader(content)

	nBytesWritten, err := io.Copy(ioWriter, contentReader)
	if err != nil {
		return tracederrors.TracedErrorf(
			"Write tar body failed: '%w'", err,
		)
	}

	requiredPaddingBytes := 512 - (len(content) % 512)
	writtenPaddingBytes, err := ioWriter.Write(make([]byte, requiredPaddingBytes))
	if err != nil {
		return tracederrors.TracedErrorf(
			"Write padding bytes failed: '%w'", err,
		)
	}

	if requiredPaddingBytes != writtenPaddingBytes {
		return tracederrors.TracedErrorf(
			"writting tar padding bytes failed. Expected to write '%d' bytes but '%d' were written",
			requiredPaddingBytes,
			writtenPaddingBytes,
		)
	}

	if nBytesWritten != fileSize {
		return tracederrors.TracedErrorf(
			"writing '%s' to tar archive failed. Content to write has len '%d' but '%d' bytes were written",
			fileName,
			len(content),
			nBytesWritten,
		)
	}

	return nil
}
