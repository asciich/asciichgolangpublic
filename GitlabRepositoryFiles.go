package asciichgolangpublic

import (
	"sort"

	aslices "github.com/asciich/asciichgolangpublic/datatypes/slices"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type GitlabRepositoryFiles struct {
	gitlabProject *GitlabProject
}

func NewGitlabRepositoryFiles() (g *GitlabRepositoryFiles) {
	return new(GitlabRepositoryFiles)
}

func (g *GitlabRepositoryFiles) CreateEmptyFile(fileName string, ref string, verbose bool) (createdFile *GitlabRepositoryFile, err error) {
	if fileName == "" {
		return nil, TracedErrorEmptyString("fileName")
	}

	if ref == "" {
		return nil, TracedErrorEmptyString("ref")
	}

	createdFile, err = g.WriteFileContent(
		&GitlabWriteFileOptions{
			Path:          fileName,
			Content:       []byte{},
			CommitMessage: "Create empty file",
			BranchName:    ref,
			Verbose:       verbose,
		},
	)
	if err != nil {
		return nil, err
	}

	return createdFile, nil
}

func (g *GitlabRepositoryFiles) DeleteAllRepositoryFiles(branchName string, verbose bool) (err error) {
	if branchName == "" {
		return TracedErrorEmptyString("branchName")
	}

	projectUrl, err := g.GetProjectUrl()
	if err != nil {
		return err
	}

	if verbose {
		LogInfof("Delete all repository files in '%s' started.", projectUrl)
	}

	files, err := g.GetFiles(branchName, verbose)
	if err != nil {
		return err
	}

	for _, toDelete := range files {
		err = toDelete.Delete("Delete all repository files", verbose)
		if err != nil {
			return err
		}
	}

	if len(files) > 0 {
		if verbose {
			LogChangedf("Deleted '%d' files in branch '%s' of gitlab project %s .", len(files), branchName, projectUrl)
		}
	} else {
		LogInfof(
			"Gitlab project '%s' is already empty, no files in branch '%s' to delete.",
			projectUrl,
			branchName,
		)
	}

	if verbose {
		LogInfof("Delete all repository files in %s finished.", projectUrl)
	}

	return nil
}

func (g *GitlabRepositoryFiles) GetDirectoryNames(ref string, verbose bool) (directoryNames []string, err error) {
	if ref == "" {
		return nil, TracedErrorEmptyString("ref")
	}

	fileAndDirectoryNames, err := g.GetFileAndDirectoryNames(ref, verbose)
	if err != nil {
		return nil, err
	}

	directoryNames = []string{}
	for _, toCheck := range fileAndDirectoryNames {
		toCheckWithAppendix := toCheck + "/"
		if aslices.AtLeastOneElementStartsWith(fileAndDirectoryNames, toCheckWithAppendix) {
			directoryNames = append(directoryNames, toCheck)
		}
	}

	directoryNames = aslices.RemoveDuplicatedStrings(directoryNames)

	sort.Strings(directoryNames)

	return directoryNames, nil
}

func (g *GitlabRepositoryFiles) GetFileAndDirectoryNames(ref string, verbose bool) (fileNames []string, err error) {
	if ref == "" {
		return nil, TracedErrorEmptyString("ref")
	}

	nativeClient, projectId, err := g.GetNativeRepositoriesClientAndProjectid()
	if err != nil {
		return nil, err
	}

	trueBoolean := true
	nextPage := 1

	for {
		if nextPage <= 0 {
			break
		}

		nativeList, response, err := nativeClient.ListTree(
			projectId,
			&gitlab.ListTreeOptions{
				Ref:       &ref,
				Recursive: &trueBoolean,
				ListOptions: gitlab.ListOptions{
					PerPage: 20,
					Page:    nextPage,
				},
			},
		)
		if err != nil {
			return nil, TracedErrorf("ListTree failed: '%w'", err)
		}

		for _, entry := range nativeList {
			toAdd := entry.Path
			fileNames = append(fileNames, toAdd)
		}

		nextPage = response.NextPage
	}

	return fileNames, nil
}

func (g *GitlabRepositoryFiles) GetFileNames(ref string, verbose bool) (fileNames []string, err error) {
	if ref == "" {
		return nil, TracedErrorEmptyString("ref")
	}

	fileAndDirectoryNames, err := g.GetFileAndDirectoryNames(ref, verbose)
	if err != nil {
		return nil, err
	}

	fileNames = []string{}
	for _, toCheck := range fileAndDirectoryNames {
		toCheckWithAppendix := toCheck + "/"
		if aslices.AtLeastOneElementStartsWith(fileAndDirectoryNames, toCheckWithAppendix) {
			continue
		}

		fileNames = append(fileNames, toCheck)
	}

	fileNames = aslices.RemoveDuplicatedStrings(fileNames)
	
	sort.Strings(fileNames)

	return fileNames, nil
}

func (g *GitlabRepositoryFiles) GetFiles(ref string, verbose bool) (files []*GitlabRepositoryFile, err error) {
	if ref == "" {
		return nil, TracedErrorEmptyString("ref")
	}

	gitlabProject, err := g.GetGitlabProject()
	if err != nil {
		return nil, err
	}

	fileNames, err := g.GetFileNames(ref, verbose)
	if err != nil {
		return nil, err
	}

	files = []*GitlabRepositoryFile{}
	for _, name := range fileNames {
		toAdd := NewGitlabRepositoryFile()

		err = toAdd.SetGitlabProject(gitlabProject)
		if err != nil {
			return nil, err
		}

		err = toAdd.SetBranchName(ref)
		if err != nil {
			return nil, err
		}

		err = toAdd.SetPath(name)
		if err != nil {
			return nil, err
		}

		files = append(files, toAdd)
	}

	return files, nil
}

func (g *GitlabRepositoryFiles) GetGitlab() (gitlab *GitlabInstance, err error) {
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

func (g *GitlabRepositoryFiles) GetGitlabProject() (gitlabProject *GitlabProject, err error) {
	if g.gitlabProject == nil {
		return nil, TracedErrorf("gitlabProject not set")
	}

	return g.gitlabProject, nil
}

func (g *GitlabRepositoryFiles) GetNativeRepositoriesClient() (nativeRepositoriesClient *gitlab.RepositoriesService, err error) {
	gitlab, err := g.GetGitlab()
	if err != nil {
		return nil, err
	}

	nativeRepositoriesClient, err = gitlab.GetNativeRepositoriesClient()
	if err != nil {
		return nil, err
	}

	return nativeRepositoriesClient, nil
}

func (g *GitlabRepositoryFiles) GetNativeRepositoriesClientAndProjectid() (nativeRepositoriesClient *gitlab.RepositoriesService, projectid int, err error) {
	nativeRepositoriesClient, err = g.GetNativeRepositoriesClient()
	if err != nil {
		return nil, -1, err
	}

	projectId, err := g.GetProjectId()
	if err != nil {
		return nil, -1, err
	}

	return nativeRepositoriesClient, projectId, nil
}

func (g *GitlabRepositoryFiles) GetNativeRepositoryFilesClient() (nativeRepositoryFilesClient *gitlab.RepositoryFilesService, err error) {
	gitlab, err := g.GetGitlab()
	if err != nil {
		return nil, err
	}

	nativeRepositoryFilesClient, err = gitlab.GetNativeRepositoryFilesClient()
	if err != nil {
		return nil, err
	}

	return nativeRepositoryFilesClient, nil
}

func (g *GitlabRepositoryFiles) GetNativeRepositoryFilesClientAndProjectId() (nativeRepositoryFilesClient *gitlab.RepositoryFilesService, projectId int, err error) {
	nativeRepositoryFilesClient, err = g.GetNativeRepositoryFilesClient()
	if err != nil {
		return nil, -1, err
	}

	projectId, err = g.GetProjectId()
	if err != nil {
		return nil, -1, err
	}

	return nativeRepositoryFilesClient, projectId, nil
}

func (g *GitlabRepositoryFiles) GetProjectId() (projectId int, err error) {
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

func (g *GitlabRepositoryFiles) GetProjectUrl() (projectUrl string, err error) {
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

func (g *GitlabRepositoryFiles) GetRepositoryFile(options *GitlabGetRepositoryFileOptions) (repositoryFile *GitlabRepositoryFile, err error) {
	if options == nil {
		return nil, TracedErrorNil("options")
	}

	repositoryFile = NewGitlabRepositoryFile()

	gitlabProject, err := g.GetGitlabProject()
	if err != nil {
		return nil, err
	}

	err = repositoryFile.SetGitlabProject(gitlabProject)
	if err != nil {
		return nil, err
	}

	path, err := options.GetPath()
	if err != nil {
		return nil, err
	}

	err = repositoryFile.SetPath(path)
	if err != nil {
		return nil, err
	}

	if options.IsBranchNameSet() {
		branchName, err := options.GetBranchName()
		if err != nil {
			return nil, err
		}

		err = repositoryFile.SetBranchName(branchName)
		if err != nil {
			return nil, err
		}
	}

	return repositoryFile, nil
}

func (g *GitlabRepositoryFiles) HasNoRepositoryFiles(branchName string, verbose bool) (hasNoRepositoryFiles bool, err error) {
	if branchName == "" {
		return false, TracedErrorEmptyString("branchName")
	}

	hasRepositoryFile, err := g.HasRepositoryFiles(branchName, false)
	if err != nil {
		return false, err
	}

	hasNoRepositoryFiles = !hasRepositoryFile

	return hasNoRepositoryFiles, nil
}

func (g *GitlabRepositoryFiles) HasRepositoryFiles(branchName string, verbose bool) (hasRepositoryFiles bool, err error) {
	if branchName == "" {
		return false, TracedErrorEmptyString("branchName")
	}

	fileNameList, err := g.GetFileNames(branchName, verbose)
	if err != nil {
		return false, err
	}

	hasRepositoryFiles = len(fileNameList) > 0

	return hasRepositoryFiles, nil
}

func (g *GitlabRepositoryFiles) MustCreateEmptyFile(fileName string, ref string, verbose bool) (createdFile *GitlabRepositoryFile) {
	createdFile, err := g.CreateEmptyFile(fileName, ref, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return createdFile
}

func (g *GitlabRepositoryFiles) MustDeleteAllRepositoryFiles(branchName string, verbose bool) {
	err := g.DeleteAllRepositoryFiles(branchName, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabRepositoryFiles) MustGetDirectoryNames(ref string, verbose bool) (directoryNames []string) {
	directoryNames, err := g.GetDirectoryNames(ref, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return directoryNames
}

func (g *GitlabRepositoryFiles) MustGetFileAndDirectoryNames(ref string, verbose bool) (fileNames []string) {
	fileNames, err := g.GetFileAndDirectoryNames(ref, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return fileNames
}

func (g *GitlabRepositoryFiles) MustGetFileNames(ref string, verbose bool) (fileNames []string) {
	fileNames, err := g.GetFileNames(ref, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return fileNames
}

func (g *GitlabRepositoryFiles) MustGetFiles(ref string, verbose bool) (files []*GitlabRepositoryFile) {
	files, err := g.GetFiles(ref, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return files
}

func (g *GitlabRepositoryFiles) MustGetGitlab() (gitlab *GitlabInstance) {
	gitlab, err := g.GetGitlab()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlab
}

func (g *GitlabRepositoryFiles) MustGetGitlabProject() (gitlabProject *GitlabProject) {
	gitlabProject, err := g.GetGitlabProject()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlabProject
}

func (g *GitlabRepositoryFiles) MustGetNativeRepositoriesClient() (nativeRepositoriesClient *gitlab.RepositoriesService) {
	nativeRepositoriesClient, err := g.GetNativeRepositoriesClient()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return nativeRepositoriesClient
}

func (g *GitlabRepositoryFiles) MustGetNativeRepositoriesClientAndProjectid() (nativeRepositoriesClient *gitlab.RepositoriesService, projectid int) {
	nativeRepositoriesClient, projectid, err := g.GetNativeRepositoriesClientAndProjectid()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return nativeRepositoriesClient, projectid
}

func (g *GitlabRepositoryFiles) MustGetNativeRepositoryFilesClient() (nativeRepositoryFilesClient *gitlab.RepositoryFilesService) {
	nativeRepositoryFilesClient, err := g.GetNativeRepositoryFilesClient()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return nativeRepositoryFilesClient
}

func (g *GitlabRepositoryFiles) MustGetNativeRepositoryFilesClientAndProjectId() (nativeRepositoryFilesClient *gitlab.RepositoryFilesService, projectId int) {
	nativeRepositoryFilesClient, projectId, err := g.GetNativeRepositoryFilesClientAndProjectId()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return nativeRepositoryFilesClient, projectId
}

func (g *GitlabRepositoryFiles) MustGetProjectId() (projectId int) {
	projectId, err := g.GetProjectId()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return projectId
}

func (g *GitlabRepositoryFiles) MustGetProjectUrl() (projectUrl string) {
	projectUrl, err := g.GetProjectUrl()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return projectUrl
}

func (g *GitlabRepositoryFiles) MustGetRepositoryFile(options *GitlabGetRepositoryFileOptions) (repositoryFile *GitlabRepositoryFile) {
	repositoryFile, err := g.GetRepositoryFile(options)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return repositoryFile
}

func (g *GitlabRepositoryFiles) MustHasNoRepositoryFiles(branchName string, verbose bool) (hasNoRepositoryFiles bool) {
	hasNoRepositoryFiles, err := g.HasNoRepositoryFiles(branchName, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return hasNoRepositoryFiles
}

func (g *GitlabRepositoryFiles) MustHasRepositoryFiles(branchName string, verbose bool) (hasRepositoryFiles bool) {
	hasRepositoryFiles, err := g.HasRepositoryFiles(branchName, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return hasRepositoryFiles
}

func (g *GitlabRepositoryFiles) MustReadFileContentAsString(options *GitlabReadFileOptions) (content string) {
	content, err := g.ReadFileContentAsString(options)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return content
}

func (g *GitlabRepositoryFiles) MustSetGitlabProject(gitlabProject *GitlabProject) {
	err := g.SetGitlabProject(gitlabProject)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabRepositoryFiles) MustWriteFileContent(options *GitlabWriteFileOptions) (gitlabRepositoryFile *GitlabRepositoryFile) {
	gitlabRepositoryFile, err := g.WriteFileContent(options)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlabRepositoryFile
}

func (g *GitlabRepositoryFiles) ReadFileContentAsString(options *GitlabReadFileOptions) (content string, err error) {
	if options == nil {
		return "", TracedErrorNil("options")
	}

	getFileOptions, err := options.GetGitlabGetRepositoryFileOptions()
	if err != nil {
		return "", err
	}

	repositoryFile, err := g.GetRepositoryFile(getFileOptions)
	if err != nil {
		return "", err
	}

	content, err = repositoryFile.GetContentAsString(options.Verbose)
	if err != nil {
		return "", err
	}

	return content, nil
}

func (g *GitlabRepositoryFiles) SetGitlabProject(gitlabProject *GitlabProject) (err error) {
	if gitlabProject == nil {
		return TracedErrorf("gitlabProject is nil")
	}

	g.gitlabProject = gitlabProject

	return nil
}

func (g *GitlabRepositoryFiles) WriteFileContent(options *GitlabWriteFileOptions) (gitlabRepositoryFile *GitlabRepositoryFile, err error) {
	if options == nil {
		return nil, TracedErrorNil("options")
	}

	getFileOptions, err := options.GetGitlabGetRepositoryFileOptions()
	if err != nil {
		return nil, err
	}

	repositoryFile, err := g.GetRepositoryFile(getFileOptions)
	if err != nil {
		return nil, err
	}

	err = repositoryFile.WriteFileContentByBytes(options.Content, options.CommitMessage, options.Verbose)
	if err != nil {
		return nil, err
	}

	return repositoryFile, nil
}
