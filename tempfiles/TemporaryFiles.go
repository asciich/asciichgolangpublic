package tempfiles

import (
	"os"
	"strings"

	"github.com/asciich/asciichgolangpublic/files"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

func CreateEmptyTemporaryFile(verbose bool) (temporaryfile files.File, err error) {
	temporaryfile, err = CreateNamedTemporaryFile("emptyFile", verbose)
	if err != nil {
		return nil, err
	}

	return temporaryfile, nil
}

func CreateEmptyTemporaryFileAndGetPath(verbose bool) (temporaryFilePath string, err error) {
	temporaryFile, err := CreateEmptyTemporaryFile(verbose)
	if err != nil {
		return "", err
	}

	temporaryFilePath, err = temporaryFile.GetLocalPath()
	if err != nil {
		return "", err
	}

	return temporaryFilePath, nil
}

func CreateFromBytes(content []byte, verbose bool) (temporaryFile files.File, err error) {
	if content == nil {
		return nil, tracederrors.TracedErrorNil("content")
	}

	temporaryFile, err = CreateEmptyTemporaryFile(verbose)
	if err != nil {
		return nil, err
	}

	err = temporaryFile.WriteBytes(content, verbose)
	if err != nil {
		return nil, err
	}

	return temporaryFile, nil
}

func CreateFromString(content string, verbose bool) (temporaryFile files.File, err error) {
	temporaryFile, err = CreateEmptyTemporaryFile(verbose)
	if err != nil {
		return nil, err
	}

	err = temporaryFile.WriteString(content, verbose)
	if err != nil {
		return nil, err
	}

	return temporaryFile, nil
}

func CreateFromStringAndGetPath(content string, verbose bool) (temporaryFilePath string, err error) {
	temporaryFile, err := CreateFromString(content, verbose)
	if err != nil {
		return "", err
	}

	temporaryFilePath, err = temporaryFile.GetPath()
	if err != nil {
		return "", err
	}

	return temporaryFilePath, nil
}

func CreateNamedTemporaryFile(fileName string, verbose bool) (temporaryfile files.File, err error) {
	if fileName == "" {
		return nil, tracederrors.TracedErrorEmptyString("fileName")
	}

	osFile, err := os.CreateTemp("", fileName)
	if err != nil {
		return nil, err
	}

	temporaryfile, err = files.GetFileByOsFile(osFile)
	if err != nil {
		return nil, err
	}

	return temporaryfile, nil
}

func CreateTemporaryFileFromBytes(content []byte, verbose bool) (temporaryFile files.File, err error) {
	if content == nil {
		return nil, tracederrors.TracedErrorNil("content")
	}

	temporaryFile, err = CreateTemporaryFileFromString(string(content), verbose)
	if err != nil {
		return nil, err
	}

	return temporaryFile, nil
}

func CreateTemporaryFileFromFile(fileToCopyAsTemporaryFile files.File, verbose bool) (temporaryFile files.File, err error) {
	if fileToCopyAsTemporaryFile == nil {
		return nil, tracederrors.TracedErrorNil("fileToCopyAsTemporaryFile")
	}

	content, err := fileToCopyAsTemporaryFile.ReadAsBytes()
	if err != nil {
		return nil, err
	}

	temporaryFile, err = CreateTemporaryFileFromBytes(content, verbose)
	if err != nil {
		return nil, err
	}

	return temporaryFile, nil
}

func CreateTemporaryFileFromPath(verbose bool, filePathToCopyAsTemporaryFile ...string) (temporaryFile files.File, err error) {
	if len(filePathToCopyAsTemporaryFile) <= 0 {
		return nil, tracederrors.TracedError("filePathToCopyAsTemporaryFile")
	}

	pathToUse := strings.Join(filePathToCopyAsTemporaryFile, "/")

	fileToCopy, err := files.GetLocalFileByPath(pathToUse)
	if err != nil {
		return nil, err
	}

	temporaryFile, err = CreateTemporaryFileFromFile(fileToCopy, verbose)
	if err != nil {
		return nil, err
	}

	return temporaryFile, nil
}

func CreateTemporaryFileFromString(content string, verbose bool) (temporaryFile files.File, err error) {
	temporaryFile, err = CreateNamedTemporaryFile("tempFile", verbose)
	if err != nil {
		return nil, err
	}

	err = temporaryFile.WriteString(content, verbose)
	if err != nil {
		return nil, err
	}

	return temporaryFile, nil
}

func MustCreateEmptyTemporaryFile(verbose bool) (temporaryfile files.File) {
	temporaryfile, err := CreateEmptyTemporaryFile(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
	return temporaryfile
}

func MustCreateEmptyTemporaryFileAndGetPath(verbose bool) (temporaryFilePath string) {
	temporaryFilePath, err := CreateEmptyTemporaryFileAndGetPath(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return temporaryFilePath
}

func MustCreateFromBytes(content []byte, verbose bool) (temporaryFile files.File) {
	temporaryFile, err := CreateFromBytes(content, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return temporaryFile
}

func MustCreateFromString(content string, verbose bool) (temporaryFile files.File) {
	file, err := CreateFromString(content, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return file
}

func MustCreateFromStringAndGetPath(content string, verbose bool) (temporaryFilePath string) {
	temporaryFilePath, err := CreateFromStringAndGetPath(content, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return temporaryFilePath
}

func MustCreateNamedTemporaryFile(fileName string, verbose bool) (temporaryfile files.File) {
	temporaryfile, err := CreateNamedTemporaryFile(fileName, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return temporaryfile
}

func MustCreateTemporaryFileFromBytes(content []byte, verbose bool) (temporaryFile files.File) {
	temporaryFile, err := CreateTemporaryFileFromBytes(content, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return temporaryFile
}

func MustCreateTemporaryFileFromFile(fileToCopyAsTemporaryFile files.File, verbose bool) (temporaryFile files.File) {
	temporaryFile, err := CreateTemporaryFileFromFile(fileToCopyAsTemporaryFile, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return temporaryFile
}

func MustCreateTemporaryFileFromPath(verbose bool, filePathToCopyAsTemporaryFile ...string) (temporaryFile files.File) {
	temporaryFile, err := CreateTemporaryFileFromPath(verbose, filePathToCopyAsTemporaryFile...)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return temporaryFile
}

func MustCreateTemporaryFileFromString(content string, verbose bool) (temporaryFile files.File) {
	temporaryFile, err := CreateTemporaryFileFromString(content, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return temporaryFile
}
