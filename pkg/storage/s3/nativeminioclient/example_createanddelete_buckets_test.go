package nativeminioclient_test

import (
	"context"
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

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func Test_Example_CreateAndDeleteBuckets_test(t *testing.T) {
	ctx := getCtx()

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
	const bucketName = "test-create-bucket"

	// Get the minio client:
	client, err := nativeminioclient.NewClient("localhost:9000", minioAdminUser, minioAdminPassword, &s3options.NewS3ClientOptions{})
	require.NoError(t, err)

	// Delete the bucket to ensure a clear defined test setup:
	err = nativeminioclient.DeleteBucket(ctx, client, bucketName)
	require.NoError(t, err)

	exists, err := nativeminioclient.BucketExists(ctx, client, bucketName)
	require.NoError(t, err)
	require.False(t, exists)

	bucketNames, err := nativeminioclient.ListBucketNames(ctx, client)
	require.NoError(t, err)
	require.NotContains(t, bucketNames, bucketName)

	// Create the bucket:
	err = nativeminioclient.CreateBucket(ctx, client, bucketName)
	require.NoError(t, err)

	exists, err = nativeminioclient.BucketExists(ctx, client, bucketName)
	require.NoError(t, err)
	require.True(t, exists)

	bucketNames, err = nativeminioclient.ListBucketNames(ctx, client)
	require.NoError(t, err)
	require.Contains(t, bucketNames, bucketName)

	// Create the bucket again to check idempotence:
	err = nativeminioclient.CreateBucket(ctx, client, bucketName)
	require.NoError(t, err)

	exists, err = nativeminioclient.BucketExists(ctx, client, bucketName)
	require.NoError(t, err)
	require.True(t, exists)

	bucketNames, err = nativeminioclient.ListBucketNames(ctx, client)
	require.NoError(t, err)
	require.Contains(t, bucketNames, bucketName)

	// Delete the bucket:
	err = nativeminioclient.DeleteBucket(ctx, client, bucketName)
	require.NoError(t, err)

	exists, err = nativeminioclient.BucketExists(ctx, client, bucketName)
	require.NoError(t, err)
	require.False(t, exists)

	bucketNames, err = nativeminioclient.ListBucketNames(ctx, client)
	require.NoError(t, err)
	require.NotContains(t, bucketNames, bucketName)

	// Delete the bucket again to check idempotence:
	err = nativeminioclient.DeleteBucket(ctx, client, bucketName)
	require.NoError(t, err)

	exists, err = nativeminioclient.BucketExists(ctx, client, bucketName)
	require.NoError(t, err)
	require.False(t, exists)

	bucketNames, err = nativeminioclient.ListBucketNames(ctx, client)
	require.NoError(t, err)
	require.NotContains(t, bucketNames, bucketName)
}
