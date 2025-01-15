package asciichgolangpublic

import (
	"errors"

	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"

	gitlab "gitlab.com/gitlab-org/api/client-go"
)

var ErrGitlabTagNotFound = errors.New("gitlab tag not found")

type GitlabTag struct {
	GitTagBase
	gitlabTags *GitlabTags
	name       string
}

func NewGitlabTag() (g *GitlabTag) {
	g = new(GitlabTag)

	g.MustSetParentGitTagForBaseClass(g)

	return g
}

func (g *GitlabTag) Delete(verbose bool) (err error) {
	name, err := g.GetName()
	if err != nil {
		return err
	}

	projectId, projectUrl, err := g.GetProjectIdAndUrl()
	if err != nil {
		return err
	}

	exists, err := g.Exists(verbose)
	if err != nil {
		return err
	}

	if exists {
		nativeClient, err := g.GetNativeTagsService()
		if err != nil {
			return err
		}

		_, err = nativeClient.DeleteTag(
			projectId,
			name,
			nil,
		)
		if err != nil {
			return tracederrors.TracedErrorf(
				"Delete tag '%s' in gitlab project %s failed: %w",
				name,
				projectUrl,
				err,
			)
		}

		if verbose {
			logging.LogChangedf(
				"Deleted tag '%s' in gitlab project %s .",
				name,
				projectUrl,
			)
		}
	} else {
		if verbose {
			logging.LogInfof(
				"Tag '%s' in gitlab project %s already absent. Skip delete.",
				name,
				projectUrl,
			)
		}
	}

	return nil
}

func (g *GitlabTag) Exists(verbose bool) (exists bool, err error) {
	projectUrl, err := g.GetProjectUrl()
	if err != nil {
		return false, err
	}

	exists = true
	_, err = g.GetRawResponse()
	if err != nil {
		if errors.Is(err, ErrGitlabTagNotFound) {
			exists = false
		} else {
			return false, err
		}
	}

	tagName, err := g.GetName()
	if err != nil {
		return false, err
	}

	if verbose {
		if exists {
			logging.LogInfof(
				"Tag '%s' in gitlab project %s exists.",
				tagName,
				projectUrl,
			)
		} else {
			logging.LogInfof(
				"Tag '%s' in gitlab project %s does not exist.",
				tagName,
				projectUrl,
			)
		}
	}

	return exists, nil
}

func (g *GitlabTag) GetGitRepository() (gitRepo GitRepository, err error) {
	// TODO: This should return the gitlab project which
	// should implement everything a git repsository does so it
	// fullfils the GitRepository interface:
	return nil, tracederrors.TracedErrorNotImplemented()
}

func (g *GitlabTag) GetGitlabProject() (gitlabProject *GitlabProject, err error) {
	gitlabTags, err := g.GetGitlabTags()
	if err != nil {
		return nil, err
	}

	gitlabProject, err = gitlabTags.GetGitlabProject()
	if err != nil {
		return nil, err
	}

	return gitlabProject, nil
}

func (g *GitlabTag) GetGitlabTags() (gitlabTags *GitlabTags, err error) {
	if g.gitlabTags == nil {
		return nil, tracederrors.TracedErrorf("gitlabTags not set")
	}

	return g.gitlabTags, nil
}

func (g *GitlabTag) GetHash() (hash string, err error) {
	raw, err := g.GetRawResponse()
	if err != nil {
		return "", err
	}

	commit := raw.Commit
	if commit == nil {
		return "", tracederrors.TracedErrorf(
			"raw.Commit is nil after evaluation.",
		)
	}

	hash = commit.ID

	if hash == "" {
		return "", tracederrors.TracedErrorf(
			"hash is empty string after evaluation.",
		)
	}

	return hash, nil
}

func (g *GitlabTag) GetName() (name string, err error) {
	if g.name == "" {
		return "", tracederrors.TracedErrorf("name not set")
	}

	return g.name, nil
}

func (g *GitlabTag) GetNativeTagsService() (nativeTagsService *gitlab.TagsService, err error) {
	gitlabTags, err := g.GetGitlabTags()
	if err != nil {
		return nil, err
	}

	nativeTagsService, err = gitlabTags.GetNativeTagsService()
	if err != nil {
		return nil, err
	}

	return nativeTagsService, nil
}

func (g *GitlabTag) GetProjectId() (projectId int, err error) {
	project, err := g.GetGitlabProject()
	if err != nil {
		return -1, err
	}

	projectId, err = project.GetId()
	if err != nil {
		return -1, err
	}

	return projectId, nil
}

func (g *GitlabTag) GetProjectIdAndUrl() (projectId int, projectUrl string, err error) {
	projectId, err = g.GetProjectId()
	if err != nil {
		return -1, "", err
	}

	projectUrl, err = g.GetProjectUrl()
	if err != nil {
		return -1, "", err
	}

	return projectId, projectUrl, err
}

func (g *GitlabTag) GetProjectUrl() (projectUrl string, err error) {
	gitlabProject, err := g.GetGitlabProject()
	if err != nil {
		return "", err
	}

	gitlabUrl, err := gitlabProject.GetProjectUrl()
	if err != nil {
		return "", err
	}

	return gitlabUrl, nil
}

func (g *GitlabTag) GetRawResponse() (rawResponse *gitlab.Tag, err error) {
	nativeClient, err := g.GetNativeTagsService()
	if err != nil {
		return nil, err
	}

	projectId, projectUrl, err := g.GetProjectIdAndUrl()
	if err != nil {
		return nil, err
	}

	name, err := g.GetName()
	if err != nil {
		return nil, err
	}

	rawResponse, _, err = nativeClient.GetTag(
		projectId,
		name,
		nil,
	)
	if err != nil {
		if err.Error() == "404 Not Found" {
			return nil, tracederrors.TracedErrorf(
				"%w: tag '%s' in gitlab project %s: %w",
				ErrGitlabTagNotFound,
				name,
				projectUrl,
				err,
			)
		}

		return nil, tracederrors.TracedErrorf(
			"Get raw response for tag '%s' in gitlab project %s failed: %w",
			name,
			projectUrl,
			err,
		)
	}

	if rawResponse == nil {
		return nil, tracederrors.TracedError("rawResponse is nil after evaluation")
	}

	return rawResponse, nil
}

func (g *GitlabTag) IsVersionTag() (isVersionTag bool, err error) {
	tagName, err := g.GetName()
	if err != nil {
		return false, err
	}

	return Versions().IsVersionString(tagName), nil
}

func (g *GitlabTag) MustDelete(verbose bool) {
	err := g.Delete(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabTag) MustExists(verbose bool) (exists bool) {
	exists, err := g.Exists(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return exists
}

func (g *GitlabTag) MustGetGitRepository() (gitRepo GitRepository) {
	gitRepo, err := g.GetGitRepository()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitRepo
}

func (g *GitlabTag) MustGetGitlabProject() (gitlabProject *GitlabProject) {
	gitlabProject, err := g.GetGitlabProject()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlabProject
}

func (g *GitlabTag) MustGetGitlabTags() (gitlabTags *GitlabTags) {
	gitlabTags, err := g.GetGitlabTags()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlabTags
}

func (g *GitlabTag) MustGetHash() (hash string) {
	hash, err := g.GetHash()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return hash
}

func (g *GitlabTag) MustGetName() (name string) {
	name, err := g.GetName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return name
}

func (g *GitlabTag) MustGetNativeTagsService() (nativeTagsService *gitlab.TagsService) {
	nativeTagsService, err := g.GetNativeTagsService()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return nativeTagsService
}

func (g *GitlabTag) MustGetProjectId() (projectId int) {
	projectId, err := g.GetProjectId()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return projectId
}

func (g *GitlabTag) MustGetProjectIdAndUrl() (projectId int, projectUrl string) {
	projectId, projectUrl, err := g.GetProjectIdAndUrl()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return projectId, projectUrl
}

func (g *GitlabTag) MustGetProjectUrl() (projectUrl string) {
	projectUrl, err := g.GetProjectUrl()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return projectUrl
}

func (g *GitlabTag) MustGetRawResponse() (rawResponse *gitlab.Tag) {
	rawResponse, err := g.GetRawResponse()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return rawResponse
}

func (g *GitlabTag) MustIsVersionTag() (isVersionTag bool) {
	isVersionTag, err := g.IsVersionTag()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return isVersionTag
}

func (g *GitlabTag) MustSetGitlabTags(gitlabTags *GitlabTags) {
	err := g.SetGitlabTags(gitlabTags)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabTag) MustSetName(name string) {
	err := g.SetName(name)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabTag) SetGitlabTags(gitlabTags *GitlabTags) (err error) {
	if gitlabTags == nil {
		return tracederrors.TracedErrorf("gitlabTags is nil")
	}

	g.gitlabTags = gitlabTags

	return nil
}

func (g *GitlabTag) SetName(name string) (err error) {
	if name == "" {
		return tracederrors.TracedErrorf("name is empty string")
	}

	g.name = name

	return nil
}
