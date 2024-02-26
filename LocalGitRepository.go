package asciichgolangpublic

import (
	"errors"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type LocalGitRepository struct {
	LocalDirectory
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

func (l *LocalGitRepository) GetAsGoGitRepository() (goGitRepository *git.Repository, err error) {
	repoPath, err := l.GetLocalPath()
	if err != nil {
		return nil, err
	}

	goGitRepository, err = git.PlainOpen(repoPath)
	if err != nil {
		return nil, TracedErrorf("%w", err)
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

func (l *LocalGitRepository) Init(verbose bool) (err error) {
	repoPath, err := l.GetLocalPath()
	if err != nil {
		return err
	}

	isInitialized, err := l.IsInitialized()
	if err != nil {
		return err
	}

	if isInitialized {
		if verbose {
			LogInfof("Local git repository '%s' is already initialized.", repoPath)
		}
	} else {
		const isBare = false
		_, err = git.PlainInit(repoPath, isBare)
		if err != nil {
			return TracedErrorf("%w", err)
		}
		if verbose {
			LogChangedf("Local git repository '%s' is initialized.", repoPath)
		}
	}

	return nil
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

func (l *LocalGitRepository) MustInit(verbose bool) {
	err := l.Init(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
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
