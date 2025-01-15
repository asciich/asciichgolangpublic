package asciichgolangpublic

import (
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type TemporaryGitRepositoriesService struct {
}

func NewTemporaryGitRepositoriesService() (temporaryGitRepositories *TemporaryGitRepositoriesService) {
	return new(TemporaryGitRepositoriesService)
}

func TemporaryGitRepositories() (temporaryDirectoriesService *TemporaryGitRepositoriesService) {
	return NewTemporaryGitRepositoriesService()
}

func (g *TemporaryGitRepositoriesService) CreateTemporaryGitRepository(verbose bool) (temporaryGitRepository GitRepository, err error) {
	tempDir, err := TemporaryDirectories().CreateEmptyTemporaryDirectory(verbose)
	if err != nil {
		return nil, err
	}

	localRepository, err := GetLocalGitReposioryFromLocalDirectory(tempDir)
	if err != nil {
		return nil, err
	}

	err = localRepository.Init(
		&CreateRepositoryOptions{
			Verbose: verbose,
		},
	)
	if err != nil {
		return nil, err
	}

	temporaryGitRepository = localRepository

	return temporaryGitRepository, nil
}

func (g *TemporaryGitRepositoriesService) CreateTemporaryGitRepositoryAndAddDataFromDirectory(dataToAdd Directory, verbose bool) (temporaryRepository GitRepository, err error) {
	if dataToAdd == nil {
		return nil, tracederrors.TracedError("dataToAdd is nil")
	}

	temporaryRepository, err = g.CreateTemporaryGitRepository(verbose)
	if err != nil {
		return nil, err
	}

	destDir, err := temporaryRepository.GetAsLocalDirectory()
	if err != nil {
		return nil, err
	}

	err = dataToAdd.CopyContentToDirectory(destDir, verbose)
	if err != nil {
		return nil, err
	}

	return temporaryRepository, nil
}

func (g *TemporaryGitRepositoriesService) MustCreateTemporaryGitRepositoryAndAddDataFromDirectory(dataToAdd Directory, verbose bool) (temporaryRepository GitRepository) {
	temporaryRepository, err := g.CreateTemporaryGitRepositoryAndAddDataFromDirectory(dataToAdd, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return temporaryRepository
}

func (g TemporaryGitRepositoriesService) MustCreateTemporaryGitRepository(verbose bool) (temporaryGitRepository GitRepository) {
	temporaryGitRepository, err := g.CreateTemporaryGitRepository(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return temporaryGitRepository
}
