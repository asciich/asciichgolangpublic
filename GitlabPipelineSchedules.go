package asciichgolangpublic

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type GitlabPipelineSchedules struct {
	gitlabProject *GitlabProject
}

func NewGitlabPipelineSchedules() (g *GitlabPipelineSchedules) {
	return new(GitlabPipelineSchedules)
}

func (g *GitlabPipelineSchedules) GetGitlabProject() (gitlab *GitlabProject, err error) {
	if g.gitlabProject == nil {
		return nil, tracederrors.TracedError("gitlabProject not set")
	}

	return g.gitlabProject, nil
}

func (g *GitlabPipelineSchedules) SetGitlabProject(gitlabProject *GitlabProject) (err error) {
	if gitlabProject == nil {
		return tracederrors.TracedErrorNil("gitlabProject")
	}

	g.gitlabProject = gitlabProject

	return nil
}

func (g *GitlabPipelineSchedules) GetNativePipelineSchedulesClient() (nativeClient *gitlab.PipelineSchedulesService, err error) {
	gitlabProject, err := g.GetGitlabProject()
	if err != nil {
		return nil, err
	}

	return gitlabProject.GetNativePipelineSchedulesClient()
}

func (g *GitlabPipelineSchedules) GetProjectId(ctx context.Context) (projectId int, err error) {
	project, err := g.GetGitlabProject()
	if err != nil {
		return 0, err
	}

	return project.GetId(ctx)
}

func (g *GitlabPipelineSchedules) GetProjectUrl(ctx context.Context) (projectUrl string, err error) {
	p, err := g.GetGitlabProject()
	if err != nil {
		return "", err
	}

	return p.GetProjectUrl(ctx)
}

func (g *GitlabPipelineSchedules) GetPipelineScheduleById(id int) (pipelineSchedule *GitlabPipelineSchedule, err error) {
	if id <= 0 {
		return nil, tracederrors.TracedErrorf("invalid id = '%d'", id)
	}

	pipelineSchedule = NewGitlabPipelineSchedule()

	err = pipelineSchedule.SetId(id)
	if err != nil {
		return nil, err
	}

	gitlabProject, err := g.GetGitlabProject()
	if err != nil {
		return nil, err
	}

	err = pipelineSchedule.SetGitlabProject(gitlabProject)
	if err != nil {
		return nil, err
	}

	return pipelineSchedule, err
}

func (g *GitlabPipelineSchedules) ListScheduledPipelineNames(ctx context.Context) (scheduledPipelineNames []string, err error) {
	pipelines, err := g.ListPipelineSchedules(ctx)
	if err != nil {
		return nil, err
	}

	scheduledPipelineNames = []string{}
	for _, p := range pipelines {
		toAdd, err := p.GetCachedName()
		if err != nil {
			return nil, err
		}

		scheduledPipelineNames = append(scheduledPipelineNames, toAdd)
	}

	return scheduledPipelineNames, nil
}

func (g *GitlabPipelineSchedules) MustListScheduledPipelineNames(ctx context.Context) (scheduledPipelineNames []string) {
	scheduledPipelineNames, err := g.ListScheduledPipelineNames(ctx)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return
}

func (g *GitlabPipelineSchedules) ListPipelineSchedules(ctx context.Context) (pipelineSchedules []*GitlabPipelineSchedule, err error) {
	nativeClient, err := g.GetNativePipelineSchedulesClient()
	if err != nil {
		return nil, err
	}

	projectId, err := g.GetProjectId(ctx)
	if err != nil {
		return nil, err
	}

	projectUrl, err := g.GetProjectUrl(ctx)
	if err != nil {
		return nil, err
	}

	pipelineSchedules = []*GitlabPipelineSchedule{}

	nextPage := 1
	for {
		if nextPage == 0 {
			break
		}

		rawList, response, err := nativeClient.ListPipelineSchedules(projectId, &gitlab.ListPipelineSchedulesOptions{Page: nextPage})
		if err != nil {
			return nil, tracederrors.TracedErrorf("Failed to load raw list of pipeline schedules in project %s: %w", projectUrl, err)
		}

		for _, r := range rawList {
			toAdd, err := g.GetPipelineScheduleById(r.ID)
			if err != nil {
				return nil, err
			}

			err = toAdd.SetCachedName(r.Description)
			if err != nil {
				return nil, err
			}

			pipelineSchedules = append(pipelineSchedules, toAdd)
		}

		nextPage = response.NextPage
	}

	logging.LogInfoByCtxf(ctx, "Collected '%d' scheduled pipelines for gitlab project '%s'", len(pipelineSchedules), projectUrl)

	return pipelineSchedules, nil
}
