package asciichgolangpublic

type GitTag interface {
	GetHash() (hash string, err error)
	GetName() (name string, err error)
	GetGitRepository() (repo GitRepository, err error)
	IsVersionTag() (isVersionTag bool, err error)
	SetName(name string) (err error)
	MustGetHash() (hash string)
	MustGetName() (name string)
	MustGetGitRepository() (repo GitRepository)
	MustIsVersionTag() (isVersionTag bool)
	MustSetName(name string)

	// These function can be implemented by embedding the GitTagBase struct:
	GetVersion() (version Version, err error)
	MustGetVersion() (version Version)
}
