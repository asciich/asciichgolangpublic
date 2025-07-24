package asciichgolangpublic

import (
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type GitlabGetRepositoryFileOptions struct {
	Path       string
	BranchName string
	Verbose    bool
}

func NewGitlabGetRepositoryFileOptions() (g *GitlabGetRepositoryFileOptions) {
	return new(GitlabGetRepositoryFileOptions)
}

func (g *GitlabGetRepositoryFileOptions) GetBranchName() (branchName string, err error) {
	if g.BranchName == "" {
		return "", tracederrors.TracedErrorf("BranchName not set")
	}

	return g.BranchName, nil
}

func (g *GitlabGetRepositoryFileOptions) GetPath() (path string, err error) {
	if g.Path == "" {
		return "", tracederrors.TracedErrorf("Path not set")
	}

	return g.Path, nil
}

func (g *GitlabGetRepositoryFileOptions) GetVerbose() (verbose bool) {

	return g.Verbose
}

func (g *GitlabGetRepositoryFileOptions) IsBranchNameSet() (isSet bool) {
	return g.BranchName != ""
}

func (g *GitlabGetRepositoryFileOptions) MustGetBranchName() (branchName string) {
	branchName, err := g.GetBranchName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return branchName
}

func (g *GitlabGetRepositoryFileOptions) MustGetPath() (path string) {
	path, err := g.GetPath()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return path
}

func (g *GitlabGetRepositoryFileOptions) MustSetBranchName(branchName string) {
	err := g.SetBranchName(branchName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabGetRepositoryFileOptions) MustSetPath(path string) {
	err := g.SetPath(path)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabGetRepositoryFileOptions) SetBranchName(branchName string) (err error) {
	if branchName == "" {
		return tracederrors.TracedErrorf("branchName is empty string")
	}

	g.BranchName = branchName

	return nil
}

func (g *GitlabGetRepositoryFileOptions) SetPath(path string) (err error) {
	if path == "" {
		return tracederrors.TracedErrorf("path is empty string")
	}

	g.Path = path

	return nil
}

func (g *GitlabGetRepositoryFileOptions) SetVerbose(verbose bool) {
	g.Verbose = verbose
}
