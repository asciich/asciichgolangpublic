package containerimagehandler_test

import (
	"archive/tar"
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/asciich/asciichgolangpublic/pkg/containerutils/containerimagehandler"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/containeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/datatypes/pointerutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
)

const testImageNameAndTag = "single-file-test:v1"
const testPathInImage = "hello.txt"
const testFileContent = "hello world\n"

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

// createTestFile writes content into a file inside a temporary directory which is removed automatically after the test.
func createTestFile(t *testing.T, content string) string {
	t.Helper()

	ret := filepath.Join(t.TempDir(), "hello.txt")
	err := os.WriteFile(ret, []byte(content), 0644)
	require.NoError(t, err)

	return ret
}

// getTestArchivePath returns a not yet existing archive path inside a temporary directory which is removed automatically after the test.
func getTestArchivePath(t *testing.T) string {
	t.Helper()

	return filepath.Join(t.TempDir(), "image.tar")
}

func getSingleFileOptions(srcPath string) *containeroptions.CreateSingleFileArchiveOptions {
	return &containeroptions.CreateSingleFileArchiveOptions{
		SourceFilePath:      srcPath,
		PathInImage:         testPathInImage,
		NewImageNameAndTag:  testImageNameAndTag,
		Mode:                pointerutils.ToInt64Pointer(0644),
		Architecture:        "amd64",
		SkipStaticLinkCheck: true, // The test file is a text file and not a statically linked binary.
	}
}

// createTestArchive creates a single file archive containing testFileContent as testPathInImage.
func createTestArchive(ctx context.Context, t *testing.T) string {
	t.Helper()

	ret := getTestArchivePath(t)
	err := containerimagehandler.CreateSingleFileArchive(ctx, ret, getSingleFileOptions(createTestFile(t, testFileContent)))
	require.NoError(t, err)

	return ret
}

func TestListImageAndTagsInArchvie(t *testing.T) {
	t.Run("alpine", func(t *testing.T) {
		ctx := getCtx()

		archivePath, err := containerimagehandler.DownloadImageAsTeporaryArchive(ctx, "alpine")
		require.NoError(t, err)
		defer nativefiles.Delete(ctx, archivePath, &filesoptions.DeleteOptions{})

		imageNamesAndTags, err := containerimagehandler.ListImageNamesAndTagsInArchive(ctx, archivePath)
		require.NoError(t, err)
		require.EqualValues(t, []string{"alpine:latest"}, imageNamesAndTags)
	})

	t.Run("single file archive", func(t *testing.T) {
		ctx := getCtx()

		archivePath := createTestArchive(ctx, t)

		imageNamesAndTags, err := containerimagehandler.ListImageNamesAndTagsInArchive(ctx, archivePath)
		require.NoError(t, err)
		require.EqualValues(t, []string{testImageNameAndTag}, imageNamesAndTags)
	})

	t.Run("empty archive path", func(t *testing.T) {
		ctx := getCtx()

		_, err := containerimagehandler.ListImageNamesAndTagsInArchive(ctx, "")
		require.Error(t, err)
	})

	t.Run("non existing archive", func(t *testing.T) {
		ctx := getCtx()

		_, err := containerimagehandler.ListImageNamesAndTagsInArchive(ctx, getTestArchivePath(t))
		require.Error(t, err)
	})
}

func TestAlpineArchive(t *testing.T) {
	ctx := getCtx()

	// Download only once for all subtests to keep the test fast.
	archivePath, err := containerimagehandler.DownloadImageAsTeporaryArchive(ctx, "alpine")
	require.NoError(t, err)
	defer nativefiles.Delete(ctx, archivePath, &filesoptions.DeleteOptions{})

	t.Run("list files contains alpine-release", func(t *testing.T) {
		files, err := containerimagehandler.ListFilesInArchive(ctx, archivePath)
		require.NoError(t, err)
		require.Contains(t, files, "etc/alpine-release")
	})

	t.Run("alpine-release exists", func(t *testing.T) {
		exists, err := containerimagehandler.FileInArchiveExists(ctx, archivePath, "etc/alpine-release")
		require.NoError(t, err)
		require.True(t, exists)
	})

	t.Run("non existing file does not exist", func(t *testing.T) {
		exists, err := containerimagehandler.FileInArchiveExists(ctx, archivePath, "this/file/does/not/exist.txt")
		require.NoError(t, err)
		require.False(t, exists)
	})

	t.Run("read alpine-release", func(t *testing.T) {
		content, err := containerimagehandler.ReadFileFromArchiveAsString(ctx, archivePath, "etc/alpine-release")
		require.NoError(t, err)
		require.NotEmpty(t, content)
	})

	t.Run("load image with auto detected tag", func(t *testing.T) {
		image, tag, err := containerimagehandler.LoadImageFromArchive(ctx, archivePath, "")
		require.NoError(t, err)
		require.NotNil(t, image)
		require.NotNil(t, tag)

		configFile, err := image.ConfigFile()
		require.NoError(t, err)
		require.EqualValues(t, "linux", configFile.OS)
	})
}

func TestCreateSingleFileImage(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		ctx := getCtx()

		image, tag, err := containerimagehandler.CreateSingleFileImage(ctx, getSingleFileOptions(createTestFile(t, testFileContent)))
		require.NoError(t, err)
		require.NotNil(t, image)
		require.NotNil(t, tag)

		t.Run("tag", func(t *testing.T) {
			require.EqualValues(t, testImageNameAndTag, tag.String())
		})

		t.Run("platform", func(t *testing.T) {
			configFile, err := image.ConfigFile()
			require.NoError(t, err)
			require.EqualValues(t, "amd64", configFile.Architecture)
			require.EqualValues(t, "linux", configFile.OS)
		})

		t.Run("single layer with single file", func(t *testing.T) {
			layers, err := image.Layers()
			require.NoError(t, err)
			require.Len(t, layers, 1)

			reader, err := layers[0].Uncompressed()
			require.NoError(t, err)
			defer reader.Close()

			tr := tar.NewReader(reader)

			header, err := tr.Next()
			require.NoError(t, err)
			require.EqualValues(t, testPathInImage, header.Name)
			require.EqualValues(t, tar.TypeReg, header.Typeflag)
			require.EqualValues(t, 0644, header.Mode)
			require.EqualValues(t, len(testFileContent), header.Size)

			content, err := io.ReadAll(tr)
			require.NoError(t, err)
			require.EqualValues(t, testFileContent, string(content))

			// There must be no other entry in the layer.
			_, err = tr.Next()
			require.ErrorIs(t, err, io.EOF)
		})
	})

	t.Run("architecture arm64", func(t *testing.T) {
		ctx := getCtx()

		options := getSingleFileOptions(createTestFile(t, testFileContent))
		options.Architecture = "arm64"

		image, _, err := containerimagehandler.CreateSingleFileImage(ctx, options)
		require.NoError(t, err)

		configFile, err := image.ConfigFile()
		require.NoError(t, err)
		require.EqualValues(t, "arm64", configFile.Architecture)
	})

	t.Run("custom mode", func(t *testing.T) {
		ctx := getCtx()

		options := getSingleFileOptions(createTestFile(t, testFileContent))
		options.Mode = pointerutils.ToInt64Pointer(0755)

		image, _, err := containerimagehandler.CreateSingleFileImage(ctx, options)
		require.NoError(t, err)

		layers, err := image.Layers()
		require.NoError(t, err)
		require.Len(t, layers, 1)

		reader, err := layers[0].Uncompressed()
		require.NoError(t, err)
		defer reader.Close()

		header, err := tar.NewReader(reader).Next()
		require.NoError(t, err)
		require.EqualValues(t, 0755, header.Mode)
	})

	t.Run("nil options", func(t *testing.T) {
		ctx := getCtx()

		_, _, err := containerimagehandler.CreateSingleFileImage(ctx, nil)
		require.Error(t, err)
	})

	t.Run("empty source file path", func(t *testing.T) {
		ctx := getCtx()

		options := getSingleFileOptions("")

		_, _, err := containerimagehandler.CreateSingleFileImage(ctx, options)
		require.Error(t, err)
	})

	t.Run("empty path in image", func(t *testing.T) {
		ctx := getCtx()

		options := getSingleFileOptions(createTestFile(t, testFileContent))
		options.PathInImage = ""

		_, _, err := containerimagehandler.CreateSingleFileImage(ctx, options)
		require.Error(t, err)
	})

	t.Run("empty new image name and tag", func(t *testing.T) {
		ctx := getCtx()

		options := getSingleFileOptions(createTestFile(t, testFileContent))
		options.NewImageNameAndTag = ""

		_, _, err := containerimagehandler.CreateSingleFileImage(ctx, options)
		require.Error(t, err)
	})

	t.Run("invalid new image name and tag", func(t *testing.T) {
		ctx := getCtx()

		options := getSingleFileOptions(createTestFile(t, testFileContent))
		options.NewImageNameAndTag = "Invalid Name With Spaces:v1"

		_, _, err := containerimagehandler.CreateSingleFileImage(ctx, options)
		require.Error(t, err)
	})

	t.Run("empty architecture", func(t *testing.T) {
		ctx := getCtx()

		options := getSingleFileOptions(createTestFile(t, testFileContent))
		options.Architecture = ""

		_, _, err := containerimagehandler.CreateSingleFileImage(ctx, options)
		require.Error(t, err)
	})

	t.Run("non existing source file", func(t *testing.T) {
		ctx := getCtx()

		options := getSingleFileOptions(filepath.Join(t.TempDir(), "does-not-exist.txt"))

		_, _, err := containerimagehandler.CreateSingleFileImage(ctx, options)
		require.Error(t, err)
	})

	t.Run("source is a directory", func(t *testing.T) {
		ctx := getCtx()

		options := getSingleFileOptions(t.TempDir())

		_, _, err := containerimagehandler.CreateSingleFileImage(ctx, options)
		require.Error(t, err)
	})

	t.Run("text file fails static link check", func(t *testing.T) {
		ctx := getCtx()

		options := getSingleFileOptions(createTestFile(t, testFileContent))
		options.SkipStaticLinkCheck = false

		_, _, err := containerimagehandler.CreateSingleFileImage(ctx, options)
		require.Error(t, err)
	})
}

func TestCreateSingleFileArchive(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		ctx := getCtx()

		archivePath := createTestArchive(ctx, t)

		t.Run("tags", func(t *testing.T) {
			imageNamesAndTags, err := containerimagehandler.ListImageNamesAndTagsInArchive(ctx, archivePath)
			require.NoError(t, err)
			require.EqualValues(t, []string{testImageNameAndTag}, imageNamesAndTags)
		})

		t.Run("files", func(t *testing.T) {
			files, err := containerimagehandler.ListFilesInArchive(ctx, archivePath)
			require.NoError(t, err)
			require.EqualValues(t, []string{testPathInImage}, files)
		})

		t.Run("content", func(t *testing.T) {
			content, err := containerimagehandler.ReadFileFromArchiveAsString(ctx, archivePath, testPathInImage)
			require.NoError(t, err)
			require.EqualValues(t, testFileContent, content)
		})

		t.Run("platform", func(t *testing.T) {
			image, _, err := containerimagehandler.LoadImageFromArchive(ctx, archivePath, "")
			require.NoError(t, err)

			configFile, err := image.ConfigFile()
			require.NoError(t, err)
			require.EqualValues(t, "amd64", configFile.Architecture)
			require.EqualValues(t, "linux", configFile.OS)
		})
	})

	t.Run("idempotent overwrite with new content", func(t *testing.T) {
		ctx := getCtx()

		archivePath := getTestArchivePath(t)

		err := containerimagehandler.CreateSingleFileArchive(ctx, archivePath, getSingleFileOptions(createTestFile(t, "first\n")))
		require.NoError(t, err)

		err = containerimagehandler.CreateSingleFileArchive(ctx, archivePath, getSingleFileOptions(createTestFile(t, "second\n")))
		require.NoError(t, err)

		content, err := containerimagehandler.ReadFileFromArchiveAsString(ctx, archivePath, testPathInImage)
		require.NoError(t, err)
		require.EqualValues(t, "second\n", content)

		imageNamesAndTags, err := containerimagehandler.ListImageNamesAndTagsInArchive(ctx, archivePath)
		require.NoError(t, err)
		require.EqualValues(t, []string{testImageNameAndTag}, imageNamesAndTags)
	})

	t.Run("no temporary archive left behind", func(t *testing.T) {
		ctx := getCtx()

		archivePath := createTestArchive(ctx, t)

		_, err := os.Stat(archivePath + ".tmp")
		require.ErrorIs(t, err, os.ErrNotExist)
	})

	t.Run("empty output path", func(t *testing.T) {
		ctx := getCtx()

		err := containerimagehandler.CreateSingleFileArchive(ctx, "", getSingleFileOptions(createTestFile(t, testFileContent)))
		require.Error(t, err)
	})

	t.Run("nil options", func(t *testing.T) {
		ctx := getCtx()

		err := containerimagehandler.CreateSingleFileArchive(ctx, getTestArchivePath(t), nil)
		require.Error(t, err)
	})

	t.Run("failed creation does not create archive", func(t *testing.T) {
		ctx := getCtx()

		archivePath := getTestArchivePath(t)
		options := getSingleFileOptions(filepath.Join(t.TempDir(), "does-not-exist.txt"))

		err := containerimagehandler.CreateSingleFileArchive(ctx, archivePath, options)
		require.Error(t, err)

		_, err = os.Stat(archivePath)
		require.ErrorIs(t, err, os.ErrNotExist)
	})
}

func TestLoadImageFromArchive(t *testing.T) {
	t.Run("auto detect tag", func(t *testing.T) {
		ctx := getCtx()

		archivePath := createTestArchive(ctx, t)

		image, tag, err := containerimagehandler.LoadImageFromArchive(ctx, archivePath, "")
		require.NoError(t, err)
		require.NotNil(t, image)
		require.EqualValues(t, testImageNameAndTag, tag.String())
	})

	t.Run("explicit tag", func(t *testing.T) {
		ctx := getCtx()

		archivePath := createTestArchive(ctx, t)

		image, tag, err := containerimagehandler.LoadImageFromArchive(ctx, archivePath, testImageNameAndTag)
		require.NoError(t, err)
		require.NotNil(t, image)
		require.EqualValues(t, testImageNameAndTag, tag.String())
	})

	t.Run("tag not in archive", func(t *testing.T) {
		ctx := getCtx()

		archivePath := createTestArchive(ctx, t)

		_, _, err := containerimagehandler.LoadImageFromArchive(ctx, archivePath, "not-in-archive:v1")
		require.Error(t, err)
	})

	t.Run("empty archive path", func(t *testing.T) {
		ctx := getCtx()

		_, _, err := containerimagehandler.LoadImageFromArchive(ctx, "", "")
		require.Error(t, err)
	})
}

func TestOverwriteArchive(t *testing.T) {
	t.Run("write in memory image", func(t *testing.T) {
		ctx := getCtx()

		image, tag, err := containerimagehandler.CreateSingleFileImage(ctx, getSingleFileOptions(createTestFile(t, testFileContent)))
		require.NoError(t, err)

		archivePath := getTestArchivePath(t)
		err = containerimagehandler.OverwriteArchive(ctx, archivePath, tag, image)
		require.NoError(t, err)

		content, err := containerimagehandler.ReadFileFromArchiveAsString(ctx, archivePath, testPathInImage)
		require.NoError(t, err)
		require.EqualValues(t, testFileContent, content)
	})

	t.Run("empty archive path", func(t *testing.T) {
		ctx := getCtx()

		image, tag, err := containerimagehandler.CreateSingleFileImage(ctx, getSingleFileOptions(createTestFile(t, testFileContent)))
		require.NoError(t, err)

		err = containerimagehandler.OverwriteArchive(ctx, "", tag, image)
		require.Error(t, err)
	})

	t.Run("nil tag", func(t *testing.T) {
		ctx := getCtx()

		image, _, err := containerimagehandler.CreateSingleFileImage(ctx, getSingleFileOptions(createTestFile(t, testFileContent)))
		require.NoError(t, err)

		err = containerimagehandler.OverwriteArchive(ctx, getTestArchivePath(t), nil, image)
		require.Error(t, err)
	})

	t.Run("nil image", func(t *testing.T) {
		ctx := getCtx()

		_, tag, err := containerimagehandler.CreateSingleFileImage(ctx, getSingleFileOptions(createTestFile(t, testFileContent)))
		require.NoError(t, err)

		err = containerimagehandler.OverwriteArchive(ctx, getTestArchivePath(t), tag, nil)
		require.Error(t, err)
	})
}

func TestFileInArchiveExists(t *testing.T) {
	t.Run("existing file", func(t *testing.T) {
		ctx := getCtx()

		archivePath := createTestArchive(ctx, t)

		exists, err := containerimagehandler.FileInArchiveExists(ctx, archivePath, testPathInImage)
		require.NoError(t, err)
		require.True(t, exists)
	})

	t.Run("non existing file", func(t *testing.T) {
		ctx := getCtx()

		archivePath := createTestArchive(ctx, t)

		exists, err := containerimagehandler.FileInArchiveExists(ctx, archivePath, "does-not-exist.txt")
		require.NoError(t, err)
		require.False(t, exists)
	})

	t.Run("empty archive path", func(t *testing.T) {
		ctx := getCtx()

		_, err := containerimagehandler.FileInArchiveExists(ctx, "", testPathInImage)
		require.Error(t, err)
	})

	t.Run("empty path to check", func(t *testing.T) {
		ctx := getCtx()

		archivePath := createTestArchive(ctx, t)

		_, err := containerimagehandler.FileInArchiveExists(ctx, archivePath, "")
		require.Error(t, err)
	})
}

func TestReadFileFromArchive(t *testing.T) {
	t.Run("as string", func(t *testing.T) {
		ctx := getCtx()

		archivePath := createTestArchive(ctx, t)

		content, err := containerimagehandler.ReadFileFromArchiveAsString(ctx, archivePath, testPathInImage)
		require.NoError(t, err)
		require.EqualValues(t, testFileContent, content)
	})

	t.Run("as bytes", func(t *testing.T) {
		ctx := getCtx()

		archivePath := createTestArchive(ctx, t)

		content, err := containerimagehandler.ReadFileFromArchiveAsBytes(ctx, archivePath, testPathInImage)
		require.NoError(t, err)
		require.EqualValues(t, []byte(testFileContent), content)
	})

	t.Run("non existing file", func(t *testing.T) {
		ctx := getCtx()

		archivePath := createTestArchive(ctx, t)

		_, err := containerimagehandler.ReadFileFromArchiveAsString(ctx, archivePath, "does-not-exist.txt")
		require.Error(t, err)
	})

	t.Run("empty archive path", func(t *testing.T) {
		ctx := getCtx()

		_, err := containerimagehandler.ReadFileFromArchiveAsBytes(ctx, "", testPathInImage)
		require.Error(t, err)
	})

	t.Run("empty path to read", func(t *testing.T) {
		ctx := getCtx()

		archivePath := createTestArchive(ctx, t)

		_, err := containerimagehandler.ReadFileFromArchiveAsBytes(ctx, archivePath, "")
		require.Error(t, err)
	})
}

func TestDeleteFileInArchive(t *testing.T) {
	getDeleteOptions := func() *containeroptions.DeleteFileFromImageOptions {
		return &containeroptions.DeleteFileFromImageOptions{
			PathInImage:            testPathInImage,
			NewImageNameAndTag:     testImageNameAndTag,
			OverwriteSourceArchive: true,
		}
	}

	t.Run("delete existing file", func(t *testing.T) {
		ctx := getCtx()

		archivePath := createTestArchive(ctx, t)

		err := containerimagehandler.DeleteFileInArchive(ctx, archivePath, getDeleteOptions())
		require.NoError(t, err)

		t.Run("file does not exist anymore", func(t *testing.T) {
			exists, err := containerimagehandler.FileInArchiveExists(ctx, archivePath, testPathInImage)
			require.NoError(t, err)
			require.False(t, exists)
		})

		t.Run("reading the file fails", func(t *testing.T) {
			_, err := containerimagehandler.ReadFileFromArchiveAsString(ctx, archivePath, testPathInImage)
			require.Error(t, err)
		})

		t.Run("whiteout layer was added", func(t *testing.T) {
			image, _, err := containerimagehandler.LoadImageFromArchive(ctx, archivePath, "")
			require.NoError(t, err)

			layers, err := image.Layers()
			require.NoError(t, err)
			require.Len(t, layers, 2)
		})

		t.Run("tag is kept", func(t *testing.T) {
			imageNamesAndTags, err := containerimagehandler.ListImageNamesAndTagsInArchive(ctx, archivePath)
			require.NoError(t, err)
			require.EqualValues(t, []string{testImageNameAndTag}, imageNamesAndTags)
		})
	})

	t.Run("not overwriting source archive is not implemented", func(t *testing.T) {
		ctx := getCtx()

		archivePath := createTestArchive(ctx, t)

		options := getDeleteOptions()
		options.OverwriteSourceArchive = false

		err := containerimagehandler.DeleteFileInArchive(ctx, archivePath, options)
		require.Error(t, err)

		// The archive must be untouched.
		exists, err := containerimagehandler.FileInArchiveExists(ctx, archivePath, testPathInImage)
		require.NoError(t, err)
		require.True(t, exists)
	})

	t.Run("empty archive path", func(t *testing.T) {
		ctx := getCtx()

		err := containerimagehandler.DeleteFileInArchive(ctx, "", getDeleteOptions())
		require.Error(t, err)
	})

	t.Run("nil options", func(t *testing.T) {
		ctx := getCtx()

		archivePath := createTestArchive(ctx, t)

		err := containerimagehandler.DeleteFileInArchive(ctx, archivePath, nil)
		require.Error(t, err)
	})
}
