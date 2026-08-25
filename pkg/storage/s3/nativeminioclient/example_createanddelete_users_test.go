package nativeminioclient_test

import (
	"testing"
	"time"

	"github.com/minio/madmin-go/v3"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/nativedocker"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/randomgenerator"
	"github.com/asciich/asciichgolangpublic/pkg/storage/s3/nativeminioclient"
	"github.com/asciich/asciichgolangpublic/pkg/storage/s3/s3options"
)

func Test_Example_CreateAndDeleteUsers_test(t *testing.T) {
	// Enable verbose output:
	ctx := contextutils.ContextVerbose()

	// Define admin credentials for the test environment
	const minioAdminUser = "minioadmin"
	minioAdminPassword, err := randomgenerator.GetRandomString(10)
	require.NoError(t, err)

	// Run minio in a docker container for testing.
	const containerName = "test-nativeminioclient-users"
	err = nativedocker.RemoveContainer(ctx, containerName, &dockeroptions.RemoveOptions{Force: true})
	require.NoError(t, err)

	container, err := nativedocker.RunContainer(ctx, &dockeroptions.DockerRunContainerOptions{
		Name:      containerName,
		ImageName: "quay.io/minio/minio",
		Command:   []string{"server", "/data", "--console-address", ":9001"},
		Ports:     []string{"9000:9000"},
		AdditionalEnvVars: map[string]string{
			"MINIO_ROOT_USER":     minioAdminUser,
			"MINIO_ROOT_PASSWORD": minioAdminPassword,
		},
		WaitForPortsOpen: true,
	})
	require.NoError(t, err)
	defer container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})

	time.Sleep(time.Second * 2)

	// Define the user credentials used for this test:
	const testUserName = "test-user"
	testUserPassword, err := randomgenerator.GetRandomString(12)
	require.NoError(t, err)

	// Get the minio admin client:
	adminClient, err := madmin.New("localhost:9000", minioAdminUser, minioAdminPassword, false)
	require.NoError(t, err)

	// Delete the user to ensure a clear defined test setup:
	err = nativeminioclient.DeleteUser(ctx, adminClient, testUserName)
	require.NoError(t, err)

	exists, err := nativeminioclient.UserExists(ctx, adminClient, testUserName)
	require.NoError(t, err)
	require.False(t, exists)

	userNames, err := nativeminioclient.ListUserNames(ctx, adminClient)
	require.NoError(t, err)
	require.NotContains(t, userNames, testUserName)

	// Create the user:
	err = nativeminioclient.CreateUser(ctx, adminClient, testUserName, testUserPassword, &s3options.CreateUserOptions{})
	require.NoError(t, err)

	exists, err = nativeminioclient.UserExists(ctx, adminClient, testUserName)
	require.NoError(t, err)
	require.True(t, exists)

	userNames, err = nativeminioclient.ListUserNames(ctx, adminClient)
	require.NoError(t, err)
	require.Contains(t, userNames, testUserName)

	// Create the user again to check idempotence (password is updated):
	newPassword, err := randomgenerator.GetRandomString(12)
	require.NoError(t, err)

	err = nativeminioclient.CreateUser(ctx, adminClient, testUserName, newPassword, &s3options.CreateUserOptions{})
	require.NoError(t, err)

	exists, err = nativeminioclient.UserExists(ctx, adminClient, testUserName)
	require.NoError(t, err)
	require.True(t, exists)

	// Verify the new password works by creating a minio client with the new credentials:
	_, err = minio.New("localhost:9000", &minio.Options{
		Creds:  credentials.NewStaticV4(testUserName, newPassword, ""),
		Secure: false,
	})
	require.NoError(t, err)

	// Create the user again with KeepCurrentPasswordIfUserExists set to true:
	anotherPassword, err := randomgenerator.GetRandomString(12)
	require.NoError(t, err)

	err = nativeminioclient.CreateUser(ctx, adminClient, testUserName, anotherPassword, &s3options.CreateUserOptions{
		KeepCurrentPasswordIfUserExists: true,
	})
	require.NoError(t, err)

	// Verify the password was NOT updated (the previous newPassword still works):
	_, err = minio.New("localhost:9000", &minio.Options{
		Creds:  credentials.NewStaticV4(testUserName, newPassword, ""),
		Secure: false,
	})
	require.NoError(t, err)

	// Delete the user:
	err = nativeminioclient.DeleteUser(ctx, adminClient, testUserName)
	require.NoError(t, err)

	exists, err = nativeminioclient.UserExists(ctx, adminClient, testUserName)
	require.NoError(t, err)
	require.False(t, exists)

	userNames, err = nativeminioclient.ListUserNames(ctx, adminClient)
	require.NoError(t, err)
	require.NotContains(t, userNames, testUserName)

	// Delete the user again to check idempotence:
	err = nativeminioclient.DeleteUser(ctx, adminClient, testUserName)
	require.NoError(t, err)

	exists, err = nativeminioclient.UserExists(ctx, adminClient, testUserName)
	require.NoError(t, err)
	require.False(t, exists)

	userNames, err = nativeminioclient.ListUserNames(ctx, adminClient)
	require.NoError(t, err)
	require.NotContains(t, userNames, testUserName)
}
