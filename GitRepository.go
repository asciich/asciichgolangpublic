package asciichgolangpublic

import "time"

// A git repository can be a LocalGitRepository or
// remote repositories like Gitlab or Github.
type GitRepository interface {
	CommitHasParentCommitByCommitHash(hash string) (hasParentCommit bool, err error)
	Create(verbose bool) (err error)
	CreateAndInit(options *CreateRepositoryOptions) (err error)
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
	GetGitStatusOutput(verbose bool) (output string, err error)
	GetPath() (path string, err error)
	Init(options *CreateRepositoryOptions) (err error)
	IsInitialized(verbose bool) (isInitialited bool, err error)

	MustCommitHasParentCommitByCommitHash(hash string) (hasParentCommit bool)
	MustCreate(verbose bool)
	MustCreateAndInit(options *CreateRepositoryOptions)
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
	MustGetGitStatusOutput(verbose bool) (output string)
	MustGetPath() (path string)
	MustInit(options *CreateRepositoryOptions)
	MustIsInitialized(verbose bool) (isInitialited bool)
}
