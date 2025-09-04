package asciichgolangpublic

import (
	"context"
	"errors"

	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"github.com/asciich/asciichgolangpublic/pkg/versionutils"

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

func (g *GitlabTag) Delete(ctx context.Context) (err error) {
	name, err := g.GetName()
	if err != nil {
		return err
	}

	projectId, projectUrl, err := g.GetProjectIdAndUrl(ctx)
	if err != nil {
		return err
	}

	exists, err := g.Exists(ctx)
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
			return tracederrors.TracedErrorf("Delete tag '%s' in gitlab project %s failed: %w", name, projectUrl, err)
		}

		logging.LogChangedByCtxf(ctx, "Deleted tag '%s' in gitlab project %s .", name, projectUrl)
	} else {
		logging.LogInfoByCtxf(ctx, "Tag '%s' in gitlab project %s already absent. Skip delete.", name, projectUrl)
	}

	return nil
}

func (g *GitlabTag) Exists(ctx context.Context) (exists bool, err error) {
	projectUrl, err := g.GetProjectUrl(ctx)
	if err != nil {
		return false, err
	}

	exists = true
	_, err = g.GetRawResponse(ctx)
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

	if exists {
		logging.LogInfoByCtxf(ctx, "Tag '%s' in gitlab project %s exists.", tagName, projectUrl)
	} else {
		logging.LogInfoByCtxf(ctx, "Tag '%s' in gitlab project %s does not exist.", tagName, projectUrl)
	}

	return exists, nil
}

func (g *GitlabTag) GetGitRepository() (gitRepo gitinterfaces.GitRepository, err error) {
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

func (g *GitlabTag) GetHash(ctx context.Context) (hash string, err error) {
	raw, err := g.GetRawResponse(ctx)
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

func (g *GitlabTag) GetProjectId(ctx context.Context) (projectId int, err error) {
	project, err := g.GetGitlabProject()
	if err != nil {
		return -1, err
	}

	projectId, err = project.GetId(ctx)
	if err != nil {
		return -1, err
	}

	return projectId, nil
}

func (g *GitlabTag) GetProjectIdAndUrl(ctx context.Context) (projectId int, projectUrl string, err error) {
	projectId, err = g.GetProjectId(ctx)
	if err != nil {
		return -1, "", err
	}

	projectUrl, err = g.GetProjectUrl(ctx)
	if err != nil {
		return -1, "", err
	}

	return projectId, projectUrl, err
}

func (g *GitlabTag) GetProjectUrl(ctx context.Context) (projectUrl string, err error) {
	gitlabProject, err := g.GetGitlabProject()
	if err != nil {
		return "", err
	}

	gitlabUrl, err := gitlabProject.GetProjectUrl(ctx)
	if err != nil {
		return "", err
	}

	return gitlabUrl, nil
}

func (g *GitlabTag) GetRawResponse(ctx context.Context) (rawResponse *gitlab.Tag, err error) {
	nativeClient, err := g.GetNativeTagsService()
	if err != nil {
		return nil, err
	}

	projectId, projectUrl, err := g.GetProjectIdAndUrl(ctx)
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

	return versionutils.IsVersionString(tagName), nil
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
