package asciichgolangpublic

import (
	"path/filepath"
	"strings"

	"github.com/asciich/asciichgolangpublic/datatypes/slicesutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type GitlabCreateProjectOptions struct {
	ProjectPath string
	IsPublic    bool
	Verbose     bool
}

func NewGitlabCreateProjectOptions() (g *GitlabCreateProjectOptions) {
	return new(GitlabCreateProjectOptions)
}

func (g *GitlabCreateProjectOptions) GetIsPublic() (isPublic bool, err error) {

	return g.IsPublic, nil
}

func (g *GitlabCreateProjectOptions) GetVerbose() (verbose bool, err error) {

	return g.Verbose, nil
}

func (g *GitlabCreateProjectOptions) MustGetGroupNames(verbose bool) (groupNames []string) {
	groupNames, err := g.GetGroupNames(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return groupNames
}

func (g *GitlabCreateProjectOptions) MustGetGroupPath(verbose bool) (groupPath string) {
	groupPath, err := g.GetGroupPath(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return groupPath
}

func (g *GitlabCreateProjectOptions) MustGetIsPublic() (isPublic bool) {
	isPublic, err := g.GetIsPublic()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return isPublic
}

func (g *GitlabCreateProjectOptions) MustGetProjectName() (projectName string) {
	projectName, err := g.GetProjectName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return projectName
}

func (g *GitlabCreateProjectOptions) MustGetProjectPath() (projectPath string) {
	projectPath, err := g.GetProjectPath()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return projectPath
}

func (g *GitlabCreateProjectOptions) MustGetVerbose() (verbose bool) {
	verbose, err := g.GetVerbose()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return verbose
}

func (g *GitlabCreateProjectOptions) MustSetIsPublic(isPublic bool) {
	err := g.SetIsPublic(isPublic)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabCreateProjectOptions) MustSetProjectPath(projectPath string) {
	err := g.SetProjectPath(projectPath)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabCreateProjectOptions) MustSetVerbose(verbose bool) {
	err := g.SetVerbose(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabCreateProjectOptions) SetIsPublic(isPublic bool) (err error) {
	g.IsPublic = isPublic

	return nil
}

func (g *GitlabCreateProjectOptions) SetProjectPath(projectPath string) (err error) {
	if projectPath == "" {
		return tracederrors.TracedErrorf("projectPath is empty string")
	}

	g.ProjectPath = projectPath

	return nil
}

func (g *GitlabCreateProjectOptions) SetVerbose(verbose bool) (err error) {
	g.Verbose = verbose

	return nil
}

func (o *GitlabCreateProjectOptions) GetGroupNames(verbose bool) (groupNames []string, err error) {
	projectPath, err := o.GetProjectPath()
	if err != nil {
		return nil, err
	}

	pathOnly := filepath.Dir(projectPath)
	pathOnly = strings.TrimSpace(pathOnly)
	if len(pathOnly) <= 0 {
		return []string{}, nil
	}

	groupNames = strings.Split(pathOnly, "/")
	groupNames = slicesutils.RemoveEmptyStrings(groupNames)
	groupNames = slicesutils.RemoveMatchingStrings(groupNames, ".")

	if verbose {
		logging.LogInfof(
			"Gitlab create project options: Evaluated group names '%v' from project path '%s'",
			groupNames,
			projectPath,
		)
	}

	return groupNames, nil
}

func (o *GitlabCreateProjectOptions) GetGroupPath(verbose bool) (groupPath string, err error) {
	groupNames, err := o.GetGroupNames(verbose)
	if err != nil {
		return "", err
	}

	groupPath = strings.Join(groupNames, "/")
	return groupPath, nil
}

func (o *GitlabCreateProjectOptions) GetProjectName() (projectName string, err error) {
	projectPath, err := o.GetProjectPath()
	if err != nil {
		return "", err
	}

	projectName = filepath.Base(projectPath)

	return projectName, nil
}

func (o *GitlabCreateProjectOptions) GetProjectPath() (projectPath string, err error) {
	if len(o.ProjectPath) <= 0 {
		return "", tracederrors.TracedError("ProjectPath is not set")
	}

	return o.ProjectPath, nil
}
