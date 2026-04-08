package archivecmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/containerimagehandler"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/containeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/datatypes/pointerutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/osutils/unixfilepermissionsutils"
)

func NewAddFileCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-file",
		Short: "Add a file to a container image archive. This will create new file layer in the image archive.",
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

			mustutils.Must0(containerimagehandler.AddFileToArchive(
				ctx,
				archive,
				&containeroptions.AddFileToImageArchiveOptions{
					SourceFilePath:         srcPath,
					PathInImage:            pathInImage,
					OverwriteSourceArchive: true,
					NewImageNameAndTag:     newTag,
					Mode:                   pointerutils.ToInt64Pointer(int64(mode)),
				},
			))

			logging.LogGoodByCtxf(ctx, "File '%s' added to image archive '%s' as '%s' named and tagged as '%s'", srcPath, archive, pathInImage, newTag)
		},
	}

	cmd.Flags().String("archive", "", "Path to the archive.")
	cmd.Flags().String("src-path", "", "Path to the file to add.")
	cmd.Flags().String("path-in-image", "", "The path where the file will be added inside the image.")
	cmd.Flags().String("new-tag", "", "New image name and tag. Use '<name>:<tag>' .")
	cmd.Flags().String("mode", "u=rw,g=r,o=r", "File mode for the file to add.")

	return cmd
}
