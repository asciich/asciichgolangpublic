package asciichgolangpublic

import (
	"errors"

	"github.com/asciich/asciichgolangpublic/logging"
	aerrors "github.com/asciich/asciichgolangpublic/errors"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

var ErrGitlabGroupNotFoundError = errors.New("Gitlab group not found")

type GitlabGroups struct {
	gitlab *GitlabInstance
}

func NewGitlabGroups() (gitlabGroups *GitlabGroups) {
	return new(GitlabGroups)
}

func (g *GitlabGroups) CreateGroup(groupPath string, createOptions *GitlabCreateGroupOptions) (createdGroup *GitlabGroup, err error) {
	if groupPath == "" {
		return nil, aerrors.TracedErrorEmptyString("groupPath")
	}

	if createOptions == nil {
		return nil, aerrors.TracedError("createOptions is nil")
	}

	group, err := g.GetGroupByPath(groupPath, createOptions.Verbose)
	if err != nil {
		return nil, err
	}

	err = group.Create(createOptions)
	if err != nil {
		return nil, err
	}

	return group, nil
}

func (g *GitlabGroups) GetGroupById(id int, verbose bool) (gitlabGroup *GitlabGroup, err error) {
	if id <= 0 {
		return nil, aerrors.TracedErrorf("Invalid group id '%d'", id)
	}

	gitlabGroup = NewGitlabGroup()

	err = gitlabGroup.SetId(id)
	if err != nil {
		return nil, err
	}

	gitlab, err := g.GetGitlab()
	if err != nil {
		return nil, err
	}

	err = gitlabGroup.SetGitlab(gitlab)
	if err != nil {
		return nil, err
	}

	return gitlabGroup, nil
}

func (g *GitlabGroups) GetGroupByPath(groupPath string, verbose bool) (gitlabGroup *GitlabGroup, err error) {
	gitlabGroup = NewGitlabGroup()

	gitlab, err := g.GetGitlab()
	if err != nil {
		return nil, err
	}

	err = gitlabGroup.SetGitlab(gitlab)
	if err != nil {
		return nil, err
	}

	err = gitlabGroup.SetGroupPath(groupPath)
	if err != nil {
		return nil, err
	}

	return gitlabGroup, nil
}

func (g *GitlabGroups) GroupByGroupPathExists(groupPath string, verbose bool) (groupExists bool, err error) {
	if len(groupPath) <= 0 {
		return false, aerrors.TracedError("groupPath is empty string")
	}

	group, err := g.GetGroupByPath(groupPath, verbose)
	if err != nil {
		return false, err
	}

	groupExists, err = group.Exists(verbose)
	if err != nil {
		return false, err
	}

	return groupExists, nil
}

func (g *GitlabGroups) MustCreateGroup(groupPath string, createOptions *GitlabCreateGroupOptions) (createdGroup *GitlabGroup) {
	createdGroup, err := g.CreateGroup(groupPath, createOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return createdGroup
}

func (g *GitlabGroups) MustGetFqdn() (fqdn string) {
	fqdn, err := g.GetFqdn()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return fqdn
}

func (g *GitlabGroups) MustGetGitlab() (gitlab *GitlabInstance) {
	gitlab, err := g.GetGitlab()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlab
}

func (g *GitlabGroups) MustGetGroupById(id int, verbose bool) (gitlabGroup *GitlabGroup) {
	gitlabGroup, err := g.GetGroupById(id, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlabGroup
}

func (g *GitlabGroups) MustGetGroupByPath(groupPath string, verbose bool) (gitlabGroup *GitlabGroup) {
	gitlabGroup, err := g.GetGroupByPath(groupPath, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlabGroup
}

func (g *GitlabGroups) MustGetNativeClient() (nativeClient *gitlab.Client) {
	nativeClient, err := g.GetNativeClient()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return nativeClient
}

func (g *GitlabGroups) MustGetNativeGroupsService() (nativeGroupsService *gitlab.GroupsService) {
	nativeGroupsService, err := g.GetNativeGroupsService()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return nativeGroupsService
}

func (g *GitlabGroups) MustGroupByGroupPathExists(groupPath string, verbose bool) (groupExists bool) {
	groupExists, err := g.GroupByGroupPathExists(groupPath, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return groupExists
}

func (g *GitlabGroups) MustSetGitlab(gitlab *GitlabInstance) {
	err := g.SetGitlab(gitlab)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (p *GitlabGroups) GetFqdn() (fqdn string, err error) {
	gitlab, err := p.GetGitlab()
	if err != nil {
		return "", err
	}

	fqdn, err = gitlab.GetFqdn()
	if err != nil {
		return "", err
	}

	return fqdn, nil
}

func (p *GitlabGroups) GetGitlab() (gitlab *GitlabInstance, err error) {
	if p.gitlab == nil {
		return nil, aerrors.TracedError("gitlab is not set")
	}

	return p.gitlab, nil
}

func (p *GitlabGroups) GetNativeClient() (nativeClient *gitlab.Client, err error) {
	gitlab, err := p.GetGitlab()
	if err != nil {
		return nil, err
	}

	nativeClient, err = gitlab.GetNativeClient()
	if err != nil {
		return nil, err
	}

	return nativeClient, nil
}

func (p *GitlabGroups) GetNativeGroupsService() (nativeGroupsService *gitlab.GroupsService, err error) {
	nativeClient, err := p.GetNativeClient()
	if err != nil {
		return nil, err
	}

	nativeGroupsService = nativeClient.Groups
	if nativeGroupsService == nil {
		return nil, aerrors.TracedError("unable to get nativeGroupsService")
	}

	return nativeGroupsService, nil
}

func (p *GitlabGroups) SetGitlab(gitlab *GitlabInstance) (err error) {
	if gitlab == nil {
		return aerrors.TracedError("gitlab is nil")
	}

	p.gitlab = gitlab

	return nil
}
