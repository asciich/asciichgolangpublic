package asciichgolangpublic

import (
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

func (g *GitlabGroup) Create(createOptions *GitlabCreateGroupOptions) (err error) {
	if createOptions == nil {
		return tracederrors.TracedErrorNil("createOptions")
	}

	if !g.IsGroupPathSet() {
		return tracederrors.TracedError("Group path must be set for a group to create")
	}
	groupPath, err := g.GetGroupPath()
	if err != nil {
		return err
	}

	fqdn, err := g.GetGitlabFqdn()
	if err != nil {
		return err
	}

	if createOptions.Verbose {
		logging.LogInfof("Going to create group '%s' on gitlab '%s'.", groupPath, fqdn)
	}

	exists, err := g.Exists(createOptions.Verbose)
	if err != nil {
		return err
	}

	if exists {
		if createOptions.Verbose {
			logging.LogInfof("Group '%s' already exists on gitlab '%s'.", groupPath, fqdn)
		}
	} else {
		isSubgroup, err := g.IsSubgroup()
		if err != nil {
			return err
		}

		var parentGroup *GitlabGroup = nil

		if isSubgroup {
			parentGroup, err = g.GetParentGroup(createOptions.Verbose)
			if err != nil {
				return err
			}

			optionsToUse := createOptions.GetDeepCopy()

			err = parentGroup.Create(optionsToUse)
			if err != nil {
				return err
			}
		}

		nativeGroupsService, err := g.GetNativeGroupsService()
		if err != nil {
			return err
		}

		groupName, err := g.GetGroupName()
		if err != nil {
			return err
		}

		if len(groupName) < 2 {
			return tracederrors.TracedErrorf("Group names with len < 2 not allowed. But got '%s'", groupName)
		}

		var namespaceId *int = nil
		if parentGroup != nil {
			idToUse, err := parentGroup.GetId()
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
			if createOptions.Verbose {
				logging.LogErrorf("Creating group '%s' failed", groupPath)
			}
			return err
		}

		if createOptions.Verbose {
			logging.LogChangedf("Group '%s' on gitlab '%s' created.", groupPath, fqdn)
		}
	}

	if createOptions.Verbose {
		logging.LogInfof("Create group '%s' in gitlab '%s' finished.", groupPath, fqdn)
	}

	return nil
}

func (g *GitlabGroup) Delete(verbose bool) (err error) {
	gid, err := g.GetGroupIdOrPath()
	if err != nil {
		return err
	}

	exists, err := g.Exists(verbose)
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
			exists, err = g.Exists(verbose)
			if err != nil {
				return err
			}

			if exists {
				if verbose {
					logging.LogInfof("Wait until delete of gitlab group '%v' is finished (%d/%d).", gid, i+1, maxRetry)
				}
				time.Sleep(time.Millisecond * 500)
				continue
			}

			deleteConfirmed = true
			break
		}

		if !deleteConfirmed {
			return tracederrors.TracedErrorf("Failed to delete gitlab group '%v'. Group still exists after delete.", gid)
		}

		if verbose {
			logging.LogChangedf("Group '%s' deleted", gid)
		}
	} else {
		if verbose {
			logging.LogInfof("Group '%s' is already absent.", gid)
		}
	}

	return nil
}

func (g *GitlabGroup) Exists(verbose bool) (exists bool, err error) {
	exists = true

	_, err = g.GetRawResponse()
	if err != nil {
		if errors.Is(err, ErrGitlabGroupNotFoundError) {
			exists = false
		} else {
			return false, err
		}
	}

	if verbose {
		gid, err := g.GetGroupIdOrPathAsString()
		if err != nil {
			return false, err
		}

		if exists {
			logging.LogInfof("Gitlab group '%s' exists.", gid)
		} else {
			logging.LogInfof("Gitlab group '%s' does not exist.", gid)
		}
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

func (g *GitlabGroup) GetGroupIdOrPath() (groupIdOrPath interface{}, err error) {
	var gid interface{} = nil
	if g.IsIdSet() {
		gid, err = g.GetId()
		if err != nil {
			return nil, err
		}
	}
	if g.IsGroupPathSet() {
		gid, err = g.GetGroupPath()
		if err != nil {
			return nil, err
		}
	}
	if gid == nil {
		return nil, tracederrors.TracedError("Either group id or path must be set to get raw response")
	}

	return gid, nil
}

func (g *GitlabGroup) GetGroupIdOrPathAsString() (groupIdOrPath string, err error) {
	gid, err := g.GetGroupIdOrPath()
	if err != nil {
		return "", err
	}

	groupIdOrPath = fmt.Sprintf("%v", gid)

	return groupIdOrPath, nil
}

func (g *GitlabGroup) GetGroupName() (groupName string, err error) {
	groupPath, err := g.GetGroupPath()
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

func (g *GitlabGroup) GetGroupPath() (groupPath string, err error) {
	if g.groupPath != "" {
		return g.groupPath, nil
	} else {
		rawResponse, err := g.GetRawResponse()
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

func (g *GitlabGroup) GetGroupPathAndId() (groupPath string, groupId int, err error) {
	groupPath, err = g.GetGroupPath()
	if err != nil {
		return "", -1, err
	}

	groupId, err = g.GetId()
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

func (g *GitlabGroup) GetParentGroup(verbose bool) (parentGroup *GitlabGroup, err error) {
	parentGroupPath, err := g.GetParentGroupPath(verbose)
	if err != nil {
		return nil, err
	}

	groups, err := g.GetGitlabGroups()
	if err != nil {
		return nil, err
	}

	parentGroup, err = groups.GetGroupByPath(parentGroupPath, verbose)
	if err != nil {
		return nil, err
	}

	return parentGroup, nil
}

func (g *GitlabGroup) GetRawResponse() (rawRespoonse *gitlab.Group, err error) {
	nativeGroupsService, err := g.GetNativeGroupsService()
	if err != nil {
		return nil, err
	}

	gid, err := g.GetGroupIdOrPath()
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

func (g *GitlabGroup) IsSubgroup() (isSubgroup bool, err error) {
	path, err := g.GetGroupPath()
	if err != nil {
		return false, err
	}

	isSubgroup = strings.Contains(path, "/")

	return isSubgroup, nil
}

func (g *GitlabGroup) ListProjectPaths(options *GitlabListProjectsOptions) (projectPaths []string, err error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	nativeService, err := g.GetNativeGroupsService()
	if err != nil {
		return nil, err
	}

	gid, err := g.GetId()
	if err != nil {
		return nil, err
	}

	groupPath, err := g.GetGroupPath()
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

	if options.Verbose {
		logging.LogInfof("Collected '%d' project paths in gitlab group '%s'.", len(projectPaths), groupPath)
	}

	return projectPaths, nil
}

func (g *GitlabGroup) MustListProjects(listProjectOptions *GitlabListProjectsOptions) (projects []*GitlabProject) {
	projects, err := g.ListProjects(listProjectOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return projects
}

func (g *GitlabGroup) ListProjects(listProjectOptions *GitlabListProjectsOptions) (projects []*GitlabProject, err error) {
	if listProjectOptions == nil {
		return nil, tracederrors.TracedErrorNil("listProjectOptions")
	}

	projectPaths, err := g.ListProjectPaths(listProjectOptions)
	if err != nil {
		return nil, err
	}

	gitlab, err := g.GetGitlab()
	if err != nil {
		return nil, err
	}

	projects = []*GitlabProject{}
	for _, path := range projectPaths {
		toAdd, err := gitlab.GetGitlabProjectByPath(path, listProjectOptions.Verbose)
		if err != nil {
			return nil, err
		}

		projects = append(projects, toAdd)
	}

	return projects, nil
}

func (g *GitlabGroup) MustCreate(createOptions *GitlabCreateGroupOptions) {
	err := g.Create(createOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabGroup) MustDelete(verbose bool) {
	err := g.Delete(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabGroup) MustExists(verbose bool) (exists bool) {
	exists, err := g.Exists(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return exists
}

func (g *GitlabGroup) MustGetFqdn() (fqdn string) {
	fqdn, err := g.GetFqdn()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return fqdn
}

func (g *GitlabGroup) MustGetGitlab() (gitlab *GitlabInstance) {
	gitlab, err := g.GetGitlab()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlab
}

func (g *GitlabGroup) MustGetGitlabFqdn() (gitlabFqdn string) {
	gitlabFqdn, err := g.GetGitlabFqdn()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlabFqdn
}

func (g *GitlabGroup) MustGetGitlabGroups() (gitlabGroups *GitlabGroups) {
	gitlabGroups, err := g.GetGitlabGroups()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlabGroups
}

func (g *GitlabGroup) MustGetGroupIdOrPath() (groupIdOrPath interface{}) {
	groupIdOrPath, err := g.GetGroupIdOrPath()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return groupIdOrPath
}

func (g *GitlabGroup) MustGetGroupIdOrPathAsString() (groupIdOrPath string) {
	groupIdOrPath, err := g.GetGroupIdOrPathAsString()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return groupIdOrPath
}

func (g *GitlabGroup) MustGetGroupName() (groupName string) {
	groupName, err := g.GetGroupName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return groupName
}

func (g *GitlabGroup) MustGetGroupPath() (groupPath string) {
	groupPath, err := g.GetGroupPath()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return groupPath
}

func (g *GitlabGroup) MustGetGroupPathAndId() (groupPath string, groupId int) {
	groupPath, groupId, err := g.GetGroupPathAndId()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return groupPath, groupId
}

func (g *GitlabGroup) MustGetId(verbose bool) (id int) {
	id, err := g.GetId()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return id
}

func (g *GitlabGroup) MustGetNativeClient() (nativeClient *gitlab.Client) {
	nativeClient, err := g.GetNativeClient()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return nativeClient
}

func (g *GitlabGroup) MustGetNativeGroupsService() (nativeGroupsService *gitlab.GroupsService) {
	nativeGroupsService, err := g.GetNativeGroupsService()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return nativeGroupsService
}

func (g *GitlabGroup) MustGetParentGroup(verbose bool) (parentGroup *GitlabGroup) {
	parentGroup, err := g.GetParentGroup(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return parentGroup
}

func (g *GitlabGroup) MustGetParentGroupPath(verbose bool) (parentGroupPath string) {
	parentGroupPath, err := g.GetParentGroupPath(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return parentGroupPath
}

func (g *GitlabGroup) MustGetRawResponse() (rawRespoonse *gitlab.Group) {
	rawRespoonse, err := g.GetRawResponse()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return rawRespoonse
}

func (g *GitlabGroup) MustIsSubgroup() (isSubgroup bool) {
	isSubgroup, err := g.IsSubgroup()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return isSubgroup
}

func (g *GitlabGroup) MustListProjectPaths(options *GitlabListProjectsOptions) (projectPaths []string) {
	projectPaths, err := g.ListProjectPaths(options)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return projectPaths
}

func (g *GitlabGroup) MustSetGitlab(gitlab *GitlabInstance) {
	err := g.SetGitlab(gitlab)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabGroup) MustSetGroupPath(groupPath string) {
	err := g.SetGroupPath(groupPath)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabGroup) MustSetId(id int) {
	err := g.SetId(id)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
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

func (p *GitlabGroup) GetId() (id int, err error) {
	if p.id > 0 {
		return p.id, nil
	}

	rawResponse, err := p.GetRawResponse()
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

func (p *GitlabGroup) GetParentGroupPath(verbose bool) (parentGroupPath string, err error) {
	path, err := p.GetGroupPath()
	if err != nil {
		return "", err
	}

	parentGroupPath = filepath.Dir(path)

	if parentGroupPath == "" {
		return "", tracederrors.TracedErrorf("parentGroupPath is empty string after evaluation, path = '%s'", path)
	}

	if verbose {
		logging.LogInfof("Parent group of gitlab project '%s' is '%s'.", path, parentGroupPath)
	}

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
