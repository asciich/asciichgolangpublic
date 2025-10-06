package asciichgolangpublic

import (
	"context"
	"fmt"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	gitlab "gitlab.com/gitlab-org/api/client-go"
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
		return nil, tracederrors.TracedErrorf("gitlabPersonalAccessTokens not set")
	}

	return g.gitlabPersonalAccessTokens, nil
}

func (g *GitlabPersonalAccessToken) SetGitlabPersonalAccessTokens(gitlabPersonalAccessTokens *GitlabPersonalAccessTokenService) (err error) {
	if gitlabPersonalAccessTokens == nil {
		return tracederrors.TracedErrorf("gitlabPersonalAccessTokens is nil")
	}

	g.gitlabPersonalAccessTokens = gitlabPersonalAccessTokens

	return nil
}

func (t *GitlabPersonalAccessToken) GetCachedName() (cachedName string, err error) {
	if len(t.cachedName) <= 0 {
		return "", tracederrors.TracedError("cachedName not implemented")
	}

	return t.cachedName, nil
}

func (t *GitlabPersonalAccessToken) GetId() (id int, err error) {
	if t.id <= 0 {
		return -1, tracederrors.TracedError("id not set")
	}

	return t.id, nil
}

func (t *GitlabPersonalAccessToken) GetInfoString(ctx context.Context) (infoString string, err error) {
	rawResponse, err := t.GetTokenRawResponse(ctx)
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
		return nil, tracederrors.TracedError("gitlabPersonalAccessTokens is not set")
	}

	return t.gitlabPersonalAccessTokens, nil
}

func (t *GitlabPersonalAccessToken) GetTokenRawResponse(ctx context.Context) (nativeResponse *gitlab.PersonalAccessToken, err error) {
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
		return nil, tracederrors.TracedError(err.Error())
	}

	if nativeResponse == nil {
		return nil, tracederrors.TracedError("nativeResponse is nil")
	}

	logging.LogInfoByCtxf(ctx, "Collected personal access token id='%d' raw response.", id)

	return nativeResponse, nil
}

func (t *GitlabPersonalAccessToken) SetCachedName(cachedName string) (err error) {
	if len(cachedName) <= 0 {
		return tracederrors.TracedError("cachedName is empty string")
	}

	t.cachedName = cachedName

	return nil
}

func (t *GitlabPersonalAccessToken) SetId(id int) (err error) {
	if id <= 0 {
		return tracederrors.TracedErrorf("invalid id '%d'", id)
	}

	t.id = id

	return nil
}

func (t *GitlabPersonalAccessToken) SetPersonalAccessTokens(tokensService *GitlabPersonalAccessTokenService) (err error) {
	if tokensService == nil {
		return tracederrors.TracedError("tokenService is nil")
	}

	t.gitlabPersonalAccessTokens = tokensService

	return nil
}
