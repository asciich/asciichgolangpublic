package asciichgolangpublic

import (
	gitlab "gitlab.com/gitlab-org/api/client-go"

	"github.com/asciich/asciichgolangpublic/datatypes/slicesutils"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type GitlabCommit struct {
	gitlabProjectsCommits *GitlabProjectCommits
	commitHash            string
}

func NewGitlabCommit() (g *GitlabCommit) {
	return new(GitlabCommit)
}

func (g *GitlabCommit) CreateRelease(createReleaseOptions *GitlabCreateReleaseOptions) (createdRelease *GitlabRelease, err error) {
	if createReleaseOptions == nil {
		return nil, tracederrors.TracedErrorNil("createReleaseOptions")
	}

	releaseName, err := createReleaseOptions.GetName()
	if err != nil {
		return nil, err
	}

	createTagOptions := &GitlabCreateTagOptions{
		Name:    releaseName,
		Verbose: createReleaseOptions.Verbose,
	}

	createdTag, err := g.CreateTag(createTagOptions)
	if err != nil {
		return nil, err
	}

	createdRelease, err = createdTag.CreateRelease(createReleaseOptions)
	if err != nil {
		return nil, err
	}

	return createdRelease, nil
}

func (g *GitlabCommit) CreateTag(createTagOptions *GitlabCreateTagOptions) (createdTag *GitlabTag, err error) {
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

	createdTag, err = tags.CreateTag(optionsToUse)
	if err != nil {
		return nil, err
	}

	return createdTag, nil
}

func (g *GitlabCommit) GetAuthorEmail(verbose bool) (authorEmail string, err error) {
	rawResponse, err := g.GetRawResponse()
	if err != nil {
		return "", err
	}

	authorEmail = rawResponse.AuthorEmail

	if verbose {
		commitHash, err := g.GetCommitHash()
		if err != nil {
			return "", err
		}

		projectUrl, err := g.GetGitlabProjectUrlAsString()
		if err != nil {
			return "", err
		}

		logging.LogInfof(
			"Gitlab commit '%s' in %s has author email '%s'.",
			commitHash,
			projectUrl,
			authorEmail,
		)
	}

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

func (g *GitlabCommit) GetGitlabProjectUrlAsString() (projectUrl string, err error) {
	project, err := g.GetGitlabProject()
	if err != nil {
		return "", err
	}

	projectUrl, err = project.GetProjectUrl()
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

		logging.LogInfof("Commit '%s' has parent commit hashes '%v'.", hash, parentCommitHashes)
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
		return nil, tracederrors.TracedErrorf("Unable to get commit: %w", err)
	}

	return rawResponse, nil
}

func (g *GitlabCommit) IsMergeCommit(verbose bool) (isMergeCommit bool, err error) {
	parentCommits, err := g.GetParentCommits(verbose)
	if err != nil {
		return false, err
	}

	isMergeCommit = len(parentCommits) > 1
	if verbose {
		projectUrl, err := g.GetGitlabProjectUrlAsString()
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
	}

	return isMergeCommit, nil
}

func (g *GitlabCommit) IsParentCommitOf(childCommit *GitlabCommit, verbose bool) (isParent bool, err error) {
	if childCommit == nil {
		return false, tracederrors.TracedErrorNil("childCommit")
	}

	parentHashes, err := childCommit.GetParentCommitHashesAsString(verbose)
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

	if verbose {
		projectUrl, err := g.GetGitlabProjectUrlAsString()
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
	}

	return isParent, nil
}

func (g *GitlabCommit) MustCreateRelease(createReleaseOptions *GitlabCreateReleaseOptions) (createdRelease *GitlabRelease) {
	createdRelease, err := g.CreateRelease(createReleaseOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return createdRelease
}

func (g *GitlabCommit) MustCreateTag(createTagOptions *GitlabCreateTagOptions) (createdTag *GitlabTag) {
	createdTag, err := g.CreateTag(createTagOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return createdTag
}

func (g *GitlabCommit) MustGetAuthorEmail(verbose bool) (authorEmail string) {
	authorEmail, err := g.GetAuthorEmail(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return authorEmail
}

func (g *GitlabCommit) MustGetCommitHash() (commitHash string) {
	commitHash, err := g.GetCommitHash()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return commitHash
}

func (g *GitlabCommit) MustGetGitlab() (gitlab *GitlabInstance) {
	gitlab, err := g.GetGitlab()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlab
}

func (g *GitlabCommit) MustGetGitlabProject() (gitlabProject *GitlabProject) {
	gitlabProject, err := g.GetGitlabProject()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlabProject
}

func (g *GitlabCommit) MustGetGitlabProjectUrlAsString() (projectUrl string) {
	projectUrl, err := g.GetGitlabProjectUrlAsString()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return projectUrl
}

func (g *GitlabCommit) MustGetGitlabProjectsCommits() (gitlabProjectsCommit *GitlabProjectCommits) {
	gitlabProjectsCommit, err := g.GetGitlabProjectsCommits()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlabProjectsCommit
}

func (g *GitlabCommit) MustGetGitlabTags() (gitlabTags *GitlabTags) {
	gitlabTags, err := g.GetGitlabTags()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlabTags
}

func (g *GitlabCommit) MustGetNativeCommitsService() (nativeCommitsService *gitlab.CommitsService) {
	nativeCommitsService, err := g.GetNativeCommitsService()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return nativeCommitsService
}

func (g *GitlabCommit) MustGetParentCommitHashesAsString(verbose bool) (parentCommitHashes []string) {
	parentCommitHashes, err := g.GetParentCommitHashesAsString(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return parentCommitHashes
}

func (g *GitlabCommit) MustGetParentCommits(verbose bool) (parentCommits []*GitlabCommit) {
	parentCommits, err := g.GetParentCommits(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return parentCommits
}

func (g *GitlabCommit) MustGetProjectId() (projectId int) {
	projectId, err := g.GetProjectId()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return projectId
}

func (g *GitlabCommit) MustGetRawResponse() (rawResponse *gitlab.Commit) {
	rawResponse, err := g.GetRawResponse()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return rawResponse
}

func (g *GitlabCommit) MustIsMergeCommit(verbose bool) (isMergeCommit bool) {
	isMergeCommit, err := g.IsMergeCommit(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return isMergeCommit
}

func (g *GitlabCommit) MustIsParentCommitOf(childCommit *GitlabCommit, verbose bool) (isParent bool) {
	isParent, err := g.IsParentCommitOf(childCommit, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return isParent
}

func (g *GitlabCommit) MustSetCommitHash(commitHash string) {
	err := g.SetCommitHash(commitHash)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabCommit) MustSetGitlabProjectsCommits(gitlabProjectsCommit *GitlabProjectCommits) {
	err := g.SetGitlabProjectsCommits(gitlabProjectsCommit)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
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
