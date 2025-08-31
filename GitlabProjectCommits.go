package asciichgolangpublic

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type GitlabProjectCommits struct {
	gitlabProject *GitlabProject
}

func NewGitlabProjectCommits() (g *GitlabProjectCommits) {
	return new(GitlabProjectCommits)
}

func (g *GitlabProjectCommits) GetCommitByHashString(ctx context.Context, hash string) (commit *GitlabCommit, err error) {
	if hash == "" {
		return nil, tracederrors.TracedErrorEmptyString("hash")
	}

	commit = NewGitlabCommit()

	err = commit.SetGitlabProjectsCommits(g)
	if err != nil {
		return nil, err
	}

	err = commit.SetCommitHash(hash)
	if err != nil {
		return nil, err
	}

	return commit, nil
}

func (g *GitlabProjectCommits) GetGitlab() (gitlab *GitlabInstance, err error) {
	project, err := g.GetGitlabProject()
	if err != nil {
		return nil, err
	}

	gitlab, err = project.GetGitlab()
	if err != nil {
		return nil, err
	}

	return gitlab, nil
}

func (g *GitlabProjectCommits) GetGitlabProject() (gitlabProject *GitlabProject, err error) {
	if g.gitlabProject == nil {
		return nil, tracederrors.TracedErrorf("gitlabProject not set")
	}

	return g.gitlabProject, nil
}

func (g *GitlabProjectCommits) GetNativeCommitsService() (nativeCommitsService *gitlab.CommitsService, err error) {
	nativeClinet, err := g.GetNativeGitlabClient()
	if err != nil {
		return nil, err
	}

	nativeCommitsService = nativeClinet.Commits

	if nativeCommitsService == nil {
		return nil, tracederrors.TracedError("nativeCommitsService is nil after evaluation")
	}

	return nativeCommitsService, nil
}

func (g *GitlabProjectCommits) GetNativeGitlabClient() (nativeClient *gitlab.Client, err error) {
	gitlab, err := g.GetGitlab()
	if err != nil {
		return nil, err
	}

	nativeClient, err = gitlab.GetNativeClient()
	if err != nil {
		return nil, err
	}

	return nativeClient, nil
}

func (g *GitlabProjectCommits) SetGitlabProject(gitlabProject *GitlabProject) (err error) {
	if gitlabProject == nil {
		return tracederrors.TracedErrorf("gitlabProject is nil")
	}

	g.gitlabProject = gitlabProject

	return nil
}
