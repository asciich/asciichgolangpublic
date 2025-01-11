package containers

type Container interface {
	IsRunning(verbose bool) (isRunning bool, err error)
	Kill(verbose bool) (err error)
	MustIsRunning(verbose bool) (isRunning bool)
	MustKill(verbose bool)
}
