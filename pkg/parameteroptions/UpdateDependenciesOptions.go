package parameteroptions

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/artifacthandler"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions/authenticationoptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type UpdateDependenciesOptions struct {
	ArtifactHandlers      []artifacthandler.ArtifactHandler
	Commit                bool
	AuthenticationOptions []authenticationoptions.AuthenticationOption
}

func NewUpdateDependenciesOptions() (u *UpdateDependenciesOptions) {
	return new(UpdateDependenciesOptions)
}

func (u *UpdateDependenciesOptions) GetArtifactHandlerForSoftwareName(softwareName string) (artifactHandler artifacthandler.ArtifactHandler, err error) {
	if softwareName == "" {
		return nil, tracederrors.TracedError("softwareName is empty string")
	}

	handlers, err := u.GetArtifactHandlers()
	if err != nil {
		return nil, err
	}

	for _, artifactHandler := range handlers {
		isHandlingSoftware, err := artifactHandler.IsHandlingArtifactByName(softwareName)
		if err != nil {
			return nil, err
		}

		if isHandlingSoftware {
			return artifactHandler, nil
		}
	}

	return nil, tracederrors.TracedErrorf("No handler for softwareName '%s' found", softwareName)
}

func (u *UpdateDependenciesOptions) GetArtifactHandlers() (artifactHandlers []artifacthandler.ArtifactHandler, err error) {
	if u.ArtifactHandlers == nil {
		return nil, tracederrors.TracedErrorf("ArtifactHandlers not set")
	}

	if len(u.ArtifactHandlers) <= 0 {
		return nil, tracederrors.TracedErrorf("ArtifactHandlers has no elements")
	}

	return u.ArtifactHandlers, nil
}

func (u *UpdateDependenciesOptions) GetAuthenticationOptions() (authenticationOptions []authenticationoptions.AuthenticationOption, err error) {
	if u.AuthenticationOptions == nil {
		return nil, tracederrors.TracedErrorf("AuthenticationOptions not set")
	}

	if len(u.AuthenticationOptions) <= 0 {
		return nil, tracederrors.TracedErrorf("AuthenticationOptions has no elements")
	}

	return u.AuthenticationOptions, nil
}

func (u *UpdateDependenciesOptions) GetCommit() (commit bool, err error) {

	return u.Commit, nil
}

func (u *UpdateDependenciesOptions) GetLatestArtifactVersionAsString(ctx context.Context, softwareName string) (latestVersion string, err error) {
	if softwareName == "" {
		return "", tracederrors.TracedError("softwareName is empty string")
	}

	artifactHandler, err := u.GetArtifactHandlerForSoftwareName(softwareName)
	if err != nil {
		return "", err
	}

	latestVersion, err = artifactHandler.GetLatestArtifactVersionAsString(ctx, softwareName)
	if err != nil {
		return "", err
	}

	return latestVersion, err
}

func (u *UpdateDependenciesOptions) SetArtifactHandlers(artifactHandlers []artifacthandler.ArtifactHandler) (err error) {
	if artifactHandlers == nil {
		return tracederrors.TracedErrorf("artifactHandlers is nil")
	}

	if len(artifactHandlers) <= 0 {
		return tracederrors.TracedErrorf("artifactHandlers has no elements")
	}

	u.ArtifactHandlers = artifactHandlers

	return nil
}

func (u *UpdateDependenciesOptions) SetAuthenticationOptions(authenticationOptions []authenticationoptions.AuthenticationOption) (err error) {
	if authenticationOptions == nil {
		return tracederrors.TracedErrorf("authenticationOptions is nil")
	}

	if len(authenticationOptions) <= 0 {
		return tracederrors.TracedErrorf("authenticationOptions has no elements")
	}

	u.AuthenticationOptions = authenticationOptions

	return nil
}

func (u *UpdateDependenciesOptions) SetCommit(commit bool) (err error) {
	u.Commit = commit

	return nil
}
