package asciichgolangpublic


// A Dependency is used to implement software and other dependencies like container images...
type Dependency interface {
	AddSourceFile(File) (err error)
	GetName() (name string, err error)
	GetNewestVersionAsString(authOptions []AuthenticationOption, verbose bool) (newestVersion string, err error)
	IsUpdateAvailable(authOptions []AuthenticationOption, verbose bool) (isUpdateAvailable bool, err error)
	Update(options *UpdateDependenciesOptions) (changeSummary *ChangeSummary, err error)
}
