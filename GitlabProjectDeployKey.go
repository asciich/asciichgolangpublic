package asciichgolangpublic

import (
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
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

func (g *GitlabProjectDeployKey) MustCreateDeployKey(createOptions *GitlabCreateDeployKeyOptions) {
	err := g.CreateDeployKey(createOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabProjectDeployKey) MustDelete(verbose bool) {
	err := g.Delete(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabProjectDeployKey) MustExists() (exists bool) {
	exists, err := g.Exists()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return exists
}

func (g *GitlabProjectDeployKey) MustGetGitlab() (gitlab *GitlabInstance) {
	gitlab, err := g.GetGitlab()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlab
}

func (g *GitlabProjectDeployKey) MustGetGitlabProject() (gitlabProject *GitlabProject) {
	gitlabProject, err := g.GetGitlabProject()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlabProject
}

func (g *GitlabProjectDeployKey) MustGetGitlabProjectDeployKeys() (gitlabProjectProjectDeployKeys *GitlabProjectDeployKeys) {
	gitlabProjectProjectDeployKeys, err := g.GetGitlabProjectDeployKeys()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlabProjectProjectDeployKeys
}

func (g *GitlabProjectDeployKey) MustGetId() (id int) {
	id, err := g.GetId()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return id
}

func (g *GitlabProjectDeployKey) MustGetName() (name string) {
	name, err := g.GetName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return name
}

func (g *GitlabProjectDeployKey) MustGetNativeProjectDeployKeyService() (nativeService *gitlab.DeployKeysService) {
	nativeService, err := g.GetNativeProjectDeployKeyService()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return nativeService
}

func (g *GitlabProjectDeployKey) MustGetProjectId() (id int) {
	id, err := g.GetProjectId()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return id
}

func (g *GitlabProjectDeployKey) MustRecreateDeployKey(createOptions *GitlabCreateDeployKeyOptions) {
	err := g.RecreateDeployKey(createOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabProjectDeployKey) MustSetGitlabProjectDeployKeys(gitlabProjectDeployKeys *GitlabProjectDeployKeys) {
	err := g.SetGitlabProjectDeployKeys(gitlabProjectDeployKeys)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabProjectDeployKey) MustSetId(id int) {
	err := g.SetId(id)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabProjectDeployKey) MustSetName(name string) {
	err := g.SetName(name)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (k *GitlabProjectDeployKey) CreateDeployKey(createOptions *GitlabCreateDeployKeyOptions) (err error) {
	if createOptions == nil {
		return tracederrors.TracedError("createOptions is nil")
	}

	keyName, err := createOptions.GetName()
	if err != nil {
		return err
	}

	exists, err := k.Exists()
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

	projectId, err := k.GetProjectId()
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

	if createOptions.Verbose {
		logging.LogInfof("Created project deploy key '%s'", keyName)
	}

	return nil
}

func (k *GitlabProjectDeployKey) Delete(verbose bool) (err error) {
	deployKeyExists, err := k.Exists()
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

		projectId, err := k.GetProjectId()
		if err != nil {
			return err
		}

		keyId, err := k.GetId()
		if err != nil {
			return err
		}

		_, err = nativeProjectDeployKeyService.DeleteDeployKey(projectId, keyId)
		if err != nil {
			return tracederrors.TracedError(err.Error())
		}

		if verbose {
			logging.LogInfof("Project deploy key '%s' Deleted.", keyName)
		}
	} else {
		if verbose {
			logging.LogInfof("Project deploy key '%s' already absent. Skip deletion.", keyName)
		}
	}

	return nil
}

func (k *GitlabProjectDeployKey) Exists() (exists bool, err error) {
	gitlabProject, err := k.GetGitlabProject()
	if err != nil {
		return false, err
	}

	keyName, err := k.GetName()
	if err != nil {
		return false, err
	}

	exists, err = gitlabProject.DeployKeyByNameExists(keyName)
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

func (k *GitlabProjectDeployKey) GetId() (id int, err error) {
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

	id, err = deployKeys.GetKeyIdByKeyName(name)
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

func (k *GitlabProjectDeployKey) GetProjectId() (id int, err error) {
	project, err := k.GetGitlabProject()
	if err != nil {
		return -1, err
	}

	id, err = project.GetId()
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (k *GitlabProjectDeployKey) RecreateDeployKey(createOptions *GitlabCreateDeployKeyOptions) (err error) {
	if createOptions == nil {
		return tracederrors.TracedError("createOptions is nil")
	}

	keyName, err := createOptions.GetName()
	if err != nil {
		return err
	}

	if createOptions.Verbose {
		logging.LogInfof("Recreate gitlab project deploy key '%s' started.", keyName)
	}

	err = k.Delete(createOptions.Verbose)
	if err != nil {
		return err
	}

	err = k.CreateDeployKey(createOptions)
	if err != nil {
		return err
	}

	if createOptions.Verbose {
		logging.LogInfof("Recreate gitlab project deploy key '%s' finished.", keyName)
	}

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
