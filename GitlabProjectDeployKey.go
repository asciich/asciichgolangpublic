package asciichgolangpublic

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type GitlabProjectDeployKey struct {
	gitlabProjectDeployKeys *GitlabProjectDeployKeys
	id                      int
	name                    string
}

func NewGitlabProjectDeployKey() (projectDeployKey *GitlabProjectDeployKey) {
	return new(GitlabProjectDeployKey)
}

func (k *GitlabProjectDeployKey) CreateDeployKey(ctx context.Context, createOptions *GitlabCreateDeployKeyOptions) (err error) {
	if createOptions == nil {
		return tracederrors.TracedError("createOptions is nil")
	}

	keyName, err := createOptions.GetName()
	if err != nil {
		return err
	}

	exists, err := k.Exists(ctx)
	if err != nil {
		return err
	}

	if exists {
		return tracederrors.TracedError("Key '%s' already exists. To recreate use the RecreateDeployKey function instead")
	}

	nativeProjectDeployKeyService, err := k.GetNativeProjectDeployKeyService()
	if err != nil {
		return err
	}

	projectId, err := k.GetProjectId(ctx)
	if err != nil {
		return err
	}

	keyMaterial, err := createOptions.GetPublicKeyMaterialString()
	if err != nil {
		return err
	}

	_, _, err = nativeProjectDeployKeyService.AddDeployKey(projectId, &gitlab.AddDeployKeyOptions{
		Title:   &keyName,
		Key:     &keyMaterial,
		CanPush: &createOptions.WriteAccess,
	})
	if err != nil {
		return tracederrors.TracedError(err.Error())
	}

	logging.LogInfoByCtxf(ctx, "Created project deploy key '%s'", keyName)

	return nil
}

func (k *GitlabProjectDeployKey) Delete(ctx context.Context) (err error) {
	deployKeyExists, err := k.Exists(ctx)
	if err != nil {
		return err
	}

	keyName, err := k.GetName()
	if err != nil {
		return err
	}

	if deployKeyExists {
		nativeProjectDeployKeyService, err := k.GetNativeProjectDeployKeyService()
		if err != nil {
			return err
		}

		projectId, err := k.GetProjectId(ctx)
		if err != nil {
			return err
		}

		keyId, err := k.GetId(ctx)
		if err != nil {
			return err
		}

		_, err = nativeProjectDeployKeyService.DeleteDeployKey(projectId, keyId)
		if err != nil {
			return tracederrors.TracedError(err.Error())
		}

		logging.LogChangedByCtxf(ctx, "Project deploy key '%s' Deleted.", keyName)
	} else {
		logging.LogInfoByCtxf(ctx, "Project deploy key '%s' already absent. Skip deletion.", keyName)
	}

	return nil
}

func (k *GitlabProjectDeployKey) Exists(ctx context.Context) (exists bool, err error) {
	gitlabProject, err := k.GetGitlabProject()
	if err != nil {
		return false, err
	}

	keyName, err := k.GetName()
	if err != nil {
		return false, err
	}

	exists, err = gitlabProject.DeployKeyByNameExists(ctx, keyName)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (k *GitlabProjectDeployKey) GetGitlab() (gitlab *GitlabInstance, err error) {
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

func (k *GitlabProjectDeployKey) GetGitlabProject() (gitlabProject *GitlabProject, err error) {
	deployKeys, err := k.GetGitlabProjectDeployKeys()
	if err != nil {
		return nil, err
	}

	gitlabProject, err = deployKeys.GetGitlabProject()
	if err != nil {
		return nil, err
	}

	return gitlabProject, nil
}

func (k *GitlabProjectDeployKey) GetGitlabProjectDeployKeys() (gitlabProjectProjectDeployKeys *GitlabProjectDeployKeys, err error) {
	if k.gitlabProjectDeployKeys == nil {
		return nil, tracederrors.TracedError("gitlabProject not set")
	}

	return k.gitlabProjectDeployKeys, nil
}

func (k *GitlabProjectDeployKey) GetId(ctx context.Context) (id int, err error) {
	if k.id > 0 {
		return k.id, nil
	}

	name, err := k.GetName()
	if err != nil {
		return -1, err
	}

	deployKeys, err := k.GetGitlabProjectDeployKeys()
	if err != nil {
		return -1, err
	}

	id, err = deployKeys.GetKeyIdByKeyName(ctx, name)
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (k *GitlabProjectDeployKey) GetName() (name string, err error) {
	if len(k.name) <= 0 {
		return "", tracederrors.TracedError("name not set")
	}

	return k.name, nil
}

func (k *GitlabProjectDeployKey) GetNativeProjectDeployKeyService() (nativeService *gitlab.DeployKeysService, err error) {
	deployKeys, err := k.GetGitlabProjectDeployKeys()
	if err != nil {
		return nil, err
	}

	nativeService, err = deployKeys.GetNativeProjectDeployKeyService()
	if err != nil {
		return nil, err
	}

	return nativeService, nil
}

func (k *GitlabProjectDeployKey) GetProjectId(ctx context.Context) (id int, err error) {
	project, err := k.GetGitlabProject()
	if err != nil {
		return -1, err
	}

	id, err = project.GetId(ctx)
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (k *GitlabProjectDeployKey) RecreateDeployKey(ctx context.Context, createOptions *GitlabCreateDeployKeyOptions) (err error) {
	if createOptions == nil {
		return tracederrors.TracedError("createOptions is nil")
	}

	keyName, err := createOptions.GetName()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Recreate gitlab project deploy key '%s' started.", keyName)

	err = k.Delete(ctx)
	if err != nil {
		return err
	}

	err = k.CreateDeployKey(ctx, createOptions)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Recreate gitlab project deploy key '%s' finished.", keyName)

	return nil
}

func (k *GitlabProjectDeployKey) SetGitlabProjectDeployKeys(gitlabProjectDeployKeys *GitlabProjectDeployKeys) (err error) {
	if gitlabProjectDeployKeys == nil {
		return tracederrors.TracedError("gitlabProject is nil")
	}

	k.gitlabProjectDeployKeys = gitlabProjectDeployKeys

	return nil
}

func (k *GitlabProjectDeployKey) SetId(id int) (err error) {
	if id <= 0 {
		return tracederrors.TracedErrorf("invalid id='%d'", id)
	}

	k.id = id

	return nil
}

func (k *GitlabProjectDeployKey) SetName(name string) (err error) {
	if len(name) <= 0 {
		return tracederrors.TracedError("name is empty string")
	}

	k.name = name

	return nil
}
