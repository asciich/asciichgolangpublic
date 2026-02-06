package nativeminioclient

import (
	"context"

	"github.com/minio/minio-go/v7"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func ListBuckets(ctx context.Context, client *minio.Client) ([]minio.BucketInfo, error) {
	logging.LogInfoByCtxf(ctx, "List minio buckets started.")

	if client == nil {
		return nil, tracederrors.TracedErrorNil("client")
	}

	buckets, err := client.ListBuckets(ctx)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to list minio buckets: %w", err)
	}

	logging.LogInfoByCtxf(ctx, "List minio buckets finished. Found '%d' buckets.", len(buckets))

	return buckets, nil
}

func ListBucketNames(ctx context.Context, client *minio.Client) ([]string, error) {
	if client == nil {
		return nil, tracederrors.TracedErrorNil("client")
	}

	buckets, err := ListBuckets(ctx, client)
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(buckets))
	for _, b := range buckets {
		names = append(names, b.Name)
	}

	return names, nil
}

func BucketExists(ctx context.Context, client *minio.Client, bucketName string) (bool, error) {
	if client == nil {
		return false, tracederrors.TracedErrorNil("client")
	}

	if bucketName == "" {
		return false, tracederrors.TracedErrorEmptyString("bucketName")
	}

	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return false, tracederrors.TracedErrorf("Failed to evaluate if bucket '%s' exists.", bucketName)
	}

	if exists {
		logging.LogInfoByCtxf(ctx, "Bucket '%s' exists.", bucketName)
	} else {
		logging.LogInfoByCtxf(ctx, "Bucket '%s' does not exist.", bucketName)
	}

	return exists, nil
}

func CreateBucket(ctx context.Context, client *minio.Client, bucketName string) error {
	if client == nil {
		return tracederrors.TracedErrorNil("client")
	}

	if bucketName == "" {
		return tracederrors.TracedErrorEmptyString("bucketName")
	}

	logging.LogInfoByCtxf(ctx, "Create bucket '%s' started.", bucketName)

	var created = true
	err := client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		errResponse := minio.ToErrorResponse(err)

		if errResponse.Code == "BucketAlreadyOwnedByYou" {
			created = false
		} else {
			return tracederrors.TracedErrorf("Failed to create bucket '%s': %w", bucketName, err)
		}
	}

	if created {
		logging.LogChangedByCtxf(ctx, "Bucket '%s' created.", bucketName)
	} else {
		logging.LogInfoByCtxf(ctx, "Bucket '%s' already exists and belongs to you. Skip creation.", bucketName)
	}

	logging.LogInfoByCtxf(ctx, "Create bucket '%s' finised.", bucketName)

	return nil
}

func DeleteBucket(ctx context.Context, client *minio.Client, bucketName string) error {
	if client == nil {
		return tracederrors.TracedErrorNil("client")
	}

	if bucketName == "" {
		return tracederrors.TracedErrorEmptyString("bucketName")
	}

	logging.LogInfoByCtxf(ctx, "Delete bucket '%s' started.", bucketName)

	var deleted = true
	err := client.RemoveBucket(ctx, bucketName)
	if err != nil {
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code == "NoSuchBucket" {
			deleted = false
		} else {
			return tracederrors.TracedErrorf("Failed to delete bucket '%s': %w", bucketName, err)
		}
	}

	if deleted {
		logging.LogChangedByCtxf(ctx, "Deleted bucket '%s'.", bucketName)
	} else {
		logging.LogInfoByCtxf(ctx, "Bucket '%s' is already absent. Skip delete.", bucketName)
	}

	logging.LogInfoByCtxf(ctx, "Delete bucket '%s' finised.", bucketName)

	return nil
}
