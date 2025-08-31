package asciichgolangpublic

import (
	"context"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfilesoo"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type GitRepositoriesService struct {
}

func GitRepositories() (g *GitRepositoriesService) {
	return NewGitRepositories()
}

func NewGitRepositories() (g *GitRepositoriesService) {
	return new(GitRepositoriesService)
}

func NewGitRepositoriesService() (g *GitRepositoriesService) {
	return new(GitRepositoriesService)
}

func (g *GitRepositoriesService) CloneGitRepositoryToDirectory(ctx context.Context, toClone GitRepository, destinationPath string) (repo GitRepository, err error) {
	if toClone == nil {
		return nil, tracederrors.TracedErrorNil("toClone")
	}

	if destinationPath == "" {
		return nil, tracederrors.TracedErrorNil("checkoutPath")
	}

	localRepository, ok := toClone.(*LocalGitRepository)
	if !ok {
		return nil, tracederrors.TracedError("Only implemented for LocalGitRepository")
	}

	localPath, err := localRepository.GetLocalPath()
	if err != nil {
		return nil, err
	}

	repo, err = g.CloneToDirectoryByPath(ctx, localPath, destinationPath)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (g *GitRepositoriesService) CloneGitRepositoryToTemporaryDirectory(ctx context.Context, toClone GitRepository) (repo GitRepository, err error) {
	if toClone == nil {
		return nil, tracederrors.TracedErrorNil("toClone")
	}

	localRepository, ok := toClone.(*LocalGitRepository)
	if ok {
		localPath, err := localRepository.GetLocalPath()
		if err != nil {
			return nil, err
		}

		repo, err = g.CloneToTemporaryDirectory(ctx, localPath)
		if err != nil {
			return nil, err
		}

		clonedPath, err := repo.GetPath()
		if err != nil {
			return nil, err
		}

		logging.LogChangedByCtxf(ctx, "Cloned local git repository '%s' into temporary directory '%s'", localPath, clonedPath)
	}

	if repo == nil {
		commandExecutorRepository, ok := toClone.(*CommandExecutorGitRepository)
		if ok {
			localPath, hostDescription, err := commandExecutorRepository.GetPathAndHostDescription()
			if err != nil {
				return nil, err
			}

			if hostDescription != "localhost" {
				return nil, tracederrors.TracedErrorf(
					"Only implemented for CommandExecutorGitRepository on localhost, but hostDescription is '%s'",
					hostDescription,
				)
			}

			repo, err = g.CloneToTemporaryDirectory(ctx, localPath)
			if err != nil {
				return nil, err
			}

			clonedPath, err := repo.GetPath()
			if err != nil {
				return nil, err
			}

			logging.LogChangedByCtxf(ctx, "Cloned git repository '%s' from host '%s' into temporary directory '%s'", localPath, hostDescription, clonedPath)
		}
	}

	return repo, nil
}

func (g *GitRepositoriesService) CloneToDirectoryByPath(ctx context.Context, urlOrPath string, destinationPath string) (repo *LocalGitRepository, err error) {
	urlOrPath = strings.TrimSpace(urlOrPath)
	if urlOrPath == "" {
		return nil, tracederrors.TracedErrorEmptyString("urlOrPath")
	}

	destinationPath = strings.TrimSpace(destinationPath)
	if destinationPath == "" {
		return nil, tracederrors.TracedErrorEmptyString("destinationPath")
	}

	repo, err = GetLocalGitRepositoryByPath(destinationPath)
	if err != nil {
		return nil, err
	}

	err = repo.CloneRepositoryByPathOrUrl(ctx, urlOrPath)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (g *GitRepositoriesService) CloneToTemporaryDirectory(ctx context.Context, urlOrPath string) (repo GitRepository, err error) {
	urlOrPath = strings.TrimSpace(urlOrPath)
	if urlOrPath == "" {
		return nil, tracederrors.TracedErrorEmptyString("urlOrPath")
	}

	destinationPath, err := tempfilesoo.CreateEmptyTemporaryDirectoryAndGetPath(ctx)
	if err != nil {
		return nil, err
	}

	repo, err = g.CloneToDirectoryByPath(ctx, urlOrPath, destinationPath)
	if err != nil {
		return nil, err
	}

	logging.LogChangedByCtxf(ctx, "Cloned git repository '%s' to local directory '%s'.", urlOrPath, destinationPath)

	return repo, nil
}

func (g *GitRepositoriesService) CreateTemporaryInitializedRepository(ctx context.Context, options *parameteroptions.CreateRepositoryOptions) (repo GitRepository, err error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	repoPath, err := tempfilesoo.CreateEmptyTemporaryDirectoryAndGetPath(ctx)
	if err != nil {
		return nil, err
	}

	repo, err = GetLocalGitRepositoryByPath(repoPath)
	if err != nil {
		return nil, err
	}

	err = repo.Init(ctx, options)
	if err != nil {
		return nil, err
	}

	return repo, nil
}
