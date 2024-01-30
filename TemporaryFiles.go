package asciichgolangpublic

import (
	"os"
)

type TemporaryFilesService struct {
}

func NewTemporaryFilesService() (t *TemporaryFilesService) {
	return new(TemporaryFilesService)
}

func TemporaryFiles() (temporaryFiles *TemporaryFilesService) {
	return new(TemporaryFilesService)
}

func (t *TemporaryFilesService) CreateEmptyTemporaryFile(verbose bool) (temporaryfile File, err error) {
	osFile, err := os.CreateTemp("", "emptyFile")
	if err != nil {
		return nil, err
	}

	temporaryfile, err = GetFileByOsFile(osFile)
	if err != nil {
		return nil, err
	}

	return temporaryfile, nil
}

func (t *TemporaryFilesService) CreateEmptyTemporaryFileAndGetPath(verbose bool) (temporaryFilePath string, err error) {
	temporaryFile, err := t.CreateEmptyTemporaryFile(verbose)
	if err != nil {
		return "", err
	}

	temporaryFilePath, err = temporaryFile.GetLocalPath()
	if err != nil {
		return "", err
	}

	return temporaryFilePath, nil
}

func (t *TemporaryFilesService) CreateFromString(content string, verbose bool) (temporaryFile File, err error) {
	temporaryFile, err = t.CreateEmptyTemporaryFile(verbose)
	if err != nil {
		return nil, err
	}

	err = temporaryFile.WriteString(content, verbose)
	if err != nil {
		return nil, err
	}

	return temporaryFile, nil
}

func (t *TemporaryFilesService) MustCreateEmptyTemporaryFile(verbose bool) (temporaryfile File) {
	temporaryfile, err := t.CreateEmptyTemporaryFile(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
	return temporaryfile
}

func (t *TemporaryFilesService) MustCreateEmptyTemporaryFileAndGetPath(verbose bool) (temporaryFilePath string) {
	temporaryFilePath, err := t.CreateEmptyTemporaryFileAndGetPath(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return temporaryFilePath
}

func (t *TemporaryFilesService) MustCreateFromString(content string, verbose bool) (temporaryFile File) {
	file, err := t.CreateFromString(content, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return file
}
