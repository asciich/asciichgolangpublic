package asciichgolangpublic

import (
	"context"
	"sort"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type GitlabMergeRequest struct {
	gitlabProjectMergeRequests *GitlabProjectMergeRequests
	cachedTitle                string
	cachedSourceBranchName     string
	cachedTargetBranchName     string
	id                         int
}

func NewGitlabMergeRequest() (g *GitlabMergeRequest) {
	return new(GitlabMergeRequest)
}

func (g *GitlabMergeRequest) GetCachedSourceBranchName() (cachedSourceBranchName string, err error) {
	if g.cachedSourceBranchName == "" {
		return "", tracederrors.TracedErrorf("cachedSourceBranchName not set")
	}

	return g.cachedSourceBranchName, nil
}

func (g *GitlabMergeRequest) GetCachedTargetBranchName() (cachedTargetBranchName string, err error) {
	if g.cachedTargetBranchName == "" {
		return "", tracederrors.TracedErrorf("cachedTargetBranchName not set")
	}

	return g.cachedTargetBranchName, nil
}

func (g *GitlabMergeRequest) GetCachedTitle() (cachedTitle string, err error) {
	if g.cachedTitle == "" {
		return "", tracederrors.TracedErrorf("cachedTitle not set")
	}

	return g.cachedTitle, nil
}

func (g *GitlabMergeRequest) GetDescription(ctx context.Context) (description string, err error) {
	rawResponse, err := g.GetRawResponse(ctx)
	if err != nil {
		return "", err
	}

	description = rawResponse.Description

	return description, nil
}

func (g *GitlabMergeRequest) GetDetailedMergeStatus(ctx context.Context) (mergeStatus string, err error) {
	rawResponse, err := g.GetRawResponse(ctx)
	if err != nil {
		return "", err
	}

	mergeStatus = rawResponse.DetailedMergeStatus
	url := rawResponse.WebURL

	if mergeStatus == "" {
		return "", tracederrors.TracedErrorf(
			"mergeStatus is empty string after evaluation for merge request %s .",
			url,
		)
	}

	logging.LogInfoByCtxf(ctx, "Merge request %s has detailed merge status '%s'", url, mergeStatus)

	return mergeStatus, nil
}

func (g *GitlabMergeRequest) GetGitlabProject() (gitlabProject *GitlabProject, err error) {
	projectMergeRequests, err := g.GetGitlabProjectMergeRequests()
	if err != nil {
		return nil, err
	}

	gitlabProject, err = projectMergeRequests.GetGitlabProject()
	if err != nil {
		return nil, err
	}

	return gitlabProject, nil
}

func (g *GitlabMergeRequest) GetGitlabProjectMergeRequests() (gitlabProjectMergeRequests *GitlabProjectMergeRequests, err error) {
	if g.gitlabProjectMergeRequests == nil {
		return nil, tracederrors.TracedErrorf("gitlabProjectMergeRequests not set")
	}

	return g.gitlabProjectMergeRequests, nil
}

func (g *GitlabMergeRequest) GetId() (id int, err error) {
	if g.id <= 0 {
		return -1, tracederrors.TracedError("Id not set")
	}

	return g.id, nil
}

func (g *GitlabMergeRequest) GetLabels(ctx context.Context) (labels []string, err error) {
	rawResponse, err := g.GetRawResponse(ctx)
	if err != nil {
		return nil, err
	}

	labels = rawResponse.Labels
	if labels == nil {
		return nil, tracederrors.TracedError("labels is nil after evaluation")
	}

	sort.Strings(labels)

	return labels, nil
}

func (g *GitlabMergeRequest) GetMergeCommit(ctx context.Context) (mergeCommit *GitlabCommit, err error) {
	mergeCommitSha, err := g.GetMergeCommitSha(ctx)
	if err != nil {
		return nil, err
	}

	gitlabProject, err := g.GetGitlabProject()
	if err != nil {
		return nil, err
	}

	mergeCommit, err = gitlabProject.GetCommitByHashString(ctx, mergeCommitSha)
	if err != nil {
		return nil, err
	}

	return mergeCommit, nil
}

func (g *GitlabMergeRequest) GetMergeCommitSha(ctx context.Context) (mergeCommitSha string, err error) {
	rawResponse, err := g.GetRawResponse(ctx)
	if err != nil {
		return "", err
	}

	mergeCommitSha = rawResponse.MergeCommitSHA

	mergeRequestUrl, err := g.GetUrlAsString(ctx)
	if err != nil {
		return "", err
	}

	if mergeCommitSha == "" {
		return "", tracederrors.TracedErrorf(
			"No merge commit sha found for %s . Is merge request already merged?",
			mergeRequestUrl,
		)
	}

	logging.LogInfoByCtxf(ctx, "Merge request %s is merged and the merge commit is '%s'.", mergeRequestUrl, mergeCommitSha)

	return mergeCommitSha, nil
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

func (g *GitlabMergeRequest) GetProjectId(ctx context.Context) (projectId int, err error) {
	project, err := g.GetProject()
	if err != nil {
		return -1, err
	}

	id, err := project.GetId(ctx)
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (g *GitlabMergeRequest) GetRawResponse(ctx context.Context) (rawResponse *gitlab.MergeRequest, err error) {
	nativeService, err := g.GetNativeMergeRequestsService()
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

	rawResponse, _, err = nativeService.GetMergeRequest(
		projectId,
		id,
		&gitlab.GetMergeRequestsOptions{},
	)
	if err != nil {
		return nil, tracederrors.TracedErrorf(
			"Failed to get raw response for merge request id='%d' of gitlab projectId='%d': '%w'",
			id,
			projectId,
			err,
		)
	}

	if rawResponse == nil {
		return nil, tracederrors.TracedError("rawResponse is nil after evaluation.")
	}

	return rawResponse, nil
}

func (g *GitlabMergeRequest) GetSourceBranchName(ctx context.Context) (sourceBranchName string, err error) {
	rawResponse, err := g.GetRawResponse(ctx)
	if err != nil {
		return "", err
	}

	sourceBranchName = rawResponse.SourceBranch
	if sourceBranchName == "" {
		return "", tracederrors.TracedError(
			"sourceBranchName is empty string after evaluation",
		)
	}

	return sourceBranchName, nil
}

func (g *GitlabMergeRequest) GetTargetBranchName(ctx context.Context) (targetBranchName string, err error) {
	rawResponse, err := g.GetRawResponse(ctx)
	if err != nil {
		return "", err
	}

	targetBranchName = rawResponse.TargetBranch

	if targetBranchName == "" {
		return "", tracederrors.TracedError("TargetBranchName is empty string after evaluation")
	}

	return targetBranchName, nil
}

func (g *GitlabMergeRequest) GetUrlAsString(ctx context.Context) (url string, err error) {
	rawResponse, err := g.GetRawResponse(ctx)
	if err != nil {
		return "", err
	}

	url = rawResponse.WebURL

	if url == "" {
		return "", tracederrors.TracedError("url is empty string after evaluation")
	}

	return url, nil
}

func (g *GitlabMergeRequest) IsClosed(ctx context.Context) (isClosed bool, err error) {
	isOpen, err := g.IsOpen(ctx)
	if err != nil {
		return false, err
	}

	isClosed = !isOpen

	return isClosed, nil
}

func (g *GitlabMergeRequest) IsMerged(ctx context.Context) (isMerged bool, err error) {
	rawResponse, err := g.GetRawResponse(ctx)
	if err != nil {
		return false, err
	}

	isMerged = rawResponse.MergedAt != nil

	return isMerged, nil
}

func (g *GitlabMergeRequest) IsOpen(ctx context.Context) (isOpen bool, err error) {
	rawResponse, err := g.GetRawResponse(ctx)
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

func (g *GitlabMergeRequest) Merge(ctx context.Context) (mergeCommit *GitlabCommit, err error) {
	isMerged, err := g.IsMerged(ctx)
	if err != nil {
		return nil, err
	}

	mergeRequestUrl, err := g.GetUrlAsString(ctx)
	if err != nil {
		return nil, err
	}

	if isMerged {
		mergeCommit, err = g.GetMergeCommit(ctx)
		if err != nil {
			return nil, err
		}

		mergeCommitHash, err := mergeCommit.GetCommitHash()
		if err != nil {
			return nil, err
		}

		logging.LogInfoByCtxf(ctx, "Merge reqeuest %s is already merged. Merge Commit is '%s'.", mergeRequestUrl, mergeCommitHash)
	} else {
		// It's not possible to open and directly merger a MergeRequest
		// because Gitlab has to check if a valid merge is possible.
		// This wait function is called to wait until the "checking" by Gitlab is done.
		err := g.WaitUntilDetailedMergeStatusIsNotChecking(ctx)
		if err != nil {
			return nil, err
		}

		nativeMergeRequestsService, err := g.GetNativeMergeRequestsService()
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

		logging.LogInfoByCtxf(ctx, "Going to merge %s .", mergeRequestUrl)

		nativeMergeRequest, _, err := nativeMergeRequestsService.AcceptMergeRequest(
			projectId,
			id,
			&gitlab.AcceptMergeRequestOptions{},
		)
		if err != nil {
			return nil, tracederrors.TracedErrorf("Merging merge request %s failed: %w", mergeRequestUrl, err)
		}

		mergedCommitSha := nativeMergeRequest.MergeCommitSHA
		if mergedCommitSha == "" {
			return nil, tracederrors.TracedErrorf("mergedCommitSha is empty string after merging '%s'", mergeRequestUrl)
		}

		gitlabProject, err := g.GetProject()
		if err != nil {
			return nil, err
		}

		mergeCommit, err = gitlabProject.GetCommitByHashString(ctx, mergedCommitSha)
		if err != nil {
			return nil, err
		}

		logging.LogChangedByCtxf(ctx, "Merged %s . Merge commit is '%s'.", mergeRequestUrl, mergedCommitSha)
	}

	return mergeCommit, nil
}

func (g *GitlabMergeRequest) Close(ctx context.Context, closeMessage string) (err error) {
	if closeMessage == "" {
		return tracederrors.TracedErrorEmptyString("closeMessage")
	}

	nativeService, err := g.GetNativeMergeRequestsService()
	if err != nil {
		return err
	}

	projectId, err := g.GetProjectId(ctx)
	if err != nil {
		return err
	}

	id, err := g.GetId()
	if err != nil {
		return err
	}

	isClosed, err := g.IsClosed(ctx)
	if err != nil {
		return err
	}

	url, err := g.GetUrlAsString(ctx)
	if err != nil {
		return err
	}

	if isClosed {
		logging.LogInfoByCtxf(ctx, "Merge request %s is already closed.", url)
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
			return tracederrors.TracedErrorf("Update merge request failed: '%w'", err)
		}

		logging.LogChangedByCtxf(ctx, "Closed merge request %s", url)
	}

	return nil
}

func (g *GitlabMergeRequest) SetCachedSourceBranchName(cachedSourceBranchName string) (err error) {
	if cachedSourceBranchName == "" {
		return tracederrors.TracedErrorf("cachedSourceBranchName is empty string")
	}

	g.cachedSourceBranchName = cachedSourceBranchName

	return nil
}

func (g *GitlabMergeRequest) SetCachedTargetBranchName(cachedTargetBranchName string) (err error) {
	if cachedTargetBranchName == "" {
		return tracederrors.TracedErrorf("cachedTargetBranchName is empty string")
	}

	g.cachedTargetBranchName = cachedTargetBranchName

	return nil
}

func (g *GitlabMergeRequest) SetCachedTitle(cachedTitle string) (err error) {
	if cachedTitle == "" {
		return tracederrors.TracedErrorf("cachedTitle is empty string")
	}

	g.cachedTitle = cachedTitle

	return nil
}

func (g *GitlabMergeRequest) SetGitlabProjectMergeRequests(gitlabProjectMergeRequests *GitlabProjectMergeRequests) (err error) {
	if gitlabProjectMergeRequests == nil {
		return tracederrors.TracedErrorf("gitlabProjectMergeRequests is nil")
	}

	g.gitlabProjectMergeRequests = gitlabProjectMergeRequests

	return nil
}

func (g *GitlabMergeRequest) SetId(id int) (err error) {
	if id <= 0 {
		return tracederrors.TracedErrorf("Invalid value '%d' for id", id)
	}

	g.id = id

	return nil
}

func (g *GitlabMergeRequest) WaitUntilDetailedMergeStatusIsNotChecking(ctx context.Context) (err error) {
	mergeRequestUrl, err := g.GetUrlAsString(ctx)
	if err != nil {
		return err
	}

	for {
		mergeStatus, err := g.GetDetailedMergeStatus(ctx)
		if err != nil {
			return err
		}

		if mergeStatus == "checking" {
			logging.LogInfoByCtxf(ctx, "Waiting for merge request %s to finish status 'checking'.", mergeRequestUrl)

			time.Sleep(200 * time.Millisecond)
			continue
		}

		break
	}

	logging.LogInfoByCtxf(ctx, "Merge request %s is not in status 'checking'.", mergeRequestUrl)

	return nil
}
