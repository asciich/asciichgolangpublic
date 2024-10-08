package asciichgolangpublic

type UpdateDependenciesOptions struct {
	ArtifactHandlers      []ArtifactHandler
	Commit                bool
	Verbose               bool
	AuthenticationOptions []AuthenticationOption
}

func NewUpdateDependenciesOptions() (u *UpdateDependenciesOptions) {
	return new(UpdateDependenciesOptions)
}

func (u *UpdateDependenciesOptions) GetArtifactHandlerForSoftwareName(softwareName string) (artifactHandler ArtifactHandler, err error) {
	if softwareName == "" {
		return nil, TracedError("softwareName is empty string")
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

	return nil, TracedErrorf("No handler for softwareName '%s' found", softwareName)
}

func (u *UpdateDependenciesOptions) GetArtifactHandlers() (artifactHandlers []ArtifactHandler, err error) {
	if u.ArtifactHandlers == nil {
		return nil, TracedErrorf("ArtifactHandlers not set")
	}

	if len(u.ArtifactHandlers) <= 0 {
		return nil, TracedErrorf("ArtifactHandlers has no elements")
	}

	return u.ArtifactHandlers, nil
}

func (u *UpdateDependenciesOptions) GetAuthenticationOptions() (authenticationOptions []AuthenticationOption, err error) {
	if u.AuthenticationOptions == nil {
		return nil, TracedErrorf("AuthenticationOptions not set")
	}

	if len(u.AuthenticationOptions) <= 0 {
		return nil, TracedErrorf("AuthenticationOptions has no elements")
	}

	return u.AuthenticationOptions, nil
}

func (u *UpdateDependenciesOptions) GetCommit() (commit bool, err error) {

	return u.Commit, nil
}

func (u *UpdateDependenciesOptions) GetLatestArtifactVersionAsString(softwareName string, verbose bool) (latestVersion string, err error) {
	if softwareName == "" {
		return "", TracedError("softwareName is empty string")
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

func (u *UpdateDependenciesOptions) MustGetArtifactHandlerForSoftwareName(softwareName string) (artifactHandler ArtifactHandler) {
	artifactHandler, err := u.GetArtifactHandlerForSoftwareName(softwareName)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return artifactHandler
}

func (u *UpdateDependenciesOptions) MustGetArtifactHandlers() (artifactHandlers []ArtifactHandler) {
	artifactHandlers, err := u.GetArtifactHandlers()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return artifactHandlers
}

func (u *UpdateDependenciesOptions) MustGetAuthenticationOptions() (authenticationOptions []AuthenticationOption) {
	authenticationOptions, err := u.GetAuthenticationOptions()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return authenticationOptions
}

func (u *UpdateDependenciesOptions) MustGetCommit() (commit bool) {
	commit, err := u.GetCommit()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return commit
}

func (u *UpdateDependenciesOptions) MustGetLatestArtifactVersionAsString(softwareName string, verbose bool) (latestVersion string) {
	latestVersion, err := u.GetLatestArtifactVersionAsString(softwareName, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return latestVersion
}

func (u *UpdateDependenciesOptions) MustGetVerbose() (verbose bool) {
	verbose, err := u.GetVerbose()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return verbose
}

func (u *UpdateDependenciesOptions) MustSetArtifactHandlers(artifactHandlers []ArtifactHandler) {
	err := u.SetArtifactHandlers(artifactHandlers)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (u *UpdateDependenciesOptions) MustSetAuthenticationOptions(authenticationOptions []AuthenticationOption) {
	err := u.SetAuthenticationOptions(authenticationOptions)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (u *UpdateDependenciesOptions) MustSetCommit(commit bool) {
	err := u.SetCommit(commit)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (u *UpdateDependenciesOptions) MustSetVerbose(verbose bool) {
	err := u.SetVerbose(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (u *UpdateDependenciesOptions) SetArtifactHandlers(artifactHandlers []ArtifactHandler) (err error) {
	if artifactHandlers == nil {
		return TracedErrorf("artifactHandlers is nil")
	}

	if len(artifactHandlers) <= 0 {
		return TracedErrorf("artifactHandlers has no elements")
	}

	u.ArtifactHandlers = artifactHandlers

	return nil
}

func (u *UpdateDependenciesOptions) SetAuthenticationOptions(authenticationOptions []AuthenticationOption) (err error) {
	if authenticationOptions == nil {
		return TracedErrorf("authenticationOptions is nil")
	}

	if len(authenticationOptions) <= 0 {
		return TracedErrorf("authenticationOptions has no elements")
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
