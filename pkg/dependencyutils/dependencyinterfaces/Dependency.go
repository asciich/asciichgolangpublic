package dependencyinterfaces

import (
	"github.com/asciich/asciichgolangpublic/changesummary"
	"github.com/asciich/asciichgolangpublic/files"
	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/parameteroptions/authenticationoptions"
)

// A Dependency is used to implement software and other dependencies like container images...
type Dependency interface {
	AddSourceFile(files.File) (err error)
	GetName() (name string, err error)
	GetNewestVersionAsString(authOptions []authenticationoptions.AuthenticationOption, verbose bool) (newestVersion string, err error)
	IsUpdateAvailable(authOptions []authenticationoptions.AuthenticationOption, verbose bool) (isUpdateAvailable bool, err error)
	Update(options *parameteroptions.UpdateDependenciesOptions) (changeSummary *changesummary.ChangeSummary, err error)
}
