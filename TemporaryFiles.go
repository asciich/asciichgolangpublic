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
	temporaryfile, err = t.CreateNamedTemporaryFile("emptyFile", verbose)
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

func (t *TemporaryFilesService) CreateNamedTemporaryFile(fileName string, verbose bool) (temporaryfile File, err error) {
	if fileName == "" {
		return nil, TracedErrorEmptyString("fileName")
	}

	osFile, err := os.CreateTemp("", fileName)
	if err != nil {
		return nil, err
	}

	temporaryfile, err = GetFileByOsFile(osFile)
	if err != nil {
		return nil, err
	}

	return temporaryfile, nil
}

func (t *TemporaryFilesService) CreateTemporaryFileFromBytes(content []byte, verbose bool) (temporaryFile File, err error) {
	if content == nil {
		return nil, TracedErrorNil("content")
	}

	temporaryFile, err = t.CreateTemporaryFileFromString(string(content), verbose)
	if err != nil {
		return nil, err
	}

	return temporaryFile, nil
}

func (t *TemporaryFilesService) CreateTemporaryFileFromFile(fileToCopyAsTemporaryFile File, verbose bool) (temporaryFile File, err error) {
	if fileToCopyAsTemporaryFile == nil {
		return nil, TracedErrorNil("fileToCopyAsTemporaryFile")
	}

	content, err := fileToCopyAsTemporaryFile.ReadAsBytes()
	if err != nil {
		return nil, err
	}

	temporaryFile, err = t.CreateTemporaryFileFromBytes(content, verbose)
	if err != nil {
		return nil, err
	}

	return temporaryFile, nil
}

func (t *TemporaryFilesService) CreateTemporaryFileFromString(content string, verbose bool) (temporaryFile File, err error) {
	temporaryFile, err = t.CreateNamedTemporaryFile("tempFile", verbose)
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

func (t *TemporaryFilesService) MustCreateNamedTemporaryFile(fileName string, verbose bool) (temporaryfile File) {
	temporaryfile, err := t.CreateNamedTemporaryFile(fileName, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return temporaryfile
}

func (t *TemporaryFilesService) MustCreateTemporaryFileFromBytes(content []byte, verbose bool) (temporaryFile File) {
	temporaryFile, err := t.CreateTemporaryFileFromBytes(content, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return temporaryFile
}

func (t *TemporaryFilesService) MustCreateTemporaryFileFromFile(fileToCopyAsTemporaryFile File, verbose bool) (temporaryFile File) {
	temporaryFile, err := t.CreateTemporaryFileFromFile(fileToCopyAsTemporaryFile, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return temporaryFile
}

func (t *TemporaryFilesService) MustCreateTemporaryFileFromString(content string, verbose bool) (temporaryFile File) {
	temporaryFile, err := t.CreateTemporaryFileFromString(content, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return temporaryFile
}
