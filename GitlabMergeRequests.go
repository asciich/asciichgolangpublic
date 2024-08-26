package asciichgolangpublic

import (
	"errors"

	"github.com/xanzy/go-gitlab"
)

var ErrNoMergeRequestWithTitleFound = errors.New("No merge request with given title found")

type GitlabMergeRequests struct {
	gitlabProject *GitlabProject
}

func NewGitlabMergeRequests() (g *GitlabMergeRequests) {
	return new(GitlabMergeRequests)
}

func (g *GitlabMergeRequests) CreateMergeRequest(options *GitlabCreateMergeRequestOptions) (createdMergeRequest *GitlabMergeRequest, err error) {
	if options == nil {
		return nil, TracedErrorNil("options")
	}

	nativeMergeRequests, err := g.GetNativeMergeRequestsService()
	if err != nil {
		return nil, err
	}

	projectId, err := g.GetProjectId()
	if err != nil {
		return nil, err
	}

	title, err := options.GetTitle()
	if err != nil {
		return nil, err
	}

	sourceBranch, err := options.GetSourceBranchName()
	if err != nil {
		return nil, err
	}

	targetBranch := ""
	if options.IsTargetBranchSet() {
		targetBranch, err = options.GetTargetBranchName()
		if err != nil {
			return nil, err
		}
	} else {
		targetBranch, err = g.GetDefaultBranchName()
		if err != nil {
			return nil, err
		}
	}

	projectUrl, err := g.GetProjectUrlAsString()
	if err != nil {
		return nil, err
	}

	createdMergeRequest, err = g.GetOpenMergeRequestByTitleOrNilIfNotPresent(title, options.Verbose)
	if err != nil {
		return nil, err
	}

	if createdMergeRequest != nil {
		url, err := createdMergeRequest.GetUrlAsString()
		if err != nil {
			return nil, err
		}

		if options.Verbose {
			LogChangedf(
				"Merge request '%s' already exists: %s .",
				title,
				url,
			)
		}
	} else {
		nativeMergeRequest, _, err := nativeMergeRequests.CreateMergeRequest(
			projectId,
			&gitlab.CreateMergeRequestOptions{
				Title:        &title,
				TargetBranch: &targetBranch,
				SourceBranch: &sourceBranch,
			},
		)
		if err != nil {
			return nil, TracedErrorf(
				"Create gitlab merge in project %s request failed: '%w'",
				projectUrl,
				err,
			)
		}

		createdMergeRequest, err = g.GetMergeRequestByNativeMergeRequest(nativeMergeRequest)
		if err != nil {
			return nil, err
		}

		url, err := createdMergeRequest.GetUrlAsString()
		if err != nil {
			return nil, err
		}

		if options.Verbose {
			LogChangedf(
				"Created merge request '%s' from branch '%s' to '%s': %s .",
				title,
				sourceBranch,
				targetBranch,
				url,
			)
		}

	}

	return createdMergeRequest, nil
}

func (g *GitlabMergeRequests) GetDefaultBranchName() (defaultBranchName string, err error) {
	gitlabProject, err := g.GetGitlabProject()
	if err != nil {
		return "", err
	}

	defaultBranchName, err = gitlabProject.GetDefaultBranchName()
	if err != nil {
		return "", err
	}

	return defaultBranchName, nil
}

func (g *GitlabMergeRequests) GetGitlab() (gitlab *GitlabInstance, err error) {
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

func (g *GitlabMergeRequests) GetGitlabProject() (gitlabProject *GitlabProject, err error) {
	if g.gitlabProject == nil {
		return nil, TracedErrorf("gitlabProject not set")
	}

	return g.gitlabProject, nil
}

func (g *GitlabMergeRequests) GetMergeRequestByNativeMergeRequest(nativeMergeRequest *gitlab.MergeRequest) (mergeRequest *GitlabMergeRequest, err error) {
	if nativeMergeRequest == nil {
		return nil, TracedErrorNil("nativeMergeRequest")
	}

	mergeRequest = NewGitlabMergeRequest()
	err = mergeRequest.SetGitlabMergeRequests(g)
	if err != nil {
		return nil, err
	}

	err = mergeRequest.SetId(nativeMergeRequest.ID)
	if err != nil {
		return nil, err
	}

	err = mergeRequest.SetCachedTitle(nativeMergeRequest.Title)
	if err != nil {
		return nil, err
	}

	return mergeRequest, nil
}

func (g *GitlabMergeRequests) GetNativeMergeRequestsService() (nativeService *gitlab.MergeRequestsService, err error) {
	gitlab, err := g.GetGitlab()
	if err != nil {
		return nil, err
	}

	nativeService, err = gitlab.GetNativeMergeRequestsService()
	if err != nil {
		return nil, err
	}

	return nativeService, nil
}

func (g *GitlabMergeRequests) GetOpenMergeRequestByTitle(title string, verbose bool) (mergeRequest *GitlabMergeRequest, err error) {
	if title == "" {
		return nil, TracedErrorEmptyString("title")
	}

	openMergeRequests, err := g.GetOpenMergeRequests(verbose)
	if err != nil {
		return nil, err
	}

	foundCounter := 0
	for _, request := range openMergeRequests {
		currentTilte, err := request.GetCachedTitle()
		if err != nil {
			return nil, err
		}

		if currentTilte == title {
			mergeRequest = request
			foundCounter += 1
		}
	}

	projectUrl, err := g.GetProjectUrlAsString()
	if err != nil {
		return nil, err
	}

	if foundCounter <= 0 {
		return nil, TracedErrorf("%w: '%s' in project %s .", ErrNoMergeRequestWithTitleFound, title, projectUrl)
	} else if foundCounter > 1 {
		return nil, TracedErrorf(
			"Found '%d' merge requests matching title '%s' in project %s but only 1 is supported.",
			foundCounter,
			title,
			projectUrl,
		)
	} else {
		if verbose {
			LogInfof("Found merge request by title '%s': %s", title, projectUrl)
		}
	}

	return mergeRequest, nil
}

func (g *GitlabMergeRequests) GetOpenMergeRequestByTitleOrNilIfNotPresent(title string, verbose bool) (mergeRequest *GitlabMergeRequest, err error) {
	mergeRequest, err = g.GetOpenMergeRequestByTitle(title, verbose)
	if err != nil {
		if errors.Is(err, ErrNoMergeRequestWithTitleFound) {
			return nil, nil
		}
		return nil, err
	}

	return mergeRequest, nil
}

func (g *GitlabMergeRequests) GetOpenMergeRequests(verbose bool) (openMergeRequest []*GitlabMergeRequest, err error) {
	var stateStringOpen string = "opened"

	rawMergeRequest, err := g.GetRawMergeRequests(&gitlab.ListProjectMergeRequestsOptions{
		State: &stateStringOpen,
	})
	if err != nil {
		return nil, err
	}

	openMergeRequest = []*GitlabMergeRequest{}
	for _, request := range rawMergeRequest {
		toAdd, err := g.GetMergeRequestByNativeMergeRequest(request)
		if err != nil {
			return nil, err
		}

		openMergeRequest = append(openMergeRequest, toAdd)
	}

	return openMergeRequest, nil
}

func (g *GitlabMergeRequests) GetProjectId() (projectId int, err error) {
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

func (g *GitlabMergeRequests) GetProjectUrlAsString() (projectUrl string, err error) {
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

func (g *GitlabMergeRequests) GetRawMergeRequests(options *gitlab.ListProjectMergeRequestsOptions) (rawMergeRequests []*gitlab.MergeRequest, err error) {
	if options == nil {
		return nil, TracedErrorNil("options")
	}

	projectId, err := g.GetProjectId()
	if err != nil {
		return nil, err
	}

	nativeService, err := g.GetNativeMergeRequestsService()
	if err != nil {
		return nil, err
	}
	rawMergeRequests, _, err = nativeService.ListProjectMergeRequests(
		projectId,
		options,
	)
	if err != nil {
		return nil, TracedErrorf("Get raw merge requests failed: '%s'", err)
	}

	return rawMergeRequests, nil
}

func (g *GitlabMergeRequests) MustCreateMergeRequest(options *GitlabCreateMergeRequestOptions) (createdMergeRequest *GitlabMergeRequest) {
	createdMergeRequest, err := g.CreateMergeRequest(options)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return createdMergeRequest
}

func (g *GitlabMergeRequests) MustGetDefaultBranchName() (defaultBranchName string) {
	defaultBranchName, err := g.GetDefaultBranchName()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return defaultBranchName
}

func (g *GitlabMergeRequests) MustGetGitlab() (gitlab *GitlabInstance) {
	gitlab, err := g.GetGitlab()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlab
}

func (g *GitlabMergeRequests) MustGetGitlabProject() (gitlabProject *GitlabProject) {
	gitlabProject, err := g.GetGitlabProject()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlabProject
}

func (g *GitlabMergeRequests) MustGetMergeRequestByNativeMergeRequest(nativeMergeRequest *gitlab.MergeRequest) (mergeRequest *GitlabMergeRequest) {
	mergeRequest, err := g.GetMergeRequestByNativeMergeRequest(nativeMergeRequest)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return mergeRequest
}

func (g *GitlabMergeRequests) MustGetNativeMergeRequestsService() (nativeService *gitlab.MergeRequestsService) {
	nativeService, err := g.GetNativeMergeRequestsService()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return nativeService
}

func (g *GitlabMergeRequests) MustGetOpenMergeRequestByTitle(title string, verbose bool) (mergeRequest *GitlabMergeRequest) {
	mergeRequest, err := g.GetOpenMergeRequestByTitle(title, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return mergeRequest
}

func (g *GitlabMergeRequests) MustGetOpenMergeRequestByTitleOrNilIfNotPresent(title string, verbose bool) (mergeRequest *GitlabMergeRequest) {
	mergeRequest, err := g.GetOpenMergeRequestByTitleOrNilIfNotPresent(title, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return mergeRequest
}

func (g *GitlabMergeRequests) MustGetOpenMergeRequests(verbose bool) (openMergeRequest []*GitlabMergeRequest) {
	openMergeRequest, err := g.GetOpenMergeRequests(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return openMergeRequest
}

func (g *GitlabMergeRequests) MustGetProjectId() (projectId int) {
	projectId, err := g.GetProjectId()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return projectId
}

func (g *GitlabMergeRequests) MustGetProjectUrlAsString() (projectUrl string) {
	projectUrl, err := g.GetProjectUrlAsString()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return projectUrl
}

func (g *GitlabMergeRequests) MustGetRawMergeRequests(options *gitlab.ListProjectMergeRequestsOptions) (rawMergeRequests []*gitlab.MergeRequest) {
	rawMergeRequests, err := g.GetRawMergeRequests(options)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return rawMergeRequests
}

func (g *GitlabMergeRequests) MustSetGitlabProject(gitlabProject *GitlabProject) {
	err := g.SetGitlabProject(gitlabProject)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabMergeRequests) SetGitlabProject(gitlabProject *GitlabProject) (err error) {
	if gitlabProject == nil {
		return TracedErrorf("gitlabProject is nil")
	}

	g.gitlabProject = gitlabProject

	return nil
}
