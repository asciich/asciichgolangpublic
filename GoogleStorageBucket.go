package asciichgolangpublic

import (
	"context"

	"cloud.google.com/go/storage"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type GoogleStorageBucket struct {
	name string

	nativeClient *storage.Client
}

func GetGoogleStorageBucketByName(bucketName string) (g *GoogleStorageBucket, err error) {
	if bucketName == "" {
		return nil, tracederrors.TracedErrorEmptyString("bucketName")
	}

	g = NewGoogleStorageBucket()

	err = g.SetName(bucketName)
	if err != nil {
		return nil, err
	}

	return g, nil
}

func MustGetGoogleStorageBucketByName(bucketName string) (g *GoogleStorageBucket) {
	g, err := GetGoogleStorageBucketByName(bucketName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return g
}

func NewGoogleStorageBucket() (g *GoogleStorageBucket) {
	return new(GoogleStorageBucket)
}

func (g *GoogleStorageBucket) Exists() (bucketExists bool, err error) {
	nativeBucket, err := g.GetNativeBucket()
	if err != nil {
		return false, err
	}

	_, err = nativeBucket.Attrs(context.Background())
	if err != nil {
		if err.Error() == "storage: bucket doesn't exist" {
			return false, nil
		}
		return false, tracederrors.TracedErrorf("Unable to get bucket Attrs: '%w'", err)
	}

	return true, nil
}

func (g *GoogleStorageBucket) GetName() (name string, err error) {
	if g.name == "" {
		return "", tracederrors.TracedErrorf("name not set")
	}

	return g.name, nil
}

func (g *GoogleStorageBucket) GetNativeBucket() (nativeBucket *storage.BucketHandle, err error) {
	nativeClient, err := g.GetNativeClient()
	if err != nil {
		return nil, err
	}

	bucketName, err := g.GetName()
	if err != nil {
		return nil, err
	}

	nativeBucket = nativeClient.Bucket(bucketName)

	return nativeBucket, nil
}

func (g *GoogleStorageBucket) GetNativeClient() (nativeClient *storage.Client, err error) {
	if g.nativeClient == nil {
		clientToAdd, err := storage.NewClient(context.Background())
		if err != nil {
			return nil, tracederrors.TracedErrorf("Unable to create native storage client: %w", err)
		}
		g.nativeClient = clientToAdd
	}
	return g.nativeClient, nil
}

func (g *GoogleStorageBucket) MustExists() (bucketExists bool) {
	bucketExists, err := g.Exists()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return bucketExists
}

func (g *GoogleStorageBucket) MustGetName() (name string) {
	name, err := g.GetName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return name
}

func (g *GoogleStorageBucket) MustGetNativeBucket() (nativeBucket *storage.BucketHandle) {
	nativeBucket, err := g.GetNativeBucket()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return nativeBucket
}

func (g *GoogleStorageBucket) MustGetNativeClient() (nativeClient *storage.Client) {
	nativeClient, err := g.GetNativeClient()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return nativeClient
}

func (g *GoogleStorageBucket) MustSetName(name string) {
	err := g.SetName(name)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GoogleStorageBucket) MustSetNativeClient(nativeClient *storage.Client) {
	err := g.SetNativeClient(nativeClient)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GoogleStorageBucket) SetName(name string) (err error) {
	if name == "" {
		return tracederrors.TracedErrorf("name is empty string")
	}

	g.name = name

	return nil
}

func (g *GoogleStorageBucket) SetNativeClient(nativeClient *storage.Client) (err error) {
	g.nativeClient = nativeClient

	return nil
}
