package parameteroptions

import (
	"github.com/asciich/asciichgolangpublic/artifacthandler"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/parameteroptions/authenticationoptions"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type UpdateDependenciesOptions struct {
	ArtifactHandlers      []artifacthandler.ArtifactHandler
	Commit                bool
	Verbose               bool
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

func (u *UpdateDependenciesOptions) GetLatestArtifactVersionAsString(softwareName string, verbose bool) (latestVersion string, err error) {
	if softwareName == "" {
		return "", tracederrors.TracedError("softwareName is empty string")
	}

	artifactHandler, err := u.GetArtifactHandlerForSoftwareName(softwareName)
	if err != nil {
		return "", err
	}

	latestVersion, err = artifactHandler.GetLatestArtifactVersionAsString(softwareName, verbose)
	if err != nil {
		return "", err
	}

	return latestVersion, err
}

func (u *UpdateDependenciesOptions) GetVerbose() (verbose bool, err error) {

	return u.Verbose, nil
}

func (u *UpdateDependenciesOptions) MustGetArtifactHandlerForSoftwareName(softwareName string) (artifactHandler artifacthandler.ArtifactHandler) {
	artifactHandler, err := u.GetArtifactHandlerForSoftwareName(softwareName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return artifactHandler
}

func (u *UpdateDependenciesOptions) MustGetArtifactHandlers() (artifactHandlers []artifacthandler.ArtifactHandler) {
	artifactHandlers, err := u.GetArtifactHandlers()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return artifactHandlers
}

func (u *UpdateDependenciesOptions) MustGetAuthenticationOptions() (authenticationOptions []authenticationoptions.AuthenticationOption) {
	authenticationOptions, err := u.GetAuthenticationOptions()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return authenticationOptions
}

func (u *UpdateDependenciesOptions) MustGetCommit() (commit bool) {
	commit, err := u.GetCommit()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return commit
}

func (u *UpdateDependenciesOptions) MustGetLatestArtifactVersionAsString(softwareName string, verbose bool) (latestVersion string) {
	latestVersion, err := u.GetLatestArtifactVersionAsString(softwareName, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return latestVersion
}

func (u *UpdateDependenciesOptions) MustGetVerbose() (verbose bool) {
	verbose, err := u.GetVerbose()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return verbose
}

func (u *UpdateDependenciesOptions) MustSetArtifactHandlers(artifactHandlers []artifacthandler.ArtifactHandler) {
	err := u.SetArtifactHandlers(artifactHandlers)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (u *UpdateDependenciesOptions) MustSetAuthenticationOptions(authenticationOptions []authenticationoptions.AuthenticationOption) {
	err := u.SetAuthenticationOptions(authenticationOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (u *UpdateDependenciesOptions) MustSetCommit(commit bool) {
	err := u.SetCommit(commit)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (u *UpdateDependenciesOptions) MustSetVerbose(verbose bool) {
	err := u.SetVerbose(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
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

func (u *UpdateDependenciesOptions) SetVerbose(verbose bool) (err error) {
	u.Verbose = verbose

	return nil
}
