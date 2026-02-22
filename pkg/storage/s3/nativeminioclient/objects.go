package nativeminioclient

import (
	"bytes"
	"context"
	"io"
	"sort"

	"github.com/minio/minio-go/v7"
	"github.com/asciich/asciichgolangpublic/pkg/checksumutils"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
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
