package asciichgolangpublic

import (
	"os"

	"github.com/asciich/asciichgolangpublic/errors"
	"github.com/asciich/asciichgolangpublic/logging"
)

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

func (t *TemporaryDirectoriesService) CreateEmptyTemporaryGitRepository(createRepoOptions *CreateRepositoryOptions) (temporaryGitRepository GitRepository, err error) {
	if createRepoOptions == nil {
		return nil, errors.TracedErrorNil("createRepoOptions")
	}

	tempDirectory, err := t.CreateEmptyTemporaryDirectory(createRepoOptions.Verbose)
	if err != nil {
		return nil, err
	}

	temporaryGitRepository, err = GetLocalGitReposioryFromDirectory(tempDirectory)
	if err != nil {
		return nil, err
	}

	err = temporaryGitRepository.Init(createRepoOptions)
	if err != nil {
		return nil, err
	}

	repoPath, err := tempDirectory.GetLocalPath()
	if err != nil {
		return nil, err
	}

	if createRepoOptions.Verbose {
		logging.LogInfof("Created temporary local git repository '%s'.", repoPath)
	}

	return temporaryGitRepository, err
}

func (t *TemporaryDirectoriesService) MustCreateEmptyTemporaryDirectory(verbose bool) (temporaryDirectory *LocalDirectory) {
	temporaryDirectory, err := t.CreateEmptyTemporaryDirectory(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
	return temporaryDirectory
}

func (t *TemporaryDirectoriesService) MustCreateEmptyTemporaryDirectoryAndGetPath(verbose bool) (TemporaryDirectoryPath string) {
	TemporaryDirectoryPath, err := t.CreateEmptyTemporaryDirectoryAndGetPath(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return TemporaryDirectoryPath
}

func (t *TemporaryDirectoriesService) MustCreateEmptyTemporaryGitRepository(createRepoOptions *CreateRepositoryOptions) (temporaryGitRepository GitRepository) {
	temporaryGitRepository, err := t.CreateEmptyTemporaryGitRepository(createRepoOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return temporaryGitRepository
}
