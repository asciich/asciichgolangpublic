package nativeminioclient_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/dockeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/nativedocker"
	"github.com/asciich/asciichgolangpublic/pkg/randomgenerator"
	"github.com/asciich/asciichgolangpublic/pkg/storage/s3/nativeminioclient"
	"github.com/asciich/asciichgolangpublic/pkg/storage/s3/s3options"
)

func Test_Example_ListObjects_test(t *testing.T) {
	// enable verbose output
	ctx := contextutils.ContextVerbose()

	// Define admin credentials for the test environment
	const minioAdminUser = "minioadmin"
	minioAdminPassword, err := randomgenerator.GetRandomString(10)
	require.NoError(t, err)

	// Run minio in a docker container for testing.
	const containerName = "test-nativeminioclient"
	err = nativedocker.RemoveContainer(ctx, containerName, &dockeroptions.RemoveOptions{Force: true})
	require.NoError(t, err)

	_, err = nativedocker.RunContainer(ctx, &dockeroptions.DockerRunContainerOptions{
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

	time.Sleep(time.Second * 2)
	//defer container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})

	// Define the bucket name used for this test:
	const bucketName = "test-bucket"

	// Get the minio client:
	client, err := nativeminioclient.NewClient("localhost:9000", minioAdminUser, minioAdminPassword, &s3options.NewS3ClientOptions{})
	require.NoError(t, err)

	// Delete the bucket to ensure a clear defined test setup:
	err = nativeminioclient.DeleteBucket(ctx, client, bucketName)
	require.NoError(t, err)

	exists, err := nativeminioclient.BucketExists(ctx, client, bucketName)
	require.NoError(t, err)
	require.False(t, exists)

	// Create the bucket which is now empty:
	err = nativeminioclient.CreateBucket(ctx, client, bucketName)
	require.NoError(t, err)

	exists, err = nativeminioclient.BucketExists(ctx, client, bucketName)
	require.NoError(t, err)
	require.True(t, exists)

	objectList, err := nativeminioclient.ListObjectNames(ctx, client, bucketName)
	require.NoError(t, err)
	require.Len(t, objectList, 0)

	// Generate one object:
	err = nativeminioclient.CreateObjectFromString(ctx, client, bucketName, "test.txt", "hello world")
	require.NoError(t, err)

	objectList, err = nativeminioclient.ListObjectNames(ctx, client, bucketName)
	require.NoError(t, err)
	require.EqualValues(t, objectList, []string{"test.txt"})

	// Generate a second object:
	err = nativeminioclient.CreateObjectFromString(ctx, client, bucketName, "abc.txt", "hello world")
	require.NoError(t, err)

	objectList, err = nativeminioclient.ListObjectNames(ctx, client, bucketName)
	require.NoError(t, err)
	require.EqualValues(t, objectList, []string{"abc.txt", "test.txt"})

	// Generate a second object again:
	err = nativeminioclient.CreateObjectFromString(ctx, client, bucketName, "abc.txt", "hello world2")
	require.NoError(t, err)

	objectList, err = nativeminioclient.ListObjectNames(ctx, client, bucketName)
	require.NoError(t, err)
	require.EqualValues(t, objectList, []string{"abc.txt", "test.txt"})

	// Generate a third object:
	err = nativeminioclient.CreateObjectFromString(ctx, client, bucketName, "aaa.txt", "hello world")
	require.NoError(t, err)

	objectList, err = nativeminioclient.ListObjectNames(ctx, client, bucketName)
	require.NoError(t, err)
	require.EqualValues(t, objectList, []string{"aaa.txt", "abc.txt", "test.txt"})
}
