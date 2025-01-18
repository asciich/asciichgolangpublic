package asciichgolangpublic

import (
	"github.com/asciich/asciichgolangpublic/changesummary"
	"github.com/asciich/asciichgolangpublic/files"
)

// A Dependency is used to implement software and other dependencies like container images...
type Dependency interface {
	AddSourceFile(files.File) (err error)
	GetName() (name string, err error)
	GetNewestVersionAsString(authOptions []AuthenticationOption, verbose bool) (newestVersion string, err error)
	IsUpdateAvailable(authOptions []AuthenticationOption, verbose bool) (isUpdateAvailable bool, err error)
	Update(options *UpdateDependenciesOptions) (changeSummary *changesummary.ChangeSummary, err error)
}
