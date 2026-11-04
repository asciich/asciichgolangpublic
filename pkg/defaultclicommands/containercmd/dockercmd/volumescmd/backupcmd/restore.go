package backupcmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewRestoreCmd() *cobra.Command {
	const short = "Restore a docker volume from a tar archive."

	cmd := &cobra.Command{
		Use:   "restore",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` container docker volumes restore --volume <volume-name> --archive <archive-path>

Usage example:
    ` + os.Args[0] + ` container docker volumes restore --volume my-volume --archive /backups/my-volume.tar.gz
`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			volumeName := mustutils.Must(cmd.Flags().GetString("volume"))
			archivePath := mustutils.Must(cmd.Flags().GetString("archive"))

			mustutils.Must0(dockerutils.RestoreVolume(ctx, volumeName, archivePath))
		},
	}

	cmd.Flags().String("volume", "", "Name of the docker volume to restore")
	cmd.Flags().String("archive", "", "Path to the tar archive file to restore from")

	mustutils.Must0(cmd.MarkFlagRequired("volume"))
	mustutils.Must0(cmd.MarkFlagRequired("archive"))

	return cmd
}
