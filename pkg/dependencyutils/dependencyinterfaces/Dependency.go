package dependencyinterfaces

import (
	"github.com/asciich/asciichgolangpublic/pkg/changesummary"
	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions/authenticationoptions"
)

// A Dependency is used to implement software and other dependencies like container images...
type Dependency interface {
	AddSourceFile(files.File) (err error)
	GetName() (name string, err error)
	GetNewestVersionAsString(authOptions []authenticationoptions.AuthenticationOption, verbose bool) (newestVersion string, err error)
	IsUpdateAvailable(authOptions []authenticationoptions.AuthenticationOption, verbose bool) (isUpdateAvailable bool, err error)
	Update(options *parameteroptions.UpdateDependenciesOptions) (changeSummary *changesummary.ChangeSummary, err error)
}
