package asciichgolangpublic

import (
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type GitlabAuthenticationOptions struct {
	AccessToken            string
	AccessTokensFromGopass []string
	GitlabUrl              string
}

func NewGitlabAuthenticationOptions() (g *GitlabAuthenticationOptions) {
	return new(GitlabAuthenticationOptions)
}

func (g *GitlabAuthenticationOptions) GetAccessToken() (accessToken string, err error) {
	if g.AccessToken == "" {
		return "", tracederrors.TracedErrorf("AccessToken not set")
	}

	return g.AccessToken, nil
}

func (g *GitlabAuthenticationOptions) GetAccessTokensFromGopass() (accessTokensFromGopass []string, err error) {
	if g.AccessTokensFromGopass == nil {
		return nil, tracederrors.TracedErrorf("AccessTokensFromGopass not set")
	}

	if len(g.AccessTokensFromGopass) <= 0 {
		return nil, tracederrors.TracedErrorf("AccessTokensFromGopass has no elements")
	}

	return g.AccessTokensFromGopass, nil
}

func (g *GitlabAuthenticationOptions) GetGitlabUrl() (gitlabUrl string, err error) {
	if g.GitlabUrl == "" {
		return "", tracederrors.TracedErrorf("GitlabUrl not set")
	}

	return g.GitlabUrl, nil
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

func (g *GitlabAuthenticationOptions) SetAccessToken(accessToken string) (err error) {
	if accessToken == "" {
		return tracederrors.TracedErrorf("accessToken is empty string")
	}

	g.AccessToken = accessToken

	return nil
}

func (g *GitlabAuthenticationOptions) SetAccessTokensFromGopass(accessTokensFromGopass []string) (err error) {
	if accessTokensFromGopass == nil {
		return tracederrors.TracedErrorf("accessTokensFromGopass is nil")
	}

	if len(accessTokensFromGopass) <= 0 {
		return tracederrors.TracedErrorf("accessTokensFromGopass has no elements")
	}

	g.AccessTokensFromGopass = accessTokensFromGopass

	return nil
}

func (g *GitlabAuthenticationOptions) SetGitlabUrl(gitlabUrl string) (err error) {
	if gitlabUrl == "" {
		return tracederrors.TracedErrorf("gitlabUrl is empty string")
	}

	g.GitlabUrl = gitlabUrl

	return nil
}

