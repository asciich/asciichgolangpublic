package tarutils

import (
	"archive/tar"
	"bufio"
	"bytes"
	"errors"
	"io"
	"sort"
	"time"

	"gitlab.asciich.ch/tools/asciichgolangpublic.git/datatypes/slicesutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/tracederrors"
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
