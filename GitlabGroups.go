package asciichgolangpublic

import (
	"errors"
	"strings"

	"github.com/xanzy/go-gitlab"
)

var ErrGitlabGroupNotFoundError = errors.New("Gitlab group not found")

type GitlabGroups struct {
	gitlab *GitlabInstance
}

func NewGitlabGroups() (gitlabGroups *GitlabGroups) {
	return new(GitlabGroups)
}

func (g *GitlabGroups) MustCreateGroup(createOptions *GitlabCreateGroupOptions) (createdGroup *GitlabGroup) {
	createdGroup, err := g.CreateGroup(createOptions)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return createdGroup
}

func (g *GitlabGroups) MustGetFqdn() (fqdn string) {
	fqdn, err := g.GetFqdn()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return fqdn
}

func (g *GitlabGroups) MustGetGitlab() (gitlab *GitlabInstance) {
	gitlab, err := g.GetGitlab()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlab
}

func (g *GitlabGroups) MustGetGroupByGroupPath(groupPath string) (gitlabGroup *GitlabGroup) {
	gitlabGroup, err := g.GetGroupByGroupPath(groupPath)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlabGroup
}

func (g *GitlabGroups) MustGetNativeClient() (nativeClient *gitlab.Client) {
	nativeClient, err := g.GetNativeClient()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return nativeClient
}

func (g *GitlabGroups) MustGetNativeGroupsService() (nativeGroupsService *gitlab.GroupsService) {
	nativeGroupsService, err := g.GetNativeGroupsService()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return nativeGroupsService
}

func (g *GitlabGroups) MustGroupByGroupPathExists(groupPath string) (groupExists bool) {
	groupExists, err := g.GroupByGroupPathExists(groupPath)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return groupExists
}

func (g *GitlabGroups) MustSetGitlab(gitlab *GitlabInstance) {
	err := g.SetGitlab(gitlab)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (p *GitlabGroups) CreateGroup(createOptions *GitlabCreateGroupOptions) (createdGroup *GitlabGroup, err error) {
	if createOptions == nil {
		return nil, TracedError("createOptions is nil")
	}

	groupPath, err := createOptions.GetGroupPath()
	if err != nil {
		return nil, err
	}

	fqdn, err := p.GetFqdn()
	if err != nil {
		return nil, err
	}

	if createOptions.Verbose {
		LogInfof("Create group '%s' in gitlab '%s' started.", groupPath, fqdn)
	}

	groupExists, err := p.GroupByGroupPathExists(groupPath)
	if err != nil {
		return nil, err
	}

	if groupExists {
		if createOptions.Verbose {
			LogInfof("Group '%s' already exists on gitlab '%s'.", groupPath, fqdn)
		}
	} else {
		if createOptions.Verbose {
			LogInfof("Going to create group '%s' on gitlab '%s'.", groupPath, fqdn)
		}
		isSubgroup, err := createOptions.IsSubgroup()
		if err != nil {
			return nil, err
		}

		var parentGroup *GitlabGroup = nil

		if isSubgroup {
			parentGroupPath, err := createOptions.GetParentGroupPath()
			if err != nil {
				return nil, err
			}

			optionsToUse := createOptions.GetDeepCopy()
			optionsToUse.GroupPath = parentGroupPath

			parentGroup, err = p.CreateGroup(optionsToUse)
			if err != nil {
				return nil, err
			}
		}

		nativeGroupsService, err := p.GetNativeGroupsService()
		if err != nil {
			return nil, err
		}

		groupName, err := createOptions.GetGroupName()
		if err != nil {
			return nil, err
		}

		if len(groupName) < 2 {
			return nil, TracedErrorf("Group names with len < 2 not allowed. But got '%s'", groupName)
		}

		var namespaceId *int = nil
		if parentGroup != nil {
			idToUse, err := parentGroup.GetId()
			if err != nil {
				return nil, err
			}
			namespaceId = &idToUse
		}

		groupPath, err := createOptions.GetGroupPath()
		if err != nil {
			return nil,
				err
		}

		defautVisibility := gitlab.PublicVisibility
		_, _, err = nativeGroupsService.CreateGroup(&gitlab.CreateGroupOptions{
			Name:       &groupName,
			ParentID:   namespaceId,
			Path:       &groupName,
			Visibility: &defautVisibility,
		})
		if err != nil {
			if createOptions.Verbose {
				LogErrorf("Creating group '%s' failed", groupPath)
			}
			return nil, err
		}

		if createOptions.Verbose {
			LogChangedf("Group '%s' on gitlab '%s' created.", groupPath, fqdn)
		}
	}

	createdGroup, err = p.GetGroupByGroupPath(groupPath)
	if err != nil {
		return nil, err
	}

	if createOptions.Verbose {
		LogInfof("Create group '%s' in gitlab '%s' finished.", groupPath, fqdn)
	}

	return createdGroup, nil
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
		return nil, TracedError("gitlab is not set")
	}

	return p.gitlab, nil
}

func (p *GitlabGroups) GetGroupByGroupPath(groupPath string) (gitlabGroup *GitlabGroup, err error) {
	if len(groupPath) <= 0 {
		return nil, TracedError("groupPath is empty string")
	}

	nativeGroupsService, err := p.GetNativeGroupsService()
	if err != nil {
		return nil, err
	}

	nativeGroup, _, err := nativeGroupsService.GetGroup(groupPath, &gitlab.GetGroupOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "404 {message: 404 Group Not Found}") {
			return nil, TracedErrorf("%w: %s", ErrGitlabGroupNotFoundError, groupPath)
		}
		return nil, err
	}

	gitlabGroup = NewGitlabGroup()

	gitlab, err := p.GetGitlab()
	if err != nil {
		return nil, err
	}

	err = gitlabGroup.SetGitlab(gitlab)
	if err != nil {
		return nil, err
	}

	err = gitlabGroup.SetId(nativeGroup.ID)
	if err != nil {
		return nil, err
	}

	return gitlabGroup, nil
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
		return nil, TracedError("unable to get nativeGroupsService")
	}

	return nativeGroupsService, nil
}

func (p *GitlabGroups) GroupByGroupPathExists(groupPath string) (groupExists bool, err error) {
	if len(groupPath) <= 0 {
		return false, TracedError("groupPath is empty string")
	}

	_, err = p.GetGroupByGroupPath(groupPath)
	if err != nil {
		if errors.Is(err, ErrGitlabGroupNotFoundError) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (p *GitlabGroups) SetGitlab(gitlab *GitlabInstance) (err error) {
	if gitlab == nil {
		return TracedError("gitlab is nil")
	}

	p.gitlab = gitlab

	return nil
}
