package archivecmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/containerimagehandler"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/containeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/datatypes/pointerutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/osutils/unixfilepermissionsutils"
)

func NewCreateSingleFileArchive() *cobra.Command {
	const shortDescription = "Create a single file container image archive by adding only one file."
	const architectureHelp = "Usually 'amd64' or 'arm' (32bit) or 'arm64'"

	cmd := &cobra.Command{
		Use:   "create-single-file-archive",
		Short: shortDescription,
		Long: shortDescription + `

Can be used to pack a single statically linked binary into a container image with absolutely no other files included.

Example packing this binary itself into a container:
  1. Pack the binary into a container:
	` + os.Args[0] + ` container images archive create-single-file-archive --verbose --archive=pack-example-latest.tar --new-tag=pack-example:latest  --src-path="$(which ` + os.Args[0] + ` )" --path-in-image=/` + os.Args[0] + ` --mode="u=rwx,g=rx,o=rx" --architecture="amd64"
  2. Load the container:
	cat pack-example-latest.tar | docker load
  3. Run the binary in the container:
	docker run --rm -it pack-example:latest /` + os.Args[0] + `
`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			archive, err := cmd.Flags().GetString("archive")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			if archive == "" {
				logging.LogFatal("Please specify --archive .")
			}

			srcPath, err := cmd.Flags().GetString("src-path")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			if srcPath == "" {
				logging.LogFatal("Please specify --src-path of the file to add.")
			}

			pathInImage, err := cmd.Flags().GetString("path-in-image")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			if pathInImage == "" {
				logging.LogFatal("Please specify the --path-in-image .")
			}

			newTag, err := cmd.Flags().GetString("new-tag")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			if newTag == "" {
				logging.LogFatal("Please specify the --new-tag .")
			}

			modeString, err := cmd.Flags().GetString("mode")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			mode, err := unixfilepermissionsutils.GetPermissionsValue(modeString)
			if err != nil {
				logging.LogGoErrorFatal(err)
			}

			architecture, err := cmd.Flags().GetString("architecture")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			if architecture == "" {
				logging.LogFatal("Please specify --architecture. " + architectureHelp)
			}

			mustutils.Must0(containerimagehandler.CreateSingleFileArchive(
				ctx,
				archive,
				&containeroptions.CreateSingleFileArchiveOptions{
					SourceFilePath:     srcPath,
					PathInImage:        pathInImage,
					NewImageNameAndTag: newTag,
					Mode:               pointerutils.ToInt64Pointer(int64(mode)),
					Architecture:       architecture,
				},
			))

			logging.LogGoodByCtxf(ctx, "File '%s' added to image archive '%s' as '%s' named and tagged as '%s' for arch ='%s'", srcPath, archive, pathInImage, newTag, architecture)
		},
	}

	cmd.Flags().String("archive", "", "Path to the archive to create.")
	cmd.Flags().String("src-path", "", "Path to the file to add.")
	cmd.Flags().String("path-in-image", "", "The path where the file will be added inside the image.")
	cmd.Flags().String("new-tag", "", "New image name and tag. Use '<name>:<tag>' .")
	cmd.Flags().String("mode", "u=rw,g=r,o=r", "File mode for the file to add.")
	cmd.Flags().String("architecture", "", "Architecture of the container image. "+architectureHelp)

	return cmd
}
