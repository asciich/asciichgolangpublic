package asciichgolangpublic

import (
	"errors"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
)

type LocalGitRepository struct {
	LocalDirectory
}

func GetLocalGitReposioryFromDirectory(directory Directory) (l *LocalGitRepository, err error) {
	if directory == nil {
		return nil, TracedErrorNil("directory")
	}

	if !directory.IsLocalDirectory() {
		return nil, TracedError("Only local directories are supported.")
	}

	localPath, err := directory.GetLocalPath()
	if err != nil {
		return nil, err
	}

	l, err = GetLocalGitRepositoryByPath(localPath)
	if err != nil {
		return nil, err
	}

	return l, nil
}

func GetLocalGitRepositoryByPath(path string) (l *LocalGitRepository, err error) {
	if path == "" {
		return nil, TracedErrorEmptyString("path")
	}

	l = NewLocalGitRepository()

	err = l.SetLocalPath(path)
	if err != nil {
		return nil, err
	}

	return l, nil
}

func MustGetLocalGitReposioryFromDirectory(directory Directory) (l *LocalGitRepository) {
	l, err := GetLocalGitReposioryFromDirectory(directory)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return l
}

func MustGetLocalGitRepositoryByPath(path string) (l *LocalGitRepository) {
	l, err := GetLocalGitRepositoryByPath(path)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return l
}

func NewLocalGitRepository() (l *LocalGitRepository) {
	return new(LocalGitRepository)
}

func (l *LocalGitRepository) Add(path string) (err error) {
	if path == "" {
		return TracedErrorNil("path")
	}

	worktree, err := l.GetGoGitWorktree()
	if err != nil {
		return err
	}

	_, err = worktree.Add(path)
	if err != nil {
		return TracedErrorf("%w", err)
	}

	return nil
}

func (l *LocalGitRepository) Commit(commitOptions *GitCommitOptions) (createdCommit *GitCommit, err error) {
	if commitOptions == nil {
		return nil, TracedErrorNil("commitOptions")
	}

	worktree, err := l.GetGoGitWorktree()
	if err != nil {
		return nil, err
	}

	commitMessage, err := commitOptions.GetMessage()
	if err != nil {
		return nil, err
	}

	hash, err := worktree.Commit(
		commitMessage,
		&git.CommitOptions{
			AllowEmptyCommits: commitOptions.GetAllowEmpty(),
		},
	)
	if err != nil {
		return nil, err
	}

	createdCommit, err = l.GetCommitByGoGitHash(&hash)
	if err != nil {
		return nil, err
	}

	path, err := l.GetLocalPath()
	if err != nil {
		return nil, err
	}

	if commitOptions.Verbose {
		LogChangedf("Created commit '%s' in git repository '%s'.", commitMessage, path)
	}

	return createdCommit, nil
}

func (l *LocalGitRepository) CommitAndPush(commitOptions *GitCommitOptions) (createdCommit *GitCommit, err error) {
	if commitOptions == nil {
		return nil, TracedErrorNil("commitOptions")
	}

	createdCommit, err = l.Commit(commitOptions)
	if err != nil {
		return nil, err
	}

	err = l.Push(commitOptions.Verbose)
	if err != nil {
		return nil, err
	}

	return createdCommit, nil
}

func (l *LocalGitRepository) GetAsGoGitRepository() (goGitRepository *git.Repository, err error) {
	repoPath, err := l.GetLocalPath()
	if err != nil {
		return nil, err
	}

	goGitRepository, err = git.PlainOpenWithOptions(repoPath, &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		return nil, TracedErrorf("%w: repoPath='%s'", err, repoPath)
	}

	if goGitRepository == nil {
		return nil, TracedError("goGitRepository is nil after evaluation.")
	}

	return goGitRepository, nil
}

func (l *LocalGitRepository) GetCommitByGoGitHash(goGitHash *plumbing.Hash) (gitCommit *GitCommit, err error) {
	if goGitHash == nil {
		return nil, TracedErrorNil("goGitHash")
	}

	gitCommit = NewGitCommit()

	err = gitCommit.SetGitRepo(l)
	if err != nil {
		return nil, err
	}

	goGitRepo, err := l.GetAsGoGitRepository()
	if err != nil {
		return nil, err
	}

	commitObject, err := goGitRepo.CommitObject(*goGitHash)
	if err != nil {
		return nil, TracedErrorf("%w", err)
	}

	err = gitCommit.SetHash(commitObject.Hash.String())
	if err != nil {
		return nil, err
	}

	return gitCommit, nil
}

func (l *LocalGitRepository) GetCommitByGoGitReference(goGitReference *plumbing.Reference) (gitCommit *GitCommit, err error) {
	if goGitReference == nil {
		return nil, TracedErrorNil("goGitReference")
	}

	hash := goGitReference.Hash()

	gitCommit, err = l.GetCommitByGoGitHash(&hash)
	if err != nil {
		return nil, err
	}

	return gitCommit, nil
}

func (l *LocalGitRepository) GetCurrentCommit() (gitCommit *GitCommit, err error) {
	head, err := l.GetGoGitHead()
	if err != nil {
		return nil, err
	}

	gitCommit, err = l.GetCommitByGoGitReference(head)
	if err != nil {
		return nil, err
	}

	return gitCommit, nil
}

func (l *LocalGitRepository) GetCurrentCommitHash() (commitHash string, err error) {
	commit, err := l.GetCurrentCommit()
	if err != nil {
		return "", err
	}

	commitHash, err = commit.GetHash()
	if err != nil {
		return "", err
	}

	return commitHash, nil
}

func (l *LocalGitRepository) GetGoGitConfig() (config *config.Config, err error) {
	goGitRepo, err := l.GetAsGoGitRepository()
	if err != nil {
		return nil, err
	}

	config, err = goGitRepo.Config()
	if err != nil {
		return nil, TracedErrorf("%w", err)
	}

	if config == nil {
		return nil, TracedError("config is nil after evaluation")
	}

	return config, nil
}

func (l *LocalGitRepository) GetGoGitHead() (head *plumbing.Reference, err error) {
	goGitRepo, err := l.GetAsGoGitRepository()
	if err != nil {
		return nil, err
	}

	head, err = goGitRepo.Head()
	if err != nil {
		return nil, TracedErrorf("%w", err)
	}

	return head, nil
}

func (l *LocalGitRepository) GetGoGitWorktree() (worktree *git.Worktree, err error) {
	goGitRepo, err := l.GetAsGoGitRepository()
	if err != nil {
		return nil, err
	}

	worktree, err = goGitRepo.Worktree()
	if err != nil {
		return nil, TracedErrorf("%w", err)
	}

	return worktree, nil
}

func (l *LocalGitRepository) HasNoUncommittedChanges() (hasUncommittedChanges bool, err error) {
	hasUncommitedChanges, err := l.HasUncommittedChanges()
	if err != nil {
		return false, err
	}

	return !hasUncommitedChanges, nil
}

func (l *LocalGitRepository) HasUncommittedChanges() (hasUncommittedChanges bool, err error) {
	worktree, err := l.GetGoGitWorktree()
	if err != nil {
		return false, err
	}

	gitStatus, err := worktree.Status()
	if err != nil {
		return false, err
	}

	if gitStatus.IsClean() {
		return false, nil
	}

	return true, nil
}

func (l *LocalGitRepository) Init(options *CreateRepositoryOptions) (err error) {
	if options == nil {
		return TracedErrorNil("options")
	}

	repoPath, err := l.GetLocalPath()
	if err != nil {
		return err
	}

	isInitialized, err := l.IsInitialized()
	if err != nil {
		return err
	}

	if isInitialized {
		if options.Verbose {
			LogInfof("Local git repository '%s' is already initialized.", repoPath)
		}
	} else {
		_, err = git.PlainInit(repoPath, options.BareRepository)
		if err != nil {
			return TracedErrorf("%w", err)
		}
		if options.Verbose {
			LogChangedf("Local git repository '%s' is initialized.", repoPath)
		}

		if options.InitializeWithEmptyCommit {
			temporaryRepository, err := GitRepositories().CloneGitRepositoryToTemporaryDirectory(
				l,
				options.Verbose,
			)
			if err != nil {
				return err
			}
			defer temporaryRepository.Delete(options.Verbose)

			err = temporaryRepository.SetGitConfig(
				&GitConfigSetOptions{
					Name:    "asciichgolangpublic git repo initilaizer",
					Email:   "asciichgolangpublic@example.net",
					Verbose: options.Verbose,
				},
			)
			if err != nil {
				return err
			}

			_, err = temporaryRepository.CommitAndPush(
				&GitCommitOptions{
					Message:    "Initial empty commit during repo initialization",
					AllowEmpty: true,
					Verbose:    true,
				},
			)
			if err != nil {
				return err
			}

			if options.Verbose {
				LogChangedf("Initialized repository '%s' with an empty commit.", repoPath)
			}
		}
	}

	return nil
}

func (l *LocalGitRepository) IsBareRepository(verbose bool) (isBareRepository bool, err error) {
	config, err := l.GetGoGitConfig()
	if err != nil {
		return false, err
	}

	isBareRepository = config.Core.IsBare

	if verbose {
		repoRoot, err := l.GetLocalPath()
		if err != nil {
			return false, err
		}

		if isBareRepository {
			LogInfof("Git repository '%s' is a bare repository.", repoRoot)
		} else {
			LogInfof("Git repository '%s' is not a bare repository.", repoRoot)
		}
	}

	return isBareRepository, nil
}

func (l *LocalGitRepository) IsInitialized() (isInitialized bool, err error) {
	_, err = l.GetAsGoGitRepository()
	if err != nil {
		if errors.Is(err, git.ErrRepositoryNotExists) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (l *LocalGitRepository) MustAdd(path string) {
	err := l.Add(path)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (l *LocalGitRepository) MustCommit(commitOptions *GitCommitOptions) (createdCommit *GitCommit) {
	createdCommit, err := l.Commit(commitOptions)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return createdCommit
}

func (l *LocalGitRepository) MustCommitAndPush(commitOptions *GitCommitOptions) (createdCommit *GitCommit) {
	createdCommit, err := l.CommitAndPush(commitOptions)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return createdCommit
}

func (l *LocalGitRepository) MustGetAsGoGitRepository() (goGitRepository *git.Repository) {
	goGitRepository, err := l.GetAsGoGitRepository()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return goGitRepository
}

func (l *LocalGitRepository) MustGetCommitByGoGitHash(goGitHash *plumbing.Hash) (gitCommit *GitCommit) {
	gitCommit, err := l.GetCommitByGoGitHash(goGitHash)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitCommit
}

func (l *LocalGitRepository) MustGetCommitByGoGitReference(goGitReference *plumbing.Reference) (gitCommit *GitCommit) {
	gitCommit, err := l.GetCommitByGoGitReference(goGitReference)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitCommit
}

func (l *LocalGitRepository) MustGetCurrentCommit() (gitCommit *GitCommit) {
	gitCommit, err := l.GetCurrentCommit()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitCommit
}

func (l *LocalGitRepository) MustGetCurrentCommitHash() (commitHash string) {
	commitHash, err := l.GetCurrentCommitHash()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return commitHash
}

func (l *LocalGitRepository) MustGetGoGitConfig() (config *config.Config) {
	config, err := l.GetGoGitConfig()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return config
}

func (l *LocalGitRepository) MustGetGoGitHead() (head *plumbing.Reference) {
	head, err := l.GetGoGitHead()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return head
}

func (l *LocalGitRepository) MustGetGoGitWorktree() (worktree *git.Worktree) {
	worktree, err := l.GetGoGitWorktree()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return worktree
}

func (l *LocalGitRepository) MustHasNoUncommittedChanges() (hasUncommittedChanges bool) {
	hasUncommittedChanges, err := l.HasNoUncommittedChanges()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return hasUncommittedChanges
}

func (l *LocalGitRepository) MustHasUncommittedChanges() (hasUncommittedChanges bool) {
	hasUncommittedChanges, err := l.HasUncommittedChanges()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return hasUncommittedChanges
}

func (l *LocalGitRepository) MustInit(options *CreateRepositoryOptions) {
	err := l.Init(options)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (l *LocalGitRepository) MustIsBareRepository(verbose bool) (isBareRepository bool) {
	isBareRepository, err := l.IsBareRepository(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return isBareRepository
}

func (l *LocalGitRepository) MustIsInitialized() (isInitialized bool) {
	isInitialized, err := l.IsInitialized()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return isInitialized
}

func (l *LocalGitRepository) MustPull(verbose bool) {
	err := l.Pull(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (l *LocalGitRepository) MustPush(verbose bool) {
	err := l.Push(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (l *LocalGitRepository) MustSetGitConfig(options *GitConfigSetOptions) {
	err := l.SetGitConfig(options)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (l *LocalGitRepository) MustSetGitConfigByGoGitConfig(config *config.Config, verbose bool) {
	err := l.SetGitConfigByGoGitConfig(config, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (l *LocalGitRepository) MustSetRemote(remoteName string, remotUrl string, verbose bool) (remote *LocalGitRemote) {
	remote, err := l.SetRemote(remoteName, remotUrl, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return remote
}

func (l *LocalGitRepository) Pull(verbose bool) (err error) {
	worktree, err := l.GetGoGitWorktree()
	if err != nil {
		return err
	}

	err = worktree.Pull(&git.PullOptions{})
	if err != nil {
		return TracedErrorf("%w", err)
	}

	return nil
}

func (l *LocalGitRepository) Push(verbose bool) (err error) {
	goGitRepo, err := l.GetAsGoGitRepository()
	if err != nil {
		return err
	}

	err = goGitRepo.Push(&git.PushOptions{})
	if err != nil {
		return TracedErrorf("%w", err)
	}

	return nil
}

func (l *LocalGitRepository) SetGitConfig(options *GitConfigSetOptions) (err error) {
	if options == nil {
		return TracedErrorNil("options")
	}

	repoPath, err := l.GetLocalPath()
	if err != nil {
		return err
	}

	config, err := l.GetGoGitConfig()
	if err != nil {
		return err
	}

	rewriteNeeded := false
	if options.IsEmailSet() {
		email, err := options.GetEmail()
		if err != nil {
			return err
		}

		if config.Author.Email == email {
			LogInfof("Email in git config of local repository '%s' is already '%s'.", repoPath, email)
		} else {
			config.Author.Email = email
			rewriteNeeded = true
			LogChangedf("Set email in git config of local repository '%s' to '%s'.", repoPath, email)
		}
	}

	if options.IsNameSet() {
		name, err := options.GetName()
		if err != nil {
			return err
		}

		if config.Author.Name == name {
			LogInfof("Author name in git config of local repository '%s' is already '%s'.", repoPath, name)
		} else {
			config.Author.Name = name
			rewriteNeeded = true
			LogChangedf("Set author name in git config of local repository '%s' to '%s'.", repoPath, name)
		}
	}

	if rewriteNeeded {
		err = l.SetGitConfigByGoGitConfig(config, options.Verbose)
		if err != nil {
			return err
		}
	}

	return nil
}

func (l *LocalGitRepository) SetGitConfigByGoGitConfig(config *config.Config, verbose bool) (err error) {
	if config == nil {
		return TracedErrorNil("config")
	}

	outFile, err := l.GetFileInDirectory(".git", "config")
	if err != nil {
		return err
	}

	configData, err := config.Marshal()
	if err != nil {
		return TracedErrorf("%w", err)
	}

	const verboseWrite bool = false
	err = outFile.WriteBytes(configData, verboseWrite)
	if err != nil {
		return err
	}

	path, err := l.GetLocalPath()
	if err != nil {
		return err
	}

	if verbose {
		LogChangedf("Wrote git config of local git repository '%s'.", path)
	}

	return nil
}

func (l *LocalGitRepository) SetRemote(remoteName string, remotUrl string, verbose bool) (remote *LocalGitRemote, err error) {
	if remoteName == "" {
		return nil, TracedErrorEmptyString("remoteName")
	}

	if remotUrl == "" {
		return nil, TracedErrorEmptyString("remotUrl")
	}

	goGitRepo, err := l.GetAsGoGitRepository()
	if err != nil {
		return nil, err
	}

	nativeRemote, err := goGitRepo.CreateRemote(&config.RemoteConfig{
		Name: remoteName,
		URLs: []string{remotUrl},
	})
	if err != nil {
		return nil, TracedErrorf("Create remote failed: '%w'", err)
	}

	remote, err = NewLocalGitRemoteByNativeGoGitRemote(nativeRemote)
	if err != nil {
		return nil, err
	}

	return remote, err
}
