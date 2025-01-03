package asciichgolangpublic

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
		LogGoErrorFatal(err)
	}

	return remoteName
}

func (g *GitRemoteAddOptions) MustGetRemoteUrl() (remoteUrl string) {
	remoteUrl, err := g.GetRemoteUrl()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return remoteUrl
}

func (g *GitRemoteAddOptions) MustGetVerbose() (verbose bool) {
	verbose, err := g.GetVerbose()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return verbose
}

func (g *GitRemoteAddOptions) MustSetRemoteName(remoteName string) {
	err := g.SetRemoteName(remoteName)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitRemoteAddOptions) MustSetRemoteUrl(remoteUrl string) {
	err := g.SetRemoteUrl(remoteUrl)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitRemoteAddOptions) MustSetVerbose(verbose bool) {
	err := g.SetVerbose(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitRemoteAddOptions) SetRemoteName(remoteName string) (err error) {
	if remoteName == "" {
		return TracedErrorf("remoteName is empty string")
	}

	g.RemoteName = remoteName

	return nil
}

func (g *GitRemoteAddOptions) SetRemoteUrl(remoteUrl string) (err error) {
	if remoteUrl == "" {
		return TracedErrorf("remoteUrl is empty string")
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
		return "", TracedError("RemoteName not set")
	}

	return o.RemoteName, nil
}

func (o *GitRemoteAddOptions) GetRemoteUrl() (remoteUrl string, err error) {
	if len(o.RemoteUrl) <= 0 {
		return "", TracedError("RemoteUrl not set")
	}

	return o.RemoteUrl, nil
}
