package asciichgolangpublic

import (
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type GitlabDeleteProjectOptions struct {
	ProjectPath string
	Verbose     bool
}

func NewGitlabDeleteProjectOptions() (g *GitlabDeleteProjectOptions) {
	return new(GitlabDeleteProjectOptions)
}

func (g *GitlabDeleteProjectOptions) GetProjectPath() (projectPath string, err error) {
	if g.ProjectPath == "" {
		return "", tracederrors.TracedErrorf("ProjectPath not set")
	}

	return g.ProjectPath, nil
}

func (g *GitlabDeleteProjectOptions) GetVerbose() (verbose bool) {

	return g.Verbose
}

func (g *GitlabDeleteProjectOptions) MustGetProjectPath() (projectPath string) {
	projectPath, err := g.GetProjectPath()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return projectPath
}

func (g *GitlabDeleteProjectOptions) MustSetProjectPath(projectPath string) {
	err := g.SetProjectPath(projectPath)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabDeleteProjectOptions) SetProjectPath(projectPath string) (err error) {
	if projectPath == "" {
		return tracederrors.TracedErrorf("projectPath is empty string")
	}

	g.ProjectPath = projectPath

	return nil
}

func (g *GitlabDeleteProjectOptions) SetVerbose(verbose bool) {
	g.Verbose = verbose
}
