package asciichgolangpublic

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	gitlab "gitlab.com/gitlab-org/api/client-go"

	"github.com/asciich/asciichgolangpublic/pkg/datatypes/slicesutils"
	"github.com/asciich/asciichgolangpublic/pkg/datatypes/stringsutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type GitlabGroup struct {
	gitlab    *GitlabInstance
	id        int
	groupPath string
}

func NewGitlabGroup() (gitlabGroup *GitlabGroup) {
	return new(GitlabGroup)
}

func (g *GitlabGroup) Create(ctx context.Context) (err error) {
	if !g.IsGroupPathSet() {
		return tracederrors.TracedError("Group path must be set for a group to create")
	}
	groupPath, err := g.GetGroupPath(ctx)
	if err != nil {
		return err
	}

	fqdn, err := g.GetGitlabFqdn()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Going to create group '%s' on gitlab '%s'.", groupPath, fqdn)

	exists, err := g.Exists(ctx)
	if err != nil {
		return err
	}

	if exists {
		logging.LogInfoByCtxf(ctx, "Group '%s' already exists on gitlab '%s'.", groupPath, fqdn)
	} else {
		isSubgroup, err := g.IsSubgroup(ctx)
		if err != nil {
			return err
		}

		var parentGroup *GitlabGroup = nil

		if isSubgroup {
			parentGroup, err = g.GetParentGroup(ctx)
			if err != nil {
				return err
			}

			err = parentGroup.Create(ctx)
			if err != nil {
				return err
			}
		}

		nativeGroupsService, err := g.GetNativeGroupsService()
		if err != nil {
			return err
		}

		groupName, err := g.GetGroupName(ctx)
		if err != nil {
			return err
		}

		if len(groupName) < 2 {
			return tracederrors.TracedErrorf("Group names with len < 2 not allowed. But got '%s'", groupName)
		}

		var namespaceId *int = nil
		if parentGroup != nil {
			idToUse, err := parentGroup.GetId(ctx)
			if err != nil {
				return err
			}
			namespaceId = &idToUse
		}

		defautVisibility := gitlab.PublicVisibility
		_, _, err = nativeGroupsService.CreateGroup(&gitlab.CreateGroupOptions{
			Name:       &groupName,
			ParentID:   namespaceId,
			Path:       &groupName,
			Visibility: &defautVisibility,
		})
		if err != nil {
			logging.LogErrorByCtxf(ctx, "Creating group '%s' failed", groupPath)
			return err
		}

		logging.LogChangedByCtxf(ctx, "Group '%s' on gitlab '%s' created.", groupPath, fqdn)
	}

	logging.LogInfoByCtxf(ctx, "Create group '%s' in gitlab '%s' finished.", groupPath, fqdn)

	return nil
}

func (g *GitlabGroup) Delete(ctx context.Context) (err error) {
	gid, err := g.GetGroupIdOrPath(ctx)
	if err != nil {
		return err
	}

	exists, err := g.Exists(ctx)
	if err != nil {
		return err
	}

	if exists {
		nativeService, err := g.GetNativeGroupsService()
		if err != nil {
			return err
		}

		_, err = nativeService.DeleteGroup(gid, &gitlab.DeleteGroupOptions{}, nil)
		if err != nil {
			return tracederrors.TracedErrorf("Delete gitlab group failed: '%w'", err)
		}

		deleteConfirmed := false
		maxRetry := 10
		for i := 0; i < maxRetry; i++ {
			exists, err = g.Exists(ctx)
			if err != nil {
				return err
			}

			if exists {
				logging.LogInfoByCtxf(ctx, "Wait until delete of gitlab group '%v' is finished (%d/%d).", gid, i+1, maxRetry)
				time.Sleep(time.Millisecond * 500)
				continue
			}

			deleteConfirmed = true
			break
		}

		if !deleteConfirmed {
			return tracederrors.TracedErrorf("Failed to delete gitlab group '%v'. Group still exists after delete.", gid)
		}

		logging.LogChangedByCtxf(ctx, "Group '%s' deleted", gid)
	} else {
		logging.LogInfoByCtxf(ctx, "Group '%s' is already absent.", gid)
	}

	return nil
}

func (g *GitlabGroup) Exists(ctx context.Context) (exists bool, err error) {
	exists = true

	_, err = g.GetRawResponse(ctx)
	if err != nil {
		if errors.Is(err, ErrGitlabGroupNotFoundError) {
			exists = false
		} else {
			return false, err
		}
	}

	gid, err := g.GetGroupIdOrPathAsString(ctx)
	if err != nil {
		return false, err
	}

	if exists {
		logging.LogInfoByCtxf(ctx, "Gitlab group '%s' exists.", gid)
	} else {
		logging.LogInfoByCtxf(ctx, "Gitlab group '%s' does not exist.", gid)
	}

	return exists, nil
}

func (g *GitlabGroup) GetGitlabFqdn() (gitlabFqdn string, err error) {
	gitlab, err := g.GetGitlab()
	if err != nil {
		return "", err
	}

	gitlabFqdn, err = gitlab.GetFqdn()
	if err != nil {
		return "", err
	}

	return gitlabFqdn, nil
}

func (g *GitlabGroup) GetGitlabGroups() (gitlabGroups *GitlabGroups, err error) {
	gitlab, err := g.GetGitlab()
	if err != nil {
		return nil, err
	}

	gitlabGroups, err = gitlab.GetGitlabGroups()
	if err != nil {
		return nil, err
	}

	return gitlabGroups, nil
}

func (g *GitlabGroup) GetGroupIdOrPath(ctx context.Context) (groupIdOrPath interface{}, err error) {
	var gid interface{} = nil
	if g.IsIdSet() {
		gid, err = g.GetId(ctx)
		if err != nil {
			return nil, err
		}
	}
	if g.IsGroupPathSet() {
		gid, err = g.GetGroupPath(ctx)
		if err != nil {
			return nil, err
		}
	}
	if gid == nil {
		return nil, tracederrors.TracedError("Either group id or path must be set to get raw response")
	}

	return gid, nil
}

func (g *GitlabGroup) GetGroupIdOrPathAsString(ctx context.Context) (groupIdOrPath string, err error) {
	gid, err := g.GetGroupIdOrPath(ctx)
	if err != nil {
		return "", err
	}

	groupIdOrPath = fmt.Sprintf("%v", gid)

	return groupIdOrPath, nil
}

func (g *GitlabGroup) GetGroupName(ctx context.Context) (groupName string, err error) {
	groupPath, err := g.GetGroupPath(ctx)
	if err != nil {
		return "", err
	}

	groupName = filepath.Base(groupPath)
	groupName = strings.TrimSpace(groupName)

	if groupName == "" {
		return "", tracederrors.TracedError("groupName is empty string after evaluation")
	}

	return groupName, nil
}

func (g *GitlabGroup) GetGroupPath(ctx context.Context) (groupPath string, err error) {
	if g.groupPath != "" {
		return g.groupPath, nil
	} else {
		rawResponse, err := g.GetRawResponse(ctx)
		if err != nil {
			return "", err
		}

		groupPath = rawResponse.FullPath
		if groupPath == "" {
			return "", tracederrors.TracedError("groupPath is empty string after evaluation.")
		}

		return groupPath, nil
	}
}

func (g *GitlabGroup) GetGroupPathAndId(ctx context.Context) (groupPath string, groupId int, err error) {
	groupPath, err = g.GetGroupPath(ctx)
	if err != nil {
		return "", -1, err
	}

	groupId, err = g.GetId(ctx)
	if err != nil {
		return "", -1, err
	}

	return groupPath, groupId, nil
}

func (g *GitlabGroup) GetGroupPathAndIdOrEmptyIfUnset() (groupPath string, groupId int) {
	groupPath = g.GetGroupPathOrEmptyStringIfUnset()

	groupId = g.GetIdOrMinusOneIfUnset()

	return groupPath, groupId
}

func (g *GitlabGroup) GetGroupPathOrEmptyStringIfUnset() (groupPath string) {
	return g.groupPath
}

func (g *GitlabGroup) GetIdOrMinusOneIfUnset() (id int) {
	if g.IsIdSet() {
		return g.id
	}

	return -1
}

func (g *GitlabGroup) GetNativeGroupsService() (nativeGroupsService *gitlab.GroupsService, err error) {
	gitlabGroups, err := g.GetGitlabGroups()
	if err != nil {
		return nil, err
	}

	nativeGroupsService, err = gitlabGroups.GetNativeGroupsService()
	if err != nil {
		return nil, err
	}

	return nativeGroupsService, nil
}

func (g *GitlabGroup) GetParentGroup(ctx context.Context) (parentGroup *GitlabGroup, err error) {
	parentGroupPath, err := g.GetParentGroupPath(ctx)
	if err != nil {
		return nil, err
	}

	groups, err := g.GetGitlabGroups()
	if err != nil {
		return nil, err
	}

	parentGroup, err = groups.GetGroupByPath(ctx, parentGroupPath)
	if err != nil {
		return nil, err
	}

	return parentGroup, nil
}

func (g *GitlabGroup) GetRawResponse(ctx context.Context) (rawRespoonse *gitlab.Group, err error) {
	nativeGroupsService, err := g.GetNativeGroupsService()
	if err != nil {
		return nil, err
	}

	gid, err := g.GetGroupIdOrPath(ctx)
	if err != nil {
		return nil, err
	}

	nativeGroup, _, err := nativeGroupsService.GetGroup(gid, &gitlab.GetGroupOptions{})
	if err != nil {
		if slicesutils.ContainsStringIgnoreCase([]string{"404 Not Found", "404 {message: 404 Group Not Found}"}, err.Error()) {
			return nil, tracederrors.TracedErrorf("%w: %v", ErrGitlabGroupNotFoundError, gid)
		}
		return nil, err
	}

	if nativeGroup == nil {
		return nil, tracederrors.TracedErrorf("NativeGroup is nil after evaluation: gid = '%v'", gid)
	}

	return nativeGroup, err
}

func (g *GitlabGroup) IsGroupPathSet() (isSet bool) {
	return g.groupPath != ""
}

func (g *GitlabGroup) IsIdSet() (isSet bool) {
	return g.id > 0
}

func (g *GitlabGroup) IsSubgroup(ctx context.Context) (isSubgroup bool, err error) {
	path, err := g.GetGroupPath(ctx)
	if err != nil {
		return false, err
	}

	isSubgroup = strings.Contains(path, "/")

	return isSubgroup, nil
}

func (g *GitlabGroup) ListProjectPaths(ctx context.Context, options *GitlabListProjectsOptions) (projectPaths []string, err error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	nativeService, err := g.GetNativeGroupsService()
	if err != nil {
		return nil, err
	}

	gid, err := g.GetId(ctx)
	if err != nil {
		return nil, err
	}

	groupPath, err := g.GetGroupPath(ctx)
	if err != nil {
		return nil, err
	}

	nextPage := 1
	falseBoolean := false

	projectPaths = []string{}

	recursive := options.Recursive

	for {
		if nextPage <= 0 {
			break
		}

		nativeProjects, response, err := nativeService.ListGroupProjects(
			gid,
			&gitlab.ListGroupProjectsOptions{
				Archived:         &falseBoolean,
				IncludeSubGroups: &recursive,
				ListOptions: gitlab.ListOptions{
					Page: nextPage,
				},
			},
		)
		if err != nil {
			return nil, tracederrors.TracedErrorf("Failed to get project paths of group '%s', groupId='%d': %w", groupPath, gid, err)
		}

		for _, toAdd := range nativeProjects {
			pathToAdd := toAdd.PathWithNamespace
			projectPaths = append(projectPaths, pathToAdd)
		}

		nextPage = response.NextPage
	}

	logging.LogInfoByCtxf(ctx, "Collected '%d' project paths in gitlab group '%s'.", len(projectPaths), groupPath)

	return projectPaths, nil
}

func (g *GitlabGroup) ListProjects(ctx context.Context, listProjectOptions *GitlabListProjectsOptions) (projects []*GitlabProject, err error) {
	if listProjectOptions == nil {
		return nil, tracederrors.TracedErrorNil("listProjectOptions")
	}

	projectPaths, err := g.ListProjectPaths(ctx, listProjectOptions)
	if err != nil {
		return nil, err
	}

	gitlab, err := g.GetGitlab()
	if err != nil {
		return nil, err
	}

	projects = []*GitlabProject{}
	for _, path := range projectPaths {
		toAdd, err := gitlab.GetGitlabProjectByPath(ctx, path)
		if err != nil {
			return nil, err
		}

		projects = append(projects, toAdd)
	}

	return projects, nil
}

func (g *GitlabGroup) SetGroupPath(groupPath string) (err error) {
	if groupPath == "" {
		return tracederrors.TracedErrorf("groupPath is empty string")
	}

	trimmed := stringsutils.TrimPrefixAndSuffix(groupPath, "/", "/")
	trimmed = strings.TrimSpace(trimmed)

	if trimmed == "" {
		return tracederrors.TracedErrorf("groupPath '%s' is invalid. It results in an empty string after avaluation.", groupPath)
	}

	g.groupPath = trimmed

	return nil
}

func (p *GitlabGroup) GetFqdn() (fqdn string, err error) {
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

func (p *GitlabGroup) GetGitlab() (gitlab *GitlabInstance, err error) {
	if p.gitlab == nil {
		return nil, tracederrors.TracedError("gitlab is not set")
	}

	return p.gitlab, nil
}

func (p *GitlabGroup) GetId(ctx context.Context) (id int, err error) {
	if p.id > 0 {
		return p.id, nil
	}

	rawResponse, err := p.GetRawResponse(ctx)
	if err != nil {
		return -1, err
	}

	id = rawResponse.ID
	if id <= 0 {
		return -1, tracederrors.TracedErrorf("id '%d' is invalid after evaluation.", id)
	}

	return id, nil
}

func (p *GitlabGroup) GetNativeClient() (nativeClient *gitlab.Client, err error) {
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

func (p *GitlabGroup) GetParentGroupPath(ctx context.Context) (parentGroupPath string, err error) {
	path, err := p.GetGroupPath(ctx)
	if err != nil {
		return "", err
	}

	parentGroupPath = filepath.Dir(path)

	if parentGroupPath == "" {
		return "", tracederrors.TracedErrorf("parentGroupPath is empty string after evaluation, path = '%s'", path)
	}

	logging.LogInfoByCtxf(ctx, "Parent group of gitlab project '%s' is '%s'.", path, parentGroupPath)

	return parentGroupPath, nil
}

func (p *GitlabGroup) SetGitlab(gitlab *GitlabInstance) (err error) {
	if gitlab == nil {
		return tracederrors.TracedError("gitlab is nil")
	}

	p.gitlab = gitlab

	return nil
}

func (p *GitlabGroup) SetId(id int) (err error) {
	if id <= 0 {
		return tracederrors.TracedErrorf("invalid id = '%d'", id)
	}

	p.id = id

	return nil
}
