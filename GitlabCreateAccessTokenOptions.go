package asciichgolangpublic

import (
	"time"

)

type GitlabCreateAccessTokenOptions struct {
	UserName  string
	TokenName string
	Scopes    []string
	ExpiresAt *time.Time
	Verbose   bool
}

func NewGitlabCreateAccessTokenOptions() (g *GitlabCreateAccessTokenOptions) {
	return new(GitlabCreateAccessTokenOptions)
}

func (g *GitlabCreateAccessTokenOptions) GetExpiresAt() (expiresAt *time.Time, err error) {
	if g.ExpiresAt == nil {
		return nil, TracedErrorf("ExpiresAt not set")
	}

	return g.ExpiresAt, nil
}

func (g *GitlabCreateAccessTokenOptions) GetVerbose() (verbose bool, err error) {

	return g.Verbose, nil
}

func (g *GitlabCreateAccessTokenOptions) MustGetExipiresAtOrDefaultIfUnset() (expiresAt *time.Time) {
	expiresAt, err := g.GetExipiresAtOrDefaultIfUnset()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return expiresAt
}

func (g *GitlabCreateAccessTokenOptions) MustGetExpiresAt() (expiresAt *time.Time) {
	expiresAt, err := g.GetExpiresAt()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return expiresAt
}

func (g *GitlabCreateAccessTokenOptions) MustGetScopes() (scopes []string) {
	scopes, err := g.GetScopes()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return scopes
}

func (g *GitlabCreateAccessTokenOptions) MustGetTokenName() (tokenName string) {
	tokenName, err := g.GetTokenName()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return tokenName
}

func (g *GitlabCreateAccessTokenOptions) MustGetUserName() (userName string) {
	userName, err := g.GetUserName()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return userName
}

func (g *GitlabCreateAccessTokenOptions) MustGetVerbose() (verbose bool) {
	verbose, err := g.GetVerbose()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return verbose
}

func (g *GitlabCreateAccessTokenOptions) MustSetExpiresAt(expiresAt *time.Time) {
	err := g.SetExpiresAt(expiresAt)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabCreateAccessTokenOptions) MustSetScopes(scopes []string) {
	err := g.SetScopes(scopes)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabCreateAccessTokenOptions) MustSetTokenName(tokenName string) {
	err := g.SetTokenName(tokenName)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabCreateAccessTokenOptions) MustSetUserName(userName string) {
	err := g.SetUserName(userName)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabCreateAccessTokenOptions) MustSetVerbose(verbose bool) {
	err := g.SetVerbose(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabCreateAccessTokenOptions) SetExpiresAt(expiresAt *time.Time) (err error) {
	if expiresAt == nil {
		return TracedErrorf("expiresAt is nil")
	}

	g.ExpiresAt = expiresAt

	return nil
}

func (g *GitlabCreateAccessTokenOptions) SetScopes(scopes []string) (err error) {
	if scopes == nil {
		return TracedErrorf("scopes is nil")
	}

	if len(scopes) <= 0 {
		return TracedErrorf("scopes has no elements")
	}

	g.Scopes = scopes

	return nil
}

func (g *GitlabCreateAccessTokenOptions) SetTokenName(tokenName string) (err error) {
	if tokenName == "" {
		return TracedErrorf("tokenName is empty string")
	}

	g.TokenName = tokenName

	return nil
}

func (g *GitlabCreateAccessTokenOptions) SetUserName(userName string) (err error) {
	if userName == "" {
		return TracedErrorf("userName is empty string")
	}

	g.UserName = userName

	return nil
}

func (g *GitlabCreateAccessTokenOptions) SetVerbose(verbose bool) (err error) {
	g.Verbose = verbose

	return nil
}

func (o *GitlabCreateAccessTokenOptions) GetExipiresAtOrDefaultIfUnset() (expiresAt *time.Time, err error) {
	if o.ExpiresAt != nil {
		return o.ExpiresAt, nil
	}

	defaultTime := time.Now()
	defaultTime = defaultTime.Add(
		*DurationParser().MustToSecondsAsTimeDuration("1month"),
	)

	return &defaultTime, nil
}

func (o *GitlabCreateAccessTokenOptions) GetScopes() (scopes []string, err error) {
	if len(o.Scopes) <= 0 {
		return nil, TracedError("Scopes not set")
	}

	return Slices().GetDeepCopyOfStringsSlice(o.Scopes), nil
}

func (o *GitlabCreateAccessTokenOptions) GetTokenName() (tokenName string, err error) {
	if len(o.TokenName) <= 0 {
		return "", TracedError("TokenName not set")
	}

	return o.TokenName, nil
}

func (o *GitlabCreateAccessTokenOptions) GetUserName() (userName string, err error) {
	if len(o.UserName) <= 0 {
		return "", TracedError("UserName not set")
	}

	return o.UserName, nil
}
