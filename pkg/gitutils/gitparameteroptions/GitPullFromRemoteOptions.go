package gitparameteroptions

import (
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type GitPullFromRemoteOptions struct {
	RemoteName string
	BranchName string
}

func NewGitPullFromRemoteOptions() (g *GitPullFromRemoteOptions) {
	return new(GitPullFromRemoteOptions)
}

func (g *GitPullFromRemoteOptions) SetBranchName(branchName string) (err error) {
	if branchName == "" {
		return tracederrors.TracedErrorf("branchName is empty string")
	}

	g.BranchName = branchName

	return nil
}

func (g *GitPullFromRemoteOptions) SetRemoteName(remoteName string) (err error) {
	if remoteName == "" {
		return tracederrors.TracedErrorf("remoteName is empty string")
	}

	g.RemoteName = remoteName

	return nil
}

func (o *GitPullFromRemoteOptions) GetBranchName() (branchName string, err error) {
	if len(o.BranchName) <= 0 {
		return "", tracederrors.TracedError("BranchName not set")
	}

	return o.BranchName, nil
}

func (o *GitPullFromRemoteOptions) GetRemoteName() (remoteName string, err error) {
	if len(o.RemoteName) <= 0 {
		return "", tracederrors.TracedError("RemoteName not set")
	}

	return o.RemoteName, nil
}
