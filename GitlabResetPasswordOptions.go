package asciichgolangpublic

type GitlabResetPasswordOptions struct {
	Username                        string
	GitlabContainerNameOnGitlabHost string
	GopassPathToStoreNewPassword    string
	SshUserNameForGitlabHost        string
	Verbose                         bool
}

func NewGitlabResetPasswordOptions() (g *GitlabResetPasswordOptions) {
	return new(GitlabResetPasswordOptions)
}

func (g *GitlabResetPasswordOptions) GetGitlabContainerNameOnGitlabHost() (gitlabContainerNameOnGitlabHost string, err error) {
	if g.GitlabContainerNameOnGitlabHost == "" {
		return "", TracedErrorf("GitlabContainerNameOnGitlabHost not set")
	}

	return g.GitlabContainerNameOnGitlabHost, nil
}

func (g *GitlabResetPasswordOptions) GetGopassPathToStoreNewPassword() (gopassPathToStoreNewPassword string, err error) {
	if g.GopassPathToStoreNewPassword == "" {
		return "", TracedErrorf("GopassPathToStoreNewPassword not set")
	}

	return g.GopassPathToStoreNewPassword, nil
}

func (g *GitlabResetPasswordOptions) GetVerbose() (verbose bool, err error) {

	return g.Verbose, nil
}

func (g *GitlabResetPasswordOptions) MustGetGitlabContainerNameOnGitlabHost() (gitlabContainerNameOnGitlabHost string) {
	gitlabContainerNameOnGitlabHost, err := g.GetGitlabContainerNameOnGitlabHost()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlabContainerNameOnGitlabHost
}

func (g *GitlabResetPasswordOptions) MustGetGopassPathToStoreNewPassword() (gopassPathToStoreNewPassword string) {
	gopassPathToStoreNewPassword, err := g.GetGopassPathToStoreNewPassword()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gopassPathToStoreNewPassword
}

func (g *GitlabResetPasswordOptions) MustGetSshUserNameForGitlabHost() (sshUserName string) {
	sshUserName, err := g.GetSshUserNameForGitlabHost()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return sshUserName
}

func (g *GitlabResetPasswordOptions) MustGetUsername() (username string) {
	username, err := g.GetUsername()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return username
}

func (g *GitlabResetPasswordOptions) MustGetVerbose() (verbose bool) {
	verbose, err := g.GetVerbose()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return verbose
}

func (g *GitlabResetPasswordOptions) MustSetGitlabContainerNameOnGitlabHost(gitlabContainerNameOnGitlabHost string) {
	err := g.SetGitlabContainerNameOnGitlabHost(gitlabContainerNameOnGitlabHost)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabResetPasswordOptions) MustSetGopassPathToStoreNewPassword(gopassPathToStoreNewPassword string) {
	err := g.SetGopassPathToStoreNewPassword(gopassPathToStoreNewPassword)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabResetPasswordOptions) MustSetSshUserNameForGitlabHost(sshUserNameForGitlabHost string) {
	err := g.SetSshUserNameForGitlabHost(sshUserNameForGitlabHost)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabResetPasswordOptions) MustSetUsername(username string) {
	err := g.SetUsername(username)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabResetPasswordOptions) MustSetVerbose(verbose bool) {
	err := g.SetVerbose(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabResetPasswordOptions) SetGitlabContainerNameOnGitlabHost(gitlabContainerNameOnGitlabHost string) (err error) {
	if gitlabContainerNameOnGitlabHost == "" {
		return TracedErrorf("gitlabContainerNameOnGitlabHost is empty string")
	}

	g.GitlabContainerNameOnGitlabHost = gitlabContainerNameOnGitlabHost

	return nil
}

func (g *GitlabResetPasswordOptions) SetGopassPathToStoreNewPassword(gopassPathToStoreNewPassword string) (err error) {
	if gopassPathToStoreNewPassword == "" {
		return TracedErrorf("gopassPathToStoreNewPassword is empty string")
	}

	g.GopassPathToStoreNewPassword = gopassPathToStoreNewPassword

	return nil
}

func (g *GitlabResetPasswordOptions) SetSshUserNameForGitlabHost(sshUserNameForGitlabHost string) (err error) {
	if sshUserNameForGitlabHost == "" {
		return TracedErrorf("sshUserNameForGitlabHost is empty string")
	}

	g.SshUserNameForGitlabHost = sshUserNameForGitlabHost

	return nil
}

func (g *GitlabResetPasswordOptions) SetUsername(username string) (err error) {
	if username == "" {
		return TracedErrorf("username is empty string")
	}

	g.Username = username

	return nil
}

func (g *GitlabResetPasswordOptions) SetVerbose(verbose bool) (err error) {
	g.Verbose = verbose

	return nil
}

func (o *GitlabResetPasswordOptions) GetSshUserNameForGitlabHost() (sshUserName string, err error) {
	if !o.IsSshUserNameForGitlabHostSet() {
		return "", TracedError("SshUserNameForGitlabHost is not set")
	}

	return o.SshUserNameForGitlabHost, nil
}

func (o *GitlabResetPasswordOptions) GetUsername() (username string, err error) {
	if len(o.Username) <= 0 {
		return "", TracedError("username not set")
	}

	return o.Username, nil
}

func (o *GitlabResetPasswordOptions) IsSshUserNameForGitlabHostSet() (isSet bool) {
	return len(o.SshUserNameForGitlabHost) > 0
}
