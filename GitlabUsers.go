package asciichgolangpublic

import (
	"strings"

	"github.com/xanzy/go-gitlab"
)

type GitlabUsers struct {
	gitlab *GitlabInstance
}

func NewGitlabUsers() (gitlabUsers *GitlabUsers) {
	return new(GitlabUsers)
}

func (g *GitlabUsers) MustCreateAccessToken(options *GitlabCreateAccessTokenOptions) (newToken string) {
	newToken, err := g.CreateAccessToken(options)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return newToken
}

func (g *GitlabUsers) MustCreateUser(createUserOptions *GitlabCreateUserOptions) (createdUser *GitlabUser) {
	createdUser, err := g.CreateUser(createUserOptions)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return createdUser
}

func (g *GitlabUsers) MustGetFqdn() (fqdn string) {
	fqdn, err := g.GetFqdn()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return fqdn
}

func (g *GitlabUsers) MustGetGitlab() (gitlab *GitlabInstance) {
	gitlab, err := g.GetGitlab()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlab
}

func (g *GitlabUsers) MustGetNativeGitlabClient() (nativeGitlabClient *gitlab.Client) {
	nativeGitlabClient, err := g.GetNativeGitlabClient()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return nativeGitlabClient
}

func (g *GitlabUsers) MustGetNativeUsersService() (nativeUsersService *gitlab.UsersService) {
	nativeUsersService, err := g.GetNativeUsersService()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return nativeUsersService
}

func (g *GitlabUsers) MustGetUserByUsername(username string) (gitlabUser *GitlabUser) {
	gitlabUser, err := g.GetUserByUsername(username)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlabUser
}

func (g *GitlabUsers) MustGetUserId() (userId int) {
	userId, err := g.GetUserId()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return userId
}

func (g *GitlabUsers) MustGetUserNames() (userNames []string) {
	userNames, err := g.GetUserNames()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return userNames
}

func (g *GitlabUsers) MustGetUsers() (users []*GitlabUser) {
	users, err := g.GetUsers()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return users
}

func (g *GitlabUsers) MustSetGitlab(gitlab *GitlabInstance) {
	err := g.SetGitlab(gitlab)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabUsers) MustUserByUserNameExists(username string) (userExists bool) {
	userExists, err := g.UserByUserNameExists(username)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return userExists
}

func (u *GitlabUsers) CreateAccessToken(options *GitlabCreateAccessTokenOptions) (newToken string, err error) {
	if options == nil {
		return "", TracedError("options is nil")
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
		return nil, TracedError("createUserOptions is nil")
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
			LogInfof(
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
				LogInfof("Going to use a random %d digit password for user '%s' since no password was specified.", nDigits, createUserOptions.Name)
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
			LogChangedf("User '%s' created on gitlab '%s'.", createUserOptions.Username, fqdn)
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
		return nil, TracedError("gitlab not set")
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
		return nil, TracedError("nativeUsersService was returned as nil pointer")
	}

	return nativeUsersService, nil
}

func (u *GitlabUsers) GetUserByUsername(username string) (gitlabUser *GitlabUser, err error) {
	username = strings.TrimSpace(username)

	if len(username) <= 0 {
		return nil, TracedError("username is empty string")
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

	return nil, TracedErrorf("User '%s' not found on gitlab '%s'", username, fqdn)
}

func (u *GitlabUsers) GetUserId() (userId int, err error) {
	usersService, err := u.GetNativeUsersService()
	if err != nil {
		return -1, err
	}

	nativeUser, _, err := usersService.CurrentUser()
	if err != nil {
		return -1, TracedError(err.Error())
	}

	userId = nativeUser.ID
	if userId <= 0 {
		return -1, TracedErrorf("Got invalid user id for current user: '%d'", userId)
	}

	return userId, nil
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

func (u *GitlabUsers) GetUsers() (users []*GitlabUser, err error) {
	nativeUsersService, err := u.GetNativeUsersService()
	if err != nil {
		return nil, err
	}

	nativeUsers, _, err := nativeUsersService.ListUsers(&gitlab.ListUsersOptions{})
	if err != nil {
		return nil, err
	}

	gitlab, err := u.GetGitlab()
	if err != nil {
		return nil, err
	}

	users = []*GitlabUser{}
	for _, nativeUser := range nativeUsers {
		userToAdd := NewGitlabUser()
		err = userToAdd.SetGitlab(gitlab)
		if err != nil {
			return nil, err
		}

		userId := nativeUser.ID
		userName := nativeUser.Name
		userEmail := nativeUser.Email
		userUsernamme := nativeUser.Username

		err = userToAdd.SetId(userId)
		if err != nil {
			return nil, err
		}

		err = userToAdd.SetCachedName(userName)
		if err != nil {
			return nil, err
		}

		if len(userEmail) > 0 {
			err = userToAdd.SetCachedEmail(userEmail)
			if err != nil {
				return nil, err
			}
		}

		err = userToAdd.SetCachedUsername(userUsernamme)
		if err != nil {
			return nil, err
		}

		users = append(users, userToAdd)
	}

	return users, nil
}

func (u *GitlabUsers) SetGitlab(gitlab *GitlabInstance) (err error) {
	if gitlab == nil {
		return TracedError("gitlab is nil")
	}

	u.gitlab = gitlab

	return nil
}

func (u *GitlabUsers) UserByUserNameExists(username string) (userExists bool, err error) {
	username = strings.TrimSpace(username)

	if len(username) <= 0 {
		return false, TracedError("username is empty string")
	}

	userNameList, err := u.GetUserNames()
	if err != nil {
		return false, err
	}

	userExists = Slices().ContainsString(userNameList, username)
	return userExists, nil
}
