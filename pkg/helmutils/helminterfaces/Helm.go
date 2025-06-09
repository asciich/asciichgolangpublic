package helminterfaces

type Helm interface {
	AddRepositoryByName(name string, url string, verbose bool) (err error)
	MustAddRepositoryByName(name string, url string, verbose bool)
}
