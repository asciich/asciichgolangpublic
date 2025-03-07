package asciichgolangpublic

import (
	"github.com/asciich/asciichgolangpublic/datatypes/slicesutils"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type GitlabWriteFileOptions struct {
	Path          string
	Content       []byte
	BranchName    string
	CommitMessage string
	Verbose       bool
}

func NewGitlabWriteFileOptions() (g *GitlabWriteFileOptions) {
	return new(GitlabWriteFileOptions)
}

func (g *GitlabWriteFileOptions) GetBranchName() (branchName string, err error) {
	if g.BranchName == "" {
		return "", tracederrors.TracedErrorf("BranchName not set")
	}

	return g.BranchName, nil
}

func (g *GitlabWriteFileOptions) GetCommitMessage() (commitMessage string, err error) {
	if g.CommitMessage == "" {
		return "", tracederrors.TracedErrorf("CommitMessage not set")
	}

	return g.CommitMessage, nil
}

func (g *GitlabWriteFileOptions) GetContent() (content []byte, err error) {
	if g.Content == nil {
		return nil, tracederrors.TracedErrorf("Content not set")
	}

	if len(g.Content) <= 0 {
		return nil, tracederrors.TracedErrorf("Content has no elements")
	}

	return g.Content, nil
}

func (g *GitlabWriteFileOptions) GetDeepCopy() (copy *GitlabWriteFileOptions) {
	copy = NewGitlabWriteFileOptions()
	*copy = *g

	if len(g.Content) > 0 {
		copy.Content = slicesutils.GetDeepCopyOfByteSlice(g.Content)
	}

	return copy
}

func (g *GitlabWriteFileOptions) GetGitlabGetRepositoryFileOptions() (getOptions *GitlabGetRepositoryFileOptions, err error) {
	getOptions = NewGitlabGetRepositoryFileOptions()
	getOptions.Path = g.Path
	getOptions.BranchName = g.BranchName
	getOptions.Verbose = g.Verbose
	return getOptions, nil
}

func (g *GitlabWriteFileOptions) GetPath() (path string, err error) {
	if g.Path == "" {
		return "", tracederrors.TracedErrorf("Path not set")
	}

	return g.Path, nil
}

func (g *GitlabWriteFileOptions) GetVerbose() (verbose bool) {

	return g.Verbose
}

func (g *GitlabWriteFileOptions) MustGetBranchName() (branchName string) {
	branchName, err := g.GetBranchName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return branchName
}

func (g *GitlabWriteFileOptions) MustGetCommitMessage() (commitMessage string) {
	commitMessage, err := g.GetCommitMessage()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return commitMessage
}

func (g *GitlabWriteFileOptions) MustGetContent() (content []byte) {
	content, err := g.GetContent()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return content
}

func (g *GitlabWriteFileOptions) MustGetGitlabGetRepositoryFileOptions() (getOptions *GitlabGetRepositoryFileOptions) {
	getOptions, err := g.GetGitlabGetRepositoryFileOptions()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return getOptions
}

func (g *GitlabWriteFileOptions) MustGetPath() (path string) {
	path, err := g.GetPath()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return path
}

func (g *GitlabWriteFileOptions) MustSetBranchName(branchName string) {
	err := g.SetBranchName(branchName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabWriteFileOptions) MustSetCommitMessage(commitMessage string) {
	err := g.SetCommitMessage(commitMessage)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabWriteFileOptions) MustSetContent(content []byte) {
	err := g.SetContent(content)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabWriteFileOptions) MustSetPath(path string) {
	err := g.SetPath(path)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabWriteFileOptions) SetBranchName(branchName string) (err error) {
	if branchName == "" {
		return tracederrors.TracedErrorf("branchName is empty string")
	}

	g.BranchName = branchName

	return nil
}

func (g *GitlabWriteFileOptions) SetCommitMessage(commitMessage string) (err error) {
	if commitMessage == "" {
		return tracederrors.TracedErrorf("commitMessage is empty string")
	}

	g.CommitMessage = commitMessage

	return nil
}

func (g *GitlabWriteFileOptions) SetContent(content []byte) (err error) {
	if content == nil {
		return tracederrors.TracedErrorf("content is nil")
	}

	if len(content) <= 0 {
		return tracederrors.TracedErrorf("content has no elements")
	}

	g.Content = content

	return nil
}

func (g *GitlabWriteFileOptions) SetPath(path string) (err error) {
	if path == "" {
		return tracederrors.TracedErrorf("path is empty string")
	}

	g.Path = path

	return nil
}

func (g *GitlabWriteFileOptions) SetVerbose(verbose bool) {
	g.Verbose = verbose
}
