package asciichgolangpublic

import (
	"path/filepath"
	"strings"
)

type GitlabCreateGroupOptions struct {
	GroupPath string
	Verbose   bool
}

func NewGitlabCreateGroupOptions() (createOptions *GitlabCreateGroupOptions) {
	return new(GitlabCreateGroupOptions)
}

func (g *GitlabCreateGroupOptions) GetVerbose() (verbose bool, err error) {

	return g.Verbose, nil
}

func (g *GitlabCreateGroupOptions) MustGetGroupName() (groupName string) {
	groupName, err := g.GetGroupName()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return groupName
}

func (g *GitlabCreateGroupOptions) MustGetGroupPath() (groupPath string) {
	groupPath, err := g.GetGroupPath()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return groupPath
}

func (g *GitlabCreateGroupOptions) MustGetParentGroupPath() (parentGroupPath string) {
	parentGroupPath, err := g.GetParentGroupPath()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return parentGroupPath
}

func (g *GitlabCreateGroupOptions) MustGetVerbose() (verbose bool) {
	verbose, err := g.GetVerbose()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return verbose
}

func (g *GitlabCreateGroupOptions) MustIsSubgroup() (isSubgroup bool) {
	isSubgroup, err := g.IsSubgroup()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return isSubgroup
}

func (g *GitlabCreateGroupOptions) MustSetGroupPath(groupPath string) {
	err := g.SetGroupPath(groupPath)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabCreateGroupOptions) MustSetVerbose(verbose bool) {
	err := g.SetVerbose(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabCreateGroupOptions) SetGroupPath(groupPath string) (err error) {
	if groupPath == "" {
		return TracedErrorf("groupPath is empty string")
	}

	g.GroupPath = groupPath

	return nil
}

func (g *GitlabCreateGroupOptions) SetVerbose(verbose bool) (err error) {
	g.Verbose = verbose

	return nil
}

func (o *GitlabCreateGroupOptions) GetDeepCopy() (copy *GitlabCreateGroupOptions) {
	copy = NewGitlabCreateGroupOptions()

	*copy = *o

	return copy
}

func (o *GitlabCreateGroupOptions) GetGroupName() (groupName string, err error) {
	groupPath, err := o.GetGroupPath()
	if err != nil {
		return "", err
	}

	groupName = filepath.Base(groupPath)
	return groupName, nil
}

func (o *GitlabCreateGroupOptions) GetGroupPath() (groupPath string, err error) {
	if len(o.GroupPath) <= 0 {
		return "", TracedError("GroupPath not set")
	}

	return o.GroupPath, nil
}

func (o *GitlabCreateGroupOptions) GetParentGroupPath() (parentGroupPath string, err error) {
	groupPath, err := o.GetGroupPath()
	if err != nil {
		return "", err
	}

	parentGroupPath = filepath.Dir(groupPath)
	parentGroupPath = strings.TrimPrefix(parentGroupPath, "/")
	return parentGroupPath, nil
}

func (o *GitlabCreateGroupOptions) IsSubgroup() (isSubgroup bool, err error) {
	groupPath, err := o.GetGroupPath()
	if err != nil {
		return false, err
	}

	groupPath = strings.TrimPrefix(groupPath, "/")
	isSubgroup = strings.Contains(groupPath, "/")

	return isSubgroup, nil
}
