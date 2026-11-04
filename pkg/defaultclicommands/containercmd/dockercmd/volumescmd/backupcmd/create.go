package backupcmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewCreateCmd() *cobra.Command {
	const short = "Backup a docker volume to a tar archive."

	cmd := &cobra.Command{
		Use:   "create",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` container docker volumes backup create --volume <volume-name> --archive <archive-path>

Usage example:
    ` + os.Args[0] + ` container docker volumes backup create --volume my-volume --archive /backups/my-volume.tar.gz
`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			volumeName := mustutils.Must(cmd.Flags().GetString("volume"))
			archivePath := mustutils.Must(cmd.Flags().GetString("archive"))

			mustutils.Must0(dockerutils.BackupVolume(ctx, volumeName, archivePath))
		},
	}

	cmd.Flags().String("volume", "", "Name of the docker volume to backup")
	cmd.Flags().String("archive", "", "Path to the tar archive file to create")

	mustutils.Must0(cmd.MarkFlagRequired("volume"))
	mustutils.Must0(cmd.MarkFlagRequired("archive"))

	return cmd
}
