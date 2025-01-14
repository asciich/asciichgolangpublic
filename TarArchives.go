package asciichgolangpublic

import (
	"archive/tar"
	"bufio"
	"bytes"
	"errors"
	"io"
	"sort"
	"time"

	aslices "github.com/asciich/asciichgolangpublic/datatypes/slices"
	aerrors "github.com/asciich/asciichgolangpublic/errors"
	"github.com/asciich/asciichgolangpublic/logging"
)

type TarArchivesService struct {
}

func NewTarArchivesService() (t *TarArchivesService) {
	return new(TarArchivesService)
}

func TarArchives() (t *TarArchivesService) {
	return NewTarArchivesService()
}

func (t *TarArchivesService) AddFileFromFileContentBytesToTarArchiveBytes(archiveToExtend []byte, fileName string, content []byte) (tarBytes []byte, err error) {
	if archiveToExtend == nil {
		return nil, aerrors.TracedErrorNil("archiveToExtend")
	}

	if fileName == "" {
		return nil, aerrors.TracedErrorEmptyString("fileName")
	}

	if content == nil {
		return nil, aerrors.TracedErrorNil("content")
	}

	sizeBeforeWriting := len(archiveToExtend)

	ioBuffer := bytes.NewBuffer(archiveToExtend)

	err = t.WriteFileContentBytesIntoWriter(ioBuffer, fileName, content)
	if err != nil {
		return nil, err
	}

	tarBytes = ioBuffer.Bytes()

	sizeAfterWriting := len(tarBytes)

	if sizeAfterWriting <= sizeBeforeWriting {
		return nil, aerrors.TracedError("Internal error: archive size did not grow after writing.")
	}

	return tarBytes, err
}

func (t *TarArchivesService) AddFileFromFileContentStringToTarArchiveBytes(archiveToExtend []byte, fileName string, content string) (tarBytes []byte, err error) {
	if archiveToExtend == nil {
		return nil, aerrors.TracedErrorNil("archiveToExtend")
	}

	if fileName == "" {
		return nil, aerrors.TracedErrorEmptyString("fileName")
	}

	if content == "" {
		return nil, aerrors.TracedErrorEmptyString("content")
	}

	tarBytes, err = t.AddFileFromFileContentBytesToTarArchiveBytes(archiveToExtend, fileName, []byte(content))
	if err != nil {
		return nil, err
	}

	return tarBytes, nil
}

func (t *TarArchivesService) CreateTarArchiveFromFileContentByteIntoWriter(ioWriter io.Writer, fileName string, content []byte) (err error) {
	if fileName == "" {
		return aerrors.TracedErrorEmptyString("fileName")
	}

	if content == nil {
		return aerrors.TracedErrorNil("content")
	}

	if ioWriter == nil {
		return aerrors.TracedErrorNil("tarWriter")
	}

	err = t.WriteFileContentBytesIntoWriter(ioWriter, fileName, content)
	if err != nil {
		return err
	}

	return nil
}

func (t *TarArchivesService) CreateTarArchiveFromFileContentStringAndGetAsBytes(fileName string, content string) (tarBytes []byte, err error) {
	if fileName == "" {
		return nil, aerrors.TracedErrorEmptyString("fileName")
	}

	if content == "" {
		return nil, aerrors.TracedErrorEmptyString("content")
	}

	var b bytes.Buffer
	tarWriter := bufio.NewWriter(&b)

	err = t.CreateTarArchiveFromFileContentStringIntoWriter(fileName, content, tarWriter)
	if err != nil {
		return nil, err
	}

	err = tarWriter.Flush()
	if err != nil {
		return nil, aerrors.TracedErrorf("flush tarWriter failed: '%w'", err)
	}

	bufferReader := bufio.NewReader(&b)

	bufferLen := b.Len()

	tarBytes = make([]byte, bufferLen)
	nReadBytes, err := io.ReadFull(bufferReader, tarBytes)
	if err != nil {
		return nil, aerrors.TracedErrorf("ReadFull on tarBytes failed: '%d'", err)
	}

	if bufferLen != nReadBytes {
		return nil, aerrors.TracedErrorf(
			"Internal error: bufferLen '%d' does not match nReadBytes '%d'",
			bufferLen,
			nReadBytes,
		)
	}

	return tarBytes, nil
}

func (t *TarArchivesService) CreateTarArchiveFromFileContentStringIntoWriter(fileName string, content string, ioWriter io.Writer) (err error) {
	if fileName == "" {
		return aerrors.TracedErrorEmptyString("fileName")
	}

	if content == "" {
		return aerrors.TracedErrorEmptyString("content")
	}

	if ioWriter == nil {
		return aerrors.TracedErrorNil("tarWriter")
	}

	err = t.CreateTarArchiveFromFileContentByteIntoWriter(ioWriter, fileName, []byte(content))
	if err != nil {
		return err
	}

	return nil
}

func (t *TarArchivesService) ListFileNamesFromTarArchiveBytes(archiveBytes []byte) (fileNames []string, err error) {
	if archiveBytes == nil {
		return nil, aerrors.TracedErrorNil("archiveBytes")
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

			return nil, aerrors.TracedErrorf(
				"Reading next header in tar archive failed: %w",
				err,
			)
		}

		fileNames = append(fileNames, header.Name)
	}

	fileNames = aslices.RemoveEmptyStrings(fileNames)

	sort.Strings(fileNames)

	return fileNames, nil
}

func (t *TarArchivesService) MustAddFileFromFileContentBytesToTarArchiveBytes(archiveToExtend []byte, fileName string, content []byte) (tarBytes []byte) {
	tarBytes, err := t.AddFileFromFileContentBytesToTarArchiveBytes(archiveToExtend, fileName, content)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return tarBytes
}

func (t *TarArchivesService) MustAddFileFromFileContentStringToTarArchiveBytes(archiveToExtend []byte, fileName string, content string) (tarBytes []byte) {
	tarBytes, err := t.AddFileFromFileContentStringToTarArchiveBytes(archiveToExtend, fileName, content)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return tarBytes
}

func (t *TarArchivesService) MustCreateTarArchiveFromFileContentByteIntoWriter(fileName string, content []byte, ioWriter io.Writer) {
	err := t.CreateTarArchiveFromFileContentByteIntoWriter(ioWriter, fileName, content)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (t *TarArchivesService) MustCreateTarArchiveFromFileContentStringAndGetAsBytes(fileName string, content string) (tarBytes []byte) {
	tarBytes, err := t.CreateTarArchiveFromFileContentStringAndGetAsBytes(fileName, content)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return tarBytes
}

func (t *TarArchivesService) MustCreateTarArchiveFromFileContentStringIntoWriter(fileName string, content string, ioWriter io.Writer) {
	err := t.CreateTarArchiveFromFileContentStringIntoWriter(fileName, content, ioWriter)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (t *TarArchivesService) MustListFileNamesFromTarArchiveBytes(archiveBytes []byte) (fileNames []string) {
	fileNames, err := t.ListFileNamesFromTarArchiveBytes(archiveBytes)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return fileNames
}

func (t *TarArchivesService) MustReadFileFromTarArchiveBytesAsBytes(archiveBytes []byte, fileNameToRead string) (content []byte) {
	content, err := t.ReadFileFromTarArchiveBytesAsBytes(archiveBytes, fileNameToRead)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return content
}

func (t *TarArchivesService) MustReadFileFromTarArchiveBytesAsString(archiveBytes []byte, fileNameToRead string) (content string) {
	content, err := t.ReadFileFromTarArchiveBytesAsString(archiveBytes, fileNameToRead)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return content
}

func (t *TarArchivesService) MustWriteFileContentBytesIntoWriter(ioWriter io.Writer, fileName string, content []byte) {
	err := t.WriteFileContentBytesIntoWriter(ioWriter, fileName, content)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (t *TarArchivesService) ReadFileFromTarArchiveBytesAsBytes(archiveBytes []byte, fileNameToRead string) (content []byte, err error) {
	if archiveBytes == nil {
		return nil, aerrors.TracedErrorNil("archiveBytes")
	}

	if fileNameToRead == "" {
		return nil, aerrors.TracedErrorEmptyString("fileNameToRead")
	}

	if len(archiveBytes) <= 0 {
		return nil, aerrors.TracedErrorf(
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

			return nil, aerrors.TracedErrorf(
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
					return nil, aerrors.TracedErrorf(
						"Reading file in tar archive failed: %w",
						err,
					)
				}
			}

			if nBytesRead != int(header.Size) {
				return nil, aerrors.TracedErrorf(
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
		return nil, aerrors.TracedErrorf(
			"Unable to read '%s' from given archive bytes. '%s' was not found.",
			fileNameToRead,
			fileNameToRead,
		)
	}

	return content, nil
}

func (t *TarArchivesService) ReadFileFromTarArchiveBytesAsString(archiveBytes []byte, fileNameToRead string) (content string, err error) {
	contentBytes, err := t.ReadFileFromTarArchiveBytesAsBytes(archiveBytes, fileNameToRead)
	if err != nil {
		return "", err
	}

	return string(contentBytes), nil
}

func (t *TarArchivesService) WriteFileContentBytesIntoWriter(ioWriter io.Writer, fileName string, content []byte) (err error) {
	if fileName == "" {
		return aerrors.TracedErrorEmptyString("fileName")
	}

	if content == nil {
		return aerrors.TracedErrorNil("content")
	}

	if ioWriter == nil {
		return aerrors.TracedErrorNil("tarWriter")
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
		return aerrors.TracedErrorf(
			"Write tar header failed: '%w'", err,
		)
	}

	contentReader := bytes.NewReader(content)

	nBytesWritten, err := io.Copy(ioWriter, contentReader)
	if err != nil {
		return aerrors.TracedErrorf(
			"Write tar body failed: '%w'", err,
		)
	}

	requiredPaddingBytes := 512 - (len(content) % 512)
	writtenPaddingBytes, err := ioWriter.Write(make([]byte, requiredPaddingBytes))
	if err != nil {
		return aerrors.TracedErrorf(
			"Write padding bytes failed: '%w'", err,
		)
	}

	if requiredPaddingBytes != writtenPaddingBytes {
		return aerrors.TracedErrorf(
			"writting tar padding bytes failed. Expected to write '%d' bytes but '%d' were written",
			requiredPaddingBytes,
			writtenPaddingBytes,
		)
	}

	if nBytesWritten != fileSize {
		return aerrors.TracedErrorf(
			"writing '%s' to tar archive failed. Content to write has len '%d' but '%d' bytes were written",
			fileName,
			len(content),
			nBytesWritten,
		)
	}

	return nil
}
