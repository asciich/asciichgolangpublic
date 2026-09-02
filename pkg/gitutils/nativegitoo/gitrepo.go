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
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
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
	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	remoteName, err := options.GetRemoteName()
	if err != nil {
		return err
	}

	remoteUrl, err := options.GetRemoteUrl()
	if err != nil {
		return err
	}

	path, hostDescription, err := n.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Add remote '%s' as '%s' to git repository '%s' on '%s' started.", remoteUrl, remoteName, path, hostDescription)

	repo, err := git.PlainOpen(path)
	if err != nil {
		return tracederrors.TracedErrorf("Open git repository '%s' on '%s' failed: %w", path, hostDescription, err)
	}

	// Check if remote already exists
	_, err = repo.Remote(remoteName)
	if err == nil {
		// Remote exists, remove it first to update URL
		err = repo.DeleteRemote(remoteName)
		if err != nil {
			return tracederrors.TracedErrorf("Delete existing remote '%s' in git repository '%s' on '%s' failed: %w", remoteName, path, hostDescription, err)
		}
	}

	// Add the remote
	_, err = repo.CreateRemote(&config.RemoteConfig{
		Name: remoteName,
		URLs: []string{remoteUrl},
	})
	if err != nil {
		return tracederrors.TracedErrorf("Create remote '%s' with url '%s' in git repository '%s' on '%s' failed: %w", remoteName, remoteUrl, path, hostDescription, err)
	}

	logging.LogChangedByCtxf(ctx, "Added remote '%s' as '%s' to git repository '%s' on '%s'.", remoteUrl, remoteName, path, hostDescription)

	return nil
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
	if repository == nil {
		return tracederrors.TracedErrorNil("repository")
	}

	repoHostDescription, err := repository.GetHostDescription()
	if err != nil {
		return err
	}

	hostDescription, err := n.GetHostDescription()
	if err != nil {
		return err
	}

	if hostDescription != repoHostDescription {
		return tracederrors.TracedErrorf(
			"Only implemented for two repositories on the same host. But repository from host '%s' should be cloned to host '%s'",
			repoHostDescription,
			hostDescription,
		)
	}

	pathToClone, err := repository.GetPath()
	if err != nil {
		return err
	}

	return n.CloneRepositoryByPathOrUrl(ctx, pathToClone)
}

func (n *NativeGitRepository) CloneToTemporaryRepository(ctx context.Context) (gitinterfaces.GitRepository, error) {
	path, hostDescription, err := n.GetPathAndHostDescription()
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Cloning repository '%s' to temporary repository on '%s' started.", path, hostDescription)

	// Create a temporary directory
	tempDir, err := tempfiles.CreateTempDir(ctx)
	if err != nil {
		return nil, err
	}

	// Create a new NativeGitRepository from the temporary directory
	clonedRepo, err := NewGitRepositoryFromPath(tempDir)
	if err != nil {
		return nil, err
	}

	// Clone the repository
	err = clonedRepo.CloneRepositoryByPathOrUrl(ctx, path)
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Cloning repository '%s' to temporary repository on '%s' finished.", path, hostDescription)

	return clonedRepo, nil
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
	if hash == "" {
		return false, tracederrors.TracedErrorEmptyString("hash")
	}

	path, err := n.GetPath()
	if err != nil {
		return false, err
	}

	repo, err := git.PlainOpen(path)
	if err != nil {
		return false, tracederrors.TracedErrorf("Open git repository '%s' failed: %w", path, err)
	}

	commitHash := plumbing.NewHash(hash)
	commitObject, err := repo.CommitObject(commitHash)
	if err != nil {
		return false, tracederrors.TracedErrorf("Get commit object for hash '%s' failed: %w", hash, err)
	}

	// A commit has a parent if it has at least one parent commit
	hasParentCommit = len(commitObject.ParentHashes) > 0

	return hasParentCommit, nil
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

func (n *NativeGitRepository) Delete(ctx context.Context, options *filesoptions.DeleteOptions) (err error) {
	path, hostDescription, err := n.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Delete git repository '%s' on '%s' started.", path, hostDescription)

	err = nativefiles.Delete(ctx, path, options)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Delete git repository '%s' on '%s' finished.", path, hostDescription)

	return nil
}

func (n *NativeGitRepository) DirectoryByPathExists(ctx context.Context, path ...string) (exists bool, err error) {
	return false, tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) Exists(ctx context.Context) (exists bool, err error) {
	return false, tracederrors.TracedErrorNotImplemented()
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

	return exists, nil
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
	path, hostDescription, err := n.GetPathAndHostDescription()
	if err != nil {
		return false, err
	}

	logging.LogInfoByCtxf(ctx, "Check for uncommitted changes in git repository '%s' on '%s'.", path, hostDescription)

	repo, err := git.PlainOpen(path)
	if err != nil {
		return false, tracederrors.TracedErrorf("Open git repository '%s' on '%s' failed: %w", path, hostDescription, err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return false, tracederrors.TracedErrorf("Get worktree for git repository '%s' on '%s' failed: %w", path, hostDescription, err)
	}

	status, err := worktree.Status()
	if err != nil {
		return false, tracederrors.TracedErrorf("Get status for git repository '%s' on '%s' failed: %w", path, hostDescription, err)
	}

	hasUncommitedChanges = !status.IsClean()

	if hasUncommitedChanges {
		logging.LogInfoByCtxf(ctx, "Git repository '%s' on '%s' has uncommitted changes.", path, hostDescription)
	} else {
		logging.LogInfoByCtxf(ctx, "Git repository '%s' on '%s' has no uncommitted changes.", path, hostDescription)
	}

	return hasUncommitedChanges, nil
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

func (n *NativeGitRepository) RemoteByNameExists(ctx context.Context, remoteName string) (remoteExists bool, err error) {
	if remoteName == "" {
		return false, tracederrors.TracedErrorEmptyString("remoteName")
	}

	path, hostDescription, err := n.GetPathAndHostDescription()
	if err != nil {
		return false, err
	}

	logging.LogInfoByCtxf(ctx, "Check if remote '%s' exists in git repository '%s' on '%s'.", remoteName, path, hostDescription)

	repo, err := git.PlainOpen(path)
	if err != nil {
		return false, tracederrors.TracedErrorf("Open git repository '%s' on '%s' failed: %w", path, hostDescription, err)
	}

	_, err = repo.Remote(remoteName)
	if err == nil {
		remoteExists = true
		logging.LogInfoByCtxf(ctx, "Remote '%s' exists in git repository '%s' on '%s'.", remoteName, path, hostDescription)
	} else if err == git.ErrRemoteNotFound {
		remoteExists = false
		logging.LogInfoByCtxf(ctx, "Remote '%s' does not exist in git repository '%s' on '%s'.", remoteName, path, hostDescription)
	} else {
		return false, tracederrors.TracedErrorf("Get remote '%s' in git repository '%s' on '%s' failed: %w", remoteName, path, hostDescription, err)
	}

	return remoteExists, nil
}

func (n *NativeGitRepository) RemoteConfigurationExists(ctx context.Context, config gitinterfaces.GitRemoteConfig) (exists bool, err error) {
	if config == nil {
		return false, tracederrors.TracedErrorNil("config")
	}

	remoteConfigs, err := n.GetRemoteConfigs(ctx)
	if err != nil {
		return false, err
	}

	for _, toCheck := range remoteConfigs {
		if config.Equals(toCheck) {
			return true, nil
		}
	}

	return false, nil
}

func (n *NativeGitRepository) RemoveRemoteByName(ctx context.Context, remoteName string) (err error) {
	if remoteName == "" {
		return tracederrors.TracedErrorEmptyString("remoteName")
	}

	path, hostDescription, err := n.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Remove remote '%s' from git repository '%s' on '%s' started.", remoteName, path, hostDescription)

	repo, err := git.PlainOpen(path)
	if err != nil {
		return tracederrors.TracedErrorf("Open git repository '%s' on '%s' failed: %w", path, hostDescription, err)
	}

	// Check if remote exists
	_, err = repo.Remote(remoteName)
	if err == git.ErrRemoteNotFound {
		logging.LogInfoByCtxf(ctx, "Remote '%s' does not exist in git repository '%s' on '%s'. Skip removal.", remoteName, path, hostDescription)
		return nil
	}
	if err != nil {
		return tracederrors.TracedErrorf("Get remote '%s' in git repository '%s' on '%s' failed: %w", remoteName, path, hostDescription, err)
	}

	err = repo.DeleteRemote(remoteName)
	if err != nil {
		return tracederrors.TracedErrorf("Delete remote '%s' from git repository '%s' on '%s' failed: %w", remoteName, path, hostDescription, err)
	}

	logging.LogChangedByCtxf(ctx, "Removed remote '%s' from git repository '%s' on '%s'.", remoteName, path, hostDescription)

	return nil
}

func (n *NativeGitRepository) PullFromRemote(ctx context.Context, pullOptions *gitparameteroptions.GitPullFromRemoteOptions) (err error) {
	if pullOptions == nil {
		return tracederrors.TracedErrorNil("pullOptions")
	}

	remoteName, err := pullOptions.GetRemoteName()
	if err != nil {
		return err
	}

	branchName, err := pullOptions.GetBranchName()
	if err != nil {
		return err
	}

	path, hostDescription, err := n.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Pull from remote '%s' branch '%s' in git repository '%s' on '%s' started.", remoteName, branchName, path, hostDescription)

	repo, err := git.PlainOpen(path)
	if err != nil {
		return tracederrors.TracedErrorf("Open git repository '%s' on '%s' failed: %w", path, hostDescription, err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return tracederrors.TracedErrorf("Get worktree for git repository '%s' on '%s' failed: %w", path, hostDescription, err)
	}

	err = worktree.PullContext(ctx, &git.PullOptions{
		RemoteName:    remoteName,
		ReferenceName: plumbing.NewBranchReferenceName(branchName),
	})
	if err != nil {
		if err == git.NoErrAlreadyUpToDate {
			logging.LogInfoByCtxf(ctx, "Git repository '%s' on '%s' is already up to date with remote '%s' branch '%s'.", path, hostDescription, remoteName, branchName)
		} else {
			return tracederrors.TracedErrorf("Pull from remote '%s' branch '%s' in git repository '%s' on '%s' failed: %w", remoteName, branchName, path, hostDescription, err)
		}
	}

	logging.LogInfoByCtxf(ctx, "Pull from remote '%s' branch '%s' in git repository '%s' on '%s' finished.", remoteName, branchName, path, hostDescription)

	return nil
}

func (n *NativeGitRepository) Push(ctx context.Context) (err error) {
	return tracederrors.TracedErrorNotImplemented()
}

func (n *NativeGitRepository) PushTagsToRemote(ctx context.Context, remoteName string) (err error) {
	if remoteName == "" {
		return tracederrors.TracedErrorEmptyString("remoteName")
	}

	path, hostDescription, err := n.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Push tags of git repository '%s' on '%s' to remote '%s' started.", path, hostDescription, remoteName)

	repo, err := git.PlainOpen(path)
	if err != nil {
		return tracederrors.TracedErrorf("Open git repository '%s' on '%s' failed: %w", path, hostDescription, err)
	}

	// Get all tags
	tags, err := repo.Tags()
	if err != nil {
		return tracederrors.TracedErrorf("Get tags from git repository '%s' on '%s' failed: %w", path, hostDescription, err)
	}

	// Build refspecs for all tags
	var refspecs []config.RefSpec
	err = tags.ForEach(func(ref *plumbing.Reference) error {
		refspecs = append(refspecs, config.RefSpec(ref.Name().String()+":"+ref.Name().String()))
		return nil
	})
	if err != nil {
		return tracederrors.TracedErrorf("Iterate tags in git repository '%s' on '%s' failed: %w", path, hostDescription, err)
	}

	if len(refspecs) == 0 {
		logging.LogInfoByCtxf(ctx, "No tags to push in git repository '%s' on '%s'.", path, hostDescription)
		return nil
	}

	err = repo.Push(&git.PushOptions{
		RemoteName: remoteName,
		RefSpecs:   refspecs,
	})
	if err != nil {
		if err == git.NoErrAlreadyUpToDate {
			logging.LogInfoByCtxf(ctx, "Git repository '%s' on '%s' tags are already up to date with remote '%s'.", path, hostDescription, remoteName)
			return nil
		}
		return tracederrors.TracedErrorf("Push tags to remote '%s' from git repository '%s' on '%s' failed: %w", remoteName, path, hostDescription, err)
	}

	logging.LogChangedByCtxf(ctx, "Pushed tags of git repository '%s' on '%s' to remote '%s'.", path, hostDescription, remoteName)

	return nil
}

func (n *NativeGitRepository) PushToRemote(ctx context.Context, remoteName string) (err error) {
	if remoteName == "" {
		return tracederrors.TracedErrorEmptyString("remoteName")
	}

	path, hostDescription, err := n.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Push git repository '%s' on '%s' to remote '%s' started.", path, hostDescription, remoteName)

	repo, err := git.PlainOpen(path)
	if err != nil {
		return tracederrors.TracedErrorf("Open git repository '%s' on '%s' failed: %w", path, hostDescription, err)
	}

	err = repo.Push(&git.PushOptions{
		RemoteName: remoteName,
	})
	if err != nil {
		if err == git.NoErrAlreadyUpToDate {
			logging.LogInfoByCtxf(ctx, "Git repository '%s' on '%s' is already up to date with remote '%s'.", path, hostDescription, remoteName)
			return nil
		}
		return tracederrors.TracedErrorf("Push to remote '%s' from git repository '%s' on '%s' failed: %w", remoteName, path, hostDescription, err)
	}

	logging.LogChangedByCtxf(ctx, "Pushed git repository '%s' on '%s' to remote '%s'.", path, hostDescription, remoteName)

	return nil
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
