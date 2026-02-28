package containerimagehandler

import (
	"context"
	"os"
	"strings"

	"github.com/google/go-containerregistry/pkg/legacy/tarball"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func DownloadImageAsArchive(ctx context.Context, imageNameAndTag string, outputPath string) error {
	if imageNameAndTag == "" {
		return tracederrors.TracedErrorEmptyString("imageNameAndTag")
	}

	if outputPath == "" {
		return tracederrors.TracedErrorEmptyString("outputPath")
	}

	logging.LogInfoByCtxf(ctx, "Download container image '%s' as archive '%s' started.", imageNameAndTag, outputPath)

	if !strings.Contains(imageNameAndTag, ":") {
		imageNameAndTag += ":latest"
		logging.LogInfoByCtxf(ctx, "Going to download latest: '%s'", imageNameAndTag)
	}

	ref, err := name.ParseReference(imageNameAndTag)
	if err != nil {
		return tracederrors.TracedErrorf("Failed parse reference to download container image '%s' as archive '%s': %w", imageNameAndTag, outputPath, err)
	}

	img, err := remote.Image(ref)
	if err != nil {
		return tracederrors.TracedErrorf("Failed to get remote image descriptor to download container image '%s' as archive '%s': %w", imageNameAndTag, outputPath, err)
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
