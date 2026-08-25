package nativeminioclient

import (
	"context"

	"github.com/minio/madmin-go/v3"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/storage/s3/s3options"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func ListUsers(ctx context.Context, adminClient *madmin.AdminClient) (map[string]madmin.UserInfo, error) {
	logging.LogInfoByCtxf(ctx, "List minio users started.")

	if adminClient == nil {
		return nil, tracederrors.TracedErrorNil("adminClient")
	}

	users, err := adminClient.ListUsers(ctx)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to list minio users: %w", err)
	}

	logging.LogInfoByCtxf(ctx, "List minio users finished. Found '%d' users.", len(users))

	return users, nil
}

func ListUserNames(ctx context.Context, adminClient *madmin.AdminClient) ([]string, error) {
	if adminClient == nil {
		return nil, tracederrors.TracedErrorNil("adminClient")
	}

	users, err := ListUsers(ctx, adminClient)
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(users))
	for name := range users {
		names = append(names, name)
	}

	return names, nil
}

func UserExists(ctx context.Context, adminClient *madmin.AdminClient, userName string) (bool, error) {
	if adminClient == nil {
		return false, tracederrors.TracedErrorNil("adminClient")
	}

	if userName == "" {
		return false, tracederrors.TracedErrorEmptyString("userName")
	}

	userNames, err := ListUserNames(ctx, adminClient)
	if err != nil {
		return false, err
	}

	exists := false
	for _, name := range userNames {
		if name == userName {
			exists = true
			break
		}
	}

	if exists {
		logging.LogInfoByCtxf(ctx, "User '%s' exists.", userName)
	} else {
		logging.LogInfoByCtxf(ctx, "User '%s' does not exist.", userName)
	}

	return exists, nil
}

func CreateUser(ctx context.Context, adminClient *madmin.AdminClient, userName string, password string, options *s3options.CreateUserOptions) error {
	if adminClient == nil {
		return tracederrors.TracedErrorNil("adminClient")
	}

	if userName == "" {
		return tracederrors.TracedErrorEmptyString("userName")
	}

	if password == "" {
		return tracederrors.TracedErrorEmptyString("password")
	}

	logging.LogInfoByCtxf(ctx, "Create user '%s' started.", userName)

	if options == nil {
		options = &s3options.CreateUserOptions{}
	}

	exists, err := UserExists(ctx, adminClient, userName)
	if err != nil {
		return err
	}

	if exists {
		if options.KeepCurrentPasswordIfUserExists {
			logging.LogInfoByCtxf(ctx, "User '%s' already exists. Keep current password as configured.", userName)
		} else {
			err = adminClient.AddUser(ctx, userName, password)
			if err != nil {
				return tracederrors.TracedErrorf("Failed to update password for user '%s': %w", userName, err)
			}

			logging.LogChangedByCtxf(ctx, "User '%s' already exists. Password updated.", userName)
		}
	} else {
		err = adminClient.AddUser(ctx, userName, password)
		if err != nil {
			return tracederrors.TracedErrorf("Failed to create user '%s': %w", userName, err)
		}

		logging.LogChangedByCtxf(ctx, "User '%s' created.", userName)
	}

	if options.ReadOnly {
		err = adminClient.SetPolicy(ctx, "readonly", userName, false)
		if err != nil {
			return tracederrors.TracedErrorf("Failed to set readonly policy for user '%s': %w", userName, err)
		}

		logging.LogChangedByCtxf(ctx, "Readonly policy set for user '%s'.", userName)
	}

	logging.LogInfoByCtxf(ctx, "Create user '%s' finished.", userName)

	return nil
}

func DeleteUser(ctx context.Context, adminClient *madmin.AdminClient, userName string) error {
	if adminClient == nil {
		return tracederrors.TracedErrorNil("adminClient")
	}

	if userName == "" {
		return tracederrors.TracedErrorEmptyString("userName")
	}

	logging.LogInfoByCtxf(ctx, "Delete user '%s' started.", userName)

	exists, err := UserExists(ctx, adminClient, userName)
	if err != nil {
		return err
	}

	if exists {
		err = adminClient.RemoveUser(ctx, userName)
		if err != nil {
			return tracederrors.TracedErrorf("Failed to delete user '%s': %w", userName, err)
		}

		logging.LogChangedByCtxf(ctx, "Deleted user '%s'.", userName)
	} else {
		logging.LogInfoByCtxf(ctx, "User '%s' is already absent. Skip delete.", userName)
	}

	logging.LogInfoByCtxf(ctx, "Delete user '%s' finished.", userName)

	return nil
}
