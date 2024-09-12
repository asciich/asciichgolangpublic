package asciichgolangpublic

import "github.com/xanzy/go-gitlab"

type GitlabCommit struct {
	gitlabProjectsCommits *GitlabProjectCommits
	commitHash            string
}

func NewGitlabCommit() (g *GitlabCommit) {
	return new(GitlabCommit)
}

func (g *GitlabCommit) GetCommitHash() (commitHash string, err error) {
	if g.commitHash == "" {
		return "", TracedErrorf("commitHash not set")
	}

	return g.commitHash, nil
}

func (g *GitlabCommit) GetGitlab() (gitlab *GitlabInstance, err error) {
	gitlabProjectCommits, err := g.GetGitlabProjectsCommits()
	if err != nil {
		return nil, err
	}

	gitlab, err = gitlabProjectCommits.GetGitlab()
	if err != nil {
		return nil, err
	}

	return gitlab, nil
}

func (g *GitlabCommit) GetGitlabProject() (gitlabProject *GitlabProject, err error) {
	gitlabProjectCommit, err := g.GetGitlabProjectsCommits()
	if err != nil {
		return nil, err
	}

	gitlabProject, err = gitlabProjectCommit.GetGitlabProject()
	if err != nil {
		return nil, err
	}

	return gitlabProject, nil
}

func (g *GitlabCommit) GetGitlabProjectsCommits() (gitlabProjectsCommit *GitlabProjectCommits, err error) {
	if g.gitlabProjectsCommits == nil {
		return nil, TracedErrorf("gitlabProjectsCommit not set")
	}

	return g.gitlabProjectsCommits, nil
}

func (g *GitlabCommit) GetNativeCommitsService() (nativeCommitsService *gitlab.CommitsService, err error) {
	gitlabProjectCommits, err := g.GetGitlabProjectsCommits()
	if err != nil {
		return nil, err
	}

	nativeCommitsService, err = gitlabProjectCommits.GetNativeCommitsService()
	if err != nil {
		return nil, err
	}

	return nativeCommitsService, nil
}

func (g *GitlabCommit) GetParentCommitHashesAsString(verbose bool) (parentCommitHashes []string, err error) {
	rawResponse, err := g.GetRawResponse()
	if err != nil {
		return nil, err
	}

	parentCommitHashes = rawResponse.ParentIDs

	if verbose {
		hash, err := g.GetCommitHash()
		if err != nil {
			return nil, err
		}

		LogInfof("Commit '%s' has parent commit hashes '%v'.", hash, parentCommitHashes)
	}

	return parentCommitHashes, nil
}

func (g *GitlabCommit) GetParentCommits(verbose bool) (parentCommits []*GitlabCommit, err error) {
	parentCommitHashes, err := g.GetParentCommitHashesAsString(verbose)
	if err != nil {
		return nil, err
	}

	projectCommits, err := g.GetGitlabProjectsCommits()
	if err != nil {
		return nil, err
	}

	parentCommits = []*GitlabCommit{}
	for _, hash := range parentCommitHashes {
		toAdd, err := projectCommits.GetCommitByHashString(hash, verbose)
		if err != nil {
			return nil, err
		}

		parentCommits = append(parentCommits, toAdd)
	}

	return parentCommits, nil
}

func (g *GitlabCommit) GetProjectId() (projectId int, err error) {
	gitlabProject, err := g.GetGitlabProject()
	if err != nil {
		return -1, err
	}

	projectId, err = gitlabProject.GetId()
	if err != nil {
		return -1, err
	}

	return projectId, nil
}

func (g *GitlabCommit) GetRawResponse() (rawResponse *gitlab.Commit, err error) {
	nativeCommitsService, err := g.GetNativeCommitsService()
	if err != nil {
		return nil, err
	}

	hash, err := g.GetCommitHash()
	if err != nil {
		return nil, err
	}

	projectId, err := g.GetProjectId()
	if err != nil {
		return nil, err
	}

	rawResponse, _, err = nativeCommitsService.GetCommit(projectId, hash, &gitlab.GetCommitOptions{})
	if err != nil {
		return nil, TracedErrorf("Unable to get commit: %w", err)
	}

	return rawResponse, nil
}

func (g *GitlabCommit) MustGetCommitHash() (commitHash string) {
	commitHash, err := g.GetCommitHash()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return commitHash
}

func (g *GitlabCommit) MustGetGitlab() (gitlab *GitlabInstance) {
	gitlab, err := g.GetGitlab()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlab
}

func (g *GitlabCommit) MustGetGitlabProject() (gitlabProject *GitlabProject) {
	gitlabProject, err := g.GetGitlabProject()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlabProject
}

func (g *GitlabCommit) MustGetGitlabProjectsCommits() (gitlabProjectsCommit *GitlabProjectCommits) {
	gitlabProjectsCommit, err := g.GetGitlabProjectsCommits()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlabProjectsCommit
}

func (g *GitlabCommit) MustGetNativeCommitsService() (nativeCommitsService *gitlab.CommitsService) {
	nativeCommitsService, err := g.GetNativeCommitsService()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return nativeCommitsService
}

func (g *GitlabCommit) MustGetParentCommitHashesAsString(verbose bool) (parentCommitHashes []string) {
	parentCommitHashes, err := g.GetParentCommitHashesAsString(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return parentCommitHashes
}

func (g *GitlabCommit) MustGetParentCommits(verbose bool) (parentCommits []*GitlabCommit) {
	parentCommits, err := g.GetParentCommits(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return parentCommits
}

func (g *GitlabCommit) MustGetProjectId() (projectId int) {
	projectId, err := g.GetProjectId()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return projectId
}

func (g *GitlabCommit) MustGetRawResponse() (rawResponse *gitlab.Commit) {
	rawResponse, err := g.GetRawResponse()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return rawResponse
}

func (g *GitlabCommit) MustSetCommitHash(commitHash string) {
	err := g.SetCommitHash(commitHash)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabCommit) MustSetGitlabProjectsCommits(gitlabProjectsCommit *GitlabProjectCommits) {
	err := g.SetGitlabProjectsCommits(gitlabProjectsCommit)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabCommit) SetCommitHash(commitHash string) (err error) {
	if commitHash == "" {
		return TracedErrorf("commitHash is empty string")
	}

	g.commitHash = commitHash

	return nil
}

func (g *GitlabCommit) SetGitlabProjectsCommits(gitlabProjectsCommit *GitlabProjectCommits) (err error) {
	if gitlabProjectsCommit == nil {
		return TracedErrorf("gitlabProjectsCommit is nil")
	}

	g.gitlabProjectsCommits = gitlabProjectsCommit

	return nil
}
