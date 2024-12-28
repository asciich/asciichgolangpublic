package asciichgolangpublic

type GitTag interface {
	GetName() (name string, err error)
	GetGitRepository() (repo GitRepository, err error)
	MustGetName() (name string)
	MustGetGitRepository() (repo GitRepository)
}
