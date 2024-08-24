package asciichgolangpublic

type GitlabDeleteProjectOptions struct {
	ProjectPath string
	Verbose     bool
}

func NewGitlabDeleteProjectOptions() (g *GitlabDeleteProjectOptions) {
	return new(GitlabDeleteProjectOptions)
}

func (g *GitlabDeleteProjectOptions) GetProjectPath() (projectPath string, err error) {
	if g.ProjectPath == "" {
		return "", TracedErrorf("ProjectPath not set")
	}

	return g.ProjectPath, nil
}

func (g *GitlabDeleteProjectOptions) GetVerbose() (verbose bool) {

	return g.Verbose
}

func (g *GitlabDeleteProjectOptions) MustGetProjectPath() (projectPath string) {
	projectPath, err := g.GetProjectPath()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return projectPath
}

func (g *GitlabDeleteProjectOptions) MustSetProjectPath(projectPath string) {
	err := g.SetProjectPath(projectPath)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabDeleteProjectOptions) SetProjectPath(projectPath string) (err error) {
	if projectPath == "" {
		return TracedErrorf("projectPath is empty string")
	}

	g.ProjectPath = projectPath

	return nil
}

func (g *GitlabDeleteProjectOptions) SetVerbose(verbose bool) {
	g.Verbose = verbose
}
