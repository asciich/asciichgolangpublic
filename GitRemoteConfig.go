package asciichgolangpublic

import (
	"github.com/asciich/asciichgolangpublic/errors"
	"github.com/asciich/asciichgolangpublic/logging"
)

type GitRemoteConfig struct {
	RemoteName string
	UrlFetch   string
	UrlPush    string
}

func NewGitRemoteConfig() (gitRemoteConfig *GitRemoteConfig) {
	return new(GitRemoteConfig)
}

func (c *GitRemoteConfig) Equals(other *GitRemoteConfig) (equals bool) {
	if other == nil {
		return false
	}

	if c.RemoteName != other.RemoteName {
		return false
	}

	if c.UrlFetch != other.UrlFetch {
		return false
	}

	if c.UrlPush != other.UrlPush {
		return false
	}

	return true
}

func (g *GitRemoteConfig) GetRemoteName() (remoteName string, err error) {
	if g.RemoteName == "" {
		return "", errors.TracedErrorf("RemoteName not set")
	}

	return g.RemoteName, nil
}

func (g *GitRemoteConfig) GetUrlFetch() (urlFetch string, err error) {
	if g.UrlFetch == "" {
		return "", errors.TracedErrorf("UrlFetch not set")
	}

	return g.UrlFetch, nil
}

func (g *GitRemoteConfig) GetUrlPush() (urlPush string, err error) {
	if g.UrlPush == "" {
		return "", errors.TracedErrorf("UrlPush not set")
	}

	return g.UrlPush, nil
}

func (g *GitRemoteConfig) MustGetRemoteName() (remoteName string) {
	remoteName, err := g.GetRemoteName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return remoteName
}

func (g *GitRemoteConfig) MustGetUrlFetch() (urlFetch string) {
	urlFetch, err := g.GetUrlFetch()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return urlFetch
}

func (g *GitRemoteConfig) MustGetUrlPush() (urlPush string) {
	urlPush, err := g.GetUrlPush()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return urlPush
}

func (g *GitRemoteConfig) MustSetRemoteName(remoteName string) {
	err := g.SetRemoteName(remoteName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitRemoteConfig) MustSetUrlFetch(urlFetch string) {
	err := g.SetUrlFetch(urlFetch)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitRemoteConfig) MustSetUrlPush(urlPush string) {
	err := g.SetUrlPush(urlPush)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitRemoteConfig) SetRemoteName(remoteName string) (err error) {
	if remoteName == "" {
		return errors.TracedErrorf("remoteName is empty string")
	}

	g.RemoteName = remoteName

	return nil
}

func (g *GitRemoteConfig) SetUrlFetch(urlFetch string) (err error) {
	if urlFetch == "" {
		return errors.TracedErrorf("urlFetch is empty string")
	}

	g.UrlFetch = urlFetch

	return nil
}

func (g *GitRemoteConfig) SetUrlPush(urlPush string) (err error) {
	if urlPush == "" {
		return errors.TracedErrorf("urlPush is empty string")
	}

	g.UrlPush = urlPush

	return nil
}
