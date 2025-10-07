package asciichgolangpublic

import (
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type GitlabCreateBranchOptions struct {
	SourceBranchName    string
	BranchName          string
	FailIfAlreadyExists bool
}

func NewGitlabCreateBranchOptions() (g *GitlabCreateBranchOptions) {
	return new(GitlabCreateBranchOptions)
}

func (g *GitlabCreateBranchOptions) GetBranchName() (branchName string, err error) {
	if g.BranchName == "" {
		return "", tracederrors.TracedErrorf("BranchName not set")
	}

	return g.BranchName, nil
}

func (g *GitlabCreateBranchOptions) GetFailIfAlreadyExists() (failIfAlreadyExists bool) {

	return g.FailIfAlreadyExists
}

func (g *GitlabCreateBranchOptions) GetSourceBranchName() (sourceBranchName string, err error) {
	if g.SourceBranchName == "" {
		return "", tracederrors.TracedErrorf("SourceBranchName not set")
	}

	return g.SourceBranchName, nil
}

func (g *GitlabCreateBranchOptions) SetBranchName(branchName string) (err error) {
	if branchName == "" {
		return tracederrors.TracedErrorf("branchName is empty string")
	}

	g.BranchName = branchName

	return nil
}

func (g *GitlabCreateBranchOptions) SetFailIfAlreadyExists(failIfAlreadyExists bool) {
	g.FailIfAlreadyExists = failIfAlreadyExists
}

func (g *GitlabCreateBranchOptions) SetSourceBranchName(sourceBranchName string) (err error) {
	if sourceBranchName == "" {
		return tracederrors.TracedErrorf("sourceBranchName is empty string")
	}

	g.SourceBranchName = sourceBranchName

	return nil
}
