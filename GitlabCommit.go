package asciichgolangpublic

import (
	"context"

	gitlab "gitlab.com/gitlab-org/api/client-go"

	"github.com/asciich/asciichgolangpublic/pkg/datatypes/slicesutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type GitlabCommit struct {
	gitlabProjectsCommits *GitlabProjectCommits
	commitHash            string
}

func NewGitlabCommit() (g *GitlabCommit) {
	return new(GitlabCommit)
}

func (g *GitlabCommit) CreateRelease(ctx context.Context, createReleaseOptions *GitlabCreateReleaseOptions) (createdRelease *GitlabRelease, err error) {
	if createReleaseOptions == nil {
		return nil, tracederrors.TracedErrorNil("createReleaseOptions")
	}

	releaseName, err := createReleaseOptions.GetName()
	if err != nil {
		return nil, err
	}

	createTagOptions := &GitlabCreateTagOptions{
		Name: releaseName,
	}

	createdTag, err := g.CreateTag(ctx, createTagOptions)
	if err != nil {
		return nil, err
	}

	createdRelease, err = createdTag.CreateRelease(ctx, createReleaseOptions)
	if err != nil {
		return nil, err
	}

	return createdRelease, nil
}

func (g *GitlabCommit) CreateTag(ctx context.Context, createTagOptions *GitlabCreateTagOptions) (createdTag *GitlabTag, err error) {
	if createTagOptions == nil {
		return nil, tracederrors.TracedErrorNil("createTagOptions")
	}

	tags, err := g.GetGitlabTags()
	if err != nil {
		return nil, err
	}

	optionsToUse := createTagOptions.GetDeepCopy()

	commitHash, err := g.GetCommitHash()
	if err != nil {
		return nil, err
	}

	err = optionsToUse.SetRef(commitHash)
	if err != nil {
		return nil, err
	}

	createdTag, err = tags.CreateTag(ctx, optionsToUse)
	if err != nil {
		return nil, err
	}

	return createdTag, nil
}

func (g *GitlabCommit) GetAuthorEmail(ctx context.Context) (authorEmail string, err error) {
	rawResponse, err := g.GetRawResponse(ctx)
	if err != nil {
		return "", err
	}

	authorEmail = rawResponse.AuthorEmail

	commitHash, err := g.GetCommitHash()
	if err != nil {
		return "", err
	}

	projectUrl, err := g.GetGitlabProjectUrlAsString(ctx)
	if err != nil {
		return "", err
	}

	logging.LogInfoByCtxf(ctx, "Gitlab commit '%s' in %s has author email '%s'.", commitHash, projectUrl, authorEmail)

	return authorEmail, nil
}

func (g *GitlabCommit) GetCommitHash() (commitHash string, err error) {
	if g.commitHash == "" {
		return "", tracederrors.TracedErrorf("commitHash not set")
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

func (g *GitlabCommit) GetGitlabProjectUrlAsString(ctx context.Context) (projectUrl string, err error) {
	project, err := g.GetGitlabProject()
	if err != nil {
		return "", err
	}

	projectUrl, err = project.GetProjectUrl(ctx)
	if err != nil {
		return "", err
	}

	return projectUrl, nil
}

func (g *GitlabCommit) GetGitlabProjectsCommits() (gitlabProjectsCommit *GitlabProjectCommits, err error) {
	if g.gitlabProjectsCommits == nil {
		return nil, tracederrors.TracedErrorf("gitlabProjectsCommit not set")
	}

	return g.gitlabProjectsCommits, nil
}

func (g *GitlabCommit) GetGitlabTags() (gitlabTags *GitlabTags, err error) {
	gitlabProject, err := g.GetGitlabProject()
	if err != nil {
		return nil, err
	}

	gitlabTags, err = gitlabProject.GetTags()
	if err != nil {
		return nil, err
	}

	return gitlabTags, nil
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

func (g *GitlabCommit) GetParentCommitHashesAsString(ctx context.Context) (parentCommitHashes []string, err error) {
	rawResponse, err := g.GetRawResponse(ctx)
	if err != nil {
		return nil, err
	}

	parentCommitHashes = rawResponse.ParentIDs

	hash, err := g.GetCommitHash()
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Commit '%s' has parent commit hashes '%v'.", hash, parentCommitHashes)

	return parentCommitHashes, nil
}

func (g *GitlabCommit) GetParentCommits(ctx context.Context) (parentCommits []*GitlabCommit, err error) {
	parentCommitHashes, err := g.GetParentCommitHashesAsString(ctx)
	if err != nil {
		return nil, err
	}

	projectCommits, err := g.GetGitlabProjectsCommits()
	if err != nil {
		return nil, err
	}

	parentCommits = []*GitlabCommit{}
	for _, hash := range parentCommitHashes {
		toAdd, err := projectCommits.GetCommitByHashString(ctx, hash)
		if err != nil {
			return nil, err
		}

		parentCommits = append(parentCommits, toAdd)
	}

	return parentCommits, nil
}

func (g *GitlabCommit) GetProjectId(ctx context.Context) (projectId int, err error) {
	gitlabProject, err := g.GetGitlabProject()
	if err != nil {
		return -1, err
	}

	projectId, err = gitlabProject.GetId(ctx)
	if err != nil {
		return -1, err
	}

	return projectId, nil
}

func (g *GitlabCommit) GetRawResponse(ctx context.Context) (rawResponse *gitlab.Commit, err error) {
	nativeCommitsService, err := g.GetNativeCommitsService()
	if err != nil {
		return nil, err
	}

	hash, err := g.GetCommitHash()
	if err != nil {
		return nil, err
	}

	projectId, err := g.GetProjectId(ctx)
	if err != nil {
		return nil, err
	}

	rawResponse, _, err = nativeCommitsService.GetCommit(projectId, hash, &gitlab.GetCommitOptions{})
	if err != nil {
		return nil, tracederrors.TracedErrorf("Unable to get commit: %w", err)
	}

	return rawResponse, nil
}

func (g *GitlabCommit) IsMergeCommit(ctx context.Context) (isMergeCommit bool, err error) {
	parentCommits, err := g.GetParentCommits(ctx)
	if err != nil {
		return false, err
	}

	isMergeCommit = len(parentCommits) > 1
	projectUrl, err := g.GetGitlabProjectUrlAsString(ctx)
	if err != nil {
		return false, err
	}

	commitSha, err := g.GetCommitHash()
	if err != nil {
		return false, err
	}

	if isMergeCommit {
		logging.LogInfof(
			"Commit '%s' of gitlab project %s is a merge commit",
			projectUrl,
			commitSha,
		)
	}

	return isMergeCommit, nil
}

func (g *GitlabCommit) IsParentCommitOf(ctx context.Context, childCommit *GitlabCommit) (isParent bool, err error) {
	if childCommit == nil {
		return false, tracederrors.TracedErrorNil("childCommit")
	}

	parentHashes, err := childCommit.GetParentCommitHashesAsString(ctx)
	if err != nil {
		return false, err
	}

	hash, err := g.GetCommitHash()
	if err != nil {
		return false, err
	}

	isParent = slicesutils.ContainsStringIgnoreCase(
		parentHashes,
		hash,
	)

	projectUrl, err := g.GetGitlabProjectUrlAsString(ctx)
	if err != nil {
		return false, err
	}

	childHash, err := childCommit.GetCommitHash()
	if err != nil {
		return false, err
	}

	if isParent {
		logging.LogInfof(
			"Commit '%s' is parent of '%s' in gitlab project %s .",
			hash,
			childHash,
			projectUrl,
		)
	} else {
		logging.LogInfof(
			"Commit '%s' is not parent of '%s' in gitlab project %s .",
			hash,
			childHash,
			projectUrl,
		)
	}

	return isParent, nil
}

func (g *GitlabCommit) SetCommitHash(commitHash string) (err error) {
	if commitHash == "" {
		return tracederrors.TracedErrorf("commitHash is empty string")
	}

	g.commitHash = commitHash

	return nil
}

func (g *GitlabCommit) SetGitlabProjectsCommits(gitlabProjectsCommit *GitlabProjectCommits) (err error) {
	if gitlabProjectsCommit == nil {
		return tracederrors.TracedErrorf("gitlabProjectsCommit is nil")
	}

	g.gitlabProjectsCommits = gitlabProjectsCommit

	return nil
}
