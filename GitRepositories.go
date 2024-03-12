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

func (g *GitRepositoriesService) CloneGitRepositoryToDirectory(toClone GitRepository, destinationPath string, verbose bool) (repo *LocalGitRepository, err error) {
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

func (g *GitRepositoriesService) CloneGitRepositoryToTemporaryDirectory(toClone GitRepository, verbose bool) (repo *LocalGitRepository, err error) {
	if toClone == nil {
		return nil, TracedErrorNil("toClone")
	}

	localRepository, ok := toClone.(*LocalGitRepository)
	if !ok {
		return nil, TracedError("Only implemented for LocalGitRepository")
	}

	localPath, err := localRepository.GetLocalPath()
	if err != nil {
		return nil, err
	}

	repo, err = g.CloneToTemporaryDirectory(localPath, verbose)
	if err != nil {
		return nil, err
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

	const isBare = false
	_, err = git.PlainClone(
		destinationPath,
		isBare,
		&git.CloneOptions{
			URL: urlOrPath,
		},
	)
	if err != nil {
		if err.Error() == "remote repository is empty" {
			if verbose {
				LogInfof("Remote repository '%s' is empty. Going to add remote for empty repository.", urlOrPath)
			}

			err = repo.Init(
				&CreateRepositoryOptions{
					Verbose:                   verbose,
					BareRepository:            isBare,
					InitializeWithEmptyCommit: false,
				})
			if err != nil {
				return nil, err
			}

			_, err = repo.SetRemote("origin", urlOrPath, verbose)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, TracedErrorf(
				"Clone '%s' failed: '%w'",
				urlOrPath,
				err,
			)
		}
	}

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

	err = repo.Init(options)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (g *GitRepositoriesService) MustCloneGitRepositoryToDirectory(toClone GitRepository, destinationPath string, verbose bool) (repo *LocalGitRepository) {
	repo, err := g.CloneGitRepositoryToDirectory(toClone, destinationPath, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return repo
}

func (g *GitRepositoriesService) MustCloneGitRepositoryToTemporaryDirectory(toClone GitRepository, verbose bool) (repo *LocalGitRepository) {
	repo, err := g.CloneGitRepositoryToTemporaryDirectory(toClone, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return repo
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
