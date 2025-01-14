package asciichgolangpublic

import (
	"strings"

	aslices "github.com/asciich/asciichgolangpublic/datatypes/slices"
	"github.com/asciich/asciichgolangpublic/errors"
	"github.com/asciich/asciichgolangpublic/logging"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type GitlabUsers struct {
	gitlab *GitlabInstance
}

func NewGitlabUsers() (gitlabUsers *GitlabUsers) {
	return new(GitlabUsers)
}

// Return the currently logged in user
func (g *GitlabUsers) GetUser() (gitlabUser *GitlabUser, err error) {
	id, err := g.GetUserId()
	if err != nil {
		return nil, err
	}

	gitlabUser, err = g.GetUserById(id)
	if err != nil {
		return nil, err
	}

	return gitlabUser, nil
}

// Returns the `userId` of the currently logged in user.
func (u *GitlabUsers) GetUserId() (userId int, err error) {
	usersService, err := u.GetNativeUsersService()
	if err != nil {
		return -1, err
	}

	nativeUser, _, err := usersService.CurrentUser()
	if err != nil {
		return -1, errors.TracedError(err.Error())
	}

	userId = nativeUser.ID
	if userId <= 0 {
		return -1, errors.TracedErrorf("Got invalid user id for current user: '%d'", userId)
	}

	return userId, nil
}

func (g *GitlabUsers) GetUserById(id int) (gitlabUser *GitlabUser, err error) {
	if id <= 0 {
		return nil, errors.TracedErrorf("id '%d' is invalid", id)
	}

	nativeClient, err := g.GetNativeUsersService()
	if err != nil {
		return nil, err
	}

	nativeUser, _, err := nativeClient.GetUser(id, gitlab.GetUsersOptions{})
	if err != nil {
		return nil, errors.TracedErrorf("Getting user with id '%d' failed: '%w'", id, err)
	}

	if nativeUser == nil {
		return nil, errors.TracedErrorf("nativeUser is nil after evaluation")
	}

	gitlabUser, err = g.GetUserByNativeGitlabUser(nativeUser)
	if err != nil {
		return nil, err
	}

	return gitlabUser, nil
}

func (g *GitlabUsers) GetUserByNativeGitlabUser(nativeUser *gitlab.User) (user *GitlabUser, err error) {
	if nativeUser == nil {
		return nil, errors.TracedErrorNil("nativeUser")
	}

	gitlab, err := g.GetGitlab()
	if err != nil {
		return nil, err
	}

	user = NewGitlabUser()
	err = user.SetGitlab(gitlab)
	if err != nil {
		return nil, err
	}

	userId := nativeUser.ID
	userName := nativeUser.Name
	userEmail := nativeUser.Email
	userUsernamme := nativeUser.Username

	err = user.SetId(userId)
	if err != nil {
		return nil, err
	}

	err = user.SetCachedName(userName)
	if err != nil {
		return nil, err
	}

	if len(userEmail) > 0 {
		err = user.SetCachedEmail(userEmail)
		if err != nil {
			return nil, err
		}
	}

	err = user.SetCachedUsername(userUsernamme)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (g *GitlabUsers) GetUsers() (users []*GitlabUser, err error) {
	nativeUsersService, err := g.GetNativeUsersService()
	if err != nil {
		return nil, err
	}

	nativeUsers, _, err := nativeUsersService.ListUsers(&gitlab.ListUsersOptions{})
	if err != nil {
		return nil, err
	}

	users = []*GitlabUser{}
	for _, nativeUser := range nativeUsers {
		userToAdd, err := g.GetUserByNativeGitlabUser(nativeUser)
		if err != nil {
			return nil, err
		}

		users = append(users, userToAdd)
	}

	return users, nil
}

func (g *GitlabUsers) MustCreateAccessToken(options *GitlabCreateAccessTokenOptions) (newToken string) {
	newToken, err := g.CreateAccessToken(options)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return newToken
}

func (g *GitlabUsers) MustCreateUser(createUserOptions *GitlabCreateUserOptions) (createdUser *GitlabUser) {
	createdUser, err := g.CreateUser(createUserOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return createdUser
}

func (g *GitlabUsers) MustGetFqdn() (fqdn string) {
	fqdn, err := g.GetFqdn()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return fqdn
}

func (g *GitlabUsers) MustGetGitlab() (gitlab *GitlabInstance) {
	gitlab, err := g.GetGitlab()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlab
}

func (g *GitlabUsers) MustGetNativeGitlabClient() (nativeGitlabClient *gitlab.Client) {
	nativeGitlabClient, err := g.GetNativeGitlabClient()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return nativeGitlabClient
}

func (g *GitlabUsers) MustGetNativeUsersService() (nativeUsersService *gitlab.UsersService) {
	nativeUsersService, err := g.GetNativeUsersService()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return nativeUsersService
}

func (g *GitlabUsers) MustGetUser() (gitlabUser *GitlabUser) {
	gitlabUser, err := g.GetUser()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlabUser
}

func (g *GitlabUsers) MustGetUserById(id int) (gitlabUser *GitlabUser) {
	gitlabUser, err := g.GetUserById(id)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlabUser
}

func (g *GitlabUsers) MustGetUserByNativeGitlabUser(nativeUser *gitlab.User) (user *GitlabUser) {
	user, err := g.GetUserByNativeGitlabUser(nativeUser)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return user
}

func (g *GitlabUsers) MustGetUserByUsername(username string) (gitlabUser *GitlabUser) {
	gitlabUser, err := g.GetUserByUsername(username)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlabUser
}

func (g *GitlabUsers) MustGetUserId() (userId int) {
	userId, err := g.GetUserId()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return userId
}

func (g *GitlabUsers) MustGetUserNames() (userNames []string) {
	userNames, err := g.GetUserNames()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return userNames
}

func (g *GitlabUsers) MustGetUsers() (users []*GitlabUser) {
	users, err := g.GetUsers()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return users
}

func (g *GitlabUsers) MustSetGitlab(gitlab *GitlabInstance) {
	err := g.SetGitlab(gitlab)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabUsers) MustUserByUserNameExists(username string) (userExists bool) {
	userExists, err := g.UserByUserNameExists(username)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return userExists
}

func (u *GitlabUsers) CreateAccessToken(options *GitlabCreateAccessTokenOptions) (newToken string, err error) {
	if options == nil {
		return "", errors.TracedError("options is nil")
	}

	username, err := options.GetUserName()
	if err != nil {
		return "", err
	}

	user, err := u.GetUserByUsername(username)
	if err != nil {
		return "", err
	}

	newToken, err = user.CreateAccessToken(options)
	if err != nil {
		return "", err
	}

	return newToken, nil
}

func (u *GitlabUsers) CreateUser(createUserOptions *GitlabCreateUserOptions) (createdUser *GitlabUser, err error) {
	if createUserOptions == nil {
		return nil, errors.TracedError("createUserOptions is nil")
	}

	nativeUsersService, err := u.GetNativeUsersService()
	if err != nil {
		return nil, err
	}

	fqdn, err := u.GetFqdn()
	if err != nil {
		return nil, err
	}

	userExists, err := u.UserByUserNameExists(createUserOptions.Username)
	if err != nil {
		return nil, err
	}

	if userExists {
		if createUserOptions.Verbose {
			logging.LogInfof(
				"User with username='%s' already exists on gitlab '%s'.",
				createUserOptions.Username,
				fqdn,
			)
		}
		createdUser, err = u.GetUserByUsername(createUserOptions.Username)
		if err != nil {
			return nil, err
		}
	} else {
		passwordToSet := createUserOptions.Password
		if len(passwordToSet) <= 0 {
			nDigits := 15
			if createUserOptions.Verbose {
				logging.LogInfof("Going to use a random %d digit password for user '%s' since no password was specified.", nDigits, createUserOptions.Name)
			}

			passwordToSet, err = RandomGenerator().GetRandomString(nDigits)
			if err != nil {
				return nil, err
			}
		}

		skipConfirmation := true

		nativeUser, _, err := nativeUsersService.CreateUser(
			&gitlab.CreateUserOptions{
				Name:             &createUserOptions.Name,
				Email:            &createUserOptions.Email,
				Username:         &createUserOptions.Username,
				Password:         &passwordToSet,
				SkipConfirmation: &skipConfirmation,
			},
		)
		if err != nil {
			return nil, err
		}

		gitlab, err := u.GetGitlab()
		if err != nil {
			return nil, err
		}

		createdUser = NewGitlabUser()
		err = createdUser.SetGitlab(gitlab)
		if err != nil {
			return nil, err
		}

		err = createdUser.SetId(nativeUser.ID)
		if err != nil {
			return nil, err
		}

		if createUserOptions.Verbose {
			logging.LogChangedf("User '%s' created on gitlab '%s'.", createUserOptions.Username, fqdn)
		}

	}

	return createdUser, nil
}

func (u *GitlabUsers) GetFqdn() (fqdn string, err error) {
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

func (u *GitlabUsers) GetGitlab() (gitlab *GitlabInstance, err error) {
	if u.gitlab == nil {
		return nil, errors.TracedError("gitlab not set")
	}

	return u.gitlab, nil
}

func (u *GitlabUsers) GetNativeGitlabClient() (nativeGitlabClient *gitlab.Client, err error) {
	gitlab, err := u.GetGitlab()
	if err != nil {
		return nil, err
	}

	nativeGitlabClient, err = gitlab.GetNativeClient()
	if err != nil {
		return nil, err
	}

	return nativeGitlabClient, nil
}

func (u *GitlabUsers) GetNativeUsersService() (nativeUsersService *gitlab.UsersService, err error) {
	nativeGitlabClient, err := u.GetNativeGitlabClient()
	if err != nil {
		return nil, err
	}

	nativeUsersService = nativeGitlabClient.Users
	if nativeUsersService == nil {
		return nil, errors.TracedError("nativeUsersService was returned as nil pointer")
	}

	return nativeUsersService, nil
}

func (u *GitlabUsers) GetUserByUsername(username string) (gitlabUser *GitlabUser, err error) {
	username = strings.TrimSpace(username)

	if len(username) <= 0 {
		return nil, errors.TracedError("username is empty string")
	}

	fqdn, err := u.GetFqdn()
	if err != nil {
		return nil, err
	}

	users, err := u.GetUsers()
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		currentUsername, err := user.GetChachedUsername()
		if err != nil {
			return nil, err
		}

		if currentUsername == username {
			return user, nil
		}
	}

	return nil, errors.TracedErrorf("User '%s' not found on gitlab '%s'", username, fqdn)
}

func (u *GitlabUsers) GetUserNames() (userNames []string, err error) {
	users, err := u.GetUsers()
	if err != nil {
		return nil, err
	}

	userNames = []string{}
	for _, user := range users {
		nameToAdd, err := user.GetChachedUsername()
		if err != nil {
			return nil, err
		}

		userNames = append(userNames, nameToAdd)
	}

	return userNames, nil
}

func (u *GitlabUsers) SetGitlab(gitlab *GitlabInstance) (err error) {
	if gitlab == nil {
		return errors.TracedError("gitlab is nil")
	}

	u.gitlab = gitlab

	return nil
}

func (u *GitlabUsers) UserByUserNameExists(username string) (userExists bool, err error) {
	username = strings.TrimSpace(username)

	if len(username) <= 0 {
		return false, errors.TracedError("username is empty string")
	}

	userNameList, err := u.GetUserNames()
	if err != nil {
		return false, err
	}

	userExists = aslices.ContainsString(userNameList, username)
	return userExists, nil
}
