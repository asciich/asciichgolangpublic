package asciichgolangpublic

import (
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type GitlabReadFileOptions struct {
	Path       string
	BranchName string
}

func NewGitlabReadFileOptions() (g *GitlabReadFileOptions) {
	return new(GitlabReadFileOptions)
}

func (g *GitlabReadFileOptions) GetBranchName() (branchName string, err error) {
	if g.BranchName == "" {
		return "", tracederrors.TracedErrorf("BranchName not set")
	}

	return g.BranchName, nil
}

func (g *GitlabReadFileOptions) GetDeepCopy() (deepCopy *GitlabReadFileOptions) {
	deepCopy = NewGitlabReadFileOptions()

	*deepCopy = *g

	return deepCopy
}

func (g *GitlabReadFileOptions) GetGitlabGetRepositoryFileOptions() (getOptions *GitlabGetRepositoryFileOptions, err error) {
	getOptions = NewGitlabGetRepositoryFileOptions()
	getOptions.Path = g.Path
	getOptions.BranchName = g.BranchName
	return getOptions, nil
}

func (g *GitlabReadFileOptions) GetPath() (path string, err error) {
	if g.Path == "" {
		return "", tracederrors.TracedErrorf("Path not set")
	}

	return g.Path, nil
}

func (g *GitlabReadFileOptions) MustGetBranchName() (branchName string) {
	branchName, err := g.GetBranchName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return branchName
}

func (g *GitlabReadFileOptions) MustGetGitlabGetRepositoryFileOptions() (getOptions *GitlabGetRepositoryFileOptions) {
	getOptions, err := g.GetGitlabGetRepositoryFileOptions()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return getOptions
}

func (g *GitlabReadFileOptions) MustGetPath() (path string) {
	path, err := g.GetPath()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return path
}

func (g *GitlabReadFileOptions) MustSetBranchName(branchName string) {
	err := g.SetBranchName(branchName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabReadFileOptions) MustSetPath(path string) {
	err := g.SetPath(path)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabReadFileOptions) SetBranchName(branchName string) (err error) {
	if branchName == "" {
		return tracederrors.TracedErrorf("branchName is empty string")
	}

	g.BranchName = branchName

	return nil
}

func (g *GitlabReadFileOptions) SetPath(path string) (err error) {
	if path == "" {
		return tracederrors.TracedErrorf("path is empty string")
	}

	g.Path = path

	return nil
}
