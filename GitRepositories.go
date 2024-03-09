package asciichgolangpublic

import (
	"strings"

	"github.com/go-git/go-git/v5"
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

func (g *GitRepositoriesService) CloneToDirectoryByPath(urlOrPath string, destinationPath string, verbose bool) (repo *LocalGitRepository, err error) {
	urlOrPath = strings.TrimSpace(urlOrPath)
	if urlOrPath == "" {
		return nil, TracedErrorEmptyString("urlOrPath")
	}

	destinationPath = strings.TrimSpace(destinationPath)
	if destinationPath == "" {
		return nil, TracedErrorEmptyString("destinationPath")
	}

	const isBare = false
	git.PlainClone(
		destinationPath,
		isBare,
		&git.CloneOptions{
			URL: urlOrPath,
		},
	)

	repo, err = GetLocalGitRepositoryByPath(destinationPath)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (g *GitRepositoriesService) CloneToTemporaryDirectory(urlOrPath string, verbose bool) (repo *LocalGitRepository, err error) {
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

func (g *GitRepositoriesService) CreateTemporaryInitializedRepository(options *CreateRepositoryOptions) (repo *LocalGitRepository, err error) {
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

	err = repo.Init(options.Verbose)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (g *GitRepositoriesService) MustCloneToDirectoryByPath(urlOrPath string, destinationPath string, verbose bool) (repo *LocalGitRepository) {
	repo, err := g.CloneToDirectoryByPath(urlOrPath, destinationPath, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return repo
}

func (g *GitRepositoriesService) MustCloneToTemporaryDirectory(urlOrPath string, verbose bool) (repo *LocalGitRepository) {
	repo, err := g.CloneToTemporaryDirectory(urlOrPath, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return repo
}

func (g *GitRepositoriesService) MustCreateTemporaryInitializedRepository(options *CreateRepositoryOptions) (repo *LocalGitRepository) {
	repo, err := g.CreateTemporaryInitializedRepository(options)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return repo
}
