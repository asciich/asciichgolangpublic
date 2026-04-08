package containerimagehandler

import (
	"context"
	"os"
	"strings"

	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func DownloadImage(ctx context.Context, imageNameAndTag string) (v1.Image, name.Reference, error) {
	if imageNameAndTag == "" {
		return nil, nil, tracederrors.TracedErrorEmptyString("imageNameAndTag")
	}

	logging.LogInfoByCtxf(ctx, "Download container image '%s' started.", imageNameAndTag)

	if !strings.Contains(imageNameAndTag, ":") {
		imageNameAndTag += ":latest"
		logging.LogInfoByCtxf(ctx, "Going to download latest: '%s'", imageNameAndTag)
	}

	ref, err := name.ParseReference(imageNameAndTag)
	if err != nil {
		return nil, nil, tracederrors.TracedErrorf("Failed parse reference to download container image archive '%s': %w", imageNameAndTag, err)
	}

	img, err := remote.Image(ref)
	if err != nil {
		return nil, nil, tracederrors.TracedErrorf("Failed to get remote image descriptor to download container image '%s': %w", imageNameAndTag, err)
	}

	logging.LogInfoByCtxf(ctx, "Download container image '%s' finished.", imageNameAndTag)

	return img, ref, err
}

func DownloadImageAsArchive(ctx context.Context, imageNameAndTag string, outputPath string) error {
	if imageNameAndTag == "" {
		return tracederrors.TracedErrorEmptyString("imageNameAndTag")
	}

	if outputPath == "" {
		return tracederrors.TracedErrorEmptyString("outputPath")
	}

	logging.LogInfoByCtxf(ctx, "Download container image '%s' as archive '%s' started.", imageNameAndTag, outputPath)

	img, ref, err := DownloadImage(ctx, imageNameAndTag)
	if err != nil {
		return err
	}

	f, err := os.Create(outputPath)
	if err != nil {
		return tracederrors.TracedErrorf("Failed to creat outputPath to download container image '%s' as archive '%s': %w", imageNameAndTag, outputPath, err)
	}
	defer f.Close()

	err = tarball.Write(ref, img, f)
	if err != nil {
		return tracederrors.TracedErrorf("Failed to download container image '%s' as archive '%s': %w", imageNameAndTag, outputPath, err)
	}

	logging.LogChangedByCtxf(ctx, "Downloaded container image '%s' as archive '%s'.", imageNameAndTag, outputPath)

	logging.LogInfoByCtxf(ctx, "Download container image '%s' as archive '%s' finished.", imageNameAndTag, outputPath)

	return nil
}

func DownloadImageAsTeporaryArchive(ctx context.Context, imageNameAndTag string) (string, error) {
	tempFile, err := tempfiles.CreateNamedTemporaryFile(ctx, strings.ReplaceAll(imageNameAndTag, ":", "_"))
	if err != nil {
		return "", err
	}

	logging.LogInfoByCtxf(ctx, "Going to download container image '%s' to temporary file '%s'.", imageNameAndTag, tempFile)

	err = DownloadImageAsArchive(ctx, imageNameAndTag, tempFile)
	if err != nil {
		return "", err
	}

	return tempFile, nil
}
