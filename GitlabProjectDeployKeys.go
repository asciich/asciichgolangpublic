package asciichgolangpublic

import (
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type GitlabProjectDeployKeys struct {
	gitlabProject *GitlabProject
}

func NewGitlabProjectDeployKeys() (deployKeys *GitlabProjectDeployKeys) {
	return new(GitlabProjectDeployKeys)
}

func (g *GitlabProjectDeployKeys) MustDeployKeyByNameExists(keyName string) (keyExists bool) {
	keyExists, err := g.DeployKeyByNameExists(keyName)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return keyExists
}

func (g *GitlabProjectDeployKeys) MustGetGitlab() (gitlab *GitlabInstance) {
	gitlab, err := g.GetGitlab()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlab
}

func (g *GitlabProjectDeployKeys) MustGetGitlabProject() (gitlabProject *GitlabProject) {
	gitlabProject, err := g.GetGitlabProject()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlabProject
}

func (g *GitlabProjectDeployKeys) MustGetGitlabProjectDeployKeyByName(keyName string) (deployKey *GitlabProjectDeployKey) {
	deployKey, err := g.GetGitlabProjectDeployKeyByName(keyName)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return deployKey
}

func (g *GitlabProjectDeployKeys) MustGetKeyIdByKeyName(keyName string) (id int) {
	id, err := g.GetKeyIdByKeyName(keyName)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return id
}

func (g *GitlabProjectDeployKeys) MustGetKeyNameList() (keyNames []string) {
	keyNames, err := g.GetKeyNameList()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return keyNames
}

func (g *GitlabProjectDeployKeys) MustGetKeysList() (keys []*GitlabProjectDeployKey) {
	keys, err := g.GetKeysList()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return keys
}

func (g *GitlabProjectDeployKeys) MustGetNativeGitlabClient() (nativeGitlabClient *gitlab.Client) {
	nativeGitlabClient, err := g.GetNativeGitlabClient()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return nativeGitlabClient
}

func (g *GitlabProjectDeployKeys) MustGetNativeProjectDeployKeyService() (nativeService *gitlab.DeployKeysService) {
	nativeService, err := g.GetNativeProjectDeployKeyService()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return nativeService
}

func (g *GitlabProjectDeployKeys) MustGetProjectId() (id int) {
	id, err := g.GetProjectId()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return id
}

func (g *GitlabProjectDeployKeys) MustSetGitlabProject(gitlabProject *GitlabProject) {
	err := g.SetGitlabProject(gitlabProject)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (k *GitlabProjectDeployKeys) DeployKeyByNameExists(keyName string) (keyExists bool, err error) {
	if len(keyName) <= 0 {
		return false, TracedError("keyName is empty string")
	}

	keyNameList, err := k.GetKeyNameList()
	if err != nil {
		return false, err
	}

	keyExists = Slices().ContainsString(keyNameList, keyName)

	return keyExists, nil
}

func (k *GitlabProjectDeployKeys) GetGitlab() (gitlab *GitlabInstance, err error) {
	gitlabProject, err := k.GetGitlabProject()
	if err != nil {
		return nil, err
	}

	gitlab, err = gitlabProject.GetGitlab()
	if err != nil {
		return nil, err
	}

	return gitlab, nil
}

func (k *GitlabProjectDeployKeys) GetGitlabProject() (gitlabProject *GitlabProject, err error) {
	if k.gitlabProject == nil {
		return nil, TracedError("gitalbProject not set")
	}

	return k.gitlabProject, nil
}

func (k *GitlabProjectDeployKeys) GetGitlabProjectDeployKeyByName(keyName string) (deployKey *GitlabProjectDeployKey, err error) {
	if len(keyName) <= 0 {
		return nil, TracedError("keyName is empty string")
	}

	deployKey = NewGitlabProjectDeployKey()

	err = deployKey.SetGitlabProjectDeployKeys(k)
	if err != nil {
		return nil, err
	}

	err = deployKey.SetName(keyName)
	if err != nil {
		return nil, err
	}

	return deployKey, nil
}

func (k *GitlabProjectDeployKeys) GetKeyIdByKeyName(keyName string) (id int, err error) {
	if len(keyName) <= 0 {
		return -1, TracedError("keyName is empty string")
	}

	keys, err := k.GetKeysList()
	if err != nil {
		return -1, err
	}

	for _, key := range keys {
		if key.name == keyName {
			id = key.id
			break
		}
	}

	if id <= 0 {
		return -1, TracedErrorf("Unable to get gitlab project deploy key id for '%s'", keyName)
	}

	return id, nil
}

func (k *GitlabProjectDeployKeys) GetKeyNameList() (keyNames []string, err error) {
	keys, err := k.GetKeysList()
	if err != nil {
		return nil, err
	}

	keyNames = []string{}
	for _, key := range keys {
		keyNameToAdd, err := key.GetName()
		if err != nil {
			return nil, err
		}
		keyNames = append(keyNames, keyNameToAdd)
	}

	return keyNames, nil
}

func (k *GitlabProjectDeployKeys) GetKeysList() (keys []*GitlabProjectDeployKey, err error) {
	nativeService, err := k.GetNativeProjectDeployKeyService()
	if err != nil {
		return nil, err
	}

	projectId, err := k.GetProjectId()
	if err != nil {
		return nil, err
	}

	nativeKeys, _, err := nativeService.ListProjectDeployKeys(projectId, &gitlab.ListProjectDeployKeysOptions{})
	if err != nil {
		return nil, err
	}

	keys = []*GitlabProjectDeployKey{}
	for _, nativeKey := range nativeKeys {
		keyToAdd := NewGitlabProjectDeployKey()
		err = keyToAdd.SetGitlabProjectDeployKeys(k)
		if err != nil {
			return nil, err
		}

		nameToAdd := nativeKey.Title
		err = keyToAdd.SetName(nameToAdd)
		if err != nil {
			return nil, err
		}

		err = keyToAdd.SetId(nativeKey.ID)
		if err != nil {
			return nil, err
		}

		keys = append(keys, keyToAdd)
	}

	return keys, nil
}

func (k *GitlabProjectDeployKeys) GetNativeGitlabClient() (nativeGitlabClient *gitlab.Client, err error) {
	gitlab, err := k.GetGitlab()
	if err != nil {
		return nil, err
	}

	nativeGitlabClient, err = gitlab.GetNativeClient()
	if err != nil {
		return nil, err
	}

	return nativeGitlabClient, nil
}

func (k *GitlabProjectDeployKeys) GetNativeProjectDeployKeyService() (nativeService *gitlab.DeployKeysService, err error) {
	nativeClient, err := k.GetNativeGitlabClient()
	if err != nil {
		return nil, err
	}

	nativeService = nativeClient.DeployKeys
	if nativeService == nil {
		return nil, TracedError("unable to get nativeService. nativeService from nativeClient is nil")
	}

	return nativeService, nil
}

func (k *GitlabProjectDeployKeys) GetProjectId() (id int, err error) {
	gitlabProject, err := k.GetGitlabProject()
	if err != nil {
		return -1, err
	}

	id, err = gitlabProject.GetId()
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (k *GitlabProjectDeployKeys) SetGitlabProject(gitlabProject *GitlabProject) (err error) {
	if gitlabProject == nil {
		return TracedError("gitlabProject is nil")
	}

	k.gitlabProject = gitlabProject

	return nil
}
