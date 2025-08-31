package asciichgolangpublic

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type GitlabPipelineSchedule struct {
	id            int
	cachedName    string
	gitlabProject *GitlabProject
}

func NewGitlabPipelineSchedule() (g *GitlabPipelineSchedule) {
	return new(GitlabPipelineSchedule)
}

func (g *GitlabPipelineSchedule) GetGitlabProject() (gitlab *GitlabProject, err error) {
	if g.gitlabProject == nil {
		return nil, tracederrors.TracedError("gitlab not set")
	}

	return g.gitlabProject, nil
}

func (g *GitlabPipelineSchedule) GetNativePipelineSchedulesClient() (nativeClient *gitlab.PipelineSchedulesService, err error) {
	project, err := g.GetGitlabProject()
	if err != nil {
		return nil, err
	}

	return project.GetNativePipelineSchedulesClient()
}

func (g *GitlabPipelineSchedule) SetId(id int) (err error) {
	if id <= 0 {
		return tracederrors.TracedErrorf("Invalid id: '%d'", id)
	}

	g.id = id

	return nil
}

func (g *GitlabPipelineSchedule) GetId() (id int, err error) {
	if g.id <= 0 {
		return 0, tracederrors.TracedError("id not set")
	}

	return g.id, nil
}

func (g *GitlabPipelineSchedule) GetProjectId(ctx context.Context) (projectId int, err error) {
	project, err := g.GetGitlabProject()
	if err != nil {
		return 0, err
	}

	return project.GetId(ctx)
}

func (g *GitlabPipelineSchedule) GetRawResponse(ctx context.Context) (rawResponse *gitlab.PipelineSchedule, err error) {
	native, err := g.GetNativePipelineSchedulesClient()
	if err != nil {
		return nil, err
	}

	projectId, err := g.GetProjectId(ctx)
	if err != nil {
		return nil, err
	}

	id, err := g.GetId()
	if err != nil {
		return nil, err
	}

	rawResponse, _, err = native.GetPipelineSchedule(projectId, id)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to get pipeline schedule: %w", err)
	}

	if rawResponse == nil {
		return nil, tracederrors.TracedError("rawResponse is nil after evaluation.")
	}

	return rawResponse, nil
}

func (g *GitlabPipelineSchedule) GetGitlabUrl(ctx context.Context) (url string, err error) {
	project, err := g.GetGitlabProject()
	if err != nil {
		return "", err
	}

	return project.GetProjectUrl(ctx)
}

func (g *GitlabPipelineSchedule) GetLastPipelineStatus(ctx context.Context) (status string, err error) {
	rawResponse, err := g.GetRawResponse(ctx)
	if err != nil {
		return "", err
	}

	lastPipeline := rawResponse.LastPipeline
	if lastPipeline == nil {
		return "never_run", nil
	}

	status = lastPipeline.Status
	if status == "" {
		return "", tracederrors.TracedErrorf("status is empty string after evaluation")
	}

	name, err := g.GetCachedName()
	if err != nil {
		return "", err
	}

	url, err := g.GetGitlabUrl(ctx)
	if err != nil {
		return "", err
	}

	logging.LogInfoByCtxf(ctx, "Last scheduled pipeline '%s' in project '%s' has status '%s'.", name, url, status)

	return status, nil
}

func (g *GitlabPipelineSchedule) GetCachedName() (cachedName string, err error) {
	if g.cachedName == "" {
		return "", tracederrors.TracedError("cachedName not set")
	}

	return g.cachedName, nil
}

func (g *GitlabPipelineSchedule) SetGitlabProject(gitlabProject *GitlabProject) (err error) {
	if gitlabProject == nil {
		return tracederrors.TracedErrorNil("gitlab")
	}

	g.gitlabProject = gitlabProject

	return nil
}

func (g *GitlabPipelineSchedule) SetCachedName(cachedName string) (err error) {
	if cachedName == "" {
		return tracederrors.TracedErrorEmptyString("cachedName")
	}

	g.cachedName = cachedName

	return nil
}
