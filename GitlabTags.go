package asciichgolangpublic

import (
	"github.com/asciich/asciichgolangpublic/errors"
	"github.com/asciich/asciichgolangpublic/logging"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type GitlabTags struct {
	gitlabProject *GitlabProject
}

func NewGitlabTags() (g *GitlabTags) {
	return new(GitlabTags)
}

func (g *GitlabTag) CreateRelease(createReleaseOptions *GitlabCreateReleaseOptions) (createdRelease *GitlabRelease, err error) {
	if createReleaseOptions == nil {
		return nil, errors.TracedErrorNil("createReleaseOptions")
	}

	releases, err := g.GetGitlabReleases()
	if err != nil {
		return nil, err
	}

	createdRelease, err = releases.CreateRelease(createReleaseOptions)
	if err != nil {
		return nil, err
	}

	return createdRelease, nil
}

func (g *GitlabTag) GetGitlabReleases() (gitlabReleases *GitlabReleases, err error) {
	gitlabProject, err := g.GetGitlabProject()
	if err != nil {
		return nil, err
	}

	gitlabReleases, err = gitlabProject.GetGitlabReleases()
	if err != nil {
		return nil, err
	}

	return gitlabReleases, nil
}

func (g *GitlabTag) MustCreateRelease(createReleaseOptions *GitlabCreateReleaseOptions) (createdRelease *GitlabRelease) {
	createdRelease, err := g.CreateRelease(createReleaseOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return createdRelease
}

func (g *GitlabTag) MustGetGitlabReleases() (gitlabReleases *GitlabReleases) {
	gitlabReleases, err := g.GetGitlabReleases()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlabReleases
}

func (g *GitlabTags) CreateTag(createTagOptions *GitlabCreateTagOptions) (createdTag *GitlabTag, err error) {
	if createTagOptions == nil {
		return nil, errors.TracedErrorNil("createTagOptions")
	}

	nativeClient, err := g.GetNativeTagsService()
	if err != nil {
		return nil, err
	}

	projectId, projectUrl, err := g.GetProjectIdAndUrl()
	if err != nil {
		return nil, err
	}

	tagName, err := createTagOptions.GetName()
	if err != nil {
		return nil, err
	}

	tagExists, err := g.TagByNameExists(tagName, createTagOptions.Verbose)
	if err != nil {
		return nil, err
	}

	if tagExists {
		if createTagOptions.Verbose {
			logging.LogInfof(
				"Tag '%s' in gitlab project %s already exists. Skip creation.",
				tagName,
				projectUrl,
			)
		}
	} else {
		ref, err := createTagOptions.GetRef()
		if err != nil {
			return nil, err
		}

		_, _, err = nativeClient.CreateTag(
			projectId,
			&gitlab.CreateTagOptions{
				TagName: &tagName,
				Ref:     &ref,
			},
		)
		if err != nil {
			return nil, errors.TracedErrorf(
				"Create tag '%s' in gitlab project %s failed: %w",
				tagName,
				projectUrl,
				err,
			)
		}

		if createTagOptions.Verbose {
			logging.LogChangedf(
				"Created tag '%s' on ref='%s' in gitlab project '%s'",
				tagName,
				ref,
				projectUrl,
			)
		}
	}

	createdTag, err = g.GetTagByName(tagName)
	if err != nil {
		return nil, err
	}

	return createdTag, nil
}

func (g *GitlabTags) GetGitlab() (gitlab *GitlabInstance, err error) {
	gitlabProject, err := g.GetGitlabProject()
	if err != nil {
		return nil, err
	}

	gitlab, err = gitlabProject.GetGitlab()
	if err != nil {
		return nil, err
	}

	return gitlab, nil
}

func (g *GitlabTags) GetGitlabProject() (gitlabProject *GitlabProject, err error) {
	if g.gitlabProject == nil {
		return nil, errors.TracedErrorf("gitlabProject not set")
	}

	return g.gitlabProject, nil
}

func (g *GitlabTags) GetNativeTagsService() (nativeTagsService *gitlab.TagsService, err error) {
	gitlab, err := g.GetGitlab()
	if err != nil {
		return nil, err
	}

	nativeTagsService, err = gitlab.GetNativeTagsService()
	if err != nil {
		return nil, err
	}

	return nativeTagsService, nil
}

func (g *GitlabTags) GetProjectId() (projectId int, err error) {
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

func (g *GitlabTags) GetProjectIdAndUrl() (projectId int, projectUrl string, err error) {
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

func (g *GitlabTags) GetProjectUrl() (projectUrl string, err error) {
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

func (g *GitlabTags) GetTagByName(tagName string) (tag *GitlabTag, err error) {
	if tagName == "" {
		return nil, errors.TracedErrorEmptyString("tagName")
	}

	tag = NewGitlabTag()

	err = tag.SetName(tagName)
	if err != nil {
		return nil, err
	}

	err = tag.SetGitlabTags(g)
	if err != nil {
		return nil, err
	}

	return tag, nil
}

func (g *GitlabTags) GetVersionTags(verbose bool) (versionTags []*GitlabTag, err error) {
	tags, err := g.ListTags(verbose)
	if err != nil {
		return nil, err
	}

	versionTags = []*GitlabTag{}
	for _, tag := range tags {
		isVersionTag, err := tag.IsVersionTag()
		if err != nil {
			return nil, err
		}

		if isVersionTag {
			versionTags = append(versionTags, tag)
		}
	}

	return versionTags, nil
}

func (g *GitlabTags) ListTagNames(verbose bool) (tagNames []string, err error) {
	tags, err := g.ListTags(verbose)
	if err != nil {
		return nil, err
	}

	tagNames = []string{}
	for _, tag := range tags {
		toAdd, err := tag.GetName()
		if err != nil {
			return nil, err
		}

		tagNames = append(tagNames, toAdd)
	}

	return tagNames, nil
}

func (g *GitlabTags) ListTags(verbose bool) (gitlabTags []*GitlabTag, err error) {
	nativeTagsService, err := g.GetNativeTagsService()
	if err != nil {
		return nil, err
	}

	projectId, err := g.GetProjectId()
	if err != nil {
		return nil, err
	}

	var nativeList []*gitlab.Tag
	pageNumber := 1
	for {
		listOptions := &gitlab.ListTagsOptions{}
		listOptions.Page = pageNumber

		tags, response, err := nativeTagsService.ListTags(
			projectId,
			listOptions,
		)
		if err != nil {
			return nil, errors.TracedErrorf("Unable to get gitlab native tag list: '%w'", err)
		}

		nativeList = append(nativeList, tags...)

		if response.NextPage <= 0 {
			break
		}

		pageNumber = response.NextPage
	}

	gitlabTags = []*GitlabTag{}
	for _, nativeTag := range nativeList {
		toAdd := NewGitlabTag()

		err = toAdd.SetName(nativeTag.Name)
		if err != nil {
			return nil, err
		}

		err = toAdd.SetGitlabTags(g)
		if err != nil {
			return nil, err
		}

		gitlabTags = append(gitlabTags, toAdd)
	}

	return gitlabTags, nil
}

func (g *GitlabTags) ListVersionTagNames(verbose bool) (tagNames []string, err error) {
	allTagNames, err := g.ListTagNames(verbose)
	if err != nil {
		return nil, err
	}

	tagNames = []string{}
	for _, toAdd := range allTagNames {
		if Versions().IsVersionString(toAdd) {
			tagNames = append(tagNames, toAdd)
		}
	}

	return tagNames, nil
}

func (g *GitlabTags) MustCreateTag(createTagOptions *GitlabCreateTagOptions) (createdTag *GitlabTag) {
	createdTag, err := g.CreateTag(createTagOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return createdTag
}

func (g *GitlabTags) MustGetGitlab() (gitlab *GitlabInstance) {
	gitlab, err := g.GetGitlab()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlab
}

func (g *GitlabTags) MustGetGitlabProject() (gitlabProject *GitlabProject) {
	gitlabProject, err := g.GetGitlabProject()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlabProject
}

func (g *GitlabTags) MustGetNativeTagsService() (nativeTagsService *gitlab.TagsService) {
	nativeTagsService, err := g.GetNativeTagsService()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return nativeTagsService
}

func (g *GitlabTags) MustGetProjectId() (projectId int) {
	projectId, err := g.GetProjectId()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return projectId
}

func (g *GitlabTags) MustGetProjectIdAndUrl() (projectId int, projectUrl string) {
	projectId, projectUrl, err := g.GetProjectIdAndUrl()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return projectId, projectUrl
}

func (g *GitlabTags) MustGetProjectUrl() (projectUrl string) {
	projectUrl, err := g.GetProjectUrl()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return projectUrl
}

func (g *GitlabTags) MustGetTagByName(tagName string) (tag *GitlabTag) {
	tag, err := g.GetTagByName(tagName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return tag
}

func (g *GitlabTags) MustGetVersionTags(verbose bool) (versionTags []*GitlabTag) {
	versionTags, err := g.GetVersionTags(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return versionTags
}

func (g *GitlabTags) MustListTagNames(verbose bool) (tagNames []string) {
	tagNames, err := g.ListTagNames(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return tagNames
}

func (g *GitlabTags) MustListTags(verbose bool) (gitlabTags []*GitlabTag) {
	gitlabTags, err := g.ListTags(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlabTags
}

func (g *GitlabTags) MustListVersionTagNames(verbose bool) (tagNames []string) {
	tagNames, err := g.ListVersionTagNames(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return tagNames
}

func (g *GitlabTags) MustSetGitlabProject(gitlabProject *GitlabProject) {
	err := g.SetGitlabProject(gitlabProject)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabTags) MustTagByNameExists(tagName string, verbose bool) (exists bool) {
	exists, err := g.TagByNameExists(tagName, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return exists
}

func (g *GitlabTags) SetGitlabProject(gitlabProject *GitlabProject) (err error) {
	if gitlabProject == nil {
		return errors.TracedErrorf("gitlabProject is nil")
	}

	g.gitlabProject = gitlabProject

	return nil
}

func (g *GitlabTags) TagByNameExists(tagName string, verbose bool) (exists bool, err error) {
	if tagName == "" {
		return false, errors.TracedErrorEmptyString("tagName")
	}

	tag, err := g.GetTagByName(tagName)
	if err != nil {
		return false, err
	}

	exists, err = tag.Exists(verbose)
	if err != nil {
		return false, err
	}

	return exists, nil
}
