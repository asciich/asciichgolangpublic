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
	Commit(commitOptions *GitCommitOptions) (createdCommit *GitCommit, err error)
	CommitHasParentCommitByCommitHash(hash string) (hasParentCommit bool, err error)
	Create(verbose bool) (err error)
	Exists(verbose bool) (exists bool, err error)
	Delete(verbose bool) (err error)
	GetAsLocalDirectory() (localDirectory *LocalDirectory, err error)
	GetAsLocalGitRepository() (localGitRepository *LocalGitRepository, err error)
	GetAuthorEmailByCommitHash(hash string) (authorEmail string, err error)
	GetAuthorStringByCommitHash(hash string) (authorString string, err error)
	GetCommitAgeDurationByCommitHash(hash string) (ageDuration *time.Duration, err error)
	GetCommitAgeSecondsByCommitHash(hash string) (ageSeconds float64, err error)
	GetCommitMessageByCommitHash(hash string) (commitMessage string, err error)
	GetCommitParentsByCommitHash(hash string, options *GitCommitGetParentsOptions) (commitParents []*GitCommit, err error)
	GetCommitTimeByCommitHash(hash string) (commitTime *time.Time, err error)
	GetCurrentCommit() (commit *GitCommit, err error)
	GetCurrentCommitHash() (currentCommitHash string, err error)
	GetGitStatusOutput(verbose bool) (output string, err error)
	GetPath() (path string, err error)
	HasInitialCommit(verbose bool) (hasInitialCommit bool, err error)
	Init(options *CreateRepositoryOptions) (err error)
	IsBareRepository(verbose bool) (isBareRepository bool, err error)
	IsInitialized(verbose bool) (isInitialited bool, err error)

	MustCommit(commitOptions *GitCommitOptions) (createdCommit *GitCommit)
	MustCommitHasParentCommitByCommitHash(hash string) (hasParentCommit bool)
	MustCreate(verbose bool)
	MustDelete(verbose bool)
	MustExists(verbose bool) (exists bool)
	MustGetAsLocalDirectory() (localDirectory *LocalDirectory)
	MustGetAsLocalGitRepository() (localGitRepository *LocalGitRepository)
	MustGetAuthorEmailByCommitHash(hash string) (authorEmail string)
	MustGetAuthorStringByCommitHash(hash string) (authorString string)
	MustGetCommitAgeDurationByCommitHash(hash string) (ageDuration *time.Duration)
	MustGetCommitAgeSecondsByCommitHash(hash string) (ageSeconds float64)
	MustGetCommitMessageByCommitHash(hash string) (commitMessage string)
	MustGetCommitParentsByCommitHash(hash string, options *GitCommitGetParentsOptions) (commitParents []*GitCommit)
	MustGetCommitTimeByCommitHash(hash string) (commitTime *time.Time)
	MustGetCurrentCommit() (commit *GitCommit)
	MustGetCurrentCommitHash() (currentCommitHash string)
	MustGetGitStatusOutput(verbose bool) (output string)
	MustGetPath() (path string)
	MustHasInitialCommit(verbose bool) (hasInitialCommit bool)
	MustInit(options *CreateRepositoryOptions)
	MustIsBareRepository(verbose bool) (isBareRepository bool)
	MustIsInitialized(verbose bool) (isInitialited bool)

	// All methods below this line can be implemented by embedding the `GitRepositoryBase` struct:
	CreateAndInit(options *CreateRepositoryOptions) (err error)
	MustCreateAndInit(options *CreateRepositoryOptions)
}

func GitRepositoryDefaultCommitMessageForInitializeWithEmptyCommit() (msg string) {
	return "Initial empty commit during repo initialization"
}

func GitRepositryDefualtAuthorName() (name string) {
	return "asciichgolangpublic git repo initializer"
}

func GitRepositryDefualtAuthorEmail() (email string) {
	return "asciichgolangpublic@example.net"
}
