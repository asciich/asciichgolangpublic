package nativegitoo

import "github.com/asciich/asciichgolangpublic/pkg/gitutils/gitinterfaces"

type NativeGitRemoteConfig struct {
	remoteName string
	urlFetch   string
	urlPush    string
}

func (c *NativeGitRemoteConfig) GetRemoteName() (string, error) {
	return c.remoteName, nil
}

func (c *NativeGitRemoteConfig) GetUrlFetch() (string, error) {
	return c.urlFetch, nil
}

func (c *NativeGitRemoteConfig) GetUrlPush() (string, error) {
	return c.urlPush, nil
}

func (c *NativeGitRemoteConfig) SetUrlFetch(url string) error {
	c.urlFetch = url
	return nil
}

func (c *NativeGitRemoteConfig) SetUrlPush(url string) error {
	c.urlPush = url
	return nil
}

func (c *NativeGitRemoteConfig) Equals(other gitinterfaces.GitRemoteConfig) bool {
	otherName, err := other.GetRemoteName()
	if err != nil {
		return false
	}
	otherFetch, err := other.GetUrlFetch()
	if err != nil {
		return false
	}
	otherPush, err := other.GetUrlPush()
	if err != nil {
		return false
	}
	return c.remoteName == otherName && c.urlFetch == otherFetch && c.urlPush == otherPush
}
