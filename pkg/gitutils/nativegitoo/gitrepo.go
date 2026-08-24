package nativegitoo

import (
	"context"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type NativeGitRepository struct {
	gitgeneric.GitRepositoryBase
	path string
}

func NewNativeGitRepository() (*NativeGitRepository, error) {
	ret := &NativeGitRepository{}

	err := ret.SetParentRepositoryForBaseClass(ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func NewGitRepositoryFromPath(path string) (gitinterfaces.GitRepository, error) {
	if path == "" {
		return nil, tracederrors.TracedErrorEmptyString("path")
	}

	repo, err := NewNativeGitRepository()
	if err != nil {
		return nil, err
	}

	err = repo.SetPath(path)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func NewGitRepositoryFromDirectory(directory filesinterfaces.Directory) (gitinterfaces.GitRepository, error) {
	if directory == nil {
		return nil, tracederrors.TracedErrorNil("directory")
	}

	hostDescription, err := directory.GetHostDescription()
	if err != nil {
		return nil, err
	}

	if hostDescription != "localhost" {
		return nil, tracederrors.TracedErrorf("The native git repository can only run on localhost but got host description = '%s'.", hostDescription)
	}

	path, err := directory.GetPath()
	if err != nil {
		return nil, err
	}

	return NewGitRepositoryFromPath(path)
}

func (n *NativeGitRepository) SetPath(path string) error {
	if path == "" {
		return tracederrors.TracedErrorEmptyString("path")
	}

	n.path = path

	return nil
}

func (n *NativeGitRepository) AddRemote(ctx context.Context, options *gitparameteroptions.GitRemoteAddOptions) (err error) {
	return tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) AddFileByPath(ctx context.Context, pathToAdd string) (err error) {
	return tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) CheckoutBranchByName(ctx context.Context, name string) (err error) {
	return tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) CloneRepository(ctx context.Context, repository gitinterfaces.GitRepository) (err error) {
	return tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) CloneRepositoryByPathOrUrl(ctx context.Context, pathOrUrl string) (err error) {
	return tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) CloneToTemporaryRepository(ctx context.Context) (gitinterfaces.GitRepository, error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) Commit(ctx context.Context, commitOptions *gitparameteroptions.GitCommitOptions) (createdCommit gitinterfaces.GitCommit, err error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) CommitHasParentCommitByCommitHash(hash string) (hasParentCommit bool, err error) {
	return false, tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) CreateBranch(ctx context.Context, createOptions *parameteroptions.CreateBranchOptions) (err error) {
	return tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) Create(ctx context.Context, options *filesoptions.CreateOptions) (err error) {
	return tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) CreateFileInDirectory(ctx context.Context, filePath string, options *filesoptions.CreateOptions) (createdFile filesinterfaces.File, err error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) CreateSubDirectory(ctx context.Context, subDirectoryName string, options *filesoptions.CreateOptions) (createdSubDirectory filesinterfaces.Directory, err error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) CreateTag(ctx context.Context, createOptions *gitparameteroptions.GitRepositoryCreateTagOptions) (createdTag gitinterfaces.GitTag, err error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) Delete(ctx context.Context, options *filesoptions.DeleteOptions) (err error) {
	return tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) DeleteBranchByName(ctx context.Context, name string) (err error) {
	return tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) DirectoryByPathExists(ctx context.Context, path ...string) (exists bool, err error) {
	return false, tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) Exists(ctx context.Context) (exists bool, err error) {
	return false, tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) Fetch(ctx context.Context) (err error) {
	return tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) FileByPathExists(ctx context.Context, path string) (exists bool, err error) {
	return false, tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) GetAuthorEmailByCommitHash(hash string) (authorEmail string, err error) {
	return "", tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) GetAuthorStringByCommitHash(hash string) (authorString string, err error) {
	return "", tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) GetDirectoryByPath(ctx context.Context, pathToSubDir ...string) (subDir filesinterfaces.Directory, err error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) GetCommitAgeDurationByCommitHash(hash string) (ageDuration *time.Duration, err error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) GetCommitAgeSecondsByCommitHash(hash string) (ageSeconds float64, err error) {
	return 0, tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) GetCommitMessageByCommitHash(hash string) (commitMessage string, err error) {
	return "", tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) GetCommitParentsByCommitHash(ctx context.Context, hash string, options *parameteroptions.GitCommitGetParentsOptions) (commitParents []gitinterfaces.GitCommit, err error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) GetCommitTimeByCommitHash(hash string) (commitTime *time.Time, err error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) GetCurrentBranchName(ctx context.Context) (branchName string, err error) {
	return "", tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) GetCurrentCommit(ctx context.Context) (commit gitinterfaces.GitCommit, err error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) GetCurrentCommitHash(ctx context.Context) (currentCommitHash string, err error) {
	return "", tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) GetGitStatusOutput(ctx context.Context) (output string, err error) {
	return "", tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) GetHashByTagName(tagName string) (hash string, err error) {
	return "", tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) GetHostDescription() (hostDescription string, err error) {
	return "", tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) GetPath() (path string, err error) {
	if n.path == "" {
		return "", tracederrors.TracedError("path not set")
	}

	return n.path, nil
}

func (n *NativeGitRepository) GetRemoteConfigs(ctx context.Context) (remoteConfigs []gitinterfaces.GitRemoteConfig, err error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) GetRootDirectory(ctx context.Context) (directory filesinterfaces.Directory, err error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) GetRootDirectoryPath(ctx context.Context) (path string, err error) {
	p, err := n.GetPath()
	if err != nil {
		return "", err
	}

	repo, err := git.PlainOpenWithOptions(p, &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		return "", tracederrors.TracedErrorf("failed to open git repository at '%s': %w", p, err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return "", tracederrors.TracedErrorf("failed to get worktree for repository at '%s': %w", p, err)
	}

	return worktree.Filesystem.Root(), nil
}

func (n *NativeGitRepository) HasInitialCommit(ctx context.Context) (hasInitialCommit bool, err error) {
	return false, tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) HasUncommittedChanges(ctx context.Context) (hasUncommitedChanges bool, err error) {
	return false, tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) Init(ctx context.Context, options *parameteroptions.CreateRepositoryOptions) (err error) {
	return tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) IsBareRepository(ctx context.Context) (isBareRepository bool, err error) {
	return false, tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) IsGitRepository(ctx context.Context) (isRepository bool, err error) {
	return false, tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) IsInitialized(ctx context.Context) (isInitialited bool, err error) {
	return false, tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) ListBranchNames(ctx context.Context) (branchNames []string, err error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) ListFilePaths(ctx context.Context, listFileOptions *parameteroptions.ListFileOptions) (filePaths []string, err error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) ListFiles(ctx context.Context, listFileOptions *parameteroptions.ListFileOptions) (files []filesinterfaces.File, err error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) ListTagNames(ctx context.Context) (tagNames []string, err error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) ListTags(ctx context.Context) (tags []gitinterfaces.GitTag, err error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) ListTagsForCommitHash(ctx context.Context, hash string) (tags []gitinterfaces.GitTag, err error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) RemoteByNameExists(ctx context.Context, remoteName string) (remoteExists bool, err error) {
	return false, tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) RemoteConfigurationExists(ctx context.Context, config gitinterfaces.GitRemoteConfig) (exists bool, err error) {
	return false, tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) RemoveRemoteByName(ctx context.Context, remoteName string) (err error) {
	return tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) Pull(ctx context.Context) (err error) {
	return tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) PullFromRemote(ctx context.Context, pullOptions *gitparameteroptions.GitPullFromRemoteOptions) (err error) {
	return tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) Push(ctx context.Context) (err error) {
	return tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) PushTagsToRemote(ctx context.Context, remoteName string) (err error) {
	return tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) PushToRemote(ctx context.Context, remoteName string) (err error) {
	return tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) SetGitConfig(ctx context.Context, options *gitparameteroptions.GitConfigSetOptions) (err error) {
	return tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) SetRemoteUrl(ctx context.Context, remoteUrl string) (err error) {
	return tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) WriteBytesToFile(ctx context.Context, path string, content []byte, options *filesoptions.WriteOptions) (writtenFile filesinterfaces.File, err error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}
