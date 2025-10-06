package asciichgolangpublic

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorbashoo"
	"github.com/asciich/asciichgolangpublic/pkg/datatypes/stringsutils"
	"github.com/asciich/asciichgolangpublic/pkg/fileformats/jsonutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type GitlabPersonalAccessTokenService struct {
	gitlab *GitlabInstance
}

func NewGitlabPersonalAccessTokenService() (tokens *GitlabPersonalAccessTokenService) {
	return new(GitlabPersonalAccessTokenService)
}

func (p *GitlabPersonalAccessTokenService) CreateToken(ctx context.Context, tokenOptions *GitlabCreatePersonalAccessTokenOptions) (newToken string, err error) {
	if tokenOptions == nil {
		return "", tracederrors.TracedError("tokenOptions is nil")
	}

	tokenName, err := tokenOptions.GetName()
	if err != nil {
		return "", err
	}

	exists, err := p.ExistsByName(ctx, tokenName)
	if err != nil {
		return "", err
	}

	if exists {
		return "", tracederrors.TracedErrorf(
			"Unable to create token '%s'. Token '%s' already exists. Use RecreateToken to overwrite existing tokens",
			tokenName,
			tokenName,
		)
	}

	// Official API to get a personal token as an User is not implemented yet in go-gitlab:
	// https://docs.gitlab.com/ee/api/users.html#create-a-personal-access-token-with-limited-scopes-for-the-currently-authenticated-user

	accessToken, err := p.GetCurrentlyUsedAccessToken()
	if err != nil {
		return "", err
	}

	apiV4Urt, err := p.GetApiV4Url()
	if err != nil {
		return "", err
	}

	command := []string{
		"curl",
		"-L",
		"--fail",
		"-s",
		"--request",
		"POST",
		"--header",
		fmt.Sprintf("PRIVATE-TOKEN: %s", accessToken),
		"--data",
		"name=" + tokenName,
		"--data",
		"scopes[]=k8s_proxy",
		apiV4Urt + "/user/personal_access_tokens",
	}
	stdout, err := commandexecutorbashoo.Bash().RunCommandAndGetStdoutAsString(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: command,
		},
	)
	if err != nil {
		return "", err
	}

	newToken, err = jsonutils.RunJqAgainstJsonStringAsString(stdout, ".token")
	if err != nil {
		return "", err
	}

	newToken = strings.TrimSpace(newToken)
	newToken = stringsutils.RemoveSurroundingQuotationMarks(newToken)
	if len(newToken) <= 0 {
		return "", tracederrors.TracedError("Unable to get newToken. newToken is empty string.")
	}

	return newToken, nil
}

func (p *GitlabPersonalAccessTokenService) ExistsByName(ctx context.Context, tokenName string) (exists bool, err error) {
	if len(tokenName) <= 0 {
		return false, tracederrors.TracedError("tokenName is nil")
	}

	tokenNames, err := p.GetPersonalAccessTokenNameList(ctx)
	if err != nil {
		return false, err
	}

	exists = slices.Contains(tokenNames, tokenName)

	return exists, nil
}

func (p *GitlabPersonalAccessTokenService) GetApiV4Url() (apiV4Url string, err error) {
	gitlab, err := p.GetGitlab()
	if err != nil {
		return "", err
	}

	apiV4Url, err = gitlab.GetApiV4Url()
	if err != nil {
		return "", err
	}

	return apiV4Url, nil
}

func (p *GitlabPersonalAccessTokenService) GetCurrentUserId(verbose bool) (userId int, err error) {
	gitlab, err := p.GetGitlab()
	if err != nil {
		return -1, err
	}

	userId, err = gitlab.GetUserId()
	if err != nil {
		return -1, err
	}

	if verbose {
		logging.LogInfof("Current gitlab user id is '%d'", userId)
	}

	return userId, nil
}

func (p *GitlabPersonalAccessTokenService) GetCurrentlyUsedAccessToken() (accessToken string, err error) {
	gitlab, err := p.GetGitlab()
	if err != nil {
		return "", err
	}

	accessToken, err = gitlab.GetCurrentlyUsedAccessToken()
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (p *GitlabPersonalAccessTokenService) GetGitlab() (gitlab *GitlabInstance, err error) {
	if p.gitlab == nil {
		return nil, tracederrors.TracedError("gitlab not set")
	}

	return p.gitlab, nil
}

func (p *GitlabPersonalAccessTokenService) GetGitlabUsers() (gitlabUsers *GitlabUsers, err error) {
	gitlab, err := p.GetGitlab()
	if err != nil {
		return nil, err
	}

	gitlabUsers, err = gitlab.GetGitlabUsers()
	if err != nil {
		return nil, err
	}

	return gitlabUsers, nil
}

func (p *GitlabPersonalAccessTokenService) GetNativeGitlabClient() (nativeClient *gitlab.Client, err error) {
	gitlab, err := p.GetGitlab()
	if err != nil {
		return nil, err
	}

	nativeClient, err = gitlab.GetNativeClient()
	if err != nil {
		return nil, err
	}

	return nativeClient, nil
}

func (p *GitlabPersonalAccessTokenService) GetNativePersonalTokenService() (nativeService *gitlab.PersonalAccessTokensService, err error) {
	nativeClient, err := p.GetNativeGitlabClient()
	if err != nil {
		return nil, err
	}

	nativeService = nativeClient.PersonalAccessTokens
	if nativeService == nil {
		return nil, tracederrors.TracedError("nativeService is nil")
	}

	return nativeService, nil
}

func (p *GitlabPersonalAccessTokenService) GetNativeUsersService() (nativeService *gitlab.UsersService, err error) {
	users, err := p.GetGitlabUsers()
	if err != nil {
		return nil, err
	}

	nativeService, err = users.GetNativeUsersService()
	if err != nil {
		return nil, err
	}

	return nativeService, nil
}

func (p *GitlabPersonalAccessTokenService) GetPersonalAccessTokenList(ctx context.Context) (tokens []*GitlabPersonalAccessToken, err error) {
	nativeService, err := p.GetNativePersonalTokenService()
	if err != nil {
		return nil, err
	}

	list, _, err := nativeService.ListPersonalAccessTokens(&gitlab.ListPersonalAccessTokensOptions{})
	if err != nil {
		return nil, err
	}

	tokens = []*GitlabPersonalAccessToken{}
	for _, l := range list {
		tokenToAdd := NewGitlabPersonalAccessToken()

		err = tokenToAdd.SetCachedName(l.Name)
		if err != nil {
			return nil, err
		}

		err = tokenToAdd.SetId(l.ID)
		if err != nil {
			return nil, err
		}

		err = tokenToAdd.SetPersonalAccessTokens(p)
		if err != nil {
			return nil, err
		}

		tokens = append(tokens, tokenToAdd)
	}

	return tokens, nil
}

func (p *GitlabPersonalAccessTokenService) GetPersonalAccessTokenNameList(ctx context.Context) (tokenNames []string, err error) {
	tokens, err := p.gitlab.GetPersonalAccessTokenList(ctx)
	if err != nil {
		return nil, err
	}

	tokenNames = []string{}
	for _, t := range tokens {
		nameToAdd, err := t.GetCachedName()
		if err != nil {
			return nil, err
		}

		tokenNames = append(tokenNames, nameToAdd)
	}

	return tokenNames, nil
}

func (p *GitlabPersonalAccessTokenService) GetTokenIdByName(ctx context.Context, tokenName string) (tokenId int, err error) {
	if len(tokenName) <= 0 {
		return -1, tracederrors.TracedError("tokenName is empty string")
	}

	tokens, err := p.GetPersonalAccessTokenList(ctx)
	if err != nil {
		return -1, err
	}

	for _, t := range tokens {
		nameToCheck, err := t.GetCachedName()
		if err != nil {
			return -1, err
		}

		if nameToCheck == tokenName {
			tokenId, err = t.GetId()
			if err != nil {
				return -1, err
			}

			return tokenId, nil
		}
	}

	return -1, tracederrors.TracedErrorf("No token with name '%s' found.", tokenName)
}

func (p *GitlabPersonalAccessTokenService) RecreateToken(ctx context.Context, tokenOptions *GitlabCreatePersonalAccessTokenOptions) (newToken string, err error) {
	if tokenOptions == nil {
		return "", tracederrors.TracedError("tokenOptions is nil")
	}

	tokenName, err := tokenOptions.GetName()
	if err != nil {
		return "", err
	}

	exists, err := p.ExistsByName(ctx, tokenName)
	if err != nil {
		return "", err
	}

	if exists {
		err = p.RevokeTokenByName(ctx, tokenName)
		if err != nil {
			return "", err
		}

		nativeService, err := p.GetNativePersonalTokenService()
		if err != nil {
			return "", err
		}

		tokenId, err := p.GetTokenIdByName(ctx, tokenName)
		if err != nil {
			return "", err
		}

		nativeToken, _, err := nativeService.RotatePersonalAccessToken(tokenId, nil)
		if err != nil {
			return "", err
		}

		newToken = nativeToken.Token
		if len(newToken) <= 0 {
			return "", tracederrors.TracedError("recreate personal access token failed. NewToken is empty string.")
		}

		logging.LogInfoByCtxf(ctx, "Personal access token '%s' recreated.", tokenName)

		return newToken, nil
	} else {
		newToken, err = p.CreateToken(ctx, tokenOptions)
		if err != nil {
			return "", err
		}
		logging.LogInfoByCtxf(ctx, "Personal access token '%s' created.", tokenName)

		return newToken, nil
	}
}

func (p *GitlabPersonalAccessTokenService) RevokeTokenByName(ctx context.Context, tokenName string) (err error) {
	if len(tokenName) <= 0 {
		return tracederrors.TracedError("tokenName is empty string")
	}

	exists, err := p.ExistsByName(ctx, tokenName)
	if err != nil {
		return err
	}

	if exists {
		nativeService, err := p.GetNativePersonalTokenService()
		if err != nil {
			return err
		}

		tokenId, err := p.GetTokenIdByName(ctx, tokenName)
		if err != nil {
			return err
		}

		_, err = nativeService.RevokePersonalAccessToken(tokenId)
		if err != nil {
			return err
		}

		logging.LogChangedByCtxf(ctx, "Personal gitlab access token '%s' deleted.", tokenName)
	} else {
		logging.LogInfoByCtxf(ctx, "Personal gitlab access token '%s' was already deleted.", tokenName)
	}

	return nil
}

func (p *GitlabPersonalAccessTokenService) SetGitlab(gitlab *GitlabInstance) (err error) {
	if gitlab == nil {
		return tracederrors.TracedError("gitlab is nil")
	}

	p.gitlab = gitlab

	return nil
}
