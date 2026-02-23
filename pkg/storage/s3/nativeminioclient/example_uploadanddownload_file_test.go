package nativeminioclient_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/dockeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/nativedocker"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
	"github.com/asciich/asciichgolangpublic/pkg/randomgenerator"
	"github.com/asciich/asciichgolangpublic/pkg/storage/s3/nativeminioclient"
	"github.com/asciich/asciichgolangpublic/pkg/storage/s3/s3options"
)

func Test_Example_UploadAndDownload_File_test(t *testing.T) {
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

	// Create the local file:
	srcFilePath, err := tempfiles.CreateTemporaryFileFromContentString(ctx, "This is the test data")
	require.NoError(t, err)
	defer nativefiles.Delete(ctx, srcFilePath, &filesoptions.DeleteOptions{})

	// Upload the local file as "example.txt" into the S3 bucket:
	ctxUpload := contextutils.WithChangeIndicator(ctx)
	objectKey := "example.txt"
	err = nativeminioclient.UploadFileByPath(ctxUpload, client, bucketName, objectKey, srcFilePath)
	require.NoError(t, err)
	require.True(t, contextutils.IsChanged(ctxUpload)) // Uploading the file is considered a change.

	// Upload the same file again will be skipped since the content is already up to date:
	ctxUpload = contextutils.WithChangeIndicator(ctx)
	err = nativeminioclient.UploadFileByPath(ctxUpload, client, bucketName, objectKey, srcFilePath)
	require.NoError(t, err)
	require.False(t, contextutils.IsChanged(ctxUpload))

	// Downlaod the file directly as a temporary file:
	ctxDownload := contextutils.WithChangeIndicator(ctx)
	destFilePath, err := nativeminioclient.DownloadAsTemporaryFile(ctxDownload, client, bucketName, objectKey)
	require.NoError(t, err)
	defer nativefiles.Delete(ctx, destFilePath, &filesoptions.DeleteOptions{})
	require.True(t, contextutils.IsChanged(ctxDownload))

	downloadedContent, err := nativefiles.ReadAsString(ctx, destFilePath, &filesoptions.ReadOptions{})
	require.NoError(t, err)
	require.EqualValues(t, "This is the test data", downloadedContent)

	// Redownload it again. This time we spedify the output path.
	ctxDownload = contextutils.WithChangeIndicator(ctx)
	err = nativeminioclient.DownloadAsFileByPath(ctxDownload, client, bucketName, objectKey, destFilePath)
	require.NoError(t, err)
	require.False(t, contextutils.IsChanged(ctxDownload)) // As it's the same file and the same content it will not considered as a change.
}
