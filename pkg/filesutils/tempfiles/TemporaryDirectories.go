package tempfiles

import (
	"os"

	"github.com/asciich/asciichgolangpublic/files"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
)

func CreateEmptyTemporaryDirectory(verbose bool) (temporaryDirectory *files.LocalDirectory, err error) {
	dirPath, err := os.MkdirTemp("", "empty")
	if err != nil {
		return nil, err
	}

	temporaryDirectory, err = files.GetLocalDirectoryByPath(dirPath)
	if err != nil {
		return nil, err
	}

	return temporaryDirectory, nil
}

func CreateEmptyTemporaryDirectoryAndGetPath(verbose bool) (TemporaryDirectoryPath string, err error) {
	TemporaryDirectory, err := CreateEmptyTemporaryDirectory(verbose)
	if err != nil {
		return "", err
	}

	TemporaryDirectoryPath, err = TemporaryDirectory.GetLocalPath()
	if err != nil {
		return "", err
	}

	return TemporaryDirectoryPath, nil
}

func MustCreateEmptyTemporaryDirectory(verbose bool) (temporaryDirectory *files.LocalDirectory) {
	temporaryDirectory, err := CreateEmptyTemporaryDirectory(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
	return temporaryDirectory
}

func MustCreateEmptyTemporaryDirectoryAndGetPath(verbose bool) (TemporaryDirectoryPath string) {
	TemporaryDirectoryPath, err := CreateEmptyTemporaryDirectoryAndGetPath(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return TemporaryDirectoryPath
}
