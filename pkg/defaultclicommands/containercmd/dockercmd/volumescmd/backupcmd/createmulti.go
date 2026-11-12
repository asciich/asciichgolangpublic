package backupcmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewCreateMultiCmd() *cobra.Command {
	const short = "Backup multiple docker volumes into one tar archive."

	cmd := &cobra.Command{
		Use:   "create-multi",
		Short: short,
		Long: short + `

The archive contains a manifest.json and the content of every volume.
Use '... backup restore-multi' to restore all volumes of the archive at once.

Usage:
    ` + os.Args[0] + ` container docker volumes backup create-multi (--volume <volume-name>... | --all) --archive <archive-path>

Usage example:
    ` + os.Args[0] + ` container docker volumes backup create-multi --volume vol-a --volume vol-b --archive /backups/volumes.tar.gz
    ` + os.Args[0] + ` container docker volumes backup create-multi --all --archive /backups/all-volumes.tar.gz
`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			volumeNames := mustutils.Must(cmd.Flags().GetStringArray("volume"))
			all := mustutils.Must(cmd.Flags().GetBool("all"))
			archivePath := mustutils.Must(cmd.Flags().GetString("archive"))

			if all {
				volumeNames = mustutils.Must(dockerutils.ListVolumesNames(ctx))
			}

			mustutils.Must0(dockerutils.BackupVolumes(ctx, volumeNames, archivePath))
		},
	}

	cmd.Flags().StringArray("volume", nil, "Name of a docker volume to backup. Can be specified multiple times")
	cmd.Flags().Bool("all", false, "Backup all docker volumes")
	cmd.Flags().String("archive", "", "Path to the tar archive file to create")

	mustutils.Must0(cmd.MarkFlagRequired("archive"))
	cmd.MarkFlagsMutuallyExclusive("volume", "all")
	cmd.MarkFlagsOneRequired("volume", "all")

	return cmd
}
