package asciichgolangpublic

import (
	"fmt"
	"strings"

	astrings "github.com/asciich/asciichgolangpublic/datatypes/strings"
	gitlab "gitlab.com/gitlab-org/api/client-go"
	aslices "github.com/asciich/asciichgolangpublic/datatypes/slices"
)

type GitlabPersonalAccessTokenService struct {
	gitlab *GitlabInstance
}

func NewGitlabPersonalAccessTokenService() (tokens *GitlabPersonalAccessTokenService) {
	return new(GitlabPersonalAccessTokenService)
}

func (g *GitlabPersonalAccessTokenService) MustCreateToken(tokenOptions *GitlabCreatePersonalAccessTokenOptions) (newToken string) {
	newToken, err := g.CreateToken(tokenOptions)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return newToken
}

func (g *GitlabPersonalAccessTokenService) MustExistsByName(tokenName string, verbose bool) (exists bool) {
	exists, err := g.ExistsByName(tokenName, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return exists
}

func (g *GitlabPersonalAccessTokenService) MustGetApiV4Url() (apiV4Url string) {
	apiV4Url, err := g.GetApiV4Url()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return apiV4Url
}

func (g *GitlabPersonalAccessTokenService) MustGetCurrentUserId(verbose bool) (userId int) {
	userId, err := g.GetCurrentUserId(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return userId
}

func (g *GitlabPersonalAccessTokenService) MustGetCurrentlyUsedAccessToken() (accessToken string) {
	accessToken, err := g.GetCurrentlyUsedAccessToken()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return accessToken
}

func (g *GitlabPersonalAccessTokenService) MustGetGitlab() (gitlab *GitlabInstance) {
	gitlab, err := g.GetGitlab()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlab
}

func (g *GitlabPersonalAccessTokenService) MustGetGitlabUsers() (gitlabUsers *GitlabUsers) {
	gitlabUsers, err := g.GetGitlabUsers()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlabUsers
}

func (g *GitlabPersonalAccessTokenService) MustGetNativeGitlabClient() (nativeClient *gitlab.Client) {
	nativeClient, err := g.GetNativeGitlabClient()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return nativeClient
}

func (g *GitlabPersonalAccessTokenService) MustGetNativePersonalTokenService() (nativeService *gitlab.PersonalAccessTokensService) {
	nativeService, err := g.GetNativePersonalTokenService()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return nativeService
}

func (g *GitlabPersonalAccessTokenService) MustGetNativeUsersService() (nativeService *gitlab.UsersService) {
	nativeService, err := g.GetNativeUsersService()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return nativeService
}

func (g *GitlabPersonalAccessTokenService) MustGetPersonalAccessTokenList(verbose bool) (tokens []*GitlabPersonalAccessToken) {
	tokens, err := g.GetPersonalAccessTokenList(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return tokens
}

func (g *GitlabPersonalAccessTokenService) MustGetPersonalAccessTokenNameList(verbose bool) (tokenNames []string) {
	tokenNames, err := g.GetPersonalAccessTokenNameList(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return tokenNames
}

func (g *GitlabPersonalAccessTokenService) MustGetTokenIdByName(tokenName string, verbose bool) (tokenId int) {
	tokenId, err := g.GetTokenIdByName(tokenName, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return tokenId
}

func (g *GitlabPersonalAccessTokenService) MustRecreateToken(tokenOptions *GitlabCreatePersonalAccessTokenOptions) (newToken string) {
	newToken, err := g.RecreateToken(tokenOptions)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return newToken
}

func (g *GitlabPersonalAccessTokenService) MustRevokeTokenByName(tokenName string, verbose bool) {
	err := g.RevokeTokenByName(tokenName, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabPersonalAccessTokenService) MustSetGitlab(gitlab *GitlabInstance) {
	err := g.SetGitlab(gitlab)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (p *GitlabPersonalAccessTokenService) CreateToken(tokenOptions *GitlabCreatePersonalAccessTokenOptions) (newToken string, err error) {
	if tokenOptions == nil {
		return "", TracedError("tokenOptions is nil")
	}

	tokenName, err := tokenOptions.GetName()
	if err != nil {
		return "", err
	}

	exists, err := p.ExistsByName(tokenName, tokenOptions.Verbose)
	if err != nil {
		return "", err
	}

	if exists {
		return "", TracedErrorf(
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
	stdout, err := Bash().RunCommandAndGetStdoutAsString(
		&RunCommandOptions{
			Command: command,
			Verbose: tokenOptions.Verbose,
		},
	)
	if err != nil {
		return "", err
	}

	newToken, err = Json().RunJqAgainstJsonStringAsString(stdout, ".token")
	if err != nil {
		return "", err
	}

	newToken = strings.TrimSpace(newToken)
	newToken = astrings.RemoveSurroundingQuotationMarks(newToken)
	if len(newToken) <= 0 {
		return "", TracedError("Unable to get newToken. newToken is empty string.")
	}

	return newToken, nil
}

func (p *GitlabPersonalAccessTokenService) ExistsByName(tokenName string, verbose bool) (exists bool, err error) {
	if len(tokenName) <= 0 {
		return false, TracedError("tokenName is nil")
	}

	tokenNames, err := p.GetPersonalAccessTokenNameList(verbose)
	if err != nil {
		return false, err
	}

	exists = aslices.ContainsString(tokenNames, tokenName)

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
		LogInfof("Current gitlab user id is '%d'", userId)
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
		return nil, TracedError("gitlab not set")
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
		return nil, TracedError("nativeService is nil")
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

func (p *GitlabPersonalAccessTokenService) GetPersonalAccessTokenList(verbose bool) (tokens []*GitlabPersonalAccessToken, err error) {
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

func (p *GitlabPersonalAccessTokenService) GetPersonalAccessTokenNameList(verbose bool) (tokenNames []string, err error) {
	tokens, err := p.gitlab.GetPersonalAccessTokenList(verbose)
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

func (p *GitlabPersonalAccessTokenService) GetTokenIdByName(tokenName string, verbose bool) (tokenId int, err error) {
	if len(tokenName) <= 0 {
		return -1, TracedError("tokenName is empty string")
	}

	tokens, err := p.GetPersonalAccessTokenList(verbose)
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

	return -1, TracedErrorf("No token with name '%s' found.", tokenName)
}

func (p *GitlabPersonalAccessTokenService) RecreateToken(tokenOptions *GitlabCreatePersonalAccessTokenOptions) (newToken string, err error) {
	if tokenOptions == nil {
		return "", TracedError("tokenOptions is nil")
	}

	tokenName, err := tokenOptions.GetName()
	if err != nil {
		return "", err
	}

	exists, err := p.ExistsByName(tokenName, tokenOptions.Verbose)
	if err != nil {
		return "", err
	}

	if exists {
		p.RevokeTokenByName(tokenName, tokenOptions.Verbose)

		nativeService, err := p.GetNativePersonalTokenService()
		if err != nil {
			return "", err
		}

		tokenId, err := p.GetTokenIdByName(tokenName, tokenOptions.Verbose)
		if err != nil {
			return "", err
		}

		nativeToken, _, err := nativeService.RotatePersonalAccessToken(tokenId, nil)
		if err != nil {
			return "", err
		}

		newToken = nativeToken.Token
		if len(newToken) <= 0 {
			return "", TracedError("recreate personal access token failed. NewToken is empty string.")
		}

		if tokenOptions.Verbose {
			LogInfof("Personal access token '%s' recreated.", tokenName)
		}

		return newToken, nil
	} else {
		newToken, err = p.CreateToken(tokenOptions)
		if err != nil {
			return "", err
		}
		if tokenOptions.Verbose {
			LogInfof("Personal access token '%s' created.", tokenName)
		}

		return newToken, nil
	}
}

func (p *GitlabPersonalAccessTokenService) RevokeTokenByName(tokenName string, verbose bool) (err error) {
	if len(tokenName) <= 0 {
		return TracedError("tokenName is empty string")
	}

	exists, err := p.ExistsByName(tokenName, verbose)
	if err != nil {
		return err
	}

	if exists {
		nativeService, err := p.GetNativePersonalTokenService()
		if err != nil {
			return err
		}

		tokenId, err := p.GetTokenIdByName(tokenName, verbose)
		if err != nil {
			return err
		}

		_, err = nativeService.RevokePersonalAccessToken(tokenId)
		if err != nil {
			return err
		}

		if verbose {
			LogChangedf("Personal gitlab access token '%s' deleted.", tokenName)
		}
	} else {
		if verbose {
			LogInfof("Personal gitlab access token '%s' was already deleted.", tokenName)
		}
	}

	return nil
}

func (p *GitlabPersonalAccessTokenService) SetGitlab(gitlab *GitlabInstance) (err error) {
	if gitlab == nil {
		return TracedError("gitlab is nil")
	}

	p.gitlab = gitlab

	return nil
}
