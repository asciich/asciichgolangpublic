package asciichgolangpublic

import (
	"github.com/asciich/asciichgolangpublic/files"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/tempfiles"
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
	tempDir, err := tempfiles.CreateEmptyTemporaryDirectory(verbose)
	if err != nil {
		return nil, err
	}

	localRepository, err := GetLocalGitReposioryFromLocalDirectory(tempDir)
	if err != nil {
		return nil, err
	}

	err = localRepository.Init(
		&parameteroptions.CreateRepositoryOptions{
			Verbose: verbose,
		},
	)
	if err != nil {
		return nil, err
	}

	temporaryGitRepository = localRepository

	return temporaryGitRepository, nil
}

func (g *TemporaryGitRepositoriesService) CreateTemporaryGitRepositoryAndAddDataFromDirectory(dataToAdd files.Directory, verbose bool) (temporaryRepository GitRepository, err error) {
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

func (g *TemporaryGitRepositoriesService) MustCreateTemporaryGitRepositoryAndAddDataFromDirectory(dataToAdd files.Directory, verbose bool) (temporaryRepository GitRepository) {
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

func (t TemporaryGitRepositoriesService) CreateEmptyTemporaryGitRepository(createRepoOptions *parameteroptions.CreateRepositoryOptions) (temporaryGitRepository GitRepository, err error) {
	if createRepoOptions == nil {
		return nil, tracederrors.TracedErrorNil("createRepoOptions")
	}

	tempDirectory, err := tempfiles.CreateEmptyTemporaryDirectory(createRepoOptions.Verbose)
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

func (t TemporaryGitRepositoriesService) MustCreateEmptyTemporaryGitRepository(createRepoOptions *parameteroptions.CreateRepositoryOptions) (temporaryGitRepository GitRepository) {
	temporaryGitRepository, err := t.CreateEmptyTemporaryGitRepository(createRepoOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return temporaryGitRepository
}
