package gitparameteroptions

import (
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type GitRemoteAddOptions struct {
	RemoteName string
	RemoteUrl  string
	Verbose    bool
}

func NewGitRemoteAddOptions() (g *GitRemoteAddOptions) {
	return new(GitRemoteAddOptions)
}

func (g *GitRemoteAddOptions) GetVerbose() (verbose bool, err error) {

	return g.Verbose, nil
}

func (g *GitRemoteAddOptions) MustGetRemoteName() (remoteName string) {
	remoteName, err := g.GetRemoteName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return remoteName
}

func (g *GitRemoteAddOptions) MustGetRemoteUrl() (remoteUrl string) {
	remoteUrl, err := g.GetRemoteUrl()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return remoteUrl
}

func (g *GitRemoteAddOptions) MustGetVerbose() (verbose bool) {
	verbose, err := g.GetVerbose()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return verbose
}

func (g *GitRemoteAddOptions) MustSetRemoteName(remoteName string) {
	err := g.SetRemoteName(remoteName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitRemoteAddOptions) MustSetRemoteUrl(remoteUrl string) {
	err := g.SetRemoteUrl(remoteUrl)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitRemoteAddOptions) MustSetVerbose(verbose bool) {
	err := g.SetVerbose(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitRemoteAddOptions) SetRemoteName(remoteName string) (err error) {
	if remoteName == "" {
		return tracederrors.TracedErrorf("remoteName is empty string")
	}

	g.RemoteName = remoteName

	return nil
}

func (g *GitRemoteAddOptions) SetRemoteUrl(remoteUrl string) (err error) {
	if remoteUrl == "" {
		return tracederrors.TracedErrorf("remoteUrl is empty string")
	}

	g.RemoteUrl = remoteUrl

	return nil
}

func (g *GitRemoteAddOptions) SetVerbose(verbose bool) (err error) {
	g.Verbose = verbose

	return nil
}

func (o *GitRemoteAddOptions) GetRemoteName() (remoteName string, err error) {
	if len(o.RemoteName) <= 0 {
		return "", tracederrors.TracedError("RemoteName not set")
	}

	return o.RemoteName, nil
}

func (o *GitRemoteAddOptions) GetRemoteUrl() (remoteUrl string, err error) {
	if len(o.RemoteUrl) <= 0 {
		return "", tracederrors.TracedError("RemoteUrl not set")
	}

	return o.RemoteUrl, nil
}
