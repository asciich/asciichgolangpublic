package asciichgolangpublic

import (
	"fmt"

	"github.com/xanzy/go-gitlab"
)

type GitlabPersonalAccessToken struct {
	gitlabPersonalAccessTokens *GitlabPersonalAccessTokenService
	id                         int
	cachedName                 string
}

func NewGitlabPersonalAccessToken() (accessToken *GitlabPersonalAccessToken) {
	return new(GitlabPersonalAccessToken)
}

func (g *GitlabPersonalAccessToken) GetGitlabPersonalAccessTokens() (gitlabPersonalAccessTokens *GitlabPersonalAccessTokenService, err error) {
	if g.gitlabPersonalAccessTokens == nil {
		return nil, TracedErrorf("gitlabPersonalAccessTokens not set")
	}

	return g.gitlabPersonalAccessTokens, nil
}

func (g *GitlabPersonalAccessToken) MustGetCachedName() (cachedName string) {
	cachedName, err := g.GetCachedName()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return cachedName
}

func (g *GitlabPersonalAccessToken) MustGetGitlabPersonalAccessTokens() (gitlabPersonalAccessTokens *GitlabPersonalAccessTokenService) {
	gitlabPersonalAccessTokens, err := g.GetGitlabPersonalAccessTokens()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlabPersonalAccessTokens
}

func (g *GitlabPersonalAccessToken) MustGetId() (id int) {
	id, err := g.GetId()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return id
}

func (g *GitlabPersonalAccessToken) MustGetInfoString(verbose bool) (infoString string) {
	infoString, err := g.GetInfoString(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return infoString
}

func (g *GitlabPersonalAccessToken) MustGetNativePersonalTokenService() (nativeService *gitlab.PersonalAccessTokensService) {
	nativeService, err := g.GetNativePersonalTokenService()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return nativeService
}

func (g *GitlabPersonalAccessToken) MustGetPersonalAccessTokens() (tokensService *GitlabPersonalAccessTokenService) {
	tokensService, err := g.GetPersonalAccessTokens()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return tokensService
}

func (g *GitlabPersonalAccessToken) MustGetTokenRawResponse(verbose bool) (nativeResponse *gitlab.PersonalAccessToken) {
	nativeResponse, err := g.GetTokenRawResponse(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return nativeResponse
}

func (g *GitlabPersonalAccessToken) MustSetCachedName(cachedName string) {
	err := g.SetCachedName(cachedName)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabPersonalAccessToken) MustSetGitlabPersonalAccessTokens(gitlabPersonalAccessTokens *GitlabPersonalAccessTokenService) {
	err := g.SetGitlabPersonalAccessTokens(gitlabPersonalAccessTokens)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabPersonalAccessToken) MustSetId(id int) {
	err := g.SetId(id)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabPersonalAccessToken) MustSetPersonalAccessTokens(tokensService *GitlabPersonalAccessTokenService) {
	err := g.SetPersonalAccessTokens(tokensService)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabPersonalAccessToken) SetGitlabPersonalAccessTokens(gitlabPersonalAccessTokens *GitlabPersonalAccessTokenService) (err error) {
	if gitlabPersonalAccessTokens == nil {
		return TracedErrorf("gitlabPersonalAccessTokens is nil")
	}

	g.gitlabPersonalAccessTokens = gitlabPersonalAccessTokens

	return nil
}

func (t *GitlabPersonalAccessToken) GetCachedName() (cachedName string, err error) {
	if len(t.cachedName) <= 0 {
		return "", TracedError("cachedName not implemented")
	}

	return t.cachedName, nil
}

func (t *GitlabPersonalAccessToken) GetId() (id int, err error) {
	if t.id <= 0 {
		return -1, TracedError("id not set")
	}

	return t.id, nil
}

func (t *GitlabPersonalAccessToken) GetInfoString(verbose bool) (infoString string, err error) {
	rawResponse, err := t.GetTokenRawResponse(verbose)
	if err != nil {
		return "", err
	}

	infoString += fmt.Sprintf("id=%d", rawResponse.ID)
	infoString += ", name=" + rawResponse.Name
	infoString += fmt.Sprintf(", revoked=%v", rawResponse.Revoked)
	infoString += fmt.Sprintf(", scopes=%v", rawResponse.Scopes)

	return infoString, nil
}

func (t *GitlabPersonalAccessToken) GetNativePersonalTokenService() (nativeService *gitlab.PersonalAccessTokensService, err error) {
	tokens, err := t.GetPersonalAccessTokens()
	if err != nil {
		return nil, err
	}

	nativeService, err = tokens.GetNativePersonalTokenService()
	if err != nil {
		return nil, err
	}

	return nativeService, nil
}

func (t *GitlabPersonalAccessToken) GetPersonalAccessTokens() (tokensService *GitlabPersonalAccessTokenService, err error) {
	if t.gitlabPersonalAccessTokens == nil {
		return nil, TracedError("gitlabPersonalAccessTokens is not set")
	}

	return t.gitlabPersonalAccessTokens, nil
}

func (t *GitlabPersonalAccessToken) GetTokenRawResponse(verbose bool) (nativeResponse *gitlab.PersonalAccessToken, err error) {
	nativeService, err := t.GetNativePersonalTokenService()
	if err != nil {
		return nil, err
	}

	id, err := t.GetId()
	if err != nil {
		return nil, err
	}

	nativeResponse, _, err = nativeService.GetSinglePersonalAccessTokenByID(id)
	if err != nil {
		return nil, TracedError(err.Error())
	}

	if nativeResponse == nil {
		return nil, TracedError("nativeResponse is nil")
	}

	if verbose {
		LogInfof("Collected personal access token id='%d' raw response.", id)
	}

	return nativeResponse, nil
}

func (t *GitlabPersonalAccessToken) SetCachedName(cachedName string) (err error) {
	if len(cachedName) <= 0 {
		return TracedError("cachedName is empty string")
	}

	t.cachedName = cachedName

	return nil
}

func (t *GitlabPersonalAccessToken) SetId(id int) (err error) {
	if id <= 0 {
		return TracedErrorf("invalid id '%d'", id)
	}

	t.id = id

	return nil
}

func (t *GitlabPersonalAccessToken) SetPersonalAccessTokens(tokensService *GitlabPersonalAccessTokenService) (err error) {
	if tokensService == nil {
		return TracedError("tokenService is nil")
	}

	t.gitlabPersonalAccessTokens = tokensService

	return nil
}
