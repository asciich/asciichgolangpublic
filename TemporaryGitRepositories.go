package asciichgolangpublic

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfilesoo"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type TemporaryGitRepositoriesService struct {
}

func NewTemporaryGitRepositoriesService() (temporaryGitRepositories *TemporaryGitRepositoriesService) {
	return new(TemporaryGitRepositoriesService)
}

func TemporaryGitRepositories() (temporaryDirectoriesService *TemporaryGitRepositoriesService) {
	return NewTemporaryGitRepositoriesService()
}

func (g *TemporaryGitRepositoriesService) CreateTemporaryGitRepository(ctx context.Context) (temporaryGitRepository GitRepository, err error) {
	tempDir, err := tempfilesoo.CreateEmptyTemporaryDirectory(ctx)
	if err != nil {
		return nil, err
	}

	localRepository, err := GetLocalGitReposioryFromLocalDirectory(tempDir)
	if err != nil {
		return nil, err
	}

	err = localRepository.Init(ctx, &parameteroptions.CreateRepositoryOptions{})
	if err != nil {
		return nil, err
	}

	temporaryGitRepository = localRepository

	return temporaryGitRepository, nil
}

func (g *TemporaryGitRepositoriesService) CreateTemporaryGitRepositoryAndAddDataFromDirectory(ctx context.Context, dataToAdd filesinterfaces.Directory) (temporaryRepository GitRepository, err error) {
	if dataToAdd == nil {
		return nil, tracederrors.TracedError("dataToAdd is nil")
	}

	temporaryRepository, err = g.CreateTemporaryGitRepository(ctx)
	if err != nil {
		return nil, err
	}

	destDir, err := temporaryRepository.GetAsLocalDirectory()
	if err != nil {
		return nil, err
	}

	err = dataToAdd.CopyContentToDirectory(ctx, destDir)
	if err != nil {
		return nil, err
	}

	return temporaryRepository, nil
}

func (t TemporaryGitRepositoriesService) CreateEmptyTemporaryGitRepository(ctx context.Context, createRepoOptions *parameteroptions.CreateRepositoryOptions) (temporaryGitRepository GitRepository, err error) {
	if createRepoOptions == nil {
		return nil, tracederrors.TracedErrorNil("createRepoOptions")
	}

	tempDirectory, err := tempfilesoo.CreateEmptyTemporaryDirectory(ctx)
	if err != nil {
		return nil, err
	}

	temporaryGitRepository, err = GetLocalGitReposioryFromDirectory(tempDirectory)
	if err != nil {
		return nil, err
	}

	err = temporaryGitRepository.Init(ctx, createRepoOptions)
	if err != nil {
		return nil, err
	}

	repoPath, err := tempDirectory.GetLocalPath()
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Created temporary local git repository '%s'.", repoPath)

	return temporaryGitRepository, err
}
