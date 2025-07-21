package asciichgolangpublic

import (
	"fmt"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type GitlabResetAccessTokenOptions struct {
	Username                        string
	GopassPathToStoreNewToken       string
	GitlabContainerNameOnGitlabHost string
	SshUserNameForGitlabHost        string
	Verbose                         bool
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

func (g *GitlabResetAccessTokenOptions) GetVerbose() (verbose bool, err error) {

	return g.Verbose, nil
}

func (g *GitlabResetAccessTokenOptions) MustGetGitlabContainerNameOnGitlabHost() (gitlabContainerNameOnGitlabHost string) {
	gitlabContainerNameOnGitlabHost, err := g.GetGitlabContainerNameOnGitlabHost()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlabContainerNameOnGitlabHost
}

func (g *GitlabResetAccessTokenOptions) MustGetGopassPathToStoreNewToken() (gopassPathToStoreNewToken string) {
	gopassPathToStoreNewToken, err := g.GetGopassPathToStoreNewToken()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gopassPathToStoreNewToken
}

func (g *GitlabResetAccessTokenOptions) MustGetSshUserNameForGitlabHost() (sshUserNameForGitlabHost string) {
	sshUserNameForGitlabHost, err := g.GetSshUserNameForGitlabHost()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return sshUserNameForGitlabHost
}

func (g *GitlabResetAccessTokenOptions) MustGetUsername() (username string) {
	username, err := g.GetUsername()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return username
}

func (g *GitlabResetAccessTokenOptions) MustGetVerbose() (verbose bool) {
	verbose, err := g.GetVerbose()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return verbose
}

func (g *GitlabResetAccessTokenOptions) MustSetGitlabContainerNameOnGitlabHost(gitlabContainerNameOnGitlabHost string) {
	err := g.SetGitlabContainerNameOnGitlabHost(gitlabContainerNameOnGitlabHost)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabResetAccessTokenOptions) MustSetGopassPathToStoreNewToken(gopassPathToStoreNewToken string) {
	err := g.SetGopassPathToStoreNewToken(gopassPathToStoreNewToken)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabResetAccessTokenOptions) MustSetSshUserNameForGitlabHost(sshUserNameForGitlabHost string) {
	err := g.SetSshUserNameForGitlabHost(sshUserNameForGitlabHost)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabResetAccessTokenOptions) MustSetUsername(username string) {
	err := g.SetUsername(username)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabResetAccessTokenOptions) MustSetVerbose(verbose bool) {
	err := g.SetVerbose(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
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

func (g *GitlabResetAccessTokenOptions) SetVerbose(verbose bool) (err error) {
	g.Verbose = verbose

	return nil
}

func (o *GitlabResetAccessTokenOptions) GetUsername() (username string, err error) {
	if len(o.Username) <= 0 {
		return "", fmt.Errorf("username not set")
	}

	return o.Username, nil
}
