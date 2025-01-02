package asciichgolangpublic

import (
	"errors"
	"time"
)

var ErrGitRepositoryDoesNotExist = errors.New("gitRepository does not exist")
var ErrGitRepositoryHeadNotFound = errors.New("gitRepository head not found")

// A git repository can be a LocalGitRepository or
// remote repositories like Gitlab or Github.
type GitRepository interface {
	AddFileByPath(pathToAdd string, verbose bool) (err error)
	CloneRepository(repository GitRepository, verbose bool) (err error)
	CloneRepositoryByPathOrUrl(pathOrUrl string, verbose bool) (err error)
	Commit(commitOptions *GitCommitOptions) (createdCommit *GitCommit, err error)
	CommitHasParentCommitByCommitHash(hash string) (hasParentCommit bool, err error)
	Create(verbose bool) (err error)
	CreateFileInDirectory(verbose bool, filePath ...string) (createdFile File, err error)
	CreateSubDirectory(subDirectoryName string, verbose bool) (createdSubDirectory Directory, err error)
	CreateTag(createOptions *GitRepositoryCreateTagOptions) (createdTag GitTag, err error)
	Delete(verbose bool) (err error)
	Exists(verbose bool) (exists bool, err error)
	FileByPathExists(path string, verbose bool) (exists bool, err error)
	// TODO: Will be removed as there should be no need to explitly get as local Directory:
	GetAsLocalDirectory() (localDirectory *LocalDirectory, err error)
	// TODO: Will be removed as there should be no need to explitly get as local GitRepository:
	GetAsLocalGitRepository() (localGitRepository *LocalGitRepository, err error)
	GetAuthorEmailByCommitHash(hash string) (authorEmail string, err error)
	GetAuthorStringByCommitHash(hash string) (authorString string, err error)
	GetCommitAgeDurationByCommitHash(hash string) (ageDuration *time.Duration, err error)
	GetCommitAgeSecondsByCommitHash(hash string) (ageSeconds float64, err error)
	GetCommitMessageByCommitHash(hash string) (commitMessage string, err error)
	GetCommitParentsByCommitHash(hash string, options *GitCommitGetParentsOptions) (commitParents []*GitCommit, err error)
	GetCommitTimeByCommitHash(hash string) (commitTime *time.Time, err error)
	GetCurrentCommit(verbose bool) (commit *GitCommit, err error)
	GetCurrentCommitHash(verbose bool) (currentCommitHash string, err error)
	GetGitStatusOutput(verbose bool) (output string, err error)
	GetHashByTagName(tagName string) (hash string, err error)
	GetHostDescription() (hostDescription string, err error)
	GetPath() (path string, err error)
	GetRootDirectory(verbose bool) (directory Directory, err error)
	GetRootDirectoryPath(verbose bool) (path string, err error)
	HasInitialCommit(verbose bool) (hasInitialCommit bool, err error)
	HasUncommittedChanges(verbose bool) (hasUncommitedChanges bool, err error)
	Init(options *CreateRepositoryOptions) (err error)
	IsBareRepository(verbose bool) (isBareRepository bool, err error)
	// Returns true if pointing to an existing git repository, false otherwise
	IsGitRepository(verbose bool) (isRepository bool, err error)
	IsInitialized(verbose bool) (isInitialited bool, err error)
	ListFiles(listFileOptions *ListFileOptions) (files []File, err error)
	ListTagNames(verbose bool) (tagNames []string, err error)
	ListTags(verbose bool) (tags []GitTag, err error)
	ListTagsForCommitHash(hash string, verbose bool) (tags []GitTag, err error)

	MustAddFileByPath(pathToAdd string, verbose bool)
	MustCloneRepository(repository GitRepository, verbose bool)
	MustCloneRepositoryByPathOrUrl(pathOrUrl string, verbose bool)
	MustCommit(commitOptions *GitCommitOptions) (createdCommit *GitCommit)
	MustCommitHasParentCommitByCommitHash(hash string) (hasParentCommit bool)
	MustCreate(verbose bool)
	MustCreateFileInDirectory(verbose bool, filePath ...string) (createdFile File)
	MustCreateSubDirectory(subDirectoryName string, verbose bool) (createdSubDirectory Directory)
	MustCreateTag(createOptions *GitRepositoryCreateTagOptions) (createdTag GitTag)
	MustDelete(verbose bool)
	MustExists(verbose bool) (exists bool)
	MustFileByPathExists(path string, verbose bool) (exists bool)
	// TODO: Will be removed as there should be no need to explitly get a local Directory:
	MustGetAsLocalDirectory() (localDirectory *LocalDirectory)
	// TODO: Will be removed as there should be no need to explitly get as local GitRepository:
	MustGetAsLocalGitRepository() (localGitRepository *LocalGitRepository)
	MustGetAuthorEmailByCommitHash(hash string) (authorEmail string)
	MustGetAuthorStringByCommitHash(hash string) (authorString string)
	MustGetCommitAgeDurationByCommitHash(hash string) (ageDuration *time.Duration)
	MustGetCommitAgeSecondsByCommitHash(hash string) (ageSeconds float64)
	MustGetCommitMessageByCommitHash(hash string) (commitMessage string)
	MustGetCommitParentsByCommitHash(hash string, options *GitCommitGetParentsOptions) (commitParents []*GitCommit)
	MustGetCommitTimeByCommitHash(hash string) (commitTime *time.Time)
	MustGetCurrentCommit(verbose bool) (commit *GitCommit)
	MustGetCurrentCommitHash(verbose bool) (currentCommitHash string)
	MustGetGitStatusOutput(verbose bool) (output string)
	MustGetHashByTagName(tagName string) (hash string)
	MustGetHostDescription() (hostDescription string)
	MustGetPath() (path string)
	MustGetRootDirectory(verbose bool) (directory Directory)
	MustGetRootDirectoryPath(verbose bool) (path string)
	MustHasInitialCommit(verbose bool) (hasInitialCommit bool)
	MustHasUncommittedChanges(verbose bool) (hasUncommitedChanges bool)
	MustInit(options *CreateRepositoryOptions)
	MustIsBareRepository(verbose bool) (isBareRepository bool)
	MustIsGitRepository(verbose bool) (isRepository bool)
	MustIsInitialized(verbose bool) (isInitialited bool)
	MustListFiles(listFileOptions *ListFileOptions) (files []File)
	MustListTagNames(verbose bool) (tagNames []string)
	MustListTags(verbose bool) (tags []GitTag)
	MustListTagsForCommitHash(hash string, verbose bool) (tags []GitTag)
	MustPull(verbose bool)
	MustPush(verbose bool)
	MustSetGitConfig(options *GitConfigSetOptions)
	Pull(verbose bool) (err error)
	Push(verbose bool) (err error)
	SetGitConfig(options *GitConfigSetOptions) (err error)

	// All methods below this line can be implemented by embedding the `GitRepositoryBase` struct:
	CheckHasNoUncommittedChanges(verbose bool) (err error)
	CheckIsGolangApplication(verbose bool) (err error)
	CheckIsGolangPackage(verbose bool) (err error)
	CheckIsOnLocalhost(verbose bool) (err error)
	CommitAndPush(commitOptions *GitCommitOptions) (createdCommit *GitCommit, err error)
	CreateAndInit(options *CreateRepositoryOptions) (err error)
	GetCurrentCommitsNewestVersion(verbose bool) (newestVersion Version, err error)
	GetCurrentCommitsNewestVersionOrNilIfNotPresent(verbose bool) (newestVersion Version, err error)
	GetFileByPath(path ...string) (fileInRepo File, err error)
	GetLatestTagVersion(verbose bool) (latestTagVersion Version, err error)
	GetLatestTagVersionAsString(verbose bool) (latestTagVersion string, err error)
	GetLatestTagVersionOrNilIfNotFound(verbose bool) (latestTagVersion Version, err error)
	GetPathAndHostDescription() (path string, hostDescription string, err error)
	HasNoUncommittedChanges(verbose bool) (noUncommitedChnages bool, err error)
	IsGolangApplication(verbose bool) (isGolangApplication bool, err error)
	IsGolangPackage(verbose bool) (isGolangPackage bool, err error)
	IsOnLocalhost(verbose bool) (isOnLocalhost bool, err error)
	ListVersionTags(verbose bool) (versionTags []GitTag, err error)
	MustCheckHasNoUncommittedChanges(verbose bool)
	MustCheckIsGolangApplication(verbose bool)
	MustCheckIsGolangPackage(verbose bool)
	MustCheckIsOnLocalhost(verbose bool)
	MustCommitAndPush(commitOptions *GitCommitOptions) (createdCommit *GitCommit)
	MustCreateAndInit(options *CreateRepositoryOptions)
	MustGetCurrentCommitsNewestVersion(verbose bool) (newestVersion Version)
	MustGetCurrentCommitsNewestVersionOrNilIfNotPresent(verbose bool) (newestVersion Version)
	MustGetFileByPath(path ...string) (fileInRepo File)
	MustGetLatestTagVersion(verbose bool) (latestTagVersion Version)
	MustGetLatestTagVersionAsString(verbose bool) (latestTagVersion string)
	MustGetLatestTagVersionOrNilIfNotFound(verbose bool) (latestTagVersion Version)
	MustGetPathAndHostDescription() (path string, hostDescription string)
	MustHasNoUncommittedChanges(verbose bool) (noUncommitedChnages bool)
	MustIsGolangApplication(verbose bool) (isGolangApplication bool)
	MustIsGolangPackage(verbose bool) (isGolangPackage bool)
	MustIsOnLocalhost(verbose bool) (isOnLocalhost bool)
	MustListVersionTags(verbose bool) (versionTags []GitTag)
	MustWriteBytesToFile(content []byte, verbose bool, path ...string) (writtenFile File)
	MustWriteStringToFile(content string, verbose bool, path ...string) (writtenFile File)
	WriteBytesToFile(content []byte, verbose bool, path ...string) (writtenFile File, err error)
	WriteStringToFile(content string, verbose bool, path ...string) (writtenFile File, err error)
}

func GitRepositoryDefaultCommitMessageForInitializeWithEmptyCommit() (msg string) {
	return "Initial empty commit during repo initialization"
}

func GitRepositryDefualtAuthorEmail() (email string) {
	return "asciichgolangpublic@example.net"
}

func GitRepositryDefualtAuthorName() (name string) {
	return "asciichgolangpublic git repo initializer"
}
