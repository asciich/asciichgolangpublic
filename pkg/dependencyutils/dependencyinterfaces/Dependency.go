package dependencyinterfaces

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/changesummary"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions/authenticationoptions"
)

// A Dependency is used to implement software and other dependencies like container images...
type Dependency interface {
	AddSourceFile(filesinterfaces.File) (err error)
	GetName() (name string, err error)
	GetNewestVersionAsString(ctx context.Context, authOptions []authenticationoptions.AuthenticationOption) (newestVersion string, err error)
	IsUpdateAvailable(ctx context.Context, authOptions []authenticationoptions.AuthenticationOption) (isUpdateAvailable bool, err error)
	Update(ctx context.Context, options *parameteroptions.UpdateDependenciesOptions) (changeSummary *changesummary.ChangeSummary, err error)
}
