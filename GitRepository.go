package asciichgolangpublic

import (
	"context"
	"errors"
	"time"

	"github.com/asciich/asciichgolangpublic/datatypes"
	"github.com/asciich/asciichgolangpublic/files"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/versionutils"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

var ErrGitRepositoryDoesNotExist = errors.New("gitRepository does not exist")
var ErrGitRepositoryHeadNotFound = errors.New("gitRepository head not found")

// A git repository can be a LocalGitRepository or
// remote repositories like Gitlab or Github.
type GitRepository interface {
	AddRemote(options *GitRemoteAddOptions) (err error)
	AddFileByPath(pathToAdd string, verbose bool) (err error)
	CheckoutBranchByName(name string, verbose bool) (err error)
	CloneRepository(repository GitRepository, verbose bool) (err error)
	CloneRepositoryByPathOrUrl(pathOrUrl string, verbose bool) (err error)
	Commit(commitOptions *GitCommitOptions) (createdCommit *GitCommit, err error)
	CommitHasParentCommitByCommitHash(hash string) (hasParentCommit bool, err error)
	CreateBranch(createOptions *parameteroptions.CreateBranchOptions) (err error)
	Create(verbose bool) (err error)
	CreateFileInDirectory(verbose bool, filePath ...string) (createdFile files.File, err error)
	CreateSubDirectory(subDirectoryName string, verbose bool) (createdSubDirectory files.Directory, err error)
	CreateTag(createOptions *GitRepositoryCreateTagOptions) (createdTag GitTag, err error)
	Delete(verbose bool) (err error)
	DeleteBranchByName(name string, verbose bool) (err error)
	DirectoryByPathExists(verbose bool, path ...string) (exists bool, err error)
	Exists(verbose bool) (exists bool, err error)
	Fetch(verbose bool) (err error)
	FileByPathExists(path string, verbose bool) (exists bool, err error)
	// TODO: Will be removed as there should be no need to explitly get as local Directory:
	GetAsLocalDirectory() (localDirectory *files.LocalDirectory, err error)
	// TODO: Will be removed as there should be no need to explitly get as local GitRepository:
	GetAsLocalGitRepository() (localGitRepository *LocalGitRepository, err error)
	GetAuthorEmailByCommitHash(hash string) (authorEmail string, err error)
	GetAuthorStringByCommitHash(hash string) (authorString string, err error)
	GetDirectoryByPath(pathToSubDir ...string) (subDir files.Directory, err error)
	GetCommitAgeDurationByCommitHash(hash string) (ageDuration *time.Duration, err error)
	GetCommitAgeSecondsByCommitHash(hash string) (ageSeconds float64, err error)
	GetCommitMessageByCommitHash(hash string) (commitMessage string, err error)
	GetCommitParentsByCommitHash(hash string, options *parameteroptions.GitCommitGetParentsOptions) (commitParents []*GitCommit, err error)
	GetCommitTimeByCommitHash(hash string) (commitTime *time.Time, err error)
	GetCurrentBranchName(verbose bool) (branchName string, err error)
	GetCurrentCommit(verbose bool) (commit *GitCommit, err error)
	GetCurrentCommitHash(verbose bool) (currentCommitHash string, err error)
	GetGitStatusOutput(verbose bool) (output string, err error)
	GetHashByTagName(tagName string) (hash string, err error)
	GetHostDescription() (hostDescription string, err error)
	GetPath() (path string, err error)
	GetRemoteConfigs(verbose bool) (remoteConfigs []*GitRemoteConfig, err error)
	GetRootDirectory(ctx context.Context) (directory files.Directory, err error)
	GetRootDirectoryPath(ctx context.Context) (path string, err error)
	HasInitialCommit(verbose bool) (hasInitialCommit bool, err error)
	HasUncommittedChanges(verbose bool) (hasUncommitedChanges bool, err error)
	Init(options *parameteroptions.CreateRepositoryOptions) (err error)
	IsBareRepository(ctx context.Context) (isBareRepository bool, err error)
	// Returns true if pointing to an existing git repository, false otherwise
	IsGitRepository(verbose bool) (isRepository bool, err error)
	IsInitialized(verbose bool) (isInitialited bool, err error)
	ListBranchNames(verbose bool) (branchNames []string, err error)
	ListFilePaths(ctx context.Context, listFileOptions *parameteroptions.ListFileOptions) (filePaths []string, err error)
	ListFiles(ctx context.Context, listFileOptions *parameteroptions.ListFileOptions) (files []files.File, err error)
	ListTagNames(verbose bool) (tagNames []string, err error)
	ListTags(verbose bool) (tags []GitTag, err error)
	ListTagsForCommitHash(hash string, verbose bool) (tags []GitTag, err error)
	MustAddRemote(options *GitRemoteAddOptions)
	MustAddFileByPath(pathToAdd string, verbose bool)
	MustCheckoutBranchByName(name string, verbose bool)
	MustCloneRepository(repository GitRepository, verbose bool)
	MustCloneRepositoryByPathOrUrl(pathOrUrl string, verbose bool)
	MustCommit(commitOptions *GitCommitOptions) (createdCommit *GitCommit)
	MustCommitHasParentCommitByCommitHash(hash string) (hasParentCommit bool)
	MustCreate(verbose bool)
	MustCreateBranch(createOptions *parameteroptions.CreateBranchOptions)
	MustCreateFileInDirectory(verbose bool, filePath ...string) (createdFile files.File)
	MustCreateSubDirectory(subDirectoryName string, verbose bool) (createdSubDirectory files.Directory)
	MustCreateTag(createOptions *GitRepositoryCreateTagOptions) (createdTag GitTag)
	MustDelete(verbose bool)
	MustDeleteBranchByName(name string, verbose bool)
	MustDirectoryByPathExists(verbose bool, path ...string) (exists bool)
	MustExists(verbose bool) (exists bool)
	MustFetch(verbose bool)
	MustFileByPathExists(path string, verbose bool) (exists bool)
	// TODO: Will be removed as there should be no need to explitly get a local Directory:
	MustGetAsLocalDirectory() (localDirectory *files.LocalDirectory)
	// TODO: Will be removed as there should be no need to explitly get as local GitRepository:
	MustGetAsLocalGitRepository() (localGitRepository *LocalGitRepository)
	MustGetDirectoryByPath(pathToSubDir ...string) (subDir files.Directory)
	MustGetAuthorEmailByCommitHash(hash string) (authorEmail string)
	MustGetAuthorStringByCommitHash(hash string) (authorString string)
	MustGetCommitAgeDurationByCommitHash(hash string) (ageDuration *time.Duration)
	MustGetCommitAgeSecondsByCommitHash(hash string) (ageSeconds float64)
	MustGetCommitMessageByCommitHash(hash string) (commitMessage string)
	MustGetCommitParentsByCommitHash(hash string, options *parameteroptions.GitCommitGetParentsOptions) (commitParents []*GitCommit)
	MustGetCommitTimeByCommitHash(hash string) (commitTime *time.Time)
	MustGetCurrentBranchName(verbose bool) (branchName string)
	MustGetCurrentCommit(verbose bool) (commit *GitCommit)
	MustGetCurrentCommitHash(verbose bool) (currentCommitHash string)
	MustGetGitStatusOutput(verbose bool) (output string)
	MustGetHashByTagName(tagName string) (hash string)
	MustGetHostDescription() (hostDescription string)
	MustGetPath() (path string)
	MustGetRemoteConfigs(verbose bool) (remoteConfigs []*GitRemoteConfig)
	MustHasInitialCommit(verbose bool) (hasInitialCommit bool)
	MustHasUncommittedChanges(verbose bool) (hasUncommitedChanges bool)
	MustInit(options *parameteroptions.CreateRepositoryOptions)
	MustIsGitRepository(verbose bool) (isRepository bool)
	MustIsInitialized(verbose bool) (isInitialited bool)
	MustListBranchNames(verbose bool) (branchNames []string)
	MustListTagNames(verbose bool) (tagNames []string)
	MustListTags(verbose bool) (tags []GitTag)
	MustListTagsForCommitHash(hash string, verbose bool) (tags []GitTag)
	MustRemoteByNameExists(remoteName string, verbose bool) (remoteExists bool)
	MustRemoteConfigurationExists(config *GitRemoteConfig, verbose bool) (exists bool)
	MustRemoveRemoteByName(remoteName string, verbose bool)
	MustPullFromRemote(pullOptions *GitPullFromRemoteOptions)
	MustPushTagsToRemote(remoteName string, verbose bool)
	MustPushToRemote(remoteName string, verbose bool)
	MustSetGitConfig(options *GitConfigSetOptions)
	MustSetRemoteUrl(remoteUrl string, verbose bool)
	RemoteByNameExists(remoteName string, verbose bool) (remoteExists bool, err error)
	RemoteConfigurationExists(config *GitRemoteConfig, verbose bool) (exists bool, err error)
	RemoveRemoteByName(remoteName string, verbose bool) (err error)
	Pull(ctx context.Context) (err error)
	PullFromRemote(pullOptions *GitPullFromRemoteOptions) (err error)
	Push(ctx context.Context) (err error)
	PushTagsToRemote(remoteName string, verbose bool) (err error)
	PushToRemote(remoteName string, verbose bool) (err error)
	SetGitConfig(options *GitConfigSetOptions) (err error)
	SetRemoteUrl(remoteUrl string, verbose bool) (err error)

	// All methods below this line can be implemented by embedding the `GitRepositoryBase` struct:
	AddFilesByPath(pathsToAdd []string, verbose bool) (err error)
	BranchByNameExists(branchName string, verbose bool) (branchExists bool, err error)
	CheckExists(ctx context.Context) (err error)
	CheckHasNoUncommittedChanges(verbose bool) (err error)
	CheckIsGolangApplication(verbose bool) (err error)
	CheckIsGolangPackage(verbose bool) (err error)
	CheckIsOnLocalhost(verbose bool) (err error)
	CheckIsPreCommitRepository(verbose bool) (err error)
	CommitAndPush(commitOptions *GitCommitOptions) (createdCommit *GitCommit, err error)
	CommitIfUncommittedChanges(commitOptions *GitCommitOptions) (createdCommit *GitCommit, err error)
	CreateAndInit(options *parameteroptions.CreateRepositoryOptions) (err error)
	EnsureMainReadmeMdExists(verbose bool) (err error)
	GetCurrentCommitMessage(verbose bool) (currentCommitMessage string, err error)
	GetCurrentCommitsNewestVersion(verbose bool) (newestVersion versionutils.Version, err error)
	GetCurrentCommitsNewestVersionOrNilIfNotPresent(verbose bool) (newestVersion versionutils.Version, err error)
	GetFileByPath(path ...string) (fileInRepo files.File, err error)
	GetLatestTagVersion(verbose bool) (latestTagVersion versionutils.Version, err error)
	GetLatestTagVersionAsString(verbose bool) (latestTagVersion string, err error)
	GetLatestTagVersionOrNilIfNotFound(verbose bool) (latestTagVersion versionutils.Version, err error)
	GetPathAndHostDescription() (path string, hostDescription string, err error)
	HasNoUncommittedChanges(verbose bool) (noUncommitedChnages bool, err error)
	IsGolangApplication(verbose bool) (isGolangApplication bool, err error)
	IsGolangPackage(verbose bool) (isGolangPackage bool, err error)
	IsOnLocalhost(verbose bool) (isOnLocalhost bool, err error)
	IsPreCommitRepository(verbose bool) (isPreCommitRepository bool, err error)
	ListVersionTags(verbose bool) (versionTags []GitTag, err error)
	MustAddFilesByPath(pathsToAdd []string, verbose bool)
	MustBranchByNameExists(branchName string, verbose bool) (branchExists bool)
	MustCheckExists(ctx context.Context)
	MustCheckHasNoUncommittedChanges(verbose bool)
	MustCheckIsGolangApplication(verbose bool)
	MustCheckIsGolangPackage(verbose bool)
	MustCheckIsOnLocalhost(verbose bool)
	MustCheckIsPreCommitRepository(verbose bool)
	MustCommitAndPush(commitOptions *GitCommitOptions) (createdCommit *GitCommit)
	MustCommitIfUncommittedChanges(commitOptions *GitCommitOptions) (createdCommit *GitCommit)
	MustCreateAndInit(options *parameteroptions.CreateRepositoryOptions)
	MustEnsureMainReadmeMdExists(verbose bool)
	MustGetCurrentCommitMessage(verbose bool) (currentCommitMessage string)
	MustGetCurrentCommitsNewestVersion(verbose bool) (newestVersion versionutils.Version)
	MustGetCurrentCommitsNewestVersionOrNilIfNotPresent(verbose bool) (newestVersion versionutils.Version)
	MustGetFileByPath(path ...string) (fileInRepo files.File)
	MustGetLatestTagVersion(verbose bool) (latestTagVersion versionutils.Version)
	MustGetLatestTagVersionAsString(verbose bool) (latestTagVersion string)
	MustGetLatestTagVersionOrNilIfNotFound(verbose bool) (latestTagVersion versionutils.Version)
	MustGetPathAndHostDescription() (path string, hostDescription string)
	MustHasNoUncommittedChanges(verbose bool) (noUncommitedChnages bool)
	MustIsGolangApplication(verbose bool) (isGolangApplication bool)
	MustIsGolangPackage(verbose bool) (isGolangPackage bool)
	MustIsOnLocalhost(verbose bool) (isOnLocalhost bool)
	MustIsPreCommitRepository(verbose bool) (isPreCommitRepository bool)
	MustListVersionTags(verbose bool) (versionTags []GitTag)
	MustWriteBytesToFile(content []byte, verbose bool, path ...string) (writtenFile files.File)
	MustWriteStringToFile(content string, verbose bool, path ...string) (writtenFile files.File)
	WriteBytesToFile(content []byte, verbose bool, path ...string) (writtenFile files.File, err error)
	WriteStringToFile(content string, verbose bool, path ...string) (writtenFile files.File, err error)
}

func GetGitRepositoryByDirectory(directory files.Directory) (repository GitRepository, err error) {
	if directory == nil {
		return nil, tracederrors.TracedErrorNil("directory")
	}

	localDirectory, ok := directory.(*files.LocalDirectory)
	if ok {
		return GetLocalGitReposioryFromDirectory(localDirectory)
	}

	commandExecutorDirectory, ok := directory.(*files.CommandExecutorDirectory)
	if ok {
		return GetCommandExecutorGitRepositoryFromDirectory(commandExecutorDirectory)
	}

	unknownTypeName, err := datatypes.GetTypeName(directory)
	if err != nil {
		return nil, err
	}

	return nil, tracederrors.TracedErrorf(
		"Unknown directory implementation '%s'. Unable to get GitRepository",
		unknownTypeName,
	)
}

func GitRepositoryDefaultCommitMessageForInitializeWithEmptyCommit() (msg string) {
	return "Initial empty commit during repo initialization"
}

func GitRepositryDefaultAuthorEmail() (email string) {
	return "asciichgolangpublic@example.net"
}

func GitRepositryDefaultAuthorName() (name string) {
	return "asciichgolangpublic git repo initializer"
}

func MustGetGitRepositoryByDirectory(directory files.Directory) (repository GitRepository) {
	repository, err := GetGitRepositoryByDirectory(directory)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return repository
}
