package asciichgolangpublic

import (
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

func (g *GitlabMergeRequest) GetDescription() (description string, err error) {
	rawResponse, err := g.GetRawResponse()
	if err != nil {
		return "", err
	}

	description = rawResponse.Description

	return description, nil
}

func (g *GitlabMergeRequest) GetDetailedMergeStatus(verbose bool) (mergeStatus string, err error) {
	rawResponse, err := g.GetRawResponse()
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

	if verbose {
		logging.LogInfof("Merge request %s has detailed merge status '%s'", url, mergeStatus)
	}

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

func (g *GitlabMergeRequest) GetLabels() (labels []string, err error) {
	rawResponse, err := g.GetRawResponse()
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

func (g *GitlabMergeRequest) GetMergeCommit(verbose bool) (mergeCommit *GitlabCommit, err error) {
	mergeCommitSha, err := g.GetMergeCommitSha(verbose)
	if err != nil {
		return nil, err
	}

	gitlabProject, err := g.GetGitlabProject()
	if err != nil {
		return nil, err
	}

	mergeCommit, err = gitlabProject.GetCommitByHashString(mergeCommitSha, verbose)
	if err != nil {
		return nil, err
	}

	return mergeCommit, nil
}

func (g *GitlabMergeRequest) GetMergeCommitSha(verbose bool) (mergeCommitSha string, err error) {
	rawResponse, err := g.GetRawResponse()
	if err != nil {
		return "", err
	}

	mergeCommitSha = rawResponse.MergeCommitSHA

	mergeRequestUrl, err := g.GetUrlAsString()
	if err != nil {
		return "", err
	}

	if mergeCommitSha == "" {
		return "", tracederrors.TracedErrorf(
			"No merge commit sha found for %s . Is merge request already merged?",
			mergeRequestUrl,
		)
	}

	if verbose {
		logging.LogInfof(
			"Merge request %s is merged and the merge commit is '%s'.",
			mergeRequestUrl,
			mergeCommitSha,
		)
	}

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

func (g *GitlabMergeRequest) GetSourceBranchName() (sourceBranchName string, err error) {
	rawResponse, err := g.GetRawResponse()
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

func (g *GitlabMergeRequest) GetTargetBranchName() (targetBranchName string, err error) {
	rawResponse, err := g.GetRawResponse()
	if err != nil {
		return "", err
	}

	targetBranchName = rawResponse.TargetBranch

	if targetBranchName == "" {
		return "", tracederrors.TracedError("TargetBranchName is empty string after evaluation")
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
		return "", tracederrors.TracedError("url is empty string after evaluation")
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

func (g *GitlabMergeRequest) IsMerged() (isMerged bool, err error) {
	rawResponse, err := g.GetRawResponse()
	if err != nil {
		return false, err
	}

	isMerged = rawResponse.MergedAt != nil

	return isMerged, nil
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

func (g *GitlabMergeRequest) Merge(options *GitlabMergeOptions) (mergeCommit *GitlabCommit, err error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	isMerged, err := g.IsMerged()
	if err != nil {
		return nil, err
	}

	mergeRequestUrl, err := g.GetUrlAsString()
	if err != nil {
		return nil, err
	}

	if isMerged {
		mergeCommit, err = g.GetMergeCommit(options.Verbose)
		if err != nil {
			return nil, err
		}

		mergeCommitHash, err := mergeCommit.GetCommitHash()
		if err != nil {
			return nil, err
		}

		if options.Verbose {
			logging.LogInfof(
				"Merge reqeuest %s is already merged. Merge Commit is '%s'.",
				mergeRequestUrl,
				mergeCommitHash,
			)
		}
	} else {
		// It's not possible to open and directly merger a MergeRequest
		// because Gitlab has to check if a valid merge is possible.
		// This wait function is called to wait until the "checking" by Gitlab is done.
		err := g.WaitUntilDetailedMergeStatusIsNotChecking(options.Verbose)
		if err != nil {
			return nil, err
		}

		nativeMergeRequestsService, err := g.GetNativeMergeRequestsService()
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

		if options.Verbose {
			logging.LogInfof("Going to merge %s .", mergeRequestUrl)
		}

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

		mergeCommit, err = gitlabProject.GetCommitByHashString(mergedCommitSha, options.Verbose)
		if err != nil {
			return nil, err
		}

		if options.Verbose {
			logging.LogChangedf(
				"Merged %s . Merge commit is '%s'.",
				mergeRequestUrl,
				mergedCommitSha,
			)
		}
	}

	return mergeCommit, nil
}

func (g *GitlabMergeRequest) Close(closeMessage string, verbose bool) (err error) {
	if closeMessage == "" {
		return tracederrors.TracedErrorEmptyString("closeMessage")
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
			logging.LogInfof("Merge request %s is already closed.", url)
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
			return tracederrors.TracedErrorf("Update merge request failed: '%w'", err)
		}

		if verbose {
			logging.LogChangedf("Closed merge request %s", url)
		}
	}

	return nil
}

func (g *GitlabMergeRequest) MustGetCachedSourceBranchName() (cachedSourceBranchName string) {
	cachedSourceBranchName, err := g.GetCachedSourceBranchName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return cachedSourceBranchName
}

func (g *GitlabMergeRequest) MustGetCachedTargetBranchName() (cachedTargetBranchName string) {
	cachedTargetBranchName, err := g.GetCachedTargetBranchName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return cachedTargetBranchName
}

func (g *GitlabMergeRequest) MustGetCachedTitle() (cachedTitle string) {
	cachedTitle, err := g.GetCachedTitle()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return cachedTitle
}

func (g *GitlabMergeRequest) MustGetDescription() (description string) {
	description, err := g.GetDescription()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return description
}

func (g *GitlabMergeRequest) MustGetDetailedMergeStatus(verbose bool) (mergeStatus string) {
	mergeStatus, err := g.GetDetailedMergeStatus(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return mergeStatus
}

func (g *GitlabMergeRequest) MustGetGitlabMergeRequests() (gitlabMergeRequests *GitlabProjectMergeRequests) {
	gitlabMergeRequests, err := g.GetGitlabProjectMergeRequests()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlabMergeRequests
}

func (g *GitlabMergeRequest) MustGetGitlabProject() (gitlabProject *GitlabProject) {
	gitlabProject, err := g.GetGitlabProject()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlabProject
}

func (g *GitlabMergeRequest) MustGetGitlabProjectMergeRequests() (gitlabProjectMergeRequests *GitlabProjectMergeRequests) {
	gitlabProjectMergeRequests, err := g.GetGitlabProjectMergeRequests()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlabProjectMergeRequests
}

func (g *GitlabMergeRequest) MustGetId() (id int) {
	id, err := g.GetId()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return id
}

func (g *GitlabMergeRequest) MustGetLabels() (labels []string) {
	labels, err := g.GetLabels()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return labels
}

func (g *GitlabMergeRequest) MustGetMergeCommit(verbose bool) (mergeCommit *GitlabCommit) {
	mergeCommit, err := g.GetMergeCommit(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return mergeCommit
}

func (g *GitlabMergeRequest) MustGetMergeCommitSha(verbose bool) (mergeCommitSha string) {
	mergeCommitSha, err := g.GetMergeCommitSha(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return mergeCommitSha
}

func (g *GitlabMergeRequest) MustGetNativeMergeRequestsService() (nativeService *gitlab.MergeRequestsService) {
	nativeService, err := g.GetNativeMergeRequestsService()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return nativeService
}

func (g *GitlabMergeRequest) MustGetProject() (project *GitlabProject) {
	project, err := g.GetProject()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return project
}

func (g *GitlabMergeRequest) MustGetProjectId() (projectId int) {
	projectId, err := g.GetProjectId()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return projectId
}

func (g *GitlabMergeRequest) MustGetRawResponse() (rawResponse *gitlab.MergeRequest) {
	rawResponse, err := g.GetRawResponse()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return rawResponse
}

func (g *GitlabMergeRequest) MustGetSourceBranchName() (sourceBranchName string) {
	sourceBranchName, err := g.GetSourceBranchName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return sourceBranchName
}

func (g *GitlabMergeRequest) MustGetTargetBranchName() (targetBranchName string) {
	targetBranchName, err := g.GetTargetBranchName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return targetBranchName
}

func (g *GitlabMergeRequest) MustGetUrlAsString() (url string) {
	url, err := g.GetUrlAsString()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return url
}

func (g *GitlabMergeRequest) MustIsClosed() (isClosed bool) {
	isClosed, err := g.IsClosed()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return isClosed
}

func (g *GitlabMergeRequest) MustIsMerged() (isMerged bool) {
	isMerged, err := g.IsMerged()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return isMerged
}

func (g *GitlabMergeRequest) MustIsOpen() (isOpen bool) {
	isOpen, err := g.IsOpen()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return isOpen
}

func (g *GitlabMergeRequest) MustMerge(options *GitlabMergeOptions) (mergeCommit *GitlabCommit) {
	mergeCommit, err := g.Merge(options)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return mergeCommit
}

func (g *GitlabMergeRequest) MustSetCachedSourceBranchName(cachedSourceBranchName string) {
	err := g.SetCachedSourceBranchName(cachedSourceBranchName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabMergeRequest) MustSetCachedTargetBranchName(cachedTargetBranchName string) {
	err := g.SetCachedTargetBranchName(cachedTargetBranchName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabMergeRequest) MustSetCachedTitle(cachedTitle string) {
	err := g.SetCachedTitle(cachedTitle)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabMergeRequest) MustSetGitlabProjectMergeRequests(gitlabProjectMergeRequests *GitlabProjectMergeRequests) {
	err := g.SetGitlabProjectMergeRequests(gitlabProjectMergeRequests)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabMergeRequest) MustSetId(id int) {
	err := g.SetId(id)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabMergeRequest) MustWaitUntilDetailedMergeStatusIsNotChecking(verbose bool) {
	err := g.WaitUntilDetailedMergeStatusIsNotChecking(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
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

func (g *GitlabMergeRequest) WaitUntilDetailedMergeStatusIsNotChecking(verbose bool) (err error) {
	mergeRequestUrl, err := g.GetUrlAsString()
	if err != nil {
		return err
	}

	for {
		mergeStatus, err := g.GetDetailedMergeStatus(verbose)
		if err != nil {
			return err
		}

		if mergeStatus == "checking" {
			if verbose {
				logging.LogInfof("Waiting for merge request %s to finish status 'checking'.", mergeRequestUrl)
			}

			time.Sleep(200 * time.Millisecond)
			continue
		}

		break
	}

	if verbose {
		logging.LogInfof("Merge request %s is not in status 'checking'.", mergeRequestUrl)
	}

	return nil
}
