package containerimagehandler

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
	"github.com/google/go-containerregistry/pkg/v1/types"

	"github.com/asciich/asciichgolangpublic/pkg/containerutils/containeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
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

func ListImageNamesAndTagsInArchive(ctx context.Context, archivePath string) ([]string, error) {
	if archivePath == "" {
		return nil, tracederrors.TracedErrorEmptyString("archivePath")
	}

	logging.LogInfoByCtxf(ctx, "List image names and tags in container image archive '%s' started.", archivePath)
	f, err := os.Open(archivePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var tags []string
	tr := tar.NewReader(f)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if header.Name == "manifest.json" {

			type imageManifest struct {
				Config   string   `json:"Config"`
				RepoTags []string `json:"RepoTags"`
				Layers   []string `json:"Layers"`
			}

			var data = []*imageManifest{}
			err := json.NewDecoder(tr).Decode(&data)
			if err != nil {
				return nil, tracederrors.TracedErrorf("Failed to decode manifest.json from container image archive '%s': %w", archivePath, err)
			}

			for _, d := range data {
				for _, r := range d.RepoTags {
					tags = append(tags, r)
				}
			}
		}
	}
	if tags == nil {
		return nil, tracederrors.TracedErrorf("No metadata found to list image names and tags in image archive '%s'.", archivePath)
	}

	logging.LogInfoByCtxf(ctx, "List image names and tags in contaienr image archive '%s' finished.", archivePath)

	return tags, nil
}

func ListFilesInArchive(ctx context.Context, archivePath string) ([]string, error) {
	if archivePath == "" {
		return nil, tracederrors.TracedErrorEmptyString("archivePath")
	}

	logging.LogInfoByCtxf(ctx, "List files in container image archive '%s' started.", archivePath)

	image, _, err := LoadImageFromArchive(ctx, archivePath, "")
	if err != nil {
		return nil, err
	}

	layers, err := image.Layers()
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to get image layers: %w", err)
	}

	fileSet := make(map[string]bool)
	for _, layer := range layers {
		reader, err := layer.Uncompressed()
		if err != nil {
			return nil, tracederrors.TracedErrorf("Failed to get uncompressed reader: %w", err)
		}
		defer reader.Close()

		tr := tar.NewReader(reader)
		for {
			header, err := tr.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, tracederrors.TracedErrorf("Failed to read tar data: %w", err)
			}

			// Handle Whiteouts:
			// If a file is prefixed with .wh., it indicates a deletion of the file
			// in the layers beneath it.
			if strings.Contains(header.Name, ".wh.") {
				// The actual path being deleted:
				originalPath := strings.Replace(header.Name, ".wh.", "", 1)
				delete(fileSet, originalPath)
			} else {
				// Otherwise, add/update the file in our set
				fileSet[header.Name] = true
			}
		}
	}

	var files []string
	for path := range fileSet {
		files = append(files, path)
	}

	logging.LogInfoByCtxf(ctx, "List files in container image archive '%s' finished.", archivePath)

	return files, nil
}

func FileInArchiveExists(ctx context.Context, archivePath string, pathToCheck string) (bool, error) {
	if archivePath == "" {
		return false, tracederrors.TracedErrorEmptyString("archivePath")
	}

	if pathToCheck == "" {
		return false, tracederrors.TracedErrorEmptyString("pathToCheck")
	}

	logging.LogInfoByCtxf(ctx, "File '%s' exists in container image archive '%s' started.", pathToCheck, archivePath)

	fileList, err := ListFilesInArchive(ctx, archivePath)
	if err != nil {
		return false, err
	}

	exists := slices.Contains(fileList, pathToCheck)

	if exists {
		logging.LogInfoByCtxf(ctx, "File '%s' exists in container image archive '%s'.", pathToCheck, archivePath)
	} else {
		logging.LogInfoByCtxf(ctx, "File '%s' does not exist in container image archive '%s'.", pathToCheck, archivePath)
	}

	logging.LogInfoByCtxf(ctx, "File '%s' exists in container image archive '%s' finished.", pathToCheck, archivePath)

	return exists, nil
}

// Load the image tagged as "imageNameAndTag" from a container image archive.
//
// If imageNammeAndTag is empty the tag is evaluated based on the archive content.
// This will fail when multiple tags are defined in one archive.
// In the case of multiple tags the tag has to be explicitly defined as imageNameAndTag
func LoadImageFromArchive(ctx context.Context, archivePath string, imageNameAndTag string) (v1.Image, *name.Tag, error) {
	if archivePath == "" {
		return nil, nil, tracederrors.TracedErrorEmptyString("archivePath")
	}

	logging.LogInfoByCtxf(ctx, "Load container image from archive '%s' started.", archivePath)

	if imageNameAndTag == "" {
		tags, err := ListImageNamesAndTagsInArchive(ctx, archivePath)
		if err != nil {
			return nil, nil, err
		}

		if len(tags) == 1 {
			imageNameAndTag = tags[0]
			logging.LogInfoByCtxf(ctx, "Going to load the only image name and tag '%s' found in the archive '%s'.", imageNameAndTag, archivePath)
		} else {
			return nil, nil, tracederrors.TracedErrorf("Failed to automatically detect imageNameAndTag. Found tags: '%v' in the archive '%s'.", tags, archivePath)
		}
	}

	tag, err := name.NewTag(imageNameAndTag)
	if err != nil {
		return nil, nil, tracederrors.TracedErrorf("Failed to load image name and tag '%s': %w", imageNameAndTag, err)
	}

	image, err := tarball.ImageFromPath(archivePath, &tag)
	if err != nil {
		return nil, nil, tracederrors.TracedErrorf("Failed to load '%s' from container image archive '%s'.", imageNameAndTag, archivePath)
	}

	logging.LogInfoByCtxf(ctx, "Load container image from archive '%s' finished.", archivePath)

	return image, &tag, nil
}

func OverwriteArchive(ctx context.Context, archivePath string, tag *name.Tag, image v1.Image) error {
	if archivePath == "" {
		return tracederrors.TracedErrorEmptyString("archivePath")
	}

	if tag == nil {
		return tracederrors.TracedErrorNil("tag")
	}

	if image == nil {
		return tracederrors.TracedErrorNil("image")
	}

	logging.LogInfoByCtxf(ctx, "Overwrite archive '%s' with new image '%s' started.", archivePath, tag)

	tmpPath := archivePath + ".tmp"
	logging.LogInfoByCtxf(ctx, "Write new image for '%s' into temporary archive '%s'.", tag, tmpPath)

	err := tarball.WriteToFile(tmpPath, *tag, image)
	if err != nil {
		return tracederrors.TracedErrorf("Failed to write '%s': %w", tmpPath, err)
	}

	err = os.Rename(tmpPath, archivePath)
	if err != nil {
		return tracederrors.TracedErrorf("Failed to rename '%s' to '%s': %w", tmpPath, archivePath, err)
	}

	logging.LogChangedByCtxf(ctx, "Overwritten container image archive '%s' with '%s'.", archivePath, tag)

	logging.LogInfoByCtxf(ctx, "Overwrite archive '%s' with new image '%s' finished.", archivePath, tag)

	return nil
}

func AddFileToArchive(ctx context.Context, archivePath string, options *containeroptions.AddFileToImageOptions) error {
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

	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

	file, err := os.Open(srcFilePath)
	if err != nil {
		return tracederrors.TracedErrorf("Failed to open '%s' to add it to the container image archive '%s': %w", srcFilePath, archivePath, err)
	}
	stat, err := file.Stat()
	if err != nil {
		return tracederrors.TracedErrorf("Failed to stat '%s' to add it to the container image archive '%s': %w", srcFilePath, archivePath, err)
	}

	header := &tar.Header{
		Name: pathInArchive,
		Size: stat.Size(),
		Mode: 0644,
	}
	tw.WriteHeader(header)
	io.Copy(tw, file)
	tw.Close()

	layer, err := tarball.LayerFromOpener(
		func() (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader(buf.Bytes())), nil
		},
		tarball.WithMediaType(types.DockerLayer),
	)
	if err != nil {
		return tracederrors.TracedErrorf("Failed to create new file layer to add '%s' to the container image archive '%s': %w", srcFilePath, archivePath, err)
	}

	newImage, err := mutate.AppendLayers(image, layer)
	if err != nil {
		return tracederrors.TracedErrorf("Failed to append new layer")
	}

	err = OverwriteArchive(ctx, archivePath, &tag, newImage)
	if err != nil {
		return err
	}

	logging.LogChangedByCtxf(ctx, "Added file '%s' as '%s' in container image '%s' to '%s'.", srcFilePath, pathInArchive, newImageNameAndTag, archivePath)

	logging.LogInfoByCtxf(ctx, "Add file '%s' as '%s' into container image archive '%s' finished.", srcFilePath, pathInArchive, archivePath)

	return nil
}

func ReadFileFromArchiveAsString(ctx context.Context, archivePath string, toReadPath string) (string, error) {
	content, err := ReadFileFromArchiveAsBytes(ctx, archivePath, toReadPath)
	if err != nil {
		return "", err
	}

	return string(content), err
}

func ReadFileFromArchiveAsBytes(ctx context.Context, archivePath string, toReadPath string) ([]byte, error) {
	if archivePath == "" {
		return nil, tracederrors.TracedErrorEmptyString("archivePath")
	}

	if toReadPath == "" {
		return nil, tracederrors.TracedErrorEmptyString("toReadPath")
	}

	logging.LogInfoByCtxf(ctx, "Read file '%s' from container image archive '%s' started.", archivePath, toReadPath)

	image, _, err := LoadImageFromArchive(ctx, archivePath, "")
	if err != nil {
		return nil, err
	}

	layers, err := image.Layers()
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to get image layers: %w", err)
	}

	var content []byte
	for i := len(layers) - 1; i >= 0; i-- {
		reader, err := layers[i].Uncompressed()
		if err != nil {
			return nil, tracederrors.TracedErrorf("Failed to get uncompressed reader: %w", err)
		}
		defer reader.Close()

		digest, err := layers[i].Digest()
		if err != nil {
			return nil, tracederrors.TracedErrorf("Unable to get layer digest: %w", err)
		}

		tr := tar.NewReader(reader)
		for {
			header, err := tr.Next()
			if err == io.EOF {
				break
			}

			// 1. Check for Whiteout (Deletion)
			if header.Name == ".wh."+toReadPath {
				return nil, tracederrors.TracedErrorf("File was deleted in the upper layer '%s' and is therefore not present in the image archive '%s'.", digest.Hex, archivePath)
			}

			// 2. Check for File Match
			if header.Name == toReadPath {
				content, err = io.ReadAll(tr)
				if err != nil {
					return nil, tracederrors.TracedErrorf("Failed to read '%s' in container image archive '%s'.", toReadPath, archivePath)
				}
				logging.LogInfoByCtxf(ctx, "Read file '%s' from container image archive '%s' layer '%s'.", toReadPath, archivePath, digest)
				break
			}
		}
	}

	if content == nil {
		return nil, tracederrors.TracedErrorf("Failed to read '%s' in container image archive '%s'. File was not found.", toReadPath, archivePath)
	}

	logging.LogInfoByCtxf(ctx, "Read file '%s' from container image archive '%s' finished.", archivePath, toReadPath)

	return content, nil
}

func DeleteFileInArchive(ctx context.Context, archivePath string, options *containeroptions.DeleteFileFromImageOptions) error {
	if archivePath == "" {
		return tracederrors.TracedErrorEmptyString(archivePath)
	}

	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	pathInArchive, err := options.GetPathInImage()
	if err != nil {
		return err
	}

	newImageNameAndTag, err := options.GetNewImageNameAndTag()
	if err != nil {
		return err
	}

	tag, err := name.NewTag(newImageNameAndTag)
	if err != nil {
		return tracederrors.TracedErrorf("Failed to parse tag '%s': %w", newImageNameAndTag, err)
	}

	logging.LogInfoByCtxf(ctx, "Delete file '%s' from container image archive '%s' started.", pathInArchive, archivePath)

	if !options.OverwriteSourceArchive {
		return tracederrors.TracedError("Only implemented to overwrite source archive.")
	}

	image, _, err := LoadImageFromArchive(ctx, archivePath, "")
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

	header := &tar.Header{
		Name: ".wh." + pathInArchive,
		Size: 0,
		Mode: 0644,
	}
	tw.WriteHeader(header)
	tw.Close()

	layer, err := tarball.LayerFromOpener(
		func() (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader(buf.Bytes())), nil
		},
		tarball.WithMediaType(types.DockerLayer),
	)
	if err != nil {
		return tracederrors.TracedErrorf("Failed to create new file layer to delete '%s' from the container image archive '%s': %w", pathInArchive, archivePath, err)
	}

	newImage, err := mutate.AppendLayers(image, layer)
	if err != nil {
		return tracederrors.TracedErrorf("Failed to add new layer: %w", err)
	}

	err = OverwriteArchive(ctx, archivePath, &tag, newImage)
	if err != nil {
		return err
	}

	logging.LogChangedByCtxf(ctx, "Deleted '%s' in container image archive '%s'.", pathInArchive, archivePath)

	logging.LogInfoByCtxf(ctx, "Delete file '%s' from container image archive '%s' started.", pathInArchive, archivePath)

	return nil
}
