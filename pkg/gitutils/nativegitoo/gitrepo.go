package nativegitoo

import (
	"context"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefilesoo"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
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

func (n *NativeGitRepository) GetHostDescription() (hostDescription string, err error) {
	return "localhost", nil
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
	if pathToAdd == "" {
		return tracederrors.TracedErrorEmptyString("pathToAdd")
	}

	path, hostDescription, err := n.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Add file '%s' in git repository '%s' on '%s' started.", pathToAdd, path, hostDescription)

	repo, err := git.PlainOpen(path)
	if err != nil {
		return tracederrors.TracedErrorf("Open git repository '%s' on '%s' failed: %w", path, hostDescription, err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return tracederrors.TracedErrorf("Get worktree for git repository '%s' on '%s' failed: %w", path, hostDescription, err)
	}

	_, err = worktree.Add(pathToAdd)
	if err != nil {
		return tracederrors.TracedErrorf("Add file '%s' in git repository '%s' on '%s' failed: %w", pathToAdd, path, hostDescription, err)
	}

	logging.LogInfoByCtxf(ctx, "Add file '%s' in git repository '%s' on '%s' finished.", pathToAdd, path, hostDescription)

	return nil
}

func (n *NativeGitRepository) CloneRepository(ctx context.Context, repository gitinterfaces.GitRepository) (err error) {
	return tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) CloneToTemporaryRepository(ctx context.Context) (gitinterfaces.GitRepository, error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) Commit(ctx context.Context, commitOptions *gitparameteroptions.GitCommitOptions) (createdCommit gitinterfaces.GitCommit, err error) {
	if commitOptions == nil {
		return nil, tracederrors.TracedErrorNil("commitOptions")
	}

	path, hostDescription, err := n.GetPathAndHostDescription()
	if err != nil {
		return nil, err
	}

	message, err := commitOptions.GetMessage()
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Commit in git repository '%s' on '%s' started.", path, hostDescription)

	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Open git repository '%s' on '%s' failed: %w", path, hostDescription, err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return nil, tracederrors.TracedErrorf("Get worktree for git repository '%s' on '%s' failed: %w", path, hostDescription, err)
	}

	if commitOptions.GetCommitAllChanges() {
		_, err = worktree.Add(".")
		if err != nil {
			return nil, tracederrors.TracedErrorf("Add all changes in git repository '%s' on '%s' failed: %w", path, hostDescription, err)
		}
	}

	cfg, err := repo.Config()
	if err != nil {
		return nil, tracederrors.TracedErrorf("Get config for git repository '%s' on '%s' failed: %w", path, hostDescription, err)
	}

	authorName := cfg.User.Name
	authorEmail := cfg.User.Email

	commitHash, err := worktree.Commit(message, &git.CommitOptions{
		AllowEmptyCommits: commitOptions.GetAllowEmpty(),
		Author: &object.Signature{
			Name:  authorName,
			Email: authorEmail,
			When:  time.Now(),
		},
	})
	if err != nil {
		return nil, tracederrors.TracedErrorf("Commit in git repository '%s' on '%s' failed: %w", path, hostDescription, err)
	}

	commitObject, err := repo.CommitObject(commitHash)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Get commit object for hash '%s' in git repository '%s' on '%s' failed: %w", commitHash.String(), path, hostDescription, err)
	}

	_ = commitObject

	createdCommit = NewNativeGitCommit(commitHash.String(), n)

	logging.LogInfoByCtxf(ctx, "Commit in git repository '%s' on '%s' finished. Created commit '%s'.", path, hostDescription, commitHash.String())

	return createdCommit, nil
}

func (n *NativeGitRepository) CommitHasParentCommitByCommitHash(hash string) (hasParentCommit bool, err error) {
	return false, tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) Create(ctx context.Context, options *filesoptions.CreateOptions) (err error) {
	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	path, hostDescription, err := n.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Create directory for git repository '%s' on '%s' started.", path, hostDescription)

	exists := nativefiles.Exists(ctx, path)
	if err != nil {
		return err
	}

	if exists {
		logging.LogInfoByCtxf(ctx, "Directory '%s' on '%s' already exists.", path, hostDescription)
	} else {
		err = nativefiles.CreateDirectory(ctx, path, &filesoptions.CreateOptions{})
		if err != nil {
			return tracederrors.TracedErrorf("Create directory '%s' on '%s' failed: %w", path, hostDescription, err)
		}

		logging.LogChangedByCtxf(ctx, "Created directory '%s' on '%s'.", path, hostDescription)
	}

	logging.LogInfoByCtxf(ctx, "Create directory for git repository '%s' on '%s' finished.", path, hostDescription)

	return nil
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
	if path == "" {
		return false, tracederrors.TracedErrorEmptyString("path")
	}

	repoPath, err := n.GetPath()
	if err != nil {
		return false, err
	}

	fullPath := filepath.Join(repoPath, path)

	exists = nativefiles.Exists(ctx, fullPath)
	if err != nil {
		return false, err
	}

	return exists, nil
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

func (n *NativeGitRepository) GetPath() (path string, err error) {
	if n.path == "" {
		return "", tracederrors.TracedError("path not set")
	}

	return n.path, nil
}

func (n *NativeGitRepository) GetRemoteConfigs(ctx context.Context) (remoteConfigs []gitinterfaces.GitRemoteConfig, err error) {
	path, hostDescription, err := n.GetPathAndHostDescription()
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Get remote configs for git repository '%s' on '%s' started.", path, hostDescription)

	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Open git repository '%s' on '%s' failed: %w", path, hostDescription, err)
	}

	remotes, err := repo.Remotes()
	if err != nil {
		return nil, tracederrors.TracedErrorf("Get remotes for git repository '%s' on '%s' failed: %w", path, hostDescription, err)
	}

	remoteConfigs = []gitinterfaces.GitRemoteConfig{}
	for _, remote := range remotes {
		cfg := remote.Config()
		fetchUrl := ""
		pushUrl := ""
		if len(cfg.URLs) > 0 {
			fetchUrl = cfg.URLs[0]
			pushUrl = cfg.URLs[0]
		}
		remoteConfigs = append(remoteConfigs, &NativeGitRemoteConfig{
			remoteName: cfg.Name,
			urlFetch:   fetchUrl,
			urlPush:    pushUrl,
		})
	}

	logging.LogInfoByCtxf(ctx, "Get remote configs for git repository '%s' on '%s' finished. Found '%d' remote configs.", path, hostDescription, len(remoteConfigs))

	return remoteConfigs, nil
}

func (n *NativeGitRepository) GetRootDirectory(ctx context.Context) (directory filesinterfaces.Directory, err error) {
	path, err := n.GetRootDirectoryPath(ctx)
	if err != nil {
		return nil, err
	}

	directory, err = nativefilesoo.NewDirectoryByPath(path)
	if err != nil {
		return nil, err
	}

	return directory, nil
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
	path, err := n.GetPath()
	if err != nil {
		return false, err
	}

	repo, err := git.PlainOpen(path)
	if err != nil {
		return false, tracederrors.TracedErrorf("Open git repository '%s' failed: %w", path, err)
	}

	_, err = repo.Head()
	if err != nil {
		if err == plumbing.ErrReferenceNotFound {
			return false, nil
		}
		return false, tracederrors.TracedErrorf("Get HEAD for git repository '%s' failed: %w", path, err)
	}

	return true, nil
}

func (n *NativeGitRepository) HasUncommittedChanges(ctx context.Context) (hasUncommitedChanges bool, err error) {
	return false, tracederrors.TracedErrorNotImplemented()
}
func (n *NativeGitRepository) IsBareRepository(ctx context.Context) (isBareRepository bool, err error) {
	path, err := n.GetPath()
	if err != nil {
		return false, err
	}

	repo, err := git.PlainOpen(path)
	if err != nil {
		return false, tracederrors.TracedErrorf("Open git repository '%s' failed: %w", path, err)
	}

	cfg, err := repo.Config()
	if err != nil {
		return false, tracederrors.TracedErrorf("Get config for git repository '%s' failed: %w", path, err)
	}

	return cfg.Core.IsBare, nil
}

func (n *NativeGitRepository) IsGitRepository(ctx context.Context) (isRepository bool, err error) {
	path, err := n.GetPath()
	if err != nil {
		return false, err
	}

	_, err = git.PlainOpen(path)
	if err != nil {
		if err == git.ErrRepositoryNotExists {
			return false, nil
		}
		return false, tracederrors.TracedErrorf("Check if '%s' is a git repository failed: %w", path, err)
	}

	return true, nil
}

func (n *NativeGitRepository) ListBranchNames(ctx context.Context) (branchNames []string, err error) {
	path, hostDescription, err := n.GetPathAndHostDescription()
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "List branch names in git repository '%s' on '%s' started.", path, hostDescription)

	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Open git repository '%s' on '%s' failed: %w", path, hostDescription, err)
	}

	branches, err := repo.Branches()
	if err != nil {
		return nil, tracederrors.TracedErrorf("List branches in git repository '%s' on '%s' failed: %w", path, hostDescription, err)
	}

	branchNames = []string{}
	err = branches.ForEach(func(ref *plumbing.Reference) error {
		branchNames = append(branchNames, ref.Name().Short())
		return nil
	})
	if err != nil {
		return nil, tracederrors.TracedErrorf("Iterate branches in git repository '%s' on '%s' failed: %w", path, hostDescription, err)
	}

	logging.LogInfoByCtxf(ctx, "List branch names in git repository '%s' on '%s' finished. Found '%d' branches.", path, hostDescription, len(branchNames))

	return branchNames, nil
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
	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	path, hostDescription, err := n.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Set git config for repository '%s' on '%s' started.", path, hostDescription)

	repo, err := git.PlainOpen(path)
	if err != nil {
		return tracederrors.TracedErrorf("Open git repository '%s' on '%s' failed: %w", path, hostDescription, err)
	}

	cfg, err := repo.Config()
	if err != nil {
		return tracederrors.TracedErrorf("Get config for git repository '%s' on '%s' failed: %w", path, hostDescription, err)
	}

	if options.Name != "" {
		cfg.User.Name = options.Name
	}

	if options.Email != "" {
		cfg.User.Email = options.Email
	}

	err = repo.SetConfig(cfg)
	if err != nil {
		return tracederrors.TracedErrorf("Set config for git repository '%s' on '%s' failed: %w", path, hostDescription, err)
	}

	logging.LogInfoByCtxf(ctx, "Set git config for repository '%s' on '%s' finished.", path, hostDescription)

	return nil
}

func (n *NativeGitRepository) SetRemoteUrl(ctx context.Context, remoteUrl string) (err error) {
	if remoteUrl == "" {
		return tracederrors.TracedErrorEmptyString("remoteUrl")
	}

	path, hostDescription, err := n.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Set remote url to '%s' for git repository '%s' on '%s' started.", remoteUrl, path, hostDescription)

	repo, err := git.PlainOpen(path)
	if err != nil {
		return tracederrors.TracedErrorf("Open git repository '%s' on '%s' failed: %w", path, hostDescription, err)
	}

	_, err = repo.Remote("origin")
	if err != nil {
		if err == git.ErrRemoteNotFound {
			_, err = repo.CreateRemote(&config.RemoteConfig{
				Name: "origin",
				URLs: []string{remoteUrl},
			})
			if err != nil {
				return tracederrors.TracedErrorf("Create remote 'origin' with url '%s' for git repository '%s' on '%s' failed: %w", remoteUrl, path, hostDescription, err)
			}
		} else {
			return tracederrors.TracedErrorf("Get remote 'origin' for git repository '%s' on '%s' failed: %w", path, hostDescription, err)
		}
	} else {
		cfg, err := repo.Config()
		if err != nil {
			return tracederrors.TracedErrorf("Get config for git repository '%s' on '%s' failed: %w", path, hostDescription, err)
		}

		cfg.Remotes["origin"].URLs = []string{remoteUrl}

		err = repo.SetConfig(cfg)
		if err != nil {
			return tracederrors.TracedErrorf("Set config for git repository '%s' on '%s' failed: %w", path, hostDescription, err)
		}
	}

	logging.LogInfoByCtxf(ctx, "Set remote url to '%s' for git repository '%s' on '%s' finished.", remoteUrl, path, hostDescription)

	return nil
}

func (n *NativeGitRepository) WriteBytesToFile(ctx context.Context, path string, content []byte, options *filesoptions.WriteOptions) (writtenFile filesinterfaces.File, err error) {
	if path == "" {
		return nil, tracederrors.TracedErrorEmptyString("path")
	}

	repoPath, hostDescription, err := n.GetPathAndHostDescription()
	if err != nil {
		return nil, err
	}

	fullPath := filepath.Join(repoPath, path)

	logging.LogInfoByCtxf(ctx, "Write bytes to file '%s' in git repository '%s' on '%s' started.", path, repoPath, hostDescription)

	err = nativefiles.WriteBytes(ctx, fullPath, content, options)
	if err != nil {
		return nil, err
	}

	writtenFile, err = n.GetFileByPath(path)
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Write bytes to file '%s' in git repository '%s' on '%s' finished.", path, repoPath, hostDescription)

	return writtenFile, nil
}
