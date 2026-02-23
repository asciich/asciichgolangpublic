package nativeminioclient

import (
	"bytes"
	"context"
	"io"
	"sort"

	"github.com/minio/minio-go/v7"
	"github.com/asciich/asciichgolangpublic/pkg/checksumutils"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func ObjectExists(ctx context.Context, client *minio.Client, bucketName string, objectKey string) (bool, error) {
	if client == nil {
		return false, tracederrors.TracedErrorNil("client")
	}

	if bucketName == "" {
		return false, tracederrors.TracedErrorEmptyString("bucketName")
	}

	if objectKey == "" {
		return false, tracederrors.TracedErrorEmptyString("objectKey")
	}

	var exists = true
	_, err := client.StatObject(ctx, bucketName, objectKey, minio.GetObjectOptions{})
	if err != nil {
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code == "NoSuchKey" {
			exists = false
		} else {
			return false, tracederrors.TracedErrorf("Failed to stat object: %w", err)
		}
	}

	if exists {
		logging.LogInfoByCtxf(ctx, "S3 object '%s' exists in bucket '%s'.", objectKey, bucketName)
	} else {
		logging.LogInfoByCtxf(ctx, "S3 object '%s' does not exsist in bucket '%s'.", objectKey, bucketName)
	}

	return exists, nil
}

func ListObjectNames(ctx context.Context, client *minio.Client, bucketName string) ([]string, error) {
	if client == nil {
		return nil, tracederrors.TracedErrorNil("client")
	}

	if bucketName == "" {
		return nil, tracederrors.TracedErrorEmptyString("bucketName")
	}

	logging.LogInfoByCtxf(ctx, "List objects in bucket '%s' started.", bucketName)

	objectChannel := client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Recursive: true,
	})

	objectNames := []string{}
	for object := range objectChannel {
		if object.Err != nil {
			return nil, tracederrors.TracedErrorf("There is an object.Err while listing the object in the bucket '%s': %w", bucketName, object.Err)
		}
		objectNames = append(objectNames, object.Key)
	}

	sort.Strings(objectNames)

	logging.LogInfoByCtxf(ctx, "List objects in bucket '%s' finished. There are '%d' objects in the bucket.", bucketName, len(objectNames))

	return objectNames, nil
}

func CreateObjectFromString(ctx context.Context, client *minio.Client, bucketName string, objectKey string, objectContent string) error {
	if client == nil {
		return tracederrors.TracedErrorNil("client")
	}

	if bucketName == "" {
		return tracederrors.TracedErrorEmptyString("bucketName")
	}

	if objectKey == "" {
		return tracederrors.TracedErrorEmptyString("objectKey")
	}

	return CreateObjectFromBytes(ctx, client, bucketName, objectKey, []byte(objectContent))
}

func CreateObjectFromBytes(ctx context.Context, client *minio.Client, bucketName string, objectKey string, objectContent []byte) error {
	if client == nil {
		return tracederrors.TracedErrorNil("client")
	}

	if bucketName == "" {
		return tracederrors.TracedErrorEmptyString("bucketName")
	}

	if objectKey == "" {
		return tracederrors.TracedErrorEmptyString("objectKey")
	}

	if objectContent == nil {
		return tracederrors.TracedErrorNil("objectContent")
	}

	logging.LogInfoByCtxf(ctx, "Create S3 object '%s' in bucket '%s' from bytes started.", objectKey, bucketName)

	isEqual, err := IsObjectContentEqualBytes(ctx, client, bucketName, objectKey, objectContent)
	if err != nil {
		return err
	}

	if isEqual {
		logging.LogInfoByCtxf(ctx, "Content of '%s' in bucket '%s' is already up to date. Skip reuploading.", objectKey, bucketName)
	} else {
		reader := bytes.NewReader(objectContent)

		_, err := client.PutObject(ctx, bucketName, objectKey, reader, int64(len(objectContent)), minio.PutObjectOptions{})
		if err != nil {
			return tracederrors.TracedErrorf("Failed to create object from bytes: %w", err)
		}

		logging.LogChangedByCtxf(ctx, "Created S3 object '%s' in bucket '%s'.", objectKey, bucketName)
	}
	logging.LogInfoByCtxf(ctx, "Create S3 object '%s' in bucket '%s' from bytes finished.", objectKey, bucketName)

	return nil
}

func IsObjectContentEqualString(ctx context.Context, client *minio.Client, bucketName string, objectKey string, objectContent string) (bool, error) {
	if client == nil {
		return false, tracederrors.TracedErrorNil("client")
	}

	if bucketName == "" {
		return false, tracederrors.TracedErrorEmptyString("bucketName")
	}

	if objectKey == "" {
		return false, tracederrors.TracedErrorEmptyString("objectKey")
	}

	return IsObjectContentEqualBytes(ctx, client, bucketName, objectKey, []byte(objectContent))
}

func GetObjectMd5Sum(ctx context.Context, client *minio.Client, bucketName string, objectKey string) (string, error) {
	if client == nil {
		return "", tracederrors.TracedErrorNil("client")
	}

	if bucketName == "" {
		return "", tracederrors.TracedErrorEmptyString("bucketName")
	}

	if objectKey == "" {
		return "", tracederrors.TracedErrorEmptyString("objectKey")
	}

	stat, err := client.StatObject(ctx, bucketName, objectKey, minio.StatObjectOptions{})
	if err != nil {
		return "", tracederrors.TracedErrorf("Failed to get MD5sum of S3 object '%s' in bucket '%s': %w", objectKey, bucketName, err)
	}

	md5Sum := stat.ETag

	logging.LogInfoByCtxf(ctx, "S3 object '%s' in bucket '%s' has MD5 checksum '%s'.", objectKey, bucketName, md5Sum)

	return md5Sum, nil
}

func IsObjectContentEqualBytes(ctx context.Context, client *minio.Client, bucketName string, objectKey string, objectContent []byte) (bool, error) {
	if client == nil {
		return false, tracederrors.TracedErrorNil("client")
	}

	if bucketName == "" {
		return false, tracederrors.TracedErrorEmptyString("bucketName")
	}

	if objectKey == "" {
		return false, tracederrors.TracedErrorEmptyString("objectKey")
	}

	if objectContent == nil {
		return false, tracederrors.TracedErrorNil("objectContent")
	}

	stat, err := client.StatObject(ctx, bucketName, objectKey, minio.StatObjectOptions{})
	if err != nil {
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code == "NoSuchKey" {
			return false, nil
		} else {
			return false, tracederrors.TracedErrorf("Failed to stat object '%s' in bucket '%s'.", objectKey, bucketName)
		}
	}

	gotMd5 := stat.ETag
	contentMd5 := checksumutils.GetMD5SumFromBytes(objectContent)

	return gotMd5 == contentMd5, nil
}

func GetObjectContentAsString(ctx context.Context, client *minio.Client, bucketName string, objectKey string) (string, error) {
	if client == nil {
		return "", tracederrors.TracedErrorNil("client")
	}

	if bucketName == "" {
		return "", tracederrors.TracedErrorEmptyString("bucketName")
	}

	if objectKey == "" {
		return "", tracederrors.TracedErrorEmptyString("objectKey")
	}

	content, err := GetObjectContentAsBytes(ctx, client, bucketName, objectKey)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func GetObjectContentAsBytes(ctx context.Context, client *minio.Client, bucketName string, objectKey string) ([]byte, error) {
	if client == nil {
		return nil, tracederrors.TracedErrorNil("client")
	}

	if bucketName == "" {
		return nil, tracederrors.TracedErrorEmptyString("bucketName")
	}

	if objectKey == "" {
		return nil, tracederrors.TracedErrorEmptyString("objectKey")
	}

	object, err := client.GetObject(ctx, bucketName, objectKey, minio.GetObjectOptions{})
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to get S3 object '%s' in bucket '%s'.", objectKey, bucketName)
	}
	defer object.Close()

	content, err := io.ReadAll(object)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to read S3 object '%s' in bucket '%s'.", objectKey, bucketName)
	}

	return content, nil
}

func DeleteObject(ctx context.Context, client *minio.Client, bucketName string, objectKey string) error {
	if client == nil {
		return tracederrors.TracedErrorNil("client")
	}

	if bucketName == "" {
		return tracederrors.TracedErrorEmptyString("bucketName")
	}

	if objectKey == "" {
		return tracederrors.TracedErrorEmptyString("objectKey")
	}

	exists, err := ObjectExists(contextutils.WithSilent(ctx), client, bucketName, objectKey)
	if err != nil {
		return err
	}

	if exists {
		err := client.RemoveObject(ctx, bucketName, objectKey, minio.RemoveObjectOptions{})
		if err != nil {
			return tracederrors.TracedErrorf("Failed to delete S3 object '%s' in bucket '%s': %w", objectKey, bucketName, err)
		}

		logging.LogChangedByCtxf(ctx, "Deleted S3 object '%s' in bucket '%s'.", objectKey, bucketName)
	} else {
		logging.LogInfoByCtxf(ctx, "S3 object '%s' in bucket '%s' is already absent. Skip deletion.", objectKey, bucketName)
	}

	return nil
}

func UploadFileByPath(ctx context.Context, client *minio.Client, bucketName string, objectKey string, path string) error {
	if client == nil {
		return tracederrors.TracedErrorNil("client")
	}

	if bucketName == "" {
		return tracederrors.TracedErrorEmptyString("bucketName")
	}

	if objectKey == "" {
		return tracederrors.TracedErrorEmptyString("objectKey")
	}

	if path == "" {
		return tracederrors.TracedErrorEmptyString("path")
	}

	logging.LogInfoByCtxf(ctx, "Upload '%s' to S3 bucket '%s' as '%s' started.", path, bucketName, objectKey)

	var performUpload bool
	stat, err := client.StatObject(ctx, bucketName, objectKey, minio.GetObjectOptions{})
	if err != nil {
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code == "NoSuchKey" {
			performUpload = true
		} else {
			return tracederrors.TracedErrorf("Failed to stat object: %w", err)
		}
	}
	if !performUpload {
		localHash, err := checksumutils.GetMD5SumFromFileByPath(ctx, path)
		if err != nil {
			return err
		}

		performUpload = localHash != stat.ETag
	}

	if performUpload {
		_, err = client.FPutObject(ctx, bucketName, objectKey, path, minio.PutObjectOptions{})
		if err != nil {
			return tracederrors.TracedErrorf("Failed to upload '%s' as S3 object '%s' into bucket '%s': %w", path, objectKey, bucketName, err)
		}
		logging.LogChangedByCtxf(ctx, "Uploaded '%s' as S3 object '%s' into bucket '%s'.", path, objectKey, bucketName)
	} else {
		logging.LogInfoByCtxf(ctx, "S3 object '%s' in bucket '%s' has already the same content as the file '%s'. Skip reupload.", objectKey, bucketName, path)
	}

	logging.LogInfoByCtxf(ctx, "Upload '%s' to S3 bucket '%s' as '%s' finished.", path, bucketName, objectKey)

	return nil
}

func DownloadAsTemporaryFile(ctx context.Context, client *minio.Client, bucketName string, objectKey string) (string, error) {
	if client == nil {
		return "", tracederrors.TracedErrorNil("client")
	}

	if bucketName == "" {
		return "", tracederrors.TracedErrorEmptyString("bucketName")
	}

	if objectKey == "" {
		return "", tracederrors.TracedErrorEmptyString("objectKey")
	}

	logging.LogInfoByCtxf(ctx, "Download S3 object '%s' from bucket '%s' as temporary file started.", objectKey, bucketName)

	tempPath, err := tempfiles.CreateTemporaryFile(ctx)
	if err != nil {
		return "", err
	}

	err = DownloadAsFileByPath(ctx, client, bucketName, objectKey, tempPath)
	if err != nil {
		return "", err
	}

	logging.LogInfoByCtxf(ctx, "Download S3 object '%s' from bucket '%s' as temporary file finished.", objectKey, bucketName)

	return tempPath, nil
}

func DownloadAsFileByPath(ctx context.Context, client *minio.Client, bucketName string, objectKey string, outputPath string) error {
	if client == nil {
		return tracederrors.TracedErrorNil("client")
	}

	if bucketName == "" {
		return tracederrors.TracedErrorEmptyString("bucketName")
	}

	if objectKey == "" {
		return tracederrors.TracedErrorEmptyString("objectKey")
	}

	if outputPath == "" {
		return tracederrors.TracedErrorEmptyString("outputPath")
	}

	logging.LogInfoByCtxf(ctx, "Download S3 object '%s' from bucket '%s' to '%s' started.", objectKey, bucketName, outputPath)

	var performDownload bool
	localChecksum, err := checksumutils.GetMD5SumFromFileByPath(ctx, outputPath)
	if err != nil {
		if filesgeneric.IsErrFileNotFound(err) {
			performDownload = true
		} else {
			return err
		}
	}

	if !performDownload {
		checksum, err := GetObjectMd5Sum(ctx, client, bucketName, objectKey)
		if err != nil {
			return err
		}

		if localChecksum != checksum {
			performDownload = true
		}
	}

	if performDownload {
		err = client.FGetObject(ctx, bucketName, objectKey, outputPath, minio.GetObjectOptions{})
		if err != nil {
			return tracederrors.TracedErrorf("Failed to download S3 object '%s' from bucket '%s' as '%s': %w", objectKey, bucketName, outputPath, err)
		}
		logging.LogChangedByCtxf(ctx, "S3 object '%s' from bucket '%s' downloaded as '%s'.", objectKey, bucketName, outputPath)
	} else {
		logging.LogInfoByCtxf(ctx, "Local file '%s' has already same content as the S3 object '%s' from bucket '%s'. Skip redownload.", outputPath, objectKey, bucketName)
	}

	logging.LogInfoByCtxf(ctx, "Download S3 object '%s' from bucket '%s' to '%s' finished.", objectKey, bucketName, outputPath)

	return nil
}
