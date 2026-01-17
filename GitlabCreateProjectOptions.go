package asciichgolangpublic

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/datatypes/slicesutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type GitlabCreateProjectOptions struct {
	ProjectPath string
	IsPublic    bool
}

func NewGitlabCreateProjectOptions() (g *GitlabCreateProjectOptions) {
	return new(GitlabCreateProjectOptions)
}

func (g *GitlabCreateProjectOptions) GetIsPublic() (isPublic bool, err error) {

	return g.IsPublic, nil
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

func (o *GitlabCreateProjectOptions) GetGroupNames(ctx context.Context) (groupNames []string, err error) {
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
	groupNames = slicesutils.RemoveMatchingStrings(groupNames, "\\.")

	logging.LogInfoByCtxf(ctx, "Gitlab create project options: Evaluated group names '%v' from project path '%s'", groupNames, projectPath)

	return groupNames, nil
}

func (o *GitlabCreateProjectOptions) GetGroupPath(ctx context.Context) (groupPath string, err error) {
	groupNames, err := o.GetGroupNames(ctx)
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
