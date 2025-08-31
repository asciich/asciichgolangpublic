package asciichgolangpublic

import (
	"context"
	"slices"

	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type GitlabProjectDeployKeys struct {
	gitlabProject *GitlabProject
}

func NewGitlabProjectDeployKeys() (deployKeys *GitlabProjectDeployKeys) {
	return new(GitlabProjectDeployKeys)
}

func (k *GitlabProjectDeployKeys) DeployKeyByNameExists(ctx context.Context, keyName string) (keyExists bool, err error) {
	if len(keyName) <= 0 {
		return false, tracederrors.TracedError("keyName is empty string")
	}

	keyNameList, err := k.GetKeyNameList(ctx)
	if err != nil {
		return false, err
	}

	keyExists = slices.Contains(keyNameList, keyName)

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
		return nil, tracederrors.TracedError("gitalbProject not set")
	}

	return k.gitlabProject, nil
}

func (k *GitlabProjectDeployKeys) GetGitlabProjectDeployKeyByName(keyName string) (deployKey *GitlabProjectDeployKey, err error) {
	if len(keyName) <= 0 {
		return nil, tracederrors.TracedError("keyName is empty string")
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

func (k *GitlabProjectDeployKeys) GetKeyIdByKeyName(ctx context.Context, keyName string) (id int, err error) {
	if len(keyName) <= 0 {
		return -1, tracederrors.TracedError("keyName is empty string")
	}

	keys, err := k.GetKeysList(ctx)
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
		return -1, tracederrors.TracedErrorf("Unable to get gitlab project deploy key id for '%s'", keyName)
	}

	return id, nil
}

func (k *GitlabProjectDeployKeys) GetKeyNameList(ctx context.Context) (keyNames []string, err error) {
	keys, err := k.GetKeysList(ctx)
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

func (k *GitlabProjectDeployKeys) GetKeysList(ctx context.Context) (keys []*GitlabProjectDeployKey, err error) {
	nativeService, err := k.GetNativeProjectDeployKeyService()
	if err != nil {
		return nil, err
	}

	projectId, err := k.GetProjectId(ctx)
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
		return nil, tracederrors.TracedError("unable to get nativeService. nativeService from nativeClient is nil")
	}

	return nativeService, nil
}

func (k *GitlabProjectDeployKeys) GetProjectId(ctx context.Context) (id int, err error) {
	gitlabProject, err := k.GetGitlabProject()
	if err != nil {
		return -1, err
	}

	id, err = gitlabProject.GetId(ctx)
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (k *GitlabProjectDeployKeys) SetGitlabProject(gitlabProject *GitlabProject) (err error) {
	if gitlabProject == nil {
		return tracederrors.TracedError("gitlabProject is nil")
	}

	k.gitlabProject = gitlabProject

	return nil
}
