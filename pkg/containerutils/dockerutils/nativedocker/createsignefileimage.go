package nativedocker

import (
	"context"

	"github.com/google/go-containerregistry/pkg/v1/daemon"

	"github.com/asciich/asciichgolangpublic/pkg/containerutils/containerimagehandler"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/containeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// CreateSingleFileImage creates a container image containing only the single file defined in the options
// and stores the image directly in the local docker daemon.
//
// An already existing image with the same name and tag is overwritten (re-tagged) by the docker daemon.
func CreateSingleFileImage(ctx context.Context, options *containeroptions.CreateSingleFileArchiveOptions) (*Image, error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	srcFilePath, err := options.GetSourceFilePath()
	if err != nil {
		return nil, err
	}

	pathInImage, err := options.GetPathInImage()
	if err != nil {
		return nil, err
	}

	newImageNameAndTag, err := options.GetNewImageNameAndTag()
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Create single file docker image '%s' with file '%s' as '%s' started.", newImageNameAndTag, srcFilePath, pathInImage)

	image, tag, err := containerimagehandler.CreateSingleFileImage(ctx, options)
	if err != nil {
		return nil, err
	}

	response, err := daemon.Write(*tag, image, daemon.WithContext(ctx))
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to write single file image '%s' into docker: %w", newImageNameAndTag, err)
	}

	logging.LogInfoByCtxf(ctx, "Docker daemon response for loading image '%s': %s", newImageNameAndTag, response)

	logging.LogChangedByCtxf(ctx, "Created single file docker image '%s' with file '%s' as '%s'.", newImageNameAndTag, srcFilePath, pathInImage)

	ret := NewImage()
	err = ret.SetName(newImageNameAndTag)
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Create single file docker image '%s' with file '%s' as '%s' finished.", newImageNameAndTag, srcFilePath, pathInImage)

	return ret, nil
}
