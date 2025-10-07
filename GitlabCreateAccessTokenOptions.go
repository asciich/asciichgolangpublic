package asciichgolangpublic

import (
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/datatypes/slicesutils"
	"github.com/asciich/asciichgolangpublic/pkg/datetime/durationparser"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type GitlabCreateAccessTokenOptions struct {
	UserName  string
	TokenName string
	Scopes    []string
	ExpiresAt *time.Time
}

func NewGitlabCreateAccessTokenOptions() (g *GitlabCreateAccessTokenOptions) {
	return new(GitlabCreateAccessTokenOptions)
}

func (g *GitlabCreateAccessTokenOptions) GetExpiresAt() (expiresAt *time.Time, err error) {
	if g.ExpiresAt == nil {
		return nil, tracederrors.TracedErrorf("ExpiresAt not set")
	}

	return g.ExpiresAt, nil
}

func (g *GitlabCreateAccessTokenOptions) SetExpiresAt(expiresAt *time.Time) (err error) {
	if expiresAt == nil {
		return tracederrors.TracedErrorf("expiresAt is nil")
	}

	g.ExpiresAt = expiresAt

	return nil
}

func (g *GitlabCreateAccessTokenOptions) SetScopes(scopes []string) (err error) {
	if scopes == nil {
		return tracederrors.TracedErrorf("scopes is nil")
	}

	if len(scopes) <= 0 {
		return tracederrors.TracedErrorf("scopes has no elements")
	}

	g.Scopes = scopes

	return nil
}

func (g *GitlabCreateAccessTokenOptions) SetTokenName(tokenName string) (err error) {
	if tokenName == "" {
		return tracederrors.TracedErrorf("tokenName is empty string")
	}

	g.TokenName = tokenName

	return nil
}

func (g *GitlabCreateAccessTokenOptions) SetUserName(userName string) (err error) {
	if userName == "" {
		return tracederrors.TracedErrorf("userName is empty string")
	}

	g.UserName = userName

	return nil
}

func (o *GitlabCreateAccessTokenOptions) GetExipiresAtOrDefaultIfUnset() (expiresAt *time.Time, err error) {
	if o.ExpiresAt != nil {
		return o.ExpiresAt, nil
	}

	defaultTime := time.Now()

	monthDuration, err := durationparser.ToSecondsAsTimeDuration("1month")
	if err != nil {
		return nil, err
	}
	defaultTime = defaultTime.Add(*monthDuration)

	return &defaultTime, nil
}

func (o *GitlabCreateAccessTokenOptions) GetScopes() (scopes []string, err error) {
	if len(o.Scopes) <= 0 {
		return nil, tracederrors.TracedError("Scopes not set")
	}

	return slicesutils.GetDeepCopyOfStringsSlice(o.Scopes), nil
}

func (o *GitlabCreateAccessTokenOptions) GetTokenName() (tokenName string, err error) {
	if len(o.TokenName) <= 0 {
		return "", tracederrors.TracedError("TokenName not set")
	}

	return o.TokenName, nil
}

func (o *GitlabCreateAccessTokenOptions) GetUserName() (userName string, err error) {
	if len(o.UserName) <= 0 {
		return "", tracederrors.TracedError("UserName not set")
	}

	return o.UserName, nil
}
