package gitgeneric

import (
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type GenericGitRemoteConfig struct {
	RemoteName string
	UrlFetch   string
	UrlPush    string
}

func NewGitRemoteConfig() (gitRemoteConfig *GenericGitRemoteConfig) {
	return new(GenericGitRemoteConfig)
}

func (c *GenericGitRemoteConfig) Equals(other gitinterfaces.GitRemoteConfig) (equals bool) {
	if other == nil {
		return false
	}

	otherRemoteName, _ := other.GetRemoteName()
	if c.RemoteName != otherRemoteName {
		return false
	}

	otherUrlFetch, _ := other.GetUrlFetch()
	if c.UrlFetch != otherUrlFetch {
		return false
	}

	otherUrlPush, _ := other.GetUrlPush()
	if c.UrlPush != otherUrlPush {
		return false
	}

	return true
}

func (g *GenericGitRemoteConfig) GetRemoteName() (remoteName string, err error) {
	if g.RemoteName == "" {
		return "", tracederrors.TracedErrorf("RemoteName not set")
	}

	return g.RemoteName, nil
}

func (g *GenericGitRemoteConfig) GetUrlFetch() (urlFetch string, err error) {
	if g.UrlFetch == "" {
		return "", tracederrors.TracedErrorf("UrlFetch not set")
	}

	return g.UrlFetch, nil
}

func (g *GenericGitRemoteConfig) GetUrlPush() (urlPush string, err error) {
	if g.UrlPush == "" {
		return "", tracederrors.TracedErrorf("UrlPush not set")
	}

	return g.UrlPush, nil
}

func (g *GenericGitRemoteConfig) SetRemoteName(remoteName string) (err error) {
	if remoteName == "" {
		return tracederrors.TracedErrorf("remoteName is empty string")
	}

	g.RemoteName = remoteName

	return nil
}

func (g *GenericGitRemoteConfig) SetUrlFetch(urlFetch string) (err error) {
	if urlFetch == "" {
		return tracederrors.TracedErrorf("urlFetch is empty string")
	}

	g.UrlFetch = urlFetch

	return nil
}

func (g *GenericGitRemoteConfig) SetUrlPush(urlPush string) (err error) {
	if urlPush == "" {
		return tracederrors.TracedErrorf("urlPush is empty string")
	}

	g.UrlPush = urlPush

	return nil
}
