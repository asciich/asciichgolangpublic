package asciichgolangpublic

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/parameteroptions/authenticationoptions"
	"github.com/asciich/asciichgolangpublic/pkg/versionutils"
	"github.com/asciich/asciichgolangpublic/tracederrors"
	"github.com/asciich/asciichgolangpublic/urlsutils"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type GitlabProject struct {
	gitlab     *GitlabInstance
	id         int
	cachedPath string
}

func GetGitlabProjectByUrl(url *urlsutils.URL, authOptions []authenticationoptions.AuthenticationOption, verbose bool) (gitlabProject *GitlabProject, err error) {
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

	if authOptions != nil {
		err = gitlab.Authenticate(gitlabAuthenticationOption)
		if err != nil {
			return nil, err
		}
	}

	gitlabProject, err = gitlab.GetGitlabProjectByPath(path, verbose)
	if err != nil {
		return nil, err
	}

	return gitlabProject, err
}

func GetGitlabProjectByUrlFromString(urlString string, authOptions []authenticationoptions.AuthenticationOption, verbose bool) (gitlabProject *GitlabProject, err error) {
	if urlString == "" {
		return nil, tracederrors.TracedErrorEmptyString("urlString")
	}

	url, err := urlsutils.GetUrlFromString(urlString)
	if err != nil {
		return nil, err
	}

	return GetGitlabProjectByUrl(url, authOptions, verbose)
}

func MustGetGitlabProjectByUrl(url *urlsutils.URL, authOptions []authenticationoptions.AuthenticationOption, verbose bool) (gitlabProject *GitlabProject) {
	gitlabProject, err := GetGitlabProjectByUrl(url, authOptions, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlabProject
}

func MustGetGitlabProjectByUrlFromString(urlString string, authOptions []authenticationoptions.AuthenticationOption, verbose bool) (gitlabProject *GitlabProject) {
	gitlabProject, err := GetGitlabProjectByUrlFromString(urlString, authOptions, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlabProject
}

func NewGitlabProject() (gitlabProject *GitlabProject) {
	return new(GitlabProject)
}

func (g *GitlabProject) Create(verbose bool) (err error) {
	gitlabProjects, err := g.GetGitlabProjects()
	if err != nil {
		return err
	}

	projectPath, err := g.GetCachedPath()
	if err != nil {
		return err
	}

	createdProject, err := gitlabProjects.CreateProject(
		&GitlabCreateProjectOptions{
			ProjectPath: projectPath,
			Verbose:     verbose,
		},
	)
	if err != nil {
		return err
	}

	createdProjectId, err := createdProject.GetId()
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

func (g *GitlabProject) CreateBranchFromDefaultBranch(branchName string, verbose bool) (createdBranch *GitlabBranch, err error) {
	branches, err := g.GetBranches()
	if err != nil {
		return nil, err
	}

	createdBranch, err = branches.CreateBranchFromDefaultBranch(branchName, verbose)
	if err != nil {
		return nil, err
	}

	return createdBranch, nil
}

func (g *GitlabProject) CreateEmptyFile(fileName string, ref string, verbose bool) (createdFile *GitlabRepositoryFile, err error) {
	repositoryFiles, err := g.GetRepositoryFiles()
	if err != nil {
		return nil, err
	}

	createdFile, err = repositoryFiles.CreateEmptyFile(fileName, ref, verbose)
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

func (g *GitlabProject) MustListScheduledPipelineNames(ctx context.Context) (scheduledPipelineNames []string) {
	scheduledPipelineNames, err := g.ListScheduledPipelineNames(ctx)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return scheduledPipelineNames
}

func (g *GitlabProject) CreateMergeRequest(options *GitlabCreateMergeRequestOptions) (createdMergeRequest *GitlabMergeRequest, err error) {
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

	createdMergeRequest, err = sourceBranch.CreateMergeRequest(options)
	if err != nil {
		return nil, err
	}

	return createdMergeRequest, nil
}

func (g *GitlabProject) CreateNextMajorReleaseFromLatestCommitInDefaultBranch(description string, verbose bool) (createdRelease *GitlabRelease, err error) {
	if description == "" {
		return nil, tracederrors.TracedErrorEmptyString("description")
	}

	nextPatchVersionString, err := g.GetNextMajorReleaseVersionString(verbose)
	if err != nil {
		return nil, err
	}

	createdRelease, err = g.CreateReleaseFromLatestCommitInDefaultBranch(
		&GitlabCreateReleaseOptions{
			Name:        nextPatchVersionString,
			Description: description,
			Verbose:     verbose,
		},
	)
	if err != nil {
		return nil, err
	}

	return createdRelease, nil
}

func (g *GitlabProject) CreateNextMinorReleaseFromLatestCommitInDefaultBranch(description string, verbose bool) (createdRelease *GitlabRelease, err error) {
	if description == "" {
		return nil, tracederrors.TracedErrorEmptyString("description")
	}

	nextPatchVersionString, err := g.GetNextMinorReleaseVersionString(verbose)
	if err != nil {
		return nil, err
	}

	createdRelease, err = g.CreateReleaseFromLatestCommitInDefaultBranch(
		&GitlabCreateReleaseOptions{
			Name:        nextPatchVersionString,
			Description: description,
			Verbose:     verbose,
		},
	)
	if err != nil {
		return nil, err
	}

	return createdRelease, nil
}

func (g *GitlabProject) CreateNextPatchReleaseFromLatestCommitInDefaultBranch(description string, verbose bool) (createdRelease *GitlabRelease, err error) {
	if description == "" {
		return nil, tracederrors.TracedErrorEmptyString("description")
	}

	nextPatchVersionString, err := g.GetNextPatchReleaseVersionString(verbose)
	if err != nil {
		return nil, err
	}

	createdRelease, err = g.CreateReleaseFromLatestCommitInDefaultBranch(
		&GitlabCreateReleaseOptions{
			Name:        nextPatchVersionString,
			Description: description,
			Verbose:     verbose,
		},
	)
	if err != nil {
		return nil, err
	}

	return createdRelease, nil
}

func (g *GitlabProject) Delete(verbose bool) (err error) {
	gitlabProjects, err := g.GetGitlabProjects()
	if err != nil {
		return err
	}

	projectPath, err := g.GetCachedPath()
	if err != nil {
		return err
	}

	err = gitlabProjects.DeleteProject(
		&GitlabDeleteProjectOptions{
			ProjectPath: projectPath,
			Verbose:     verbose,
		},
	)
	if err != nil {
		return err
	}

	g.id = 0

	return nil
}

func (g *GitlabProject) DeleteAllBranchesExceptDefaultBranch(verbose bool) (err error) {
	branches, err := g.GetBranches()
	if err != nil {
		return err
	}

	err = branches.DeleteAllBranchesExceptDefaultBranch(verbose)
	if err != nil {
		return err
	}

	return nil
}

func (g *GitlabProject) DeleteAllReleases(deleteOptions *GitlabDeleteReleaseOptions) (err error) {
	releases, err := g.GetGitlabReleases()
	if err != nil {
		return err
	}

	err = releases.DeleteAllReleases(deleteOptions)
	if err != nil {
		return err
	}

	return nil
}

func (g *GitlabProject) DeleteBranch(branchName string, deleteOptions *GitlabDeleteBranchOptions) (err error) {
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

	err = branch.Delete(deleteOptions)
	if err != nil {
		return err
	}

	return nil
}

func (g *GitlabProject) DeleteFileInDefaultBranch(fileName string, commitMessage string, verbose bool) (err error) {
	if fileName == "" {
		return tracederrors.TracedErrorEmptyString("fileName")
	}

	if commitMessage == "" {
		return tracederrors.TracedErrorEmptyString("commitMessage")
	}

	fileInRepo, err := g.GetFileInDefaultBranch(
		fileName,
		verbose,
	)
	if err != nil {
		return err
	}

	err = fileInRepo.Delete(commitMessage, verbose)
	if err != nil {
		return err
	}

	return nil
}

func (g *GitlabProject) Exists(verbose bool) (projectExists bool, err error) {
	gitlab, err := g.GetGitlab()
	if err != nil {
		return false, err
	}

	idSet, err := g.IsIdSet()
	if err != nil {
		return false, err
	}
	if idSet {
		projectId, err := g.GetId()
		if err != nil {
			return false, err
		}

		projectExists, err = gitlab.ProjectByProjectIdExists(projectId, verbose)
		if err != nil {
			return false, err
		}
	} else {
		projectPath, err := g.GetCachedPath()
		if err != nil {
			return false, err
		}

		projectExists, err = gitlab.ProjectByProjectPathExists(projectPath, verbose)
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

func (g *GitlabProject) GetBranchNames(verbose bool) (branchNames []string, err error) {
	branches, err := g.GetBranches()
	if err != nil {
		return nil, err
	}

	branchNames, err = branches.GetBranchNames(verbose)
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

func (g *GitlabProject) GetCachedPathForPersonalProject() (cachedPath string, err error) {
	projectName, err := g.GetCachedProjectName()
	if err != nil {
		return "", err
	}

	const verbose = false
	userName, err := g.GetCurrentUserName(verbose)
	if err != nil {
		return "", err
	}

	cachedPath = fmt.Sprintf("%s/%s", userName, projectName)

	return cachedPath, nil
}

func (g *GitlabProject) GetCachedProjectName() (projectName string, err error) {
	cachedPath, err := g.GetCachedPath()
	if err != nil {
		return "", err
	}

	projectName = filepath.Base(cachedPath)
	if projectName == "" {
		return "", tracederrors.TracedErrorf("Unable to extract project name from cachedPath = '%s'", cachedPath)
	}

	return projectName, nil
}

func (g *GitlabProject) GetCommitByHashString(hashString string, verbose bool) (commit *GitlabCommit, err error) {
	if hashString == "" {
		return nil, tracederrors.TracedErrorNil("hashString")
	}

	projectCommits, err := g.GetProjectCommits()
	if err != nil {
		return nil, err
	}

	commit, err = projectCommits.GetCommitByHashString(hashString, verbose)
	if err != nil {
		return nil, err
	}

	return commit, nil
}

func (g *GitlabProject) GetCurrentUserName(verbose bool) (userName string, err error) {
	gitlab, err := g.GetGitlab()
	if err != nil {
		return "", err
	}

	userName, err = gitlab.GetCurrentUsersName(verbose)
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

func (g *GitlabProject) GetDefaultBranch() (defaultBranch *GitlabBranch, err error) {
	defaultBranchName, err := g.GetDefaultBranchName()
	if err != nil {
		return nil, err
	}

	defaultBranch, err = g.GetBranchByName(defaultBranchName)
	if err != nil {
		return nil, err
	}

	return defaultBranch, nil
}

func (g *GitlabProject) GetDefaultBranchName() (defaultBranchName string, err error) {
	nativeProject, err := g.GetRawResponse()
	if err != nil {
		return "", err
	}

	defaultBranchName = nativeProject.DefaultBranch
	if defaultBranchName == "" {
		return "", tracederrors.TracedError("defaultBranchName is empty string after evaluation")
	}

	return defaultBranchName, nil
}

func (g *GitlabProject) GetDirectoryNames(ref string, verbose bool) (directoryNames []string, err error) {
	repositoryFiles, err := g.GetRepositoryFiles()
	if err != nil {
		return nil, err
	}

	directoryNames, err = repositoryFiles.GetDirectoryNames(ref, verbose)
	if err != nil {
		return nil, err
	}

	return directoryNames, nil
}

func (g *GitlabProject) GetFileInDefaultBranch(fileName string, verbose bool) (repositoryFile *GitlabRepositoryFile, err error) {
	if fileName == "" {
		return nil, tracederrors.TracedErrorEmptyString("fileName")
	}

	repositoryFile, err = g.GetRepositoryFile(
		&GitlabGetRepositoryFileOptions{
			Path:    fileName,
			Verbose: verbose,
		},
	)
	if err != nil {
		return nil, err
	}

	return repositoryFile, nil
}

func (g *GitlabProject) GetFilesNames(ref string, verbose bool) (fileNames []string, err error) {
	repositoryFiles, err := g.GetRepositoryFiles()
	if err != nil {
		return nil, err
	}

	fileNames, err = repositoryFiles.GetFileNames(ref, verbose)
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

func (g *GitlabProject) GetLatestCommit(branchName string, verbose bool) (latestCommit *GitlabCommit, err error) {
	if branchName == "" {
		return nil, tracederrors.TracedErrorNil("branchName")
	}

	latestHash, err := g.GetLatestCommitHashAsString(branchName, verbose)
	if err != nil {
		return nil, err
	}

	latestCommit, err = g.GetCommitByHashString(latestHash, verbose)
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

func (g *GitlabProject) GetLatestCommitHashAsString(branchName string, verbose bool) (commitHash string, err error) {
	if branchName == "" {
		return "", tracederrors.TracedErrorEmptyString("branchName")
	}

	branch, err := g.GetBranchByName(branchName)
	if err != nil {
		return "", err
	}

	commitHash, err = branch.GetLatestCommitHashAsString(verbose)
	if err != nil {
		return "", err
	}

	return commitHash, nil
}

func (g *GitlabProject) GetLatestCommitOfDefaultBranch(verbose bool) (latestCommit *GitlabCommit, err error) {
	defaultBranch, err := g.GetDefaultBranchName()
	if err != nil {
		return nil, err
	}

	latestCommit, err = g.GetLatestCommit(defaultBranch, verbose)
	if err != nil {
		return nil, err
	}

	if verbose {
		latestHash, err := latestCommit.GetCommitHash()
		if err != nil {
			return nil, err
		}

		logging.LogInfof(
			"Latest commit of default branch '%s' has hash '%s'.",
			defaultBranch,
			latestHash,
		)
	}

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

func (g *GitlabProject) GetNewestSemanticVersion(verbose bool) (newestSemanticVersion *versionutils.VersionSemanticVersion, err error) {
	semanticVersions, err := g.GetSemanticVersions(verbose)
	if err != nil {
		return nil, err
	}

	newestVersion, err := versionutils.Versions().GetLatestVersionFromSlice(semanticVersions)
	if err != nil {
		return nil, err
	}

	newestSemanticVersion, ok := newestVersion.(*versionutils.VersionSemanticVersion)
	if !ok {
		return nil, tracederrors.TracedErrorf(
			"Unable to get newest semantiv version from '%v'",
			newestVersion,
		)
	}

	return newestSemanticVersion, nil
}

func (g *GitlabProject) GetNewestVersion(verbose bool) (newestVersion versionutils.Version, err error) {
	availableVersions, err := g.GetVersions(verbose)
	if err != nil {
		return nil, err
	}

	if len(availableVersions) <= 0 {
		return nil, tracederrors.TracedError("No versionTags returned")
	}

	newestVersion, err = versionutils.Versions().GetLatestVersionFromSlice(availableVersions)
	if err != nil {
		return nil, err
	}

	return newestVersion, nil
}

func (g *GitlabProject) GetNewestVersionAsString(verbose bool) (newestVersionString string, err error) {
	newestVersion, err := g.GetNewestVersion(verbose)
	if err != nil {
		return "", err
	}

	newestVersionString, err = newestVersion.GetAsString()
	if err != nil {
		return "", err
	}

	return newestVersionString, nil
}

func (g *GitlabProject) GetNextMajorReleaseVersionString(verbose bool) (nextVersionString string, err error) {
	newestVersion, err := g.GetNewestSemanticVersion(verbose)
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

func (g *GitlabProject) GetNextMinorReleaseVersionString(verbose bool) (nextVersionString string, err error) {
	newestVersion, err := g.GetNewestSemanticVersion(verbose)
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

func (g *GitlabProject) GetNextPatchReleaseVersionString(verbose bool) (nextVersionString string, err error) {
	newestVersion, err := g.GetNewestSemanticVersion(verbose)
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

func (g *GitlabProject) GetOpenMergeRequestBySourceAndTargetBranch(sourceBranchName string, targetBranchName string, verbose bool) (mergeRequest *GitlabMergeRequest, err error) {
	mergeRequests, err := g.GetMergeRequests()
	if err != nil {
		return nil, err
	}

	mergeRequest, err = mergeRequests.GetOpenMergeRequestBySourceAndTargetBranch(sourceBranchName, targetBranchName, verbose)
	if err != nil {
		return nil, err
	}

	return mergeRequest, nil
}

func (g *GitlabProject) GetOpenMergeRequestByTitle(title string, verbose bool) (mergeRequest *GitlabMergeRequest, err error) {
	mergeRequests, err := g.GetMergeRequests()
	if err != nil {
		return nil, err
	}

	mergeRequest, err = mergeRequests.GetOpenMergeRequestByTitle(title, verbose)
	if err != nil {
		return nil, err
	}

	return mergeRequest, nil
}

func (g *GitlabProject) GetPath() (projectPath string, err error) {
	nativeProject, err := g.GetRawResponse()
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

func (g *GitlabProject) GetProjectUrl() (projectUrl string, err error) {
	fqdn, err := g.GetGitlabFqdn()
	if err != nil {
		return "", err
	}

	projectPath, err := g.GetPath()
	if err != nil {
		return "", err
	}

	projectUrl = fmt.Sprintf("https://%s/%s", fqdn, projectPath)

	return projectUrl, nil
}

func (g *GitlabProject) GetRawResponse() (nativeGitlabProject *gitlab.Project, err error) {
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
		pid, err = g.GetId()
		if err != nil {
			return nil, err
		}
	} else {
		if g.IsCachedPathSet() {
			isPersonalProject, err := g.IsPersonalProject()
			if err != nil {
				return nil, err
			}
			if isPersonalProject {
				pid, err = g.GetCachedPathForPersonalProject()
				if err != nil {
					return nil, err
				}
			} else {
				pid, err = g.GetCachedPath()
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

func (g *GitlabProject) GetRepositoryFile(options *GitlabGetRepositoryFileOptions) (repositoryFile *GitlabRepositoryFile, err error) {
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

func (g *GitlabProject) GetSemanticVersions(verbose bool) (semanticVersions []versionutils.Version, err error) {
	versions, err := g.GetVersions(verbose)
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

func (g *GitlabProject) GetVersionTags(verbose bool) (versionTags []*GitlabTag, err error) {
	tags, err := g.GetTags()
	if err != nil {
		return nil, err
	}

	versionTags, err = tags.GetVersionTags(verbose)
	if err != nil {
		return nil, err
	}

	return versionTags, nil
}

func (g *GitlabProject) GetVersions(verbose bool) (versions []versionutils.Version, err error) {
	versionTags, err := g.GetVersionTags(verbose)
	if err != nil {
		return nil, err
	}

	versions = []versionutils.Version{}

	for _, tag := range versionTags {
		versionName, err := tag.GetName()
		if err != nil {
			return nil, err
		}

		toAdd, err := versionutils.Versions().GetNewVersionByString(versionName)
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

func (g *GitlabProject) IsPersonalProject() (isPersonalProject bool, err error) {
	gitlabProjects, err := g.GetGitlabProjects()
	if err != nil {
		return false, err
	}

	projectPath, err := g.GetCachedPath()
	if err != nil {
		return false, err
	}

	isPersonalProject, err = gitlabProjects.IsProjectPathPersonalProject(projectPath)
	if err != nil {
		return false, err
	}

	return isPersonalProject, nil
}

func (g *GitlabProject) ListVersionTagNames(verbose bool) (versionTagNames []string, err error) {
	tags, err := g.GetTags()
	if err != nil {
		return nil, err
	}

	versionTagNames, err = tags.ListVersionTagNames(verbose)
	if err != nil {
		return nil, err
	}

	return versionTagNames, nil
}

func (g *GitlabProject) MustCreate(verbose bool) {
	err := g.Create(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabProject) MustCreateBranchFromDefaultBranch(branchName string, verbose bool) (createdBranch *GitlabBranch) {
	createdBranch, err := g.CreateBranchFromDefaultBranch(branchName, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return createdBranch
}

func (g *GitlabProject) MustCreateEmptyFile(fileName string, ref string, verbose bool) (createdFile *GitlabRepositoryFile) {
	createdFile, err := g.CreateEmptyFile(fileName, ref, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return createdFile
}

func (g *GitlabProject) MustCreateMergeRequest(options *GitlabCreateMergeRequestOptions) (createdMergeRequest *GitlabMergeRequest) {
	createdMergeRequest, err := g.CreateMergeRequest(options)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return createdMergeRequest
}

func (g *GitlabProject) MustCreateNextMajorReleaseFromLatestCommitInDefaultBranch(description string, verbose bool) (createdRelease *GitlabRelease) {
	createdRelease, err := g.CreateNextMajorReleaseFromLatestCommitInDefaultBranch(description, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return createdRelease
}

func (g *GitlabProject) MustCreateNextMinorReleaseFromLatestCommitInDefaultBranch(description string, verbose bool) (createdRelease *GitlabRelease) {
	createdRelease, err := g.CreateNextMinorReleaseFromLatestCommitInDefaultBranch(description, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return createdRelease
}

func (g *GitlabProject) MustCreateNextPatchReleaseFromLatestCommitInDefaultBranch(description string, verbose bool) (createdRelease *GitlabRelease) {
	createdRelease, err := g.CreateNextPatchReleaseFromLatestCommitInDefaultBranch(description, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return createdRelease
}

func (g *GitlabProject) MustCreateReleaseFromLatestCommitInDefaultBranch(createReleaseOptions *GitlabCreateReleaseOptions) (createdRelease *GitlabRelease) {
	createdRelease, err := g.CreateReleaseFromLatestCommitInDefaultBranch(createReleaseOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return createdRelease
}

func (g *GitlabProject) MustDelete(verbose bool) {
	err := g.Delete(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabProject) MustDeleteAllBranchesExceptDefaultBranch(verbose bool) {
	err := g.DeleteAllBranchesExceptDefaultBranch(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabProject) MustDeleteAllReleases(deleteOptions *GitlabDeleteReleaseOptions) {
	err := g.DeleteAllReleases(deleteOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabProject) MustDeleteBranch(branchName string, deleteOptions *GitlabDeleteBranchOptions) {
	err := g.DeleteBranch(branchName, deleteOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabProject) MustDeleteFileInDefaultBranch(fileName string, commitMessage string, verbose bool) {
	err := g.DeleteFileInDefaultBranch(fileName, commitMessage, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabProject) MustDeployKeyByNameExists(keyName string) (exists bool) {
	exists, err := g.DeployKeyByNameExists(keyName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return exists
}

func (g *GitlabProject) MustExists(verbose bool) (projectExists bool) {
	projectExists, err := g.Exists(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return projectExists
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

func (g *GitlabProject) MustGetBranchByName(branchName string) (branch *GitlabBranch) {
	branch, err := g.GetBranchByName(branchName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return branch
}

func (g *GitlabProject) MustGetBranchNames(verbose bool) (branchNames []string) {
	branchNames, err := g.GetBranchNames(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return branchNames
}

func (g *GitlabProject) MustGetBranches() (branches *GitlabBranches) {
	branches, err := g.GetBranches()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return branches
}

func (g *GitlabProject) MustGetCachedPath() (path string) {
	path, err := g.GetCachedPath()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return path
}

func (g *GitlabProject) MustGetCachedPathForPersonalProject() (cachedPath string) {
	cachedPath, err := g.GetCachedPathForPersonalProject()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return cachedPath
}

func (g *GitlabProject) MustGetCachedProjectName() (projectName string) {
	projectName, err := g.GetCachedProjectName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return projectName
}

func (g *GitlabProject) MustGetCommitByHashString(hashString string, verbose bool) (commit *GitlabCommit) {
	commit, err := g.GetCommitByHashString(hashString, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return commit
}

func (g *GitlabProject) MustGetCurrentUserName(verbose bool) (userName string) {
	userName, err := g.GetCurrentUserName(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return userName
}

func (g *GitlabProject) MustGetDefaultBranch() (defaultBranch *GitlabBranch) {
	defaultBranch, err := g.GetDefaultBranch()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return defaultBranch
}

func (g *GitlabProject) MustGetDefaultBranchName() (defaultBranchName string) {
	defaultBranchName, err := g.GetDefaultBranchName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return defaultBranchName
}

func (g *GitlabProject) MustGetDeployKeyByName(keyName string) (projectDeployKey *GitlabProjectDeployKey) {
	projectDeployKey, err := g.GetDeployKeyByName(keyName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return projectDeployKey
}

func (g *GitlabProject) MustGetDeployKeys() (deployKeys *GitlabProjectDeployKeys) {
	deployKeys, err := g.GetDeployKeys()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return deployKeys
}

func (g *GitlabProject) MustGetDirectoryNames(ref string, verbose bool) (directoryNames []string) {
	directoryNames, err := g.GetDirectoryNames(ref, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return directoryNames
}

func (g *GitlabProject) MustGetFileInDefaultBranch(fileName string, verbose bool) (repositoryFile *GitlabRepositoryFile) {
	repositoryFile, err := g.GetFileInDefaultBranch(fileName, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return repositoryFile
}

func (g *GitlabProject) MustGetFilesNames(ref string, verbose bool) (fileNames []string) {
	fileNames, err := g.GetFilesNames(ref, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return fileNames
}

func (g *GitlabProject) MustGetGitlab() (gitlab *GitlabInstance) {
	gitlab, err := g.GetGitlab()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlab
}

func (g *GitlabProject) MustGetGitlabFqdn() (fqdn string) {
	fqdn, err := g.GetGitlabFqdn()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return fqdn
}

func (g *GitlabProject) MustGetGitlabProjectDeployKeys() (projectDeployKeys *GitlabProjectDeployKeys) {
	projectDeployKeys, err := g.GetGitlabProjectDeployKeys()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return projectDeployKeys
}

func (g *GitlabProject) MustGetGitlabProjects() (projects *GitlabProjects) {
	projects, err := g.GetGitlabProjects()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return projects
}

func (g *GitlabProject) MustGetGitlabReleases() (gitlabReleases *GitlabReleases) {
	gitlabReleases, err := g.GetGitlabReleases()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlabReleases
}

func (g *GitlabProject) MustGetId() (id int) {
	id, err := g.GetId()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return id
}

func (g *GitlabProject) MustGetLatestCommit(branchName string, verbose bool) (latestCommit *GitlabCommit) {
	latestCommit, err := g.GetLatestCommit(branchName, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return latestCommit
}

func (g *GitlabProject) MustGetLatestCommitHashAsString(branchName string, verbose bool) (commitHash string) {
	commitHash, err := g.GetLatestCommitHashAsString(branchName, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return commitHash
}

func (g *GitlabProject) MustGetLatestCommitOfDefaultBranch(verbose bool) (latestCommit *GitlabCommit) {
	latestCommit, err := g.GetLatestCommitOfDefaultBranch(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return latestCommit
}

func (g *GitlabProject) MustGetMergeRequests() (mergeRequestes *GitlabProjectMergeRequests) {
	mergeRequestes, err := g.GetMergeRequests()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return mergeRequestes
}

func (g *GitlabProject) MustGetNativeGitlabProject() (nativeGitlabProject *gitlab.Project) {
	nativeGitlabProject, err := g.GetRawResponse()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return nativeGitlabProject
}

func (g *GitlabProject) MustGetNativeProjectsService() (nativeGitlabProject *gitlab.ProjectsService) {
	nativeGitlabProject, err := g.GetNativeProjectsService()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return nativeGitlabProject
}

func (g *GitlabProject) MustGetNewestSemanticVersion(verbose bool) (newestSemanticVersion *versionutils.VersionSemanticVersion) {
	newestSemanticVersion, err := g.GetNewestSemanticVersion(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return newestSemanticVersion
}

func (g *GitlabProject) MustGetNewestVersion(verbose bool) (newestVersion versionutils.Version) {
	newestVersion, err := g.GetNewestVersion(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return newestVersion
}

func (g *GitlabProject) MustGetNewestVersionAsString(verbose bool) (newestVersionString string) {
	newestVersionString, err := g.GetNewestVersionAsString(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return newestVersionString
}

func (g *GitlabProject) MustGetNextMajorReleaseVersionString(verbose bool) (nextVersionString string) {
	nextVersionString, err := g.GetNextMajorReleaseVersionString(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return nextVersionString
}

func (g *GitlabProject) MustGetNextMinorReleaseVersionString(verbose bool) (nextVersionString string) {
	nextVersionString, err := g.GetNextMinorReleaseVersionString(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return nextVersionString
}

func (g *GitlabProject) MustGetNextPatchReleaseVersionString(verbose bool) (nextVersionString string) {
	nextVersionString, err := g.GetNextPatchReleaseVersionString(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return nextVersionString
}

func (g *GitlabProject) MustGetOpenMergeRequestBySourceAndTargetBranch(sourceBranchName string, targetBranchName string, verbose bool) (mergeRequest *GitlabMergeRequest) {
	mergeRequest, err := g.GetOpenMergeRequestBySourceAndTargetBranch(sourceBranchName, targetBranchName, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return mergeRequest
}

func (g *GitlabProject) MustGetOpenMergeRequestByTitle(title string, verbose bool) (mergeRequest *GitlabMergeRequest) {
	mergeRequest, err := g.GetOpenMergeRequestByTitle(title, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return mergeRequest
}

func (g *GitlabProject) MustGetPath() (projectPath string) {
	projectPath, err := g.GetPath()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return projectPath
}

func (g *GitlabProject) MustGetProjectCommits() (projectCommits *GitlabProjectCommits) {
	projectCommits, err := g.GetProjectCommits()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return projectCommits
}

func (g *GitlabProject) MustGetProjectUrl() (projectUrl string) {
	projectUrl, err := g.GetProjectUrl()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return projectUrl
}

func (g *GitlabProject) MustGetRawResponse() (nativeGitlabProject *gitlab.Project) {
	nativeGitlabProject, err := g.GetRawResponse()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return nativeGitlabProject
}

func (g *GitlabProject) MustGetReleaseByName(releaseName string) (gitlabRelease *GitlabRelease) {
	gitlabRelease, err := g.GetReleaseByName(releaseName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlabRelease
}

func (g *GitlabProject) MustGetRepositoryFile(options *GitlabGetRepositoryFileOptions) (repositoryFile *GitlabRepositoryFile) {
	repositoryFile, err := g.GetRepositoryFile(options)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return repositoryFile
}

func (g *GitlabProject) MustGetRepositoryFiles() (repositoryFiles *GitlabRepositoryFiles) {
	repositoryFiles, err := g.GetRepositoryFiles()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return repositoryFiles
}

func (g *GitlabProject) MustGetSemanticVersions(verbose bool) (semanticVersions []versionutils.Version) {
	semanticVersions, err := g.GetSemanticVersions(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return semanticVersions
}

func (g *GitlabProject) MustGetTagByName(tagName string) (tag *GitlabTag) {
	tag, err := g.GetTagByName(tagName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return tag
}

func (g *GitlabProject) MustGetTags() (gitlabTags *GitlabTags) {
	gitlabTags, err := g.GetTags()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlabTags
}

func (g *GitlabProject) MustGetVersionTags(verbose bool) (versionTags []*GitlabTag) {
	versionTags, err := g.GetVersionTags(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return versionTags
}

func (g *GitlabProject) MustGetVersions(verbose bool) (versions []versionutils.Version) {
	versions, err := g.GetVersions(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return versions
}

func (g *GitlabProject) MustIsIdSet() (isSet bool) {
	isSet, err := g.IsIdSet()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return isSet
}

func (g *GitlabProject) MustIsPersonalProject() (isPersonalProject bool) {
	isPersonalProject, err := g.IsPersonalProject()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return isPersonalProject
}

func (g *GitlabProject) MustListVersionTagNames(verbose bool) (versionTagNames []string) {
	versionTagNames, err := g.ListVersionTagNames(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return versionTagNames
}

func (g *GitlabProject) MustMakePrivate(verbose bool) {
	err := g.MakePrivate(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabProject) MustMakePublic(verbose bool) {
	err := g.MakePublic(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabProject) MustReadFileContentAsString(options *GitlabReadFileOptions) (content string) {
	content, err := g.ReadFileContentAsString(options)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return content
}

func (g *GitlabProject) MustRecreateDeployKey(keyOptions *GitlabCreateDeployKeyOptions) {
	err := g.RecreateDeployKey(keyOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabProject) MustSetCachedPath(pathToCache string) {
	err := g.SetCachedPath(pathToCache)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabProject) MustSetGitlab(gitlab *GitlabInstance) {
	err := g.SetGitlab(gitlab)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabProject) MustSetId(id int) {
	err := g.SetId(id)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabProject) MustWriteFileContent(options *GitlabWriteFileOptions) (gitlabRepositoryFile *GitlabRepositoryFile) {
	gitlabRepositoryFile, err := g.WriteFileContent(options)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlabRepositoryFile
}

func (g *GitlabProject) MustWriteFileContentInDefaultBranch(writeOptions *GitlabWriteFileOptions) (gitlabRepositoryFile *GitlabRepositoryFile) {
	gitlabRepositoryFile, err := g.WriteFileContentInDefaultBranch(writeOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlabRepositoryFile
}

func (g *GitlabProject) ReadFileContentAsString(options *GitlabReadFileOptions) (content string, err error) {
	if options == nil {
		return "", tracederrors.TracedErrorNil("options")
	}

	repositoryFiles, err := g.GetRepositoryFiles()
	if err != nil {
		return "", err
	}

	content, err = repositoryFiles.ReadFileContentAsString(options)
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

func (g *GitlabProject) WriteFileContent(options *GitlabWriteFileOptions) (gitlabRepositoryFile *GitlabRepositoryFile, err error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	repositoryFiles, err := g.GetRepositoryFiles()
	if err != nil {
		return nil, err
	}

	gitlabRepositoryFile, err = repositoryFiles.WriteFileContent(options)
	if err != nil {
		return nil, err
	}

	return gitlabRepositoryFile, nil
}

func (g *GitlabProject) WriteFileContentInDefaultBranch(writeOptions *GitlabWriteFileOptions) (gitlabRepositoryFile *GitlabRepositoryFile, err error) {
	if writeOptions == nil {
		return nil, tracederrors.TracedErrorNil("writeOptions")
	}

	defaultBranch, err := g.GetDefaultBranch()
	if err != nil {
		return nil, err
	}

	gitlabRepositoryFile, err = defaultBranch.WriteFileContent(writeOptions)
	if err != nil {
		return nil, err
	}

	return gitlabRepositoryFile, nil
}

func (p *GitlabProject) CreateReleaseFromLatestCommitInDefaultBranch(createReleaseOptions *GitlabCreateReleaseOptions) (createdRelease *GitlabRelease, err error) {
	if createReleaseOptions == nil {
		return nil, tracederrors.TracedErrorNil("createReleaseOptions")
	}

	latestCommit, err := p.GetLatestCommitOfDefaultBranch(createReleaseOptions.Verbose)
	if err != nil {
		return nil, err
	}

	createdRelease, err = latestCommit.CreateRelease(createReleaseOptions)
	if err != nil {
		return nil, err
	}

	return createdRelease, nil
}

func (p *GitlabProject) DeployKeyByNameExists(keyName string) (exists bool, err error) {
	if len(keyName) <= 0 {
		return false, tracederrors.TracedError("keyName is empty string")
	}

	deployKeys, err := p.GetDeployKeys()
	if err != nil {
		return false, err
	}

	exists, err = deployKeys.DeployKeyByNameExists(keyName)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (p *GitlabProject) GetCachedPath() (path string, err error) {
	if !p.IsCachedPathSet() {
		_, err := p.GetPath()
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

func (p *GitlabProject) GetId() (id int, err error) {
	if p.id < 0 {
		return -1, tracederrors.TracedErrorf("id is set to invalid value of '%d'", p.id)
	}

	if p.id == 0 {
		rawResponse, err := p.GetRawResponse()
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

func (p *GitlabProject) MakePrivate(verbose bool) (err error) {
	nativeProjectsService, err := p.GetNativeProjectsService()
	if err != nil {
		return err
	}

	projectId, err := p.GetId()
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

	if verbose {
		logging.LogInfof("Gitlab project '%v' made private.", projectId)
	}

	return nil
}

func (p *GitlabProject) MakePublic(verbose bool) (err error) {
	nativeProjectsService, err := p.GetNativeProjectsService()
	if err != nil {
		return err
	}

	projectId, err := p.GetId()
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

	if verbose {
		logging.LogInfof("Gitlab project '%v' made public.", projectId)
	}

	return nil
}

func (p *GitlabProject) RecreateDeployKey(keyOptions *GitlabCreateDeployKeyOptions) (err error) {
	if keyOptions == nil {
		return tracederrors.TracedError("keyOptions is nil")
	}

	deployKey, err := p.GetDeployKeyByName(keyOptions.Name)
	if err != nil {
		return err
	}

	err = deployKey.RecreateDeployKey(keyOptions)
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
