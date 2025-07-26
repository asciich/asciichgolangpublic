package asciichgolangpublic

import (
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/sshutils"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type GitlabUser struct {
	gitlab *GitlabInstance
	id     int

	cachedName     string
	cachedEmail    string
	cachedUsername string
}

func NewGitlabUser() (gitlabUser *GitlabUser) {
	return new(GitlabUser)
}

func (g *GitlabUser) GetCachedEmail() (cachedEmail string, err error) {
	if g.cachedEmail == "" {
		return "", tracederrors.TracedErrorf("cachedEmail not set")
	}

	return g.cachedEmail, nil
}

func (g *GitlabUser) GetCachedUsername() (cachedUsername string, err error) {
	if g.cachedUsername == "" {
		return "", tracederrors.TracedErrorf("cachedUsername not set")
	}

	return g.cachedUsername, nil
}

func (u *GitlabUser) AddSshKey(sshKey *sshutils.SSHPublicKey, verbose bool) (err error) {
	if sshKey == nil {
		return tracederrors.TracedError("sshKey is nil")
	}

	nativeUsersService, err := u.GetNativeUsersService()
	if err != nil {
		return err
	}

	userAtHost, err := sshKey.GetKeyUserAtHost()
	if err != nil {
		return err
	}

	keyMaterial, err := sshKey.GetAsPublicKeyLine()
	if err != nil {
		return err
	}

	username, err := u.GetChachedUsername()
	if err != nil {
		return err
	}

	userId, err := u.GetId()
	if err != nil {
		return err
	}

	if verbose {
		logging.LogInfof("Going to add SSH key to gitlab user '%s': '%s'", username, keyMaterial)
	}

	keyExists, err := u.SshKeyExists(sshKey)
	if err != nil {
		return err
	}

	if keyExists {
		if verbose {
			logging.LogInfof("SSH key '%s' already present for gitlab user '%s'.", keyMaterial, username)
		}
	} else {
		_, _, err = nativeUsersService.AddSSHKeyForUser(
			userId,
			&gitlab.AddSSHKeyOptions{
				Title: &userAtHost,
				Key:   &keyMaterial,
			},
		)
		if err != nil {
			return tracederrors.TracedError(err.Error())
		}

		if verbose {
			logging.LogChangedf("SSH key '%s' added for gitlab user '%s'.", keyMaterial, username)
		}
	}

	return nil
}

func (u *GitlabUser) AddSshKeysFromFile(sshKeysFile filesinterfaces.File, verbose bool) (err error) {
	if sshKeysFile == nil {
		return tracederrors.TracedError("sshKeysFile is nil")
	}

	username, err := u.GetCachedName()
	if err != nil {
		return err
	}

	sshKeys, err := sshutils.LoadPublicKeysFromFile(contextutils.GetVerbosityContextByBool(verbose), sshKeysFile)
	if err != nil {
		return err
	}

	if verbose {
		logging.LogInfof("Going to add '%d' SSH keys for gitlab user '%s'.", len(sshKeys), username)
	}

	for _, keyToAdd := range sshKeys {
		err = u.AddSshKey(keyToAdd, verbose)
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *GitlabUser) AddSshKeysFromFilePath(sshKeyFilePath string, verbose bool) (err error) {
	if len(sshKeyFilePath) <= 0 {
		return tracederrors.TracedError("sshKeyFilePath is empty string")
	}

	sshKeyFile, err := files.GetLocalFileByPath(sshKeyFilePath)
	if err != nil {
		return err
	}

	err = u.AddSshKeysFromFile(sshKeyFile, verbose)
	if err != nil {
		return err
	}

	return nil
}

func (u *GitlabUser) CreateAccessToken(options *GitlabCreateAccessTokenOptions) (newToken string, err error) {
	if options == nil {
		return "", tracederrors.TracedError("options is nil")
	}

	nativeService, err := u.GetNativeUsersService()
	if err != nil {
		return "", err
	}

	userId, err := u.GetId()
	if err != nil {
		return "", err
	}

	tokenName, err := options.GetTokenName()
	if err != nil {
		return "", err
	}

	scopes, err := options.GetScopes()
	if err != nil {
		return "", err
	}

	expiresAt, err := options.GetExipiresAtOrDefaultIfUnset()
	if err != nil {
		return "", err
	}

	nativeToken, _, err := nativeService.CreatePersonalAccessToken(
		userId,
		&gitlab.CreatePersonalAccessTokenOptions{
			Name:      &tokenName,
			ExpiresAt: (*gitlab.ISOTime)(expiresAt),
			Scopes:    &scopes,
		},
	)
	if err != nil {
		return "", tracederrors.TracedError(err.Error())
	}

	newToken = nativeToken.Token
	if len(newToken) <= 0 {
		return "", tracederrors.TracedError("newToken is empty string")
	}

	return newToken, nil
}

func (u *GitlabUser) GetCachedName() (name string, err error) {
	if len(u.cachedName) <= 0 {
		return "", tracederrors.TracedError("Cached name not set")
	}

	return u.cachedName, nil
}

func (u *GitlabUser) GetChachedUsername() (username string, err error) {
	if len(u.cachedUsername) <= 0 {
		return "", tracederrors.TracedError("Cached username not set")
	}

	return u.cachedUsername, nil
}

func (u *GitlabUser) GetFqdn() (fqdn string, err error) {
	gitlab, err := u.GetGitlab()
	if err != nil {
		return "", err
	}

	fqdn, err = gitlab.GetFqdn()
	if err != nil {
		return "", err
	}

	return fqdn, nil
}

func (u *GitlabUser) GetGitlab() (gitlab *GitlabInstance, err error) {
	if u.gitlab == nil {
		return nil, tracederrors.TracedError("gitlab not set")
	}

	return u.gitlab, nil
}

func (u *GitlabUser) GetGitlabUsers() (gitlabUsers *GitlabUsers, err error) {
	gitlab, err := u.GetGitlab()
	if err != nil {
		return nil, err
	}

	gitlabUsers, err = gitlab.GetGitlabUsers()
	if err != nil {
		return nil, err
	}

	return gitlabUsers, nil
}

func (u *GitlabUser) GetId() (id int, err error) {
	if u.id <= 0 {
		return -1, tracederrors.TracedError("id not set")
	}

	return u.id, nil
}

func (u *GitlabUser) GetNativeUsersService() (nativeUsersService *gitlab.UsersService, err error) {
	gitlabUsers, err := u.GetGitlabUsers()
	if err != nil {
		return nil, err
	}

	nativeUsersService, err = gitlabUsers.GetNativeUsersService()
	if err != nil {
		return nil, err
	}

	return nativeUsersService, nil
}

func (u *GitlabUser) GetRawNativeUser() (rawUser *gitlab.User, err error) {
	nativeUserService, err := u.GetNativeUsersService()
	if err != nil {
		return nil, err
	}

	id, err := u.GetId()
	if err != nil {
		return nil, err
	}

	rawUser, _, err = nativeUserService.GetUser(id, gitlab.GetUsersOptions{})
	if err != nil {
		return nil, tracederrors.TracedError(err.Error())
	}

	if rawUser == nil {
		return nil, tracederrors.TracedError("rawUser is nil")
	}

	return rawUser, nil
}

func (u *GitlabUser) GetSshKeys() (sshKeys []*sshutils.SSHPublicKey, err error) {
	sshKeysString, err := u.GetSshKeysAsString()
	if err != nil {
		return nil, err
	}

	sshKeys = []*sshutils.SSHPublicKey{}
	for _, keyString := range sshKeysString {
		keyToAdd := sshutils.NewSSHPublicKey()
		err = keyToAdd.SetFromString(keyString)
		if err != nil {
			return nil, err
		}

		sshKeys = append(sshKeys, keyToAdd)
	}

	return sshKeys, nil
}

func (u *GitlabUser) GetSshKeysAsString() (sshKeys []string, err error) {
	nativeUsersService, err := u.GetNativeUsersService()
	if err != nil {
		return nil, err
	}

	id, err := u.GetId()
	if err != nil {
		return nil, err
	}

	nativeSshKeys, _, err := nativeUsersService.ListSSHKeysForUser(
		id,
		&gitlab.ListSSHKeysForUserOptions{},
	)
	if err != nil {
		return nil, tracederrors.TracedError(err.Error())
	}

	sshKeys = []string{}
	for _, nativeKey := range nativeSshKeys {
		keyToAdd := nativeKey.Key
		sshKeys = append(sshKeys, keyToAdd)
	}

	return sshKeys, nil
}

func (u *GitlabUser) SetCachedEmail(email string) (err error) {
	if len(email) <= 0 {
		return tracederrors.TracedError("email is empty string")
	}

	u.cachedEmail = email

	return nil
}

func (u *GitlabUser) SetCachedName(name string) (err error) {
	if len(name) <= 0 {
		return tracederrors.TracedError("name is empty string")
	}

	u.cachedName = name

	return nil
}

func (u *GitlabUser) SetCachedUsername(username string) (err error) {
	if len(username) <= 0 {
		return tracederrors.TracedError("cached usernae is empty string")
	}

	u.cachedUsername = username

	return nil
}

func (u *GitlabUser) SetGitlab(gitlab *GitlabInstance) (err error) {
	if gitlab == nil {
		return tracederrors.TracedError("gitlab is nil")
	}

	u.gitlab = gitlab

	return nil
}

func (u *GitlabUser) SetId(id int) (err error) {
	if id <= 0 {
		return tracederrors.TracedErrorf("invalid id: '%d'", id)
	}

	u.id = id

	return nil
}

func (u *GitlabUser) SshKeyExists(sshKey *sshutils.SSHPublicKey) (keyExistsForUser bool, err error) {
	return false, tracederrors.TracedErrorNotImplemented()
	/* TODO enable again
	if sshKey == nil {
		return false, tracederrors.TracedError("sshKey is nil")
	}

	existingKeys, err := u.GetSshKeys()
	if err != nil {
		return false, err
	}

	keyExistsForUser = aslices.ContainsSshPublicKeyWithSameKeyMaterial(existingKeys, sshKey)
	return keyExistsForUser, nil
	*/
}

func (u *GitlabUser) UpdatePassword(newPassword string, verbose bool) (err error) {
	if len(newPassword) <= 0 {
		return tracederrors.TracedError("newPassword is empty string")
	}

	fqdn, err := u.GetFqdn()
	if err != nil {
		return err
	}

	username, err := u.GetChachedUsername()
	if err != nil {
		return err
	}

	nativeUsersService, err := u.GetNativeUsersService()
	if err != nil {
		return err
	}

	id, err := u.GetId()
	if err != nil {
		return err
	}

	if id <= 0 {
		return tracederrors.TracedErrorf("Got invalid id: '%d'", id)
	}

	_, _, err = nativeUsersService.ModifyUser(
		id,
		&gitlab.ModifyUserOptions{
			Password: &newPassword,
		},
	)
	if err != nil {
		return err
	}

	if verbose {
		logging.LogChangedf("password for user '%s' on gitlab '%s' updated", username, fqdn)
	}

	return nil
}
