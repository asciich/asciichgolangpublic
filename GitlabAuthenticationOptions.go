package asciichgolangpublic

import (
	"strings"
)

type GitlabAuthenticationOptions struct {
	AccessToken            string
	AccessTokensFromGopass []string
	Verbose                bool
	GitlabUrl              string
}

func NewGitlabAuthenticationOptions() (g *GitlabAuthenticationOptions) {
	return new(GitlabAuthenticationOptions)
}

func (g *GitlabAuthenticationOptions) GetAccessToken() (accessToken string, err error) {
	if g.AccessToken == "" {
		return "", TracedErrorf("AccessToken not set")
	}

	return g.AccessToken, nil
}

func (g *GitlabAuthenticationOptions) GetAccessTokensFromGopass() (accessTokensFromGopass []string, err error) {
	if g.AccessTokensFromGopass == nil {
		return nil, TracedErrorf("AccessTokensFromGopass not set")
	}

	if len(g.AccessTokensFromGopass) <= 0 {
		return nil, TracedErrorf("AccessTokensFromGopass has no elements")
	}

	return g.AccessTokensFromGopass, nil
}

func (g *GitlabAuthenticationOptions) GetGitlabUrl() (gitlabUrl string, err error) {
	if g.GitlabUrl == "" {
		return "", TracedErrorf("GitlabUrl not set")
	}

	return g.GitlabUrl, nil
}

func (g *GitlabAuthenticationOptions) GetVerbose() (verbose bool, err error) {

	return g.Verbose, nil
}

func (g *GitlabAuthenticationOptions) IsAccessTokenSet() (isSet bool) {
	return g.AccessToken != ""
}

func (g *GitlabAuthenticationOptions) IsAuthenticatingAgainst(serviceName string) (isAuthenticatingAgainst bool, err error) {
	gitlabUrl, err := g.GetGitlabUrl()
	if err != nil {
		return false, err
	}

	isAuthenticatingAgainst = strings.HasPrefix(serviceName, gitlabUrl)

	return isAuthenticatingAgainst, nil
}

func (g *GitlabAuthenticationOptions) IsVerbose() (isVerbose bool) {
	return g.Verbose
}

func (g *GitlabAuthenticationOptions) MustGetAccessToken() (accessToken string) {
	accessToken, err := g.GetAccessToken()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return accessToken
}

func (g *GitlabAuthenticationOptions) MustGetAccessTokensFromGopass() (accessTokensFromGopass []string) {
	accessTokensFromGopass, err := g.GetAccessTokensFromGopass()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return accessTokensFromGopass
}

func (g *GitlabAuthenticationOptions) MustGetGitlabUrl() (gitlabUrl string) {
	gitlabUrl, err := g.GetGitlabUrl()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlabUrl
}

func (g *GitlabAuthenticationOptions) MustGetVerbose() (verbose bool) {
	verbose, err := g.GetVerbose()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return verbose
}

func (g *GitlabAuthenticationOptions) MustIsAuthenticatingAgainst(serviceName string) (isAuthenticatingAgainst bool) {
	isAuthenticatingAgainst, err := g.IsAuthenticatingAgainst(serviceName)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return isAuthenticatingAgainst
}

func (g *GitlabAuthenticationOptions) MustSetAccessToken(accessToken string) {
	err := g.SetAccessToken(accessToken)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabAuthenticationOptions) MustSetAccessTokensFromGopass(accessTokensFromGopass []string) {
	err := g.SetAccessTokensFromGopass(accessTokensFromGopass)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabAuthenticationOptions) MustSetGitlabUrl(gitlabUrl string) {
	err := g.SetGitlabUrl(gitlabUrl)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabAuthenticationOptions) MustSetVerbose(verbose bool) {
	err := g.SetVerbose(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabAuthenticationOptions) SetAccessToken(accessToken string) (err error) {
	if accessToken == "" {
		return TracedErrorf("accessToken is empty string")
	}

	g.AccessToken = accessToken

	return nil
}

func (g *GitlabAuthenticationOptions) SetAccessTokensFromGopass(accessTokensFromGopass []string) (err error) {
	if accessTokensFromGopass == nil {
		return TracedErrorf("accessTokensFromGopass is nil")
	}

	if len(accessTokensFromGopass) <= 0 {
		return TracedErrorf("accessTokensFromGopass has no elements")
	}

	g.AccessTokensFromGopass = accessTokensFromGopass

	return nil
}

func (g *GitlabAuthenticationOptions) SetGitlabUrl(gitlabUrl string) (err error) {
	if gitlabUrl == "" {
		return TracedErrorf("gitlabUrl is empty string")
	}

	g.GitlabUrl = gitlabUrl

	return nil
}

func (g *GitlabAuthenticationOptions) SetVerbose(verbose bool) (err error) {
	g.Verbose = verbose

	return nil
}
