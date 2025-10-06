package asciichgolangpublic

import (
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type GitlabResetPasswordOptions struct {
	Username                        string
	GitlabContainerNameOnGitlabHost string
	GopassPathToStoreNewPassword    string
	SshUserNameForGitlabHost        string
}

func NewGitlabResetPasswordOptions() (g *GitlabResetPasswordOptions) {
	return new(GitlabResetPasswordOptions)
}

func (g *GitlabResetPasswordOptions) GetGitlabContainerNameOnGitlabHost() (gitlabContainerNameOnGitlabHost string, err error) {
	if g.GitlabContainerNameOnGitlabHost == "" {
		return "", tracederrors.TracedErrorf("GitlabContainerNameOnGitlabHost not set")
	}

	return g.GitlabContainerNameOnGitlabHost, nil
}

func (g *GitlabResetPasswordOptions) GetGopassPathToStoreNewPassword() (gopassPathToStoreNewPassword string, err error) {
	if g.GopassPathToStoreNewPassword == "" {
		return "", tracederrors.TracedErrorf("GopassPathToStoreNewPassword not set")
	}

	return g.GopassPathToStoreNewPassword, nil
}

func (g *GitlabResetPasswordOptions) SetGitlabContainerNameOnGitlabHost(gitlabContainerNameOnGitlabHost string) (err error) {
	if gitlabContainerNameOnGitlabHost == "" {
		return tracederrors.TracedErrorf("gitlabContainerNameOnGitlabHost is empty string")
	}

	g.GitlabContainerNameOnGitlabHost = gitlabContainerNameOnGitlabHost

	return nil
}

func (g *GitlabResetPasswordOptions) SetGopassPathToStoreNewPassword(gopassPathToStoreNewPassword string) (err error) {
	if gopassPathToStoreNewPassword == "" {
		return tracederrors.TracedErrorf("gopassPathToStoreNewPassword is empty string")
	}

	g.GopassPathToStoreNewPassword = gopassPathToStoreNewPassword

	return nil
}

func (g *GitlabResetPasswordOptions) SetSshUserNameForGitlabHost(sshUserNameForGitlabHost string) (err error) {
	if sshUserNameForGitlabHost == "" {
		return tracederrors.TracedErrorf("sshUserNameForGitlabHost is empty string")
	}

	g.SshUserNameForGitlabHost = sshUserNameForGitlabHost

	return nil
}

func (g *GitlabResetPasswordOptions) SetUsername(username string) (err error) {
	if username == "" {
		return tracederrors.TracedErrorf("username is empty string")
	}

	g.Username = username

	return nil
}

func (o *GitlabResetPasswordOptions) GetSshUserNameForGitlabHost() (sshUserName string, err error) {
	if !o.IsSshUserNameForGitlabHostSet() {
		return "", tracederrors.TracedError("SshUserNameForGitlabHost is not set")
	}

	return o.SshUserNameForGitlabHost, nil
}

func (o *GitlabResetPasswordOptions) GetUsername() (username string, err error) {
	if len(o.Username) <= 0 {
		return "", tracederrors.TracedError("username not set")
	}

	return o.Username, nil
}

func (o *GitlabResetPasswordOptions) IsSshUserNameForGitlabHostSet() (isSet bool) {
	return len(o.SshUserNameForGitlabHost) > 0
}
