package asciichgolangpublic

import (
	"github.com/go-git/go-git/v5"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type LocalGitRemote struct {
	Name      string
	RemoteUrl string
}

func MustNewLocalGitRemoteByNativeGoGitRemote(goGitRemote *git.Remote) (l *LocalGitRemote) {
	l, err := NewLocalGitRemoteByNativeGoGitRemote(goGitRemote)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return l
}

func NewLocalGitRemote() (l *LocalGitRemote) {
	return new(LocalGitRemote)
}

func NewLocalGitRemoteByNativeGoGitRemote(goGitRemote *git.Remote) (l *LocalGitRemote, err error) {
	if goGitRemote == nil {
		return nil, tracederrors.TracedErrorEmptyString("goGitRemote")
	}

	l = NewLocalGitRemote()

	remoteConfig := goGitRemote.Config()
	if remoteConfig == nil {
		return nil, tracederrors.TracedErrorEmptyString("Config")
	}

	err = l.SetName(remoteConfig.Name)
	if err != nil {
		return nil, err
	}

	if len(remoteConfig.URLs) != 1 {
		return nil, tracederrors.TracedErrorf(
			"Only implemented for 1 remote URL at the moment but got '%v'",
			remoteConfig.URLs,
		)
	}

	err = l.SetRemoteUrl(remoteConfig.URLs[0])
	if err != nil {
		return nil, err
	}

	return l, nil
}

func (l *LocalGitRemote) GetName() (name string, err error) {
	if l.Name == "" {
		return "", tracederrors.TracedErrorf("Name not set")
	}

	return l.Name, nil
}

func (l *LocalGitRemote) GetRemoteUrl() (remoteUrl string, err error) {
	if l.RemoteUrl == "" {
		return "", tracederrors.TracedErrorf("RemoteUrl not set")
	}

	return l.RemoteUrl, nil
}

func (l *LocalGitRemote) MustGetName() (name string) {
	name, err := l.GetName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return name
}

func (l *LocalGitRemote) MustGetRemoteUrl() (remoteUrl string) {
	remoteUrl, err := l.GetRemoteUrl()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return remoteUrl
}

func (l *LocalGitRemote) MustSetName(name string) {
	err := l.SetName(name)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (l *LocalGitRemote) MustSetRemoteUrl(remoteUrl string) {
	err := l.SetRemoteUrl(remoteUrl)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (l *LocalGitRemote) SetName(name string) (err error) {
	if name == "" {
		return tracederrors.TracedErrorf("name is empty string")
	}

	l.Name = name

	return nil
}

func (l *LocalGitRemote) SetRemoteUrl(remoteUrl string) (err error) {
	if remoteUrl == "" {
		return tracederrors.TracedErrorf("remoteUrl is empty string")
	}

	l.RemoteUrl = remoteUrl

	return nil
}
