package asciichgolangpublic

import (
	"fmt"

	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type GitlabResetAccessTokenOptions struct {
	Username                        string
	GopassPathToStoreNewToken       string
	GitlabContainerNameOnGitlabHost string
	SshUserNameForGitlabHost        string
}

func NewGitlabResetAccessTokenOptions() (g *GitlabResetAccessTokenOptions) {
	return new(GitlabResetAccessTokenOptions)
}

func (g *GitlabResetAccessTokenOptions) GetGitlabContainerNameOnGitlabHost() (gitlabContainerNameOnGitlabHost string, err error) {
	if g.GitlabContainerNameOnGitlabHost == "" {
		return "", tracederrors.TracedErrorf("GitlabContainerNameOnGitlabHost not set")
	}

	return g.GitlabContainerNameOnGitlabHost, nil
}

func (g *GitlabResetAccessTokenOptions) GetGopassPathToStoreNewToken() (gopassPathToStoreNewToken string, err error) {
	if g.GopassPathToStoreNewToken == "" {
		return "", tracederrors.TracedErrorf("GopassPathToStoreNewToken not set")
	}

	return g.GopassPathToStoreNewToken, nil
}

func (g *GitlabResetAccessTokenOptions) GetSshUserNameForGitlabHost() (sshUserNameForGitlabHost string, err error) {
	if g.SshUserNameForGitlabHost == "" {
		return "", tracederrors.TracedErrorf("SshUserNameForGitlabHost not set")
	}

	return g.SshUserNameForGitlabHost, nil
}

func (g *GitlabResetAccessTokenOptions) SetGitlabContainerNameOnGitlabHost(gitlabContainerNameOnGitlabHost string) (err error) {
	if gitlabContainerNameOnGitlabHost == "" {
		return tracederrors.TracedErrorf("gitlabContainerNameOnGitlabHost is empty string")
	}

	g.GitlabContainerNameOnGitlabHost = gitlabContainerNameOnGitlabHost

	return nil
}

func (g *GitlabResetAccessTokenOptions) SetGopassPathToStoreNewToken(gopassPathToStoreNewToken string) (err error) {
	if gopassPathToStoreNewToken == "" {
		return tracederrors.TracedErrorf("gopassPathToStoreNewToken is empty string")
	}

	g.GopassPathToStoreNewToken = gopassPathToStoreNewToken

	return nil
}

func (g *GitlabResetAccessTokenOptions) SetSshUserNameForGitlabHost(sshUserNameForGitlabHost string) (err error) {
	if sshUserNameForGitlabHost == "" {
		return tracederrors.TracedErrorf("sshUserNameForGitlabHost is empty string")
	}

	g.SshUserNameForGitlabHost = sshUserNameForGitlabHost

	return nil
}

func (g *GitlabResetAccessTokenOptions) SetUsername(username string) (err error) {
	if username == "" {
		return tracederrors.TracedErrorf("username is empty string")
	}

	g.Username = username

	return nil
}

func (o *GitlabResetAccessTokenOptions) GetUsername() (username string, err error) {
	if len(o.Username) <= 0 {
		return "", fmt.Errorf("username not set")
	}

	return o.Username, nil
}
