package imagescmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/asciich/asciichgolangpublic/pkg/containerutils/containeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/nativedocker"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/datatypes/pointerutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/osutils/unixfilepermissionsutils"
)

func NewCreateSingleFileDockerImage() *cobra.Command {
	const short = "Create a single file docker image."
	const architectureHelp = "Usually 'amd64' or 'arm' (32bit) or 'arm64'"

	cmd := &cobra.Command{
		Use:   "create-single-file-docker-image",
		Short: short,
		Long: short + `

Creates a docker image containing only the --src-path as file.
The generated image is directly stored in docker.

Can be used to pack a single statically linked binary into a container image with absolutely no other files included.

Example:
  ` + os.Args[0] + ` docker images create-single-file-docker-image --src-path ./mybinary --path-in-image /mybinary --new-tag mybinary:latest --architecture amd64
`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

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

			skipStaticLinkCheck, err := cmd.Flags().GetBool("skip-static-link-check")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			image := mustutils.Must(nativedocker.CreateSingleFileImage(
				ctx,
				&containeroptions.CreateSingleFileArchiveOptions{
					SourceFilePath:      srcPath,
					PathInImage:         pathInImage,
					NewImageNameAndTag:  newTag,
					Mode:                pointerutils.ToInt64Pointer(int64(mode)),
					Architecture:        architecture,
					SkipStaticLinkCheck: skipStaticLinkCheck,
				},
			))

			imageName := mustutils.Must(image.GetName())

			logging.LogGoodByCtxf(ctx, "File '%s' added as '%s' to docker image '%s' for arch='%s'.", srcPath, pathInImage, imageName, architecture)
		},
	}

	cmd.Flags().String("src-path", "", "Path to the file to add.")
	cmd.Flags().String("path-in-image", "", "The path where the file will be added inside the image.")
	cmd.Flags().String("new-tag", "", "New image name and tag. Use '<name>:<tag>' .")
	cmd.Flags().String("mode", "u=rw,g=r,o=r", "File mode for the file to add.")
	cmd.Flags().String("architecture", "", "Architecture of the container image. "+architectureHelp)
	cmd.Flags().Bool("skip-static-link-check", false, "Skip the check for statically linked binaries. Only use this if you know what you are doing.")

	return cmd
}
