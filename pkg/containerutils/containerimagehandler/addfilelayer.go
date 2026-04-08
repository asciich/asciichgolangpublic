package containerimagehandler

import (
	"archive/tar"
	"bytes"
	"context"
	"io"
	"os"

	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
	"github.com/google/go-containerregistry/pkg/v1/types"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/containeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func AddFileToImage(ctx context.Context, image v1.Image, options *containeroptions.AddFileToImageOptions) (v1.Image, error) {
	if image == nil {
		return nil, tracederrors.TracedErrorNil("image")
	}

	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	sourcePath, err := options.GetSourceFilePath()
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Add '%s' to container image stared.", sourcePath)

	pathInArchive, err := options.GetPathInImage()
	if err != nil {
		return nil, err
	}

	mode, err := options.GetMode()
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

	file, err := os.Open(sourcePath)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to open '%s' to add it to the container image: %w", sourcePath, err)
	}
	stat, err := file.Stat()
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to stat '%s' to add it to the container image: %w", sourcePath, err)
	}

	header := &tar.Header{
		Name: pathInArchive,
		Size: stat.Size(),
		Mode: mode,
	}
	err = tw.WriteHeader(header)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to write tar header for '%s': %w", sourcePath, err)
	}

	_, err = io.Copy(tw, file)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to write file content for '%s' into tar: %w", sourcePath, err)
	}

	err = tw.Close()
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to close tar writer for '%s': %w", sourcePath, err)
	}

	layer, err := tarball.LayerFromOpener(
		func() (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader(buf.Bytes())), nil
		},
		tarball.WithMediaType(types.DockerLayer),
	)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to create new file layer to add '%s' to the container image: %w", sourcePath, err)
	}

	newImage, err := mutate.AppendLayers(image, layer)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to append new layer")
	}

	logging.LogInfoByCtxf(ctx, "Add '%s' as '%s' to container image finished.", sourcePath, pathInArchive)

	return newImage, nil
}

func AddFileToArchive(ctx context.Context, archivePath string, options *containeroptions.AddFileToImageArchiveOptions) error {
	if archivePath == "" {
		return tracederrors.TracedErrorEmptyString("archivePath")
	}

	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	srcFilePath, err := options.GetSourceFilePath()
	if err != nil {
		return err
	}

	pathInArchive, err := options.GetPathInImage()
	if err != nil {
		return err
	}

	newImageNameAndTag, err := options.GetNewImageNameAndTag()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Add file '%s' as '%s' into container image archive '%s' started.", srcFilePath, pathInArchive, archivePath)

	if !options.OverwriteSourceArchive {
		return tracederrors.TracedError("Only implemented to overwrite source archive.")
	}

	tag, err := name.NewTag(newImageNameAndTag)
	if err != nil {
		return tracederrors.TracedErrorf("Failed to parse new image name and tag '%s': %w", newImageNameAndTag, err)
	}

	image, _, err := LoadImageFromArchive(ctx, archivePath, "")
	if err != nil {
		return err
	}

	newImage, err := AddFileToImage(ctx, image, &containeroptions.AddFileToImageOptions{
		SourceFilePath: srcFilePath,
		PathInImage: pathInArchive,
		Mode: options.Mode,
	})
	if err != nil {
		return err
	}

	err = OverwriteArchive(ctx, archivePath, &tag, newImage)
	if err != nil {
		return err
	}

	logging.LogChangedByCtxf(ctx, "Added file '%s' as '%s' in container image '%s' to '%s'.", srcFilePath, pathInArchive, newImageNameAndTag, archivePath)

	logging.LogInfoByCtxf(ctx, "Add file '%s' as '%s' into container image archive '%s' finished.", srcFilePath, pathInArchive, archivePath)

	return nil
}
