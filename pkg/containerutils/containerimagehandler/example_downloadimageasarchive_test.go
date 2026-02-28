package containerimagehandler_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexec"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/containerimagehandler"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
)

// This example downloads the "alpine" linux container image to a local archive file.
//
// At the end we use docker load to validate the downloaded image is a container image.
func Test_DownloadImageAsArchive(t *testing.T) {
	// enable verbose output:
	ctx := contextutils.ContextVerbose()

	// define the temporary output path:
	outputPath, err := tempfiles.CreateTemporaryFile(ctx)
	require.NoError(t, err)
	defer nativefiles.Delete(ctx, outputPath, &filesoptions.DeleteOptions{})

	// Download the image as archive:
	err = containerimagehandler.DownloadImageAsArchive(ctx, "alpine:latest", outputPath)
	require.NoError(t, err)

	// Validate the downloaded archive is actually a image loadable by docker:
	stdout, err := commandexecutorexec.RunCommandAndGetStdoutAsString(ctx, &parameteroptions.RunCommandOptions{
		Command: []string{"bash", "-c", "cat '" + outputPath + "' | docker load"},
	})
	require.NoError(t, err)
	require.Contains(t, stdout, "image: alpine:latest")
}
