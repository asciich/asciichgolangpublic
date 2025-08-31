package asciichgolangpublic

import (
	"context"
	"errors"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

var ErrNoMergeRequestWithTitleFound = errors.New("no merge request with given title found")
var ErrNoMergeRequestWithSourceAndTargetBranchFound = errors.New("no merge request with given source and target branch found")

// Handle Gitlab merge requests related to a project.
type GitlabProjectMergeRequests struct {
	gitlabProject *GitlabProject
}

func NewGitlabMergeRequests() (g *GitlabProjectMergeRequests) {
	return new(GitlabProjectMergeRequests)
}

func NewGitlabProjectMergeRequests() (g *GitlabProjectMergeRequests) {
	return new(GitlabProjectMergeRequests)
}

// Returns the `userId` of the currently logged in user.
func (g *GitlabProjectMergeRequests) GetUserId() (userId int, err error) {
	gitlabInstance, err := g.GetGitlab()
	if err != nil {
		return -1, err
	}

	userId, err = gitlabInstance.GetUserId()
	if err != nil {
		return -1, err
	}

	return userId, nil
}

func (g *GitlabProjectMergeRequests) CreateMergeRequest(ctx context.Context, options *GitlabCreateMergeRequestOptions) (createdMergeRequest *GitlabMergeRequest, err error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	nativeMergeRequests, err := g.GetNativeMergeRequestsService()
	if err != nil {
		return nil, err
	}

	projectId, err := g.GetProjectId(ctx)
	if err != nil {
		return nil, err
	}

	title, err := options.GetTitle()
	if err != nil {
		return nil, err
	}

	description := options.GetDescriptionOrEmptyStringIfUnset()

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
		targetBranch, err = g.GetDefaultBranchName(ctx)
		if err != nil {
			return nil, err
		}
	}

	projectUrl, err := g.GetProjectUrlAsString(ctx)
	if err != nil {
		return nil, err
	}

	createdMergeRequest, err = g.GetOpenMergeRequestByTitleOrNilIfNotPresent(ctx, title)
	if err != nil {
		return nil, err
	}

	if createdMergeRequest != nil {
		if options.GetFailIfMergeRequestAlreadyExists() {
			return nil, tracederrors.TracedErrorf(
				"Failed to create merge request: merge request with title '%s' already exists.",
				title,
			)
		}
	}

	labels := options.GetLabelsOrEmptySliceIfUnset()
	labelOptions := gitlab.LabelOptions(labels)

	squashEnabled := options.GetSquashEnabled()
	deleteSourceBranch := options.GetDeleteSourceBranchOnMerge()

	if createdMergeRequest != nil {
		url, err := createdMergeRequest.GetUrlAsString(ctx)
		if err != nil {
			return nil, err
		}

		logging.LogChangedByCtxf(ctx, "Merge request '%s' already exists: %s .", title, url)
	} else {
		assigneIds := []int{}

		if options.GetAssignToSelf() {
			userId, err := g.GetUserId()
			if err != nil {
				return nil, err
			}

			assigneIds = append(assigneIds, userId)
		}

		nativeMergeRequest, _, err := nativeMergeRequests.CreateMergeRequest(
			projectId,
			&gitlab.CreateMergeRequestOptions{
				Title:              &title,
				Description:        &description,
				TargetBranch:       &targetBranch,
				SourceBranch:       &sourceBranch,
				Labels:             &labelOptions,
				Squash:             &squashEnabled,
				RemoveSourceBranch: &deleteSourceBranch,
				AssigneeIDs:        &assigneIds,
			},
		)
		if err != nil {
			return nil, tracederrors.TracedErrorf(
				"Create gitlab merge in project %s request failed: '%w'",
				projectUrl,
				err,
			)
		}

		createdMergeRequest, err = g.GetMergeRequestByNativeMergeRequest(nativeMergeRequest)
		if err != nil {
			return nil, err
		}

		url, err := createdMergeRequest.GetUrlAsString(ctx)
		if err != nil {
			return nil, err
		}

		logging.LogChangedByCtxf(ctx, "Created merge request '%s' from branch '%s' to '%s': %s .", title, sourceBranch, targetBranch, url)

	}

	return createdMergeRequest, nil
}

func (g *GitlabProjectMergeRequests) GetDefaultBranchName(ctx context.Context) (defaultBranchName string, err error) {
	gitlabProject, err := g.GetGitlabProject()
	if err != nil {
		return "", err
	}

	defaultBranchName, err = gitlabProject.GetDefaultBranchName(ctx)
	if err != nil {
		return "", err
	}

	return defaultBranchName, nil
}

func (g *GitlabProjectMergeRequests) GetGitlab() (gitlab *GitlabInstance, err error) {
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

func (g *GitlabProjectMergeRequests) GetGitlabProject() (gitlabProject *GitlabProject, err error) {
	if g.gitlabProject == nil {
		return nil, tracederrors.TracedErrorf("gitlabProject not set")
	}

	return g.gitlabProject, nil
}

func (g *GitlabProjectMergeRequests) GetMergeRequestByNativeMergeRequest(nativeMergeRequest *gitlab.MergeRequest) (mergeRequest *GitlabMergeRequest, err error) {
	if nativeMergeRequest == nil {
		return nil, tracederrors.TracedErrorNil("nativeMergeRequest")
	}

	mergeRequest = NewGitlabMergeRequest()
	err = mergeRequest.SetGitlabProjectMergeRequests(g)
	if err != nil {
		return nil, err
	}

	err = mergeRequest.SetId(nativeMergeRequest.IID)
	if err != nil {
		return nil, err
	}

	err = mergeRequest.SetCachedTitle(nativeMergeRequest.Title)
	if err != nil {
		return nil, err
	}

	err = mergeRequest.SetCachedSourceBranchName(nativeMergeRequest.SourceBranch)
	if err != nil {
		return nil, err
	}

	err = mergeRequest.SetCachedTargetBranchName(nativeMergeRequest.TargetBranch)
	if err != nil {
		return nil, err
	}

	return mergeRequest, nil
}

func (g *GitlabProjectMergeRequests) GetNativeMergeRequestsService() (nativeService *gitlab.MergeRequestsService, err error) {
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

func (g *GitlabProjectMergeRequests) GetOpenMergeRequestBySourceAndTargetBranch(ctx context.Context, sourceBranchName string, targetBranchName string) (mergeRequest *GitlabMergeRequest, err error) {
	if sourceBranchName == "" {
		return nil, tracederrors.TracedErrorEmptyString("sourceBranchName")
	}

	if targetBranchName == "" {
		return nil, tracederrors.TracedErrorEmptyString("targetBranchName")
	}

	openMergeRequests, err := g.GetOpenMergeRequests(ctx)
	if err != nil {
		return nil, err
	}

	foundCounter := 0
	for _, request := range openMergeRequests {
		currentSourceBranchName, err := request.GetCachedSourceBranchName()
		if err != nil {
			return nil, err
		}

		currentTargetBranchName, err := request.GetCachedTargetBranchName()
		if err != nil {
			return nil, err
		}

		if currentSourceBranchName == sourceBranchName {
			if currentTargetBranchName == targetBranchName {
				mergeRequest = request
				foundCounter += 1
			}
		}
	}

	projectUrl, err := g.GetProjectUrlAsString(ctx)
	if err != nil {
		return nil, err
	}

	if foundCounter <= 0 {
		return nil, tracederrors.TracedErrorf(
			"%w: sourceBranch '%s' and targetBranch '%s' in project %s .",
			ErrNoMergeRequestWithSourceAndTargetBranchFound,
			sourceBranchName,
			targetBranchName,
			projectUrl,
		)
	} else if foundCounter > 1 {
		return nil, tracederrors.TracedErrorf(
			"Found '%d' merge requests matching sourceBranch '%s' and targetBranch '%s' in project %s but only 1 is supported.",
			foundCounter,
			sourceBranchName,
			targetBranchName,
			projectUrl,
		)
	} else {
		title, err := mergeRequest.GetCachedTitle()
		if err != nil {
			return nil, err
		}

		logging.LogInfoByCtxf(ctx, "Found merge request by sourceBranch  '%s' and targetBranch '%s': '%s' in %s", sourceBranchName, targetBranchName, title, projectUrl)
	}

	return mergeRequest, nil
}

func (g *GitlabProjectMergeRequests) GetOpenMergeRequestByTitle(ctx context.Context, title string) (mergeRequest *GitlabMergeRequest, err error) {
	if title == "" {
		return nil, tracederrors.TracedErrorEmptyString("title")
	}

	openMergeRequests, err := g.GetOpenMergeRequests(ctx)
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

	projectUrl, err := g.GetProjectUrlAsString(ctx)
	if err != nil {
		return nil, err
	}

	if foundCounter <= 0 {
		return nil, tracederrors.TracedErrorf("%w: '%s' in project %s .", ErrNoMergeRequestWithTitleFound, title, projectUrl)
	} else if foundCounter > 1 {
		return nil, tracederrors.TracedErrorf(
			"Found '%d' merge requests matching title '%s' in project %s but only 1 is supported.",
			foundCounter,
			title,
			projectUrl,
		)
	} else {
		logging.LogInfoByCtxf(ctx, "Found merge request by title '%s': %s", title, projectUrl)
	}

	return mergeRequest, nil
}

func (g *GitlabProjectMergeRequests) GetOpenMergeRequestByTitleOrNilIfNotPresent(ctx context.Context, title string) (mergeRequest *GitlabMergeRequest, err error) {
	mergeRequest, err = g.GetOpenMergeRequestByTitle(ctx, title)
	if err != nil {
		if errors.Is(err, ErrNoMergeRequestWithTitleFound) {
			return nil, nil
		}
		return nil, err
	}

	return mergeRequest, nil
}

func (g *GitlabProjectMergeRequests) GetOpenMergeRequests(ctx context.Context) (openMergeRequest []*GitlabMergeRequest, err error) {
	var stateStringOpen string = "opened"

	rawMergeRequest, err := g.GetRawMergeRequests(ctx, &gitlab.ListProjectMergeRequestsOptions{State: &stateStringOpen})
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

func (g *GitlabProjectMergeRequests) GetProjectId(ctx context.Context) (projectId int, err error) {
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

func (g *GitlabProjectMergeRequests) GetProjectUrlAsString(ctx context.Context) (projectUrl string, err error) {
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

func (g *GitlabProjectMergeRequests) GetRawMergeRequests(ctx context.Context, options *gitlab.ListProjectMergeRequestsOptions) (rawMergeRequests []*gitlab.MergeRequest, err error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	projectId, err := g.GetProjectId(ctx)
	if err != nil {
		return nil, err
	}

	rawMergeRequests = []*gitlab.MergeRequest{}
	nextPage := 1

	for {
		if nextPage <= 0 {
			break
		}

		nativeService, err := g.GetNativeMergeRequestsService()
		if err != nil {
			return nil, err
		}
		rawMergeRequestsToAdd, response, err := nativeService.ListProjectMergeRequests(
			projectId,
			options,
		)
		if err != nil {
			return nil, tracederrors.TracedErrorf("Get raw merge requests failed: '%s'", err)
		}

		rawMergeRequests = append(rawMergeRequests, rawMergeRequestsToAdd...)

		nextPage = response.NextPage

	}
	return rawMergeRequests, nil
}

func (g *GitlabProjectMergeRequests) SetGitlabProject(gitlabProject *GitlabProject) (err error) {
	if gitlabProject == nil {
		return tracederrors.TracedErrorf("gitlabProject is nil")
	}

	g.gitlabProject = gitlabProject

	return nil
}
