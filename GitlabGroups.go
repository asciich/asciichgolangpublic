package asciichgolangpublic

import (
	"context"
	"errors"

	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

var ErrGitlabGroupNotFoundError = errors.New("Gitlab group not found")

type GitlabGroups struct {
	gitlab *GitlabInstance
}

func NewGitlabGroups() (gitlabGroups *GitlabGroups) {
	return new(GitlabGroups)
}

func (g *GitlabGroups) CreateGroup(ctx context.Context, groupPath string) (createdGroup *GitlabGroup, err error) {
	if groupPath == "" {
		return nil, tracederrors.TracedErrorEmptyString("groupPath")
	}

	group, err := g.GetGroupByPath(ctx, groupPath)
	if err != nil {
		return nil, err
	}

	err = group.Create(ctx)
	if err != nil {
		return nil, err
	}

	return group, nil
}

func (g *GitlabGroups) GetGroupById(ctx context.Context, id int) (gitlabGroup *GitlabGroup, err error) {
	if id <= 0 {
		return nil, tracederrors.TracedErrorf("Invalid group id '%d'", id)
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

func (g *GitlabGroups) GetGroupByPath(ctx context.Context, groupPath string) (gitlabGroup *GitlabGroup, err error) {
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

func (g *GitlabGroups) GroupByGroupPathExists(ctx context.Context, groupPath string) (groupExists bool, err error) {
	if len(groupPath) <= 0 {
		return false, tracederrors.TracedError("groupPath is empty string")
	}

	group, err := g.GetGroupByPath(ctx, groupPath)
	if err != nil {
		return false, err
	}

	groupExists, err = group.Exists(ctx)
	if err != nil {
		return false, err
	}

	return groupExists, nil
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
		return nil, tracederrors.TracedError("gitlab is not set")
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
		return nil, tracederrors.TracedError("unable to get nativeGroupsService")
	}

	return nativeGroupsService, nil
}

func (p *GitlabGroups) SetGitlab(gitlab *GitlabInstance) (err error) {
	if gitlab == nil {
		return tracederrors.TracedError("gitlab is nil")
	}

	p.gitlab = gitlab

	return nil
}
