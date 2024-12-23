package asciichgolangpublic

import (
	"strings"
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

func (g *GitRepositoriesService) CloneGitRepositoryToDirectory(toClone GitRepository, destinationPath string, verbose bool) (repo GitRepository, err error) {
	if toClone == nil {
		return nil, TracedErrorNil("toClone")
	}

	if destinationPath == "" {
		return nil, TracedErrorNil("checkoutPath")
	}

	localRepository, ok := toClone.(*LocalGitRepository)
	if !ok {
		return nil, TracedError("Only implemented for LocalGitRepository")
	}

	localPath, err := localRepository.GetLocalPath()
	if err != nil {
		return nil, err
	}

	repo, err = g.CloneToDirectoryByPath(localPath, destinationPath, verbose)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (g *GitRepositoriesService) CloneGitRepositoryToTemporaryDirectory(toClone GitRepository, verbose bool) (repo GitRepository, err error) {
	if toClone == nil {
		return nil, TracedErrorNil("toClone")
	}

	localRepository, ok := toClone.(*LocalGitRepository)
	if ok {
		localPath, err := localRepository.GetLocalPath()
		if err != nil {
			return nil, err
		}

		repo, err = g.CloneToTemporaryDirectory(localPath, verbose)
		if err != nil {
			return nil, err
		}

		if verbose {
			clonedPath, err := repo.GetPath()
			if err != nil {
				return nil, err
			}

			LogChangedf(
				"Cloned local git repository '%s' into temporary directory '%s'",
				localPath,
				clonedPath,
			)
		}
	}

	if repo == nil {
		commandExecutorRepository, ok := toClone.(*CommandExecutorGitRepository)
		if ok {
			localPath, hostDescription, err := commandExecutorRepository.GetPathAndHostDescription()
			if err != nil {
				return nil, err
			}

			if hostDescription != "localhost" {
				return nil, TracedErrorf(
					"Only implemented for CommandExecutorGitRepository on localhost, but hostDescription is '%s'",
					hostDescription,
				)
			}

			repo, err = g.CloneToTemporaryDirectory(localPath, verbose)
			if err != nil {
				return nil, err
			}

			if verbose {
				clonedPath, err := repo.GetPath()
				if err != nil {
					return nil, err
				}

				LogChangedf(
					"Cloned git repository '%s' from host '%s' into temporary directory '%s'",
					localPath,
					hostDescription,
					clonedPath,
				)
			}
		}
	}

	return repo, nil
}

func (g *GitRepositoriesService) CloneToDirectoryByPath(urlOrPath string, destinationPath string, verbose bool) (repo *LocalGitRepository, err error) {
	urlOrPath = strings.TrimSpace(urlOrPath)
	if urlOrPath == "" {
		return nil, TracedErrorEmptyString("urlOrPath")
	}

	destinationPath = strings.TrimSpace(destinationPath)
	if destinationPath == "" {
		return nil, TracedErrorEmptyString("destinationPath")
	}

	repo, err = GetLocalGitRepositoryByPath(destinationPath)
	if err != nil {
		return nil, err
	}

	err = repo.CloneRepositoryByPathOrUrl(urlOrPath, verbose)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (g *GitRepositoriesService) CloneToTemporaryDirectory(urlOrPath string, verbose bool) (repo GitRepository, err error) {
	urlOrPath = strings.TrimSpace(urlOrPath)
	if urlOrPath == "" {
		return nil, TracedErrorEmptyString("urlOrPath")
	}

	destinationPath, err := TemporaryDirectories().CreateEmptyTemporaryDirectoryAndGetPath(verbose)
	if err != nil {
		return nil, err
	}

	repo, err = g.CloneToDirectoryByPath(urlOrPath, destinationPath, verbose)
	if err != nil {
		return nil, err
	}

	if verbose {
		LogChangedf("Cloned git repository '%s' to local directory '%s'.", urlOrPath, destinationPath)
	}

	return repo, nil
}

func (g *GitRepositoriesService) CreateTemporaryInitializedRepository(options *CreateRepositoryOptions) (repo GitRepository, err error) {
	if options == nil {
		return nil, TracedErrorNil("options")
	}

	repoPath, err := TemporaryDirectories().CreateEmptyTemporaryDirectoryAndGetPath(options.Verbose)
	if err != nil {
		return nil, err
	}

	repo, err = GetLocalGitRepositoryByPath(repoPath)
	if err != nil {
		return nil, err
	}

	err = repo.Init(options)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (g *GitRepositoriesService) MustCloneGitRepositoryToDirectory(toClone GitRepository, destinationPath string, verbose bool) (repo GitRepository) {
	repo, err := g.CloneGitRepositoryToDirectory(toClone, destinationPath, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return repo
}

func (g *GitRepositoriesService) MustCloneGitRepositoryToTemporaryDirectory(toClone GitRepository, verbose bool) (repo GitRepository) {
	repo, err := g.CloneGitRepositoryToTemporaryDirectory(toClone, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return repo
}

func (g *GitRepositoriesService) MustCloneToDirectoryByPath(urlOrPath string, destinationPath string, verbose bool) (repo GitRepository) {
	repo, err := g.CloneToDirectoryByPath(urlOrPath, destinationPath, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return repo
}

func (g *GitRepositoriesService) MustCloneToTemporaryDirectory(urlOrPath string, verbose bool) (repo GitRepository) {
	repo, err := g.CloneToTemporaryDirectory(urlOrPath, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return repo
}

func (g *GitRepositoriesService) MustCreateTemporaryInitializedRepository(options *CreateRepositoryOptions) (repo GitRepository) {
	repo, err := g.CreateTemporaryInitializedRepository(options)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return repo
}
