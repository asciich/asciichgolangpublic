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

func Test_Example_CreateAndDelete_Objects_test(t *testing.T) {
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
	err = nativeminioclient.CreateBucket(ctx, client, bucketName, &s3options.CreateBucketOptions{})
	require.NoError(t, err)

	exists, err = nativeminioclient.BucketExists(ctx, client, bucketName)
	require.NoError(t, err)
	require.True(t, exists)

	objectList, err := nativeminioclient.ListObjectNames(ctx, client, bucketName)
	require.NoError(t, err)
	require.Len(t, objectList, 0)

	// Our test object does not exist yet:
	objectKey := "test.txt"
	exists, err = nativeminioclient.ObjectExists(ctx, client, bucketName, objectKey)
	require.NoError(t, err)
	require.False(t, exists)

	// Create the object:
	ctxCreate := contextutils.WithChangeIndicator(ctx)
	err = nativeminioclient.CreateObjectFromString(ctxCreate, client, bucketName, objectKey, "hello world")
	require.NoError(t, err)
	require.True(t, contextutils.IsChanged(ctxCreate)) // Creating the object is indicated by a change.

	exists, err = nativeminioclient.ObjectExists(ctx, client, bucketName, objectKey)
	require.NoError(t, err)
	require.True(t, exists)

	isEqual, err := nativeminioclient.IsObjectContentEqualString(ctx, client, bucketName, objectKey, "hello world")
	require.NoError(t, err)
	require.True(t, isEqual)

	content, err := nativeminioclient.GetObjectContentAsString(ctx, client, bucketName, objectKey)
	require.NoError(t, err)
	require.EqualValues(t, "hello world", content)

	// Create the object with the same content again:
	// This will not change the object as it's already correct in the S3 storage:
	ctxCreate = contextutils.WithChangeIndicator(ctx)
	err = nativeminioclient.CreateObjectFromString(ctxCreate, client, bucketName, objectKey, "hello world")
	require.NoError(t, err)
	require.False(t, contextutils.IsChanged(ctxCreate)) // Recreating the object with the same content is skipped and will not create a change.

	exists, err = nativeminioclient.ObjectExists(ctx, client, bucketName, objectKey)
	require.NoError(t, err)
	require.True(t, exists)

	content, err = nativeminioclient.GetObjectContentAsString(ctx, client, bucketName, objectKey)
	require.NoError(t, err)
	require.EqualValues(t, "hello world", content)

	// Create the object with another content again:
	// This will agin change the object as it's content differs now:
	ctxCreate = contextutils.WithChangeIndicator(ctx)
	err = nativeminioclient.CreateObjectFromString(ctxCreate, client, bucketName, objectKey, "another content")
	require.NoError(t, err)
	require.True(t, contextutils.IsChanged(ctxCreate)) // Recreating the object with differenct content results in an indicated change.

	exists, err = nativeminioclient.ObjectExists(ctx, client, bucketName, objectKey)
	require.NoError(t, err)
	require.True(t, exists)

	content, err = nativeminioclient.GetObjectContentAsString(ctx, client, bucketName, objectKey)
	require.NoError(t, err)
	require.EqualValues(t, "another content", content)

	// Delete the object:
	ctxDelete := contextutils.WithChangeIndicator(ctx)
	err = nativeminioclient.DeleteObject(ctxDelete, client, bucketName, objectKey)
	require.NoError(t, err)
	require.True(t, contextutils.IsChanged(ctxDelete))

	exists, err = nativeminioclient.ObjectExists(ctx, client, bucketName, objectKey)
	require.NoError(t, err)
	require.False(t, exists)

	// Delete the object again.
	// This will not indicate a change as the object was already absent.
	ctxDelete = contextutils.WithChangeIndicator(ctx)
	err = nativeminioclient.DeleteObject(ctxDelete, client, bucketName, objectKey)
	require.NoError(t, err)
	require.False(t, contextutils.IsChanged(ctxDelete))

	exists, err = nativeminioclient.ObjectExists(ctx, client, bucketName, objectKey)
	require.NoError(t, err)
	require.False(t, exists)
}
