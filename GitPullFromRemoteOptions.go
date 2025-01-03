package asciichgolangpublic

type GitPullFromRemoteOptions struct {
	RemoteName string
	BranchName string
	Verbose    bool
}

func NewGitPullFromRemoteOptions() (g *GitPullFromRemoteOptions) {
	return new(GitPullFromRemoteOptions)
}

func (g *GitPullFromRemoteOptions) GetVerbose() (verbose bool, err error) {

	return g.Verbose, nil
}

func (g *GitPullFromRemoteOptions) MustGetBranchName() (branchName string) {
	branchName, err := g.GetBranchName()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return branchName
}

func (g *GitPullFromRemoteOptions) MustGetRemoteName() (remoteName string) {
	remoteName, err := g.GetRemoteName()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return remoteName
}

func (g *GitPullFromRemoteOptions) MustGetVerbose() (verbose bool) {
	verbose, err := g.GetVerbose()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return verbose
}

func (g *GitPullFromRemoteOptions) MustSetBranchName(branchName string) {
	err := g.SetBranchName(branchName)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitPullFromRemoteOptions) MustSetRemoteName(remoteName string) {
	err := g.SetRemoteName(remoteName)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitPullFromRemoteOptions) MustSetVerbose(verbose bool) {
	err := g.SetVerbose(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitPullFromRemoteOptions) SetBranchName(branchName string) (err error) {
	if branchName == "" {
		return TracedErrorf("branchName is empty string")
	}

	g.BranchName = branchName

	return nil
}

func (g *GitPullFromRemoteOptions) SetRemoteName(remoteName string) (err error) {
	if remoteName == "" {
		return TracedErrorf("remoteName is empty string")
	}

	g.RemoteName = remoteName

	return nil
}

func (g *GitPullFromRemoteOptions) SetVerbose(verbose bool) (err error) {
	g.Verbose = verbose

	return nil
}

func (o *GitPullFromRemoteOptions) GetBranchName() (branchName string, err error) {
	if len(o.BranchName) <= 0 {
		return "", TracedError("BranchName not set")
	}

	return o.BranchName, nil
}

func (o *GitPullFromRemoteOptions) GetRemoteName() (remoteName string, err error) {
	if len(o.RemoteName) <= 0 {
		return "", TracedError("RemoteName not set")
	}

	return o.RemoteName, nil
}
