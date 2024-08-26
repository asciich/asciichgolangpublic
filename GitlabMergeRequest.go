package asciichgolangpublic

import "github.com/xanzy/go-gitlab"

type GitlabMergeRequest struct {
	gitlabProjectMergeRequests *GitlabProjectMergeRequests
	cachedTitle                string
	id                         int
}

func NewGitlabMergeRequest() (g *GitlabMergeRequest) {
	return new(GitlabMergeRequest)
}

func (g *GitlabMergeRequest) GetCachedTitle() (cachedTitle string, err error) {
	if g.cachedTitle == "" {
		return "", TracedErrorf("cachedTitle not set")
	}

	return g.cachedTitle, nil
}

func (g *GitlabMergeRequest) GetGitlabProjectMergeRequests() (gitlabProjectMergeRequests *GitlabProjectMergeRequests, err error) {
	if g.gitlabProjectMergeRequests == nil {
		return nil, TracedErrorf("gitlabProjectMergeRequests not set")
	}

	return g.gitlabProjectMergeRequests, nil
}

func (g *GitlabMergeRequest) GetId() (id int, err error) {
	if g.id <= 0 {
		return -1, TracedError("Id not set")
	}

	return g.id, nil
}

func (g *GitlabMergeRequest) GetNativeMergeRequestsService() (nativeService *gitlab.MergeRequestsService, err error) {
	mergeRequests, err := g.GetGitlabProjectMergeRequests()
	if err != nil {
		return nil, err
	}

	nativeService, err = mergeRequests.GetNativeMergeRequestsService()
	if err != nil {
		return nil, err
	}

	return nativeService, nil
}

func (g *GitlabMergeRequest) GetProject() (project *GitlabProject, err error) {
	mergeRequests, err := g.GetGitlabProjectMergeRequests()
	if err != nil {
		return nil, err
	}

	project, err = mergeRequests.GetGitlabProject()
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (g *GitlabMergeRequest) GetProjectId() (projectId int, err error) {
	project, err := g.GetProject()
	if err != nil {
		return -1, err
	}

	id, err := project.GetId()
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (g *GitlabMergeRequest) GetRawResponse() (rawResponse *gitlab.MergeRequest, err error) {
	nativeService, err := g.GetNativeMergeRequestsService()
	if err != nil {
		return nil, err
	}

	projectId, err := g.GetProjectId()
	if err != nil {
		return nil, err
	}

	id, err := g.GetId()
	if err != nil {
		return nil, err
	}

	rawResponse, _, err = nativeService.GetMergeRequest(
		projectId,
		id,
		&gitlab.GetMergeRequestsOptions{},
	)
	if err != nil {
		return nil, TracedErrorf("Failed to get raw resopnse for gitlab project: '%w'", err)
	}

	if rawResponse == nil {
		return nil, TracedError("rawResponse is nil after evaluation.")
	}

	return rawResponse, nil
}

func (g *GitlabMergeRequest) GetSourceBranchName() (sourceBranchName string, err error) {
	rawResponse, err := g.GetRawResponse()
	if err != nil {
		return "", err
	}

	sourceBranchName = rawResponse.SourceBranch
	if err != nil {
		return "", err
	}

	return sourceBranchName, nil
}

func (g *GitlabMergeRequest) GetTargetBranchName() (targetBranchName string, err error) {
	rawResponse, err := g.GetRawResponse()
	if err != nil {
		return "", err
	}

	targetBranchName = rawResponse.TargetBranch

	if targetBranchName == "" {
		return "", TracedError("TargetBranchName is empty string after evaluation")
	}

	return targetBranchName, nil
}

func (g *GitlabMergeRequest) GetUrlAsString() (url string, err error) {
	rawResponse, err := g.GetRawResponse()
	if err != nil {
		return "", err
	}

	url = rawResponse.WebURL

	if url == "" {
		return "", TracedError("url is empty string after evaluation")
	}

	return url, nil
}

func (g *GitlabMergeRequest) IsClosed() (isClosed bool, err error) {
	isOpen, err := g.IsOpen()
	if err != nil {
		return false, err
	}

	isClosed = !isOpen

	return isClosed, nil
}

func (g *GitlabMergeRequest) IsOpen() (isOpen bool, err error) {
	rawResponse, err := g.GetRawResponse()
	if err != nil {
		return false, err
	}

	if rawResponse.ClosedAt != nil {
		return false, nil
	}

	if rawResponse.MergedAt != nil {
		return false, nil
	}

	return true, nil
}

func (g *GitlabMergeRequest) MustClose(closeMessage string, verbose bool) (err error) {
	if closeMessage == "" {
		return TracedErrorEmptyString("closeMessage")
	}

	nativeService, err := g.GetNativeMergeRequestsService()
	if err != nil {
		return err
	}

	projectId, err := g.GetProjectId()
	if err != nil {
		return err
	}

	id, err := g.GetId()
	if err != nil {
		return err
	}

	isClosed, err := g.IsClosed()
	if err != nil {
		return err
	}

	url, err := g.GetUrlAsString()
	if err != nil {
		return err
	}

	if isClosed {
		if verbose {
			LogInfof("Merge request %s is already closed.", url)
		}
	} else {
		stateEvent := "close"

		_, _, err = nativeService.UpdateMergeRequest(
			projectId,
			id,
			&gitlab.UpdateMergeRequestOptions{
				StateEvent: &stateEvent,
			},
		)
		if err != nil {
			return TracedErrorf("Update merge request failed: '%w'", err)
		}

		if verbose {
			LogChangedf("Closed merge request %s", url)
		}
	}

	return nil
}

func (g *GitlabMergeRequest) MustGetCachedTitle() (cachedTitle string) {
	cachedTitle, err := g.GetCachedTitle()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return cachedTitle
}

func (g *GitlabMergeRequest) MustGetGitlabMergeRequests() (gitlabMergeRequests *GitlabProjectMergeRequests) {
	gitlabMergeRequests, err := g.GetGitlabProjectMergeRequests()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlabMergeRequests
}

func (g *GitlabMergeRequest) MustGetGitlabProjectMergeRequests() (gitlabProjectMergeRequests *GitlabProjectMergeRequests) {
	gitlabProjectMergeRequests, err := g.GetGitlabProjectMergeRequests()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlabProjectMergeRequests
}

func (g *GitlabMergeRequest) MustGetId() (id int) {
	id, err := g.GetId()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return id
}

func (g *GitlabMergeRequest) MustGetNativeMergeRequestsService() (nativeService *gitlab.MergeRequestsService) {
	nativeService, err := g.GetNativeMergeRequestsService()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return nativeService
}

func (g *GitlabMergeRequest) MustGetProject() (project *GitlabProject) {
	project, err := g.GetProject()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return project
}

func (g *GitlabMergeRequest) MustGetProjectId() (projectId int) {
	projectId, err := g.GetProjectId()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return projectId
}

func (g *GitlabMergeRequest) MustGetRawResponse() (rawResponse *gitlab.MergeRequest) {
	rawResponse, err := g.GetRawResponse()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return rawResponse
}

func (g *GitlabMergeRequest) MustGetSourceBranchName() (sourceBranchName string) {
	sourceBranchName, err := g.GetSourceBranchName()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return sourceBranchName
}

func (g *GitlabMergeRequest) MustGetTargetBranchName() (targetBranchName string) {
	targetBranchName, err := g.GetTargetBranchName()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return targetBranchName
}

func (g *GitlabMergeRequest) MustGetUrlAsString() (url string) {
	url, err := g.GetUrlAsString()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return url
}

func (g *GitlabMergeRequest) MustIsClosed() (isClosed bool) {
	isClosed, err := g.IsClosed()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return isClosed
}

func (g *GitlabMergeRequest) MustIsOpen() (isOpen bool) {
	isOpen, err := g.IsOpen()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return isOpen
}

func (g *GitlabMergeRequest) MustSetCachedTitle(cachedTitle string) {
	err := g.SetCachedTitle(cachedTitle)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabMergeRequest) MustSetGitlabProjectMergeRequests(gitlabProjectMergeRequests *GitlabProjectMergeRequests) {
	err := g.SetGitlabProjectMergeRequests(gitlabProjectMergeRequests)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabMergeRequest) MustSetId(id int) {
	err := g.SetId(id)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabMergeRequest) SetCachedTitle(cachedTitle string) (err error) {
	if cachedTitle == "" {
		return TracedErrorf("cachedTitle is empty string")
	}

	g.cachedTitle = cachedTitle

	return nil
}

func (g *GitlabMergeRequest) SetGitlabProjectMergeRequests(gitlabProjectMergeRequests *GitlabProjectMergeRequests) (err error) {
	if gitlabProjectMergeRequests == nil {
		return TracedErrorf("gitlabProjectMergeRequests is nil")
	}

	g.gitlabProjectMergeRequests = gitlabProjectMergeRequests

	return nil
}

func (g *GitlabMergeRequest) SetId(id int) (err error) {
	if id <= 0 {
		return TracedErrorf("Invalid value '%d' for id", id)
	}

	g.id = id

	return nil
}
