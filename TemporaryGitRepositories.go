package asciichgolangpublic

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfilesoo"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/commandexecutorgitoo"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitinterfaces"
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

func (g *TemporaryGitRepositoriesService) CreateTemporaryGitRepositoryAndAddDataFromDirectory(ctx context.Context, dataToAdd filesinterfaces.Directory) (temporaryRepository gitinterfaces.GitRepository, err error) {
	if dataToAdd == nil {
		return nil, tracederrors.TracedError("dataToAdd is nil")
	}

	temporaryRepository, err = commandexecutorgitoo.CreateLocalTemporaryRepository(ctx, &parameteroptions.CreateRepositoryOptions{})
	if err != nil {
		return nil, err
	}

	hostDescription, err := temporaryRepository.GetHostDescription()
	if err != nil {
		return nil, err
	}

	if hostDescription != "localhost" {
		return nil, tracederrors.TracedErrorf("Only implemented for localhost but got '%s'.", hostDescription)
	}

	localPath, err := temporaryRepository.GetPath()
	if err != nil {
		return nil, err
	}

	destDir, err := files.GetLocalDirectoryByPath(localPath)
	if err != nil {
		return nil, err
	}

	err = dataToAdd.CopyContentToDirectory(ctx, destDir)
	if err != nil {
		return nil, err
	}

	return temporaryRepository, nil
}

func (t TemporaryGitRepositoriesService) CreateEmptyTemporaryGitRepository(ctx context.Context, createRepoOptions *parameteroptions.CreateRepositoryOptions) (temporaryGitRepository gitinterfaces.GitRepository, err error) {
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
