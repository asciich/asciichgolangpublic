package asciichgolangpublic

import "os"

type TemporaryDirectoriesService struct {
}

func NewTemporaryDirectoriesService() (t *TemporaryDirectoriesService) {
	return new(TemporaryDirectoriesService)
}

func TemporaryDirectories() (TemporaryDirectorys *TemporaryDirectoriesService) {
	return new(TemporaryDirectoriesService)
}

func (t *TemporaryDirectoriesService) CreateEmptyTemporaryDirectory(verbose bool) (temporaryDirectory *LocalDirectory, err error) {
	dirPath, err := os.MkdirTemp("", "empty")
	if err != nil {
		return nil, err
	}

	temporaryDirectory, err = GetLocalDirectoryByPath(dirPath)
	if err != nil {
		return nil, err
	}

	return temporaryDirectory, nil
}

func (t *TemporaryDirectoriesService) CreateEmptyTemporaryDirectoryAndGetPath(verbose bool) (TemporaryDirectoryPath string, err error) {
	TemporaryDirectory, err := t.CreateEmptyTemporaryDirectory(verbose)
	if err != nil {
		return "", err
	}

	TemporaryDirectoryPath, err = TemporaryDirectory.GetLocalPath()
	if err != nil {
		return "", err
	}

	return TemporaryDirectoryPath, nil
}

func (t *TemporaryDirectoriesService) MustCreateEmptyTemporaryDirectory(verbose bool) (temporaryDirectory *LocalDirectory) {
	temporaryDirectory, err := t.CreateEmptyTemporaryDirectory(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
	return temporaryDirectory
}

func (t *TemporaryDirectoriesService) MustCreateEmptyTemporaryDirectoryAndGetPath(verbose bool) (TemporaryDirectoryPath string) {
	TemporaryDirectoryPath, err := t.CreateEmptyTemporaryDirectoryAndGetPath(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return TemporaryDirectoryPath
}
