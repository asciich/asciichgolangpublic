package asciichgolangpublic

import (
	"github.com/xanzy/go-gitlab"
)

type GitlabTags struct {
	gitlabProject *GitlabProject
}

func NewGitlabTags() (g *GitlabTags) {
	return new(GitlabTags)
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
		return nil, TracedErrorf("gitlabProject not set")
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

func (g *GitlabTags) GetTagNames(verbose bool) (tagNames []string, err error) {
	tags, err := g.GetTags(verbose)
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

func (g *GitlabTags) GetTags(verbose bool) (gitlabTags []*GitlabTag, err error) {
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
			return nil, TracedErrorf("Unable to get gitlab native tag list: '%w'", err)
		}

		nativeList = append(nativeList, tags...)

		if response.NextPage > 0 {
			break
		}

		pageNumber = response.NextPage
	}

	gitlabProject, err := g.GetGitlabProject()
	if err != nil {
		return nil, err
	}

	gitlabTags = []*GitlabTag{}
	for _, nativeTag := range nativeList {
		toAdd := NewGitlabTag()

		err = toAdd.SetName(nativeTag.Name)
		if err != nil {
			return nil, err
		}

		err = toAdd.SetGitlabProject(gitlabProject)
		if err != nil {
			return nil, err
		}

		gitlabTags = append(gitlabTags, toAdd)
	}

	return gitlabTags, nil
}

func (g *GitlabTags) GetVersionTagNames(verbose bool) (tagNames []string, err error) {
	allTagNames, err := g.GetTagNames(verbose)
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

func (g *GitlabTags) GetVersionTags(verbose bool) (versionTags []*GitlabTag, err error) {
	tags, err := g.GetTags(verbose)
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

func (g *GitlabTags) MustGetGitlab() (gitlab *GitlabInstance) {
	gitlab, err := g.GetGitlab()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlab
}

func (g *GitlabTags) MustGetGitlabProject() (gitlabProject *GitlabProject) {
	gitlabProject, err := g.GetGitlabProject()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlabProject
}

func (g *GitlabTags) MustGetNativeTagsService() (nativeTagsService *gitlab.TagsService) {
	nativeTagsService, err := g.GetNativeTagsService()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return nativeTagsService
}

func (g *GitlabTags) MustGetProjectId() (projectId int) {
	projectId, err := g.GetProjectId()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return projectId
}

func (g *GitlabTags) MustGetTagNames(verbose bool) (tagNames []string) {
	tagNames, err := g.GetTagNames(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return tagNames
}

func (g *GitlabTags) MustGetTags(verbose bool) (gitlabTags []*GitlabTag) {
	gitlabTags, err := g.GetTags(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlabTags
}

func (g *GitlabTags) MustGetVersionTagNames(verbose bool) (tagNames []string) {
	tagNames, err := g.GetVersionTagNames(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return tagNames
}

func (g *GitlabTags) MustGetVersionTags(verbose bool) (versionTags []*GitlabTag) {
	versionTags, err := g.GetVersionTags(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return versionTags
}

func (g *GitlabTags) MustSetGitlabProject(gitlabProject *GitlabProject) {
	err := g.SetGitlabProject(gitlabProject)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabTags) SetGitlabProject(gitlabProject *GitlabProject) (err error) {
	if gitlabProject == nil {
		return TracedErrorf("gitlabProject is nil")
	}

	g.gitlabProject = gitlabProject

	return nil
}
