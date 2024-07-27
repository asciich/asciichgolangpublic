package asciichgolangpublic

import "time"

// A git repository can be a LocalGitRepository or
// remote repositories like Gitlab or Github.
type GitRepository interface {
	GetAuthorEmailByCommitHash(hash string) (authorEmail string, err error)
	GetAuthorStringByCommitHash(hash string) (authorString string, err error)
	GetCommitAgeDurationByCommitHash(hash string) (ageDuration *time.Duration, err error)
	GetCommitAgeSecondsByCommitHash(hash string) (ageSeconds float64, err error)
	GetCommitMessageByCommitHash(hash string) (commitMessage string, err error)
	GetCommitTimeByCommitHash(hash string) (commitTime *time.Time, err error)

	MustGetAuthorEmailByCommitHash(hash string) (authorEmail string)
	MustGetAuthorStringByCommitHash(hash string) (authorString string)
	MustGetCommitAgeDurationByCommitHash(hash string) (ageDuration *time.Duration)
	MustGetCommitAgeSecondsByCommitHash(hash string) (ageSeconds float64)
	MustGetCommitMessageByCommitHash(hash string) (commitMessage string)
	MustGetCommitTimeByCommitHash(hash string) (commitTime *time.Time)
}
