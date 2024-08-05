package asciichgolangpublic

import "github.com/xanzy/go-gitlab"

type GitlabRepositoryFiles struct {
	gitlabProject *GitlabProject
}

func NewGitlabRepositoryFiles() (g *GitlabRepositoryFiles) {
	return new(GitlabRepositoryFiles)
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

func (g *GitlabRepositoryFiles) MustGetNativeRepositoryFilesClient() (nativeRepositoryFilesClient *gitlab.RepositoryFilesService) {
	nativeRepositoryFilesClient, err := g.GetNativeRepositoryFilesClient()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return nativeRepositoryFilesClient
}

func (g *GitlabRepositoryFiles) MustGetRepositoryFile(options *GitlabGetRepositoryFileOptions) (repositoryFile *GitlabRepositoryFile) {
	repositoryFile, err := g.GetRepositoryFile(options)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return repositoryFile
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
