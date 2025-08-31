package asciichgolangpublic

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions/authenticationoptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"github.com/asciich/asciichgolangpublic/pkg/urlsutils"
	"github.com/asciich/asciichgolangpublic/pkg/versionutils"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type GitlabProject struct {
	gitlab     *GitlabInstance
	id         int
	cachedPath string
}

func GetGitlabProjectByUrl(ctx context.Context, url *urlsutils.URL, authOptions []authenticationoptions.AuthenticationOption) (gitlabProject *GitlabProject, err error) {
	if url == nil {
		return nil, tracederrors.TracedErrorNil("url")
	}

	if authOptions == nil {
		return nil, tracederrors.TracedErrorNil("authOptions")
	}

	fqdnWithSheme, path, err := url.GetFqdnWitShemeAndPathAsString()
	if err != nil {
		return nil, err
	}

	gitlab, err := GetGitlabByFQDN(fqdnWithSheme)
	if err != nil {
		return nil, err
	}

	authOption, err := authenticationoptions.AuthenticationOptionsHandler().GetAuthenticationoptionsForServiceByUrl(authOptions, url)
	if err != nil {
		return nil, err
	}

	gitlabAuthenticationOption, ok := authOption.(*GitlabAuthenticationOptions)
	if !ok {
		return nil, tracederrors.TracedErrorf("Unable to get %v as GitlabAuthenticationOptions", authOption)
	}

	err = gitlab.Authenticate(ctx, gitlabAuthenticationOption)
	if err != nil {
		return nil, err
	}

	gitlabProject, err = gitlab.GetGitlabProjectByPath(ctx, path)
	if err != nil {
		return nil, err
	}

	return gitlabProject, err
}

func GetGitlabProjectByUrlFromString(ctx context.Context, urlString string, authOptions []authenticationoptions.AuthenticationOption) (gitlabProject *GitlabProject, err error) {
	if urlString == "" {
		return nil, tracederrors.TracedErrorEmptyString("urlString")
	}

	url, err := urlsutils.GetUrlFromString(urlString)
	if err != nil {
		return nil, err
	}

	return GetGitlabProjectByUrl(ctx, url, authOptions)
}

func NewGitlabProject() (gitlabProject *GitlabProject) {
	return new(GitlabProject)
}

func (g *GitlabProject) Create(ctx context.Context) (err error) {
	gitlabProjects, err := g.GetGitlabProjects()
	if err != nil {
		return err
	}

	projectPath, err := g.GetCachedPath(ctx)
	if err != nil {
		return err
	}

	createdProject, err := gitlabProjects.CreateProject(
		ctx,
		&GitlabCreateProjectOptions{
			ProjectPath: projectPath,
		},
	)
	if err != nil {
		return err
	}

	createdProjectId, err := createdProject.GetId(ctx)
	if err != nil {
		return err
	}

	err = g.SetId(createdProjectId)
	if err != nil {
		return err
	}

	return err
}

func (g *GitlabProject) GetNativePipelineSchedulesClient() (nativeClient *gitlab.PipelineSchedulesService, err error) {
	gitlab, err := g.GetGitlab()
	if err != nil {
		return nil, err
	}

	return gitlab.GetNativePipelineSchedulesClient()
}

func (g *GitlabProject) CreateBranchFromDefaultBranch(ctx context.Context, branchName string) (createdBranch *GitlabBranch, err error) {
	branches, err := g.GetBranches()
	if err != nil {
		return nil, err
	}

	createdBranch, err = branches.CreateBranchFromDefaultBranch(ctx, branchName)
	if err != nil {
		return nil, err
	}

	return createdBranch, nil
}

func (g *GitlabProject) CreateEmptyFile(ctx context.Context, fileName string, ref string) (createdFile *GitlabRepositoryFile, err error) {
	repositoryFiles, err := g.GetRepositoryFiles()
	if err != nil {
		return nil, err
	}

	createdFile, err = repositoryFiles.CreateEmptyFile(ctx, fileName, ref)
	if err != nil {
		return nil, err
	}

	return createdFile, nil
}

func (g *GitlabProject) ListScheduledPipelineNames(ctx context.Context) (scheduledPipelineNames []string, err error) {
	schedules, err := g.GetPipelineSchedules()
	if err != nil {
		return nil, err
	}

	return schedules.ListScheduledPipelineNames(ctx)
}

func (g *GitlabProject) CreateMergeRequest(ctx context.Context, options *GitlabCreateMergeRequestOptions) (createdMergeRequest *GitlabMergeRequest, err error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	sourceBranchName, err := options.GetSourceBranchName()
	if err != nil {
		return nil, err
	}

	sourceBranch, err := g.GetBranchByName(sourceBranchName)
	if err != nil {
		return nil, err
	}

	createdMergeRequest, err = sourceBranch.CreateMergeRequest(ctx, options)
	if err != nil {
		return nil, err
	}

	return createdMergeRequest, nil
}

func (g *GitlabProject) CreateNextMajorReleaseFromLatestCommitInDefaultBranch(ctx context.Context, description string) (createdRelease *GitlabRelease, err error) {
	if description == "" {
		return nil, tracederrors.TracedErrorEmptyString("description")
	}

	nextPatchVersionString, err := g.GetNextMajorReleaseVersionString(ctx)
	if err != nil {
		return nil, err
	}

	createdRelease, err = g.CreateReleaseFromLatestCommitInDefaultBranch(
		ctx,
		&GitlabCreateReleaseOptions{
			Name:        nextPatchVersionString,
			Description: description,
		},
	)
	if err != nil {
		return nil, err
	}

	return createdRelease, nil
}

func (g *GitlabProject) CreateNextMinorReleaseFromLatestCommitInDefaultBranch(ctx context.Context, description string) (createdRelease *GitlabRelease, err error) {
	if description == "" {
		return nil, tracederrors.TracedErrorEmptyString("description")
	}

	nextPatchVersionString, err := g.GetNextMinorReleaseVersionString(ctx)
	if err != nil {
		return nil, err
	}

	createdRelease, err = g.CreateReleaseFromLatestCommitInDefaultBranch(
		ctx,
		&GitlabCreateReleaseOptions{
			Name:        nextPatchVersionString,
			Description: description,
		},
	)
	if err != nil {
		return nil, err
	}

	return createdRelease, nil
}

func (g *GitlabProject) CreateNextPatchReleaseFromLatestCommitInDefaultBranch(ctx context.Context, description string) (createdRelease *GitlabRelease, err error) {
	if description == "" {
		return nil, tracederrors.TracedErrorEmptyString("description")
	}

	nextPatchVersionString, err := g.GetNextPatchReleaseVersionString(ctx)
	if err != nil {
		return nil, err
	}

	createdRelease, err = g.CreateReleaseFromLatestCommitInDefaultBranch(
		ctx,
		&GitlabCreateReleaseOptions{
			Name:        nextPatchVersionString,
			Description: description,
		},
	)
	if err != nil {
		return nil, err
	}

	return createdRelease, nil
}

func (g *GitlabProject) Delete(ctx context.Context) (err error) {
	gitlabProjects, err := g.GetGitlabProjects()
	if err != nil {
		return err
	}

	projectPath, err := g.GetCachedPath(ctx)
	if err != nil {
		return err
	}

	err = gitlabProjects.DeleteProject(
		ctx,
		&GitlabDeleteProjectOptions{
			ProjectPath: projectPath,
		},
	)
	if err != nil {
		return err
	}

	g.id = 0

	return nil
}

func (g *GitlabProject) DeleteAllBranchesExceptDefaultBranch(ctx context.Context) (err error) {
	branches, err := g.GetBranches()
	if err != nil {
		return err
	}

	err = branches.DeleteAllBranchesExceptDefaultBranch(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (g *GitlabProject) DeleteAllReleases(ctx context.Context, deleteOptions *GitlabDeleteReleaseOptions) (err error) {
	releases, err := g.GetGitlabReleases()
	if err != nil {
		return err
	}

	err = releases.DeleteAllReleases(ctx, deleteOptions)
	if err != nil {
		return err
	}

	return nil
}

func (g *GitlabProject) DeleteBranch(ctx context.Context, branchName string, deleteOptions *GitlabDeleteBranchOptions) (err error) {
	if branchName == "" {
		return tracederrors.TracedErrorEmptyString("branchNAme")
	}

	if deleteOptions == nil {
		return tracederrors.TracedErrorNil("deleteOptons")
	}

	branch, err := g.GetBranchByName(branchName)
	if err != nil {
		return err
	}

	err = branch.Delete(ctx, deleteOptions)
	if err != nil {
		return err
	}

	return nil
}

func (g *GitlabProject) DeleteFileInDefaultBranch(ctx context.Context, fileName string, commitMessage string) (err error) {
	if fileName == "" {
		return tracederrors.TracedErrorEmptyString("fileName")
	}

	if commitMessage == "" {
		return tracederrors.TracedErrorEmptyString("commitMessage")
	}

	fileInRepo, err := g.GetFileInDefaultBranch(
		ctx,
		fileName,
	)
	if err != nil {
		return err
	}

	err = fileInRepo.Delete(ctx, commitMessage)
	if err != nil {
		return err
	}

	return nil
}

func (g *GitlabProject) Exists(ctx context.Context) (projectExists bool, err error) {
	gitlab, err := g.GetGitlab()
	if err != nil {
		return false, err
	}

	idSet, err := g.IsIdSet()
	if err != nil {
		return false, err
	}
	if idSet {
		projectId, err := g.GetId(ctx)
		if err != nil {
			return false, err
		}

		projectExists, err = gitlab.ProjectByProjectIdExists(ctx, projectId)
		if err != nil {
			return false, err
		}
	} else {
		projectPath, err := g.GetCachedPath(ctx)
		if err != nil {
			return false, err
		}

		projectExists, err = gitlab.ProjectByProjectPathExists(ctx, projectPath)
		if err != nil {
			return false, err
		}
	}

	return projectExists, nil
}

func (g *GitlabProject) GetBranchByName(branchName string) (branch *GitlabBranch, err error) {
	if branchName == "" {
		return nil, tracederrors.TracedErrorEmptyString("branchName")
	}

	branches, err := g.GetBranches()
	if err != nil {
		return nil, err
	}

	branch, err = branches.GetBranchByName(branchName)
	if err != nil {
		return nil, err
	}

	return branch, nil
}

func (g *GitlabProject) GetBranchNames(ctx context.Context) (branchNames []string, err error) {
	branches, err := g.GetBranches()
	if err != nil {
		return nil, err
	}

	branchNames, err = branches.GetBranchNames(ctx)
	if err != nil {
		return nil, err
	}

	return branchNames, nil
}

func (g *GitlabProject) GetBranches() (branches *GitlabBranches, err error) {
	branches = NewGitlabBranches()

	err = branches.SetGitlabProject(g)
	if err != nil {
		return nil, err
	}

	return branches, nil
}

func (g *GitlabProject) GetCachedPathForPersonalProject(ctx context.Context) (cachedPath string, err error) {
	projectName, err := g.GetCachedProjectName(ctx)
	if err != nil {
		return "", err
	}

	userName, err := g.GetCurrentUserName(ctx)
	if err != nil {
		return "", err
	}

	cachedPath = fmt.Sprintf("%s/%s", userName, projectName)

	return cachedPath, nil
}

func (g *GitlabProject) GetCachedProjectName(ctx context.Context) (projectName string, err error) {
	cachedPath, err := g.GetCachedPath(ctx)
	if err != nil {
		return "", err
	}

	projectName = filepath.Base(cachedPath)
	if projectName == "" {
		return "", tracederrors.TracedErrorf("Unable to extract project name from cachedPath = '%s'", cachedPath)
	}

	return projectName, nil
}

func (g *GitlabProject) GetCommitByHashString(ctx context.Context, hashString string) (commit *GitlabCommit, err error) {
	if hashString == "" {
		return nil, tracederrors.TracedErrorNil("hashString")
	}

	projectCommits, err := g.GetProjectCommits()
	if err != nil {
		return nil, err
	}

	commit, err = projectCommits.GetCommitByHashString(ctx, hashString)
	if err != nil {
		return nil, err
	}

	return commit, nil
}

func (g *GitlabProject) GetCurrentUserName(ctx context.Context) (userName string, err error) {
	gitlab, err := g.GetGitlab()
	if err != nil {
		return "", err
	}

	userName, err = gitlab.GetCurrentUsersName(ctx)
	if err != nil {
		return "", err
	}

	return userName, nil
}

func (g *GitlabProject) GetDeepCopy() (copy *GitlabProject) {
	copy = NewGitlabProject()

	*copy = *g

	if g.gitlab != nil {
		copy.gitlab = g.gitlab.GetDeepCopy()
	}

	return copy
}

func (g *GitlabProject) GetDefaultBranch(ctx context.Context) (defaultBranch *GitlabBranch, err error) {
	defaultBranchName, err := g.GetDefaultBranchName(ctx)
	if err != nil {
		return nil, err
	}

	defaultBranch, err = g.GetBranchByName(defaultBranchName)
	if err != nil {
		return nil, err
	}

	return defaultBranch, nil
}

func (g *GitlabProject) GetDefaultBranchName(ctx context.Context) (defaultBranchName string, err error) {
	nativeProject, err := g.GetRawResponse(ctx)
	if err != nil {
		return "", err
	}

	defaultBranchName = nativeProject.DefaultBranch
	if defaultBranchName == "" {
		return "", tracederrors.TracedError("defaultBranchName is empty string after evaluation")
	}

	return defaultBranchName, nil
}

func (g *GitlabProject) GetDirectoryNames(ctx context.Context, ref string) (directoryNames []string, err error) {
	repositoryFiles, err := g.GetRepositoryFiles()
	if err != nil {
		return nil, err
	}

	directoryNames, err = repositoryFiles.GetDirectoryNames(ctx, ref)
	if err != nil {
		return nil, err
	}

	return directoryNames, nil
}

func (g *GitlabProject) GetFileInDefaultBranch(ctx context.Context, fileName string) (repositoryFile *GitlabRepositoryFile, err error) {
	if fileName == "" {
		return nil, tracederrors.TracedErrorEmptyString("fileName")
	}

	repositoryFile, err = g.GetRepositoryFile(
		ctx,
		&GitlabGetRepositoryFileOptions{
			Path: fileName,
		},
	)
	if err != nil {
		return nil, err
	}

	return repositoryFile, nil
}

func (g *GitlabProject) GetFilesNames(ctx context.Context, ref string) (fileNames []string, err error) {
	repositoryFiles, err := g.GetRepositoryFiles()
	if err != nil {
		return nil, err
	}

	fileNames, err = repositoryFiles.GetFileNames(ctx, ref)
	if err != nil {
		return nil, err
	}

	return fileNames, nil
}

func (g *GitlabProject) GetGitlabFqdn() (fqdn string, err error) {
	gitlab, err := g.GetGitlab()
	if err != nil {
		return "", err
	}

	fqdn, err = gitlab.GetFqdn()
	if err != nil {
		return "", err
	}

	return fqdn, nil
}

func (g *GitlabProject) GetLatestCommit(ctx context.Context, branchName string) (latestCommit *GitlabCommit, err error) {
	if branchName == "" {
		return nil, tracederrors.TracedErrorNil("branchName")
	}

	latestHash, err := g.GetLatestCommitHashAsString(ctx, branchName)
	if err != nil {
		return nil, err
	}

	latestCommit, err = g.GetCommitByHashString(ctx, latestHash)
	if err != nil {
		return nil, err
	}

	logging.LogInfof(
		"Latest commit of branch '%s' has hash '%s'.",
		branchName,
		latestHash,
	)

	return latestCommit, nil
}

func (g *GitlabProject) GetLatestCommitHashAsString(ctx context.Context, branchName string) (commitHash string, err error) {
	if branchName == "" {
		return "", tracederrors.TracedErrorEmptyString("branchName")
	}

	branch, err := g.GetBranchByName(branchName)
	if err != nil {
		return "", err
	}

	commitHash, err = branch.GetLatestCommitHashAsString(ctx)
	if err != nil {
		return "", err
	}

	return commitHash, nil
}

func (g *GitlabProject) GetLatestCommitOfDefaultBranch(ctx context.Context) (latestCommit *GitlabCommit, err error) {
	defaultBranch, err := g.GetDefaultBranchName(ctx)
	if err != nil {
		return nil, err
	}

	latestCommit, err = g.GetLatestCommit(ctx, defaultBranch)
	if err != nil {
		return nil, err
	}

	latestHash, err := latestCommit.GetCommitHash()
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Latest commit of default branch '%s' has hash '%s'.", defaultBranch, latestHash)

	return latestCommit, nil
}

func (g *GitlabProject) GetMergeRequests() (mergeRequestes *GitlabProjectMergeRequests, err error) {
	mergeRequestes = NewGitlabMergeRequests()

	err = mergeRequestes.SetGitlabProject(g)
	if err != nil {
		return nil, err
	}

	return mergeRequestes, nil
}

func (g *GitlabProject) GetNewestSemanticVersion(ctx context.Context) (newestSemanticVersion *versionutils.SemanticVersion, err error) {
	semanticVersions, err := g.GetSemanticVersions(ctx)
	if err != nil {
		return nil, err
	}

	newestVersion, err := versionutils.GetLatestVersionFromSlice(semanticVersions)
	if err != nil {
		return nil, err
	}

	newestSemanticVersion, ok := newestVersion.(*versionutils.SemanticVersion)
	if !ok {
		return nil, tracederrors.TracedErrorf(
			"Unable to get newest semantiv version from '%v'",
			newestVersion,
		)
	}

	return newestSemanticVersion, nil
}

func (g *GitlabProject) GetNewestVersion(ctx context.Context) (newestVersion versionutils.Version, err error) {
	availableVersions, err := g.GetVersions(ctx)
	if err != nil {
		return nil, err
	}

	if len(availableVersions) <= 0 {
		return nil, tracederrors.TracedError("No versionTags returned")
	}

	newestVersion, err = versionutils.GetLatestVersionFromSlice(availableVersions)
	if err != nil {
		return nil, err
	}

	return newestVersion, nil
}

func (g *GitlabProject) GetNewestVersionAsString(ctx context.Context) (newestVersionString string, err error) {
	newestVersion, err := g.GetNewestVersion(ctx)
	if err != nil {
		return "", err
	}

	newestVersionString, err = newestVersion.GetAsString()
	if err != nil {
		return "", err
	}

	return newestVersionString, nil
}

func (g *GitlabProject) GetNextMajorReleaseVersionString(ctx context.Context) (nextVersionString string, err error) {
	newestVersion, err := g.GetNewestSemanticVersion(ctx)
	if err != nil {
		return "", err
	}

	nextVersion, err := newestVersion.GetNextVersion("major")
	if err != nil {
		return "", err
	}

	nextVersionString, err = nextVersion.GetAsString()
	if err != nil {
		return "", err
	}

	return nextVersionString, nil
}

func (g *GitlabProject) GetNextMinorReleaseVersionString(ctx context.Context) (nextVersionString string, err error) {
	newestVersion, err := g.GetNewestSemanticVersion(ctx)
	if err != nil {
		return "", err
	}

	nextVersion, err := newestVersion.GetNextVersion("minor")
	if err != nil {
		return "", err
	}

	nextVersionString, err = nextVersion.GetAsString()
	if err != nil {
		return "", err
	}

	return nextVersionString, nil
}

func (g *GitlabProject) GetNextPatchReleaseVersionString(ctx context.Context) (nextVersionString string, err error) {
	newestVersion, err := g.GetNewestSemanticVersion(ctx)
	if err != nil {
		return "", err
	}

	nextVersion, err := newestVersion.GetNextVersion("patch")
	if err != nil {
		return "", err
	}

	nextVersionString, err = nextVersion.GetAsString()
	if err != nil {
		return "", err
	}

	return nextVersionString, nil
}

func (g *GitlabProject) GetOpenMergeRequestBySourceAndTargetBranch(ctx context.Context, sourceBranchName string, targetBranchName string) (mergeRequest *GitlabMergeRequest, err error) {
	mergeRequests, err := g.GetMergeRequests()
	if err != nil {
		return nil, err
	}

	mergeRequest, err = mergeRequests.GetOpenMergeRequestBySourceAndTargetBranch(ctx, sourceBranchName, targetBranchName)
	if err != nil {
		return nil, err
	}

	return mergeRequest, nil
}

func (g *GitlabProject) GetOpenMergeRequestByTitle(ctx context.Context, title string) (mergeRequest *GitlabMergeRequest, err error) {
	mergeRequests, err := g.GetMergeRequests()
	if err != nil {
		return nil, err
	}

	mergeRequest, err = mergeRequests.GetOpenMergeRequestByTitle(ctx, title)
	if err != nil {
		return nil, err
	}

	return mergeRequest, nil
}

func (g *GitlabProject) GetPath(ctx context.Context) (projectPath string, err error) {
	nativeProject, err := g.GetRawResponse(ctx)
	if err != nil {
		return "", err
	}

	projectPath = nativeProject.PathWithNamespace
	if projectPath == "" {
		return "", tracederrors.TracedError("projectPath is empty string after evaluation")
	}

	err = g.SetCachedPath(projectPath)
	if err != nil {
		return "", err
	}

	return projectPath, nil
}

func (g *GitlabProject) GetProjectCommits() (projectCommits *GitlabProjectCommits, err error) {
	projectCommits = NewGitlabProjectCommits()

	err = projectCommits.SetGitlabProject(g)
	if err != nil {
		return nil, err
	}

	return projectCommits, nil
}

func (g *GitlabProject) GetProjectUrl(ctx context.Context) (projectUrl string, err error) {
	fqdn, err := g.GetGitlabFqdn()
	if err != nil {
		return "", err
	}

	projectPath, err := g.GetPath(ctx)
	if err != nil {
		return "", err
	}

	projectUrl = fmt.Sprintf("https://%s/%s", fqdn, projectPath)

	return projectUrl, nil
}

func (g *GitlabProject) GetRawResponse(ctx context.Context) (nativeGitlabProject *gitlab.Project, err error) {
	gitlabProjects, err := g.GetGitlabProjects()
	if err != nil {
		return nil, err
	}

	projectsService, err := gitlabProjects.GetNativeProjectsService()
	if err != nil {
		return nil, err
	}

	isIdSet, err := g.IsIdSet()
	if err != nil {
		return nil, err
	}

	var pid interface{} = nil
	if isIdSet {
		pid, err = g.GetId(ctx)
		if err != nil {
			return nil, err
		}
	} else {
		if g.IsCachedPathSet() {
			isPersonalProject, err := g.IsPersonalProject(ctx)
			if err != nil {
				return nil, err
			}
			if isPersonalProject {
				pid, err = g.GetCachedPathForPersonalProject(ctx)
				if err != nil {
					return nil, err
				}
			} else {
				pid, err = g.GetCachedPath(ctx)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	if pid == nil {
		return nil, tracederrors.TracedErrorf("Unable to evaluate pid to get native gitlab project: '%w'", err)
	}

	nativeProject, _, err := projectsService.GetProject(pid, nil, nil)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Unable to get native project: '%w'", err)
	}

	if nativeProject == nil {
		return nil, tracederrors.TracedError("nativeProject is nil after evaluation")
	}

	return nativeProject, nil
}

func (g *GitlabProject) GetRepositoryFile(ctx context.Context, options *GitlabGetRepositoryFileOptions) (repositoryFile *GitlabRepositoryFile, err error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	repositoryFiles, err := g.GetRepositoryFiles()
	if err != nil {
		return nil, err
	}

	repositoryFile, err = repositoryFiles.GetRepositoryFile(options)
	if err != nil {
		return nil, err
	}

	return repositoryFile, nil
}

func (g *GitlabProject) GetRepositoryFiles() (repositoryFiles *GitlabRepositoryFiles, err error) {
	repositoryFiles = NewGitlabRepositoryFiles()

	err = repositoryFiles.SetGitlabProject(g)
	if err != nil {
		return nil, err
	}

	return repositoryFiles, nil
}

func (g *GitlabProject) GetSemanticVersions(ctx context.Context) (semanticVersions []versionutils.Version, err error) {
	versions, err := g.GetVersions(ctx)
	if err != nil {
		return nil, err
	}

	semanticVersions = []versionutils.Version{}
	for _, toAdd := range versions {
		if toAdd.IsSemanticVersion() {
			semanticVersions = append(semanticVersions, toAdd)
		}
	}

	return semanticVersions, nil
}

func (g *GitlabProject) GetTagByName(tagName string) (tag *GitlabTag, err error) {
	if tagName == "" {
		return nil, tracederrors.TracedErrorEmptyString("tagName")
	}

	tags, err := g.GetTags()
	if err != nil {
		return nil, err
	}

	tag, err = tags.GetTagByName(tagName)
	if err != nil {
		return nil, err
	}

	return tag, nil
}

func (g *GitlabProject) GetTags() (gitlabTags *GitlabTags, err error) {
	gitlabTags = NewGitlabTags()

	err = gitlabTags.SetGitlabProject(g)
	if err != nil {
		return nil, err
	}

	return gitlabTags, err
}

func (g *GitlabProject) GetVersionTags(ctx context.Context) (versionTags []*GitlabTag, err error) {
	tags, err := g.GetTags()
	if err != nil {
		return nil, err
	}

	versionTags, err = tags.GetVersionTags(ctx)
	if err != nil {
		return nil, err
	}

	return versionTags, nil
}

func (g *GitlabProject) GetVersions(ctx context.Context) (versions []versionutils.Version, err error) {
	versionTags, err := g.GetVersionTags(ctx)
	if err != nil {
		return nil, err
	}

	versions = []versionutils.Version{}

	for _, tag := range versionTags {
		versionName, err := tag.GetName()
		if err != nil {
			return nil, err
		}

		toAdd, err := versionutils.ReadFromString(versionName)
		if err != nil {
			return nil, err
		}

		versions = append(versions, toAdd)
	}

	return versions, nil
}

func (g *GitlabProject) IsCachedPathSet() (isSet bool) {
	return g.cachedPath != ""
}

func (g *GitlabProject) IsIdSet() (isSet bool, err error) {
	return g.id > 0, nil
}

func (g *GitlabProject) IsPersonalProject(ctx context.Context) (isPersonalProject bool, err error) {
	gitlabProjects, err := g.GetGitlabProjects()
	if err != nil {
		return false, err
	}

	projectPath, err := g.GetCachedPath(ctx)
	if err != nil {
		return false, err
	}

	isPersonalProject, err = gitlabProjects.IsProjectPathPersonalProject(projectPath)
	if err != nil {
		return false, err
	}

	return isPersonalProject, nil
}

func (g *GitlabProject) ListVersionTagNames(ctx context.Context) (versionTagNames []string, err error) {
	tags, err := g.GetTags()
	if err != nil {
		return nil, err
	}

	versionTagNames, err = tags.ListVersionTagNames(ctx)
	if err != nil {
		return nil, err
	}

	return versionTagNames, nil
}

func (g *GitlabProject) GetPipelineSchedules() (scheduledPipelines *GitlabPipelineSchedules, err error) {
	scheduledPipelines = NewGitlabPipelineSchedules()

	err = scheduledPipelines.SetGitlabProject(g)
	if err != nil {
		return nil, err
	}

	return scheduledPipelines, err
}

func (g *GitlabProject) ListScheduledPipelines(ctx context.Context) (scheduledPipelines []*GitlabPipelineSchedule, err error) {
	scheduled, err := g.GetPipelineSchedules()
	if err != nil {
		return nil, err
	}

	return scheduled.ListPipelineSchedules(ctx)
}

func (g *GitlabProject) ReadFileContentAsString(ctx context.Context, options *GitlabReadFileOptions) (content string, err error) {
	if options == nil {
		return "", tracederrors.TracedErrorNil("options")
	}

	repositoryFiles, err := g.GetRepositoryFiles()
	if err != nil {
		return "", err
	}

	content, err = repositoryFiles.ReadFileContentAsString(ctx, options)
	if err != nil {
		return "", err
	}

	return content, nil
}

func (g *GitlabProject) SetId(id int) (err error) {
	if id <= 0 {
		return tracederrors.TracedErrorf("invalid id = '%d'", id)
	}

	g.id = id

	return nil
}

func (g *GitlabProject) WriteFileContent(ctx context.Context, options *GitlabWriteFileOptions) (gitlabRepositoryFile *GitlabRepositoryFile, err error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	repositoryFiles, err := g.GetRepositoryFiles()
	if err != nil {
		return nil, err
	}

	gitlabRepositoryFile, err = repositoryFiles.WriteFileContent(ctx, options)
	if err != nil {
		return nil, err
	}

	return gitlabRepositoryFile, nil
}

func (g *GitlabProject) WriteFileContentInDefaultBranch(ctx context.Context, writeOptions *GitlabWriteFileOptions) (gitlabRepositoryFile *GitlabRepositoryFile, err error) {
	if writeOptions == nil {
		return nil, tracederrors.TracedErrorNil("writeOptions")
	}

	defaultBranch, err := g.GetDefaultBranch(ctx)
	if err != nil {
		return nil, err
	}

	gitlabRepositoryFile, err = defaultBranch.WriteFileContent(ctx, writeOptions)
	if err != nil {
		return nil, err
	}

	return gitlabRepositoryFile, nil
}

func (p *GitlabProject) CreateReleaseFromLatestCommitInDefaultBranch(ctx context.Context, createReleaseOptions *GitlabCreateReleaseOptions) (createdRelease *GitlabRelease, err error) {
	if createReleaseOptions == nil {
		return nil, tracederrors.TracedErrorNil("createReleaseOptions")
	}

	latestCommit, err := p.GetLatestCommitOfDefaultBranch(ctx)
	if err != nil {
		return nil, err
	}

	createdRelease, err = latestCommit.CreateRelease(ctx, createReleaseOptions)
	if err != nil {
		return nil, err
	}

	return createdRelease, nil
}

func (p *GitlabProject) DeployKeyByNameExists(ctx context.Context, keyName string) (exists bool, err error) {
	if len(keyName) <= 0 {
		return false, tracederrors.TracedError("keyName is empty string")
	}

	deployKeys, err := p.GetDeployKeys()
	if err != nil {
		return false, err
	}

	exists, err = deployKeys.DeployKeyByNameExists(ctx, keyName)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (p *GitlabProject) GetCachedPath(ctx context.Context) (path string, err error) {
	if !p.IsCachedPathSet() {
		_, err := p.GetPath(ctx)
		if err != nil {
			return "", err
		}
	}

	if len(p.cachedPath) <= 0 {
		return "", tracederrors.TracedError("cachedPath not set")
	}

	return p.cachedPath, nil
}

func (p *GitlabProject) GetDeployKeyByName(keyName string) (projectDeployKey *GitlabProjectDeployKey, err error) {
	if len(keyName) <= 0 {
		return nil, tracederrors.TracedError("keyName is nil")
	}

	deployKeys, err := p.GetGitlabProjectDeployKeys()
	if err != nil {
		return nil, err
	}

	projectDeployKey, err = deployKeys.GetGitlabProjectDeployKeyByName(keyName)
	if err != nil {
		return nil, err
	}

	return projectDeployKey, nil
}

func (p *GitlabProject) GetDeployKeys() (deployKeys *GitlabProjectDeployKeys, err error) {
	deployKeys = NewGitlabProjectDeployKeys()
	err = deployKeys.SetGitlabProject(p)
	if err != nil {
		return nil, err
	}

	return deployKeys, nil
}

func (p *GitlabProject) GetGitlab() (gitlab *GitlabInstance, err error) {
	if p.gitlab == nil {
		return nil, tracederrors.TracedError("gitlab is not set")
	}

	return p.gitlab, nil
}

func (p *GitlabProject) GetGitlabProjectDeployKeys() (projectDeployKeys *GitlabProjectDeployKeys, err error) {
	projectDeployKeys = NewGitlabProjectDeployKeys()

	err = projectDeployKeys.SetGitlabProject(p)
	if err != nil {
		return nil, err
	}

	return projectDeployKeys, nil
}

func (p *GitlabProject) GetGitlabProjects() (projects *GitlabProjects, err error) {
	gitlab, err := p.GetGitlab()
	if err != nil {
		return nil, err
	}

	projects, err = gitlab.GetGitlabProjects()
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (p *GitlabProject) GetGitlabReleases() (gitlabReleases *GitlabReleases, err error) {
	gitlabReleases = NewGitlabReleases()

	err = gitlabReleases.SetGitlabProject(p)
	if err != nil {
		return nil, err
	}

	return gitlabReleases, nil
}

func (p *GitlabProject) GetId(ctx context.Context) (id int, err error) {
	if p.id < 0 {
		return -1, tracederrors.TracedErrorf("id is set to invalid value of '%d'", p.id)
	}

	if p.id == 0 {
		rawResponse, err := p.GetRawResponse(ctx)
		if err != nil {
			return -1, err
		}

		id = rawResponse.ID
		if id <= 0 {
			return -1, tracederrors.TracedErrorf("GetId failed for GitlabProject: id is '%d' after evaluation", id)
		}
	}

	return p.id, nil
}

func (p *GitlabProject) GetNativeProjectsService() (nativeGitlabProject *gitlab.ProjectsService, err error) {
	projects, err := p.GetGitlabProjects()
	if err != nil {
		return nil, err
	}

	nativeGitlabProject, err = projects.GetNativeProjectsService()
	if err != nil {
		return nil, err
	}

	return nativeGitlabProject, nil
}

func (p *GitlabProject) GetReleaseByName(releaseName string) (gitlabRelease *GitlabRelease, err error) {
	if releaseName == "" {
		return nil, tracederrors.TracedErrorEmptyString("releaseName")
	}

	gitlabReleases, err := p.GetGitlabReleases()
	if err != nil {
		return nil, err
	}

	gitlabRelease, err = gitlabReleases.GetGitlabReleaseByName(releaseName)
	if err != nil {
		return nil, err
	}

	return gitlabRelease, nil
}

func (p *GitlabProject) MakePrivate(ctx context.Context) (err error) {
	nativeProjectsService, err := p.GetNativeProjectsService()
	if err != nil {
		return err
	}

	projectId, err := p.GetId(ctx)
	if err != nil {
		return err
	}

	var visibility = gitlab.PrivateVisibility

	_, _, err = nativeProjectsService.EditProject(
		projectId,
		&gitlab.EditProjectOptions{
			Visibility: &visibility,
		},
	)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Gitlab project '%v' made private.", projectId)

	return nil
}

func (p *GitlabProject) MakePublic(ctx context.Context) (err error) {
	nativeProjectsService, err := p.GetNativeProjectsService()
	if err != nil {
		return err
	}

	projectId, err := p.GetId(ctx)
	if err != nil {
		return err
	}

	var visibility = gitlab.PublicVisibility

	_, _, err = nativeProjectsService.EditProject(
		projectId,
		&gitlab.EditProjectOptions{
			Visibility: &visibility,
		},
	)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Gitlab project '%v' made public.", projectId)

	return nil
}

func (p *GitlabProject) RecreateDeployKey(ctx context.Context, keyOptions *GitlabCreateDeployKeyOptions) (err error) {
	if keyOptions == nil {
		return tracederrors.TracedError("keyOptions is nil")
	}

	deployKey, err := p.GetDeployKeyByName(keyOptions.Name)
	if err != nil {
		return err
	}

	err = deployKey.RecreateDeployKey(ctx, keyOptions)
	if err != nil {
		return err
	}

	return nil
}

func (p *GitlabProject) SetCachedPath(pathToCache string) (err error) {
	if len(pathToCache) <= 0 {
		return tracederrors.TracedError("pathToCache is empty string")
	}

	p.cachedPath = pathToCache

	return nil
}

func (p *GitlabProject) SetGitlab(gitlab *GitlabInstance) (err error) {
	if gitlab == nil {
		return tracederrors.TracedError("gitlab is nil")
	}

	p.gitlab = gitlab

	return nil
}
