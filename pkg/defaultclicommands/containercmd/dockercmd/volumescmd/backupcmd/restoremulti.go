package backupcmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewRestoreMultiCmd() *cobra.Command {
	const short = "Restore all docker volumes contained in a multi volume tar archive."

	cmd := &cobra.Command{
		Use:   "restore-multi",
		Short: short,
		Long: short + `

The volumes to restore are taken from the manifest.json inside the archive.
Volumes which do not exist yet are created. If the archive has no manifest, an error is returned and nothing is restored.

Usage:
    ` + os.Args[0] + ` container docker volumes backup restore-multi --archive <archive-path>

Usage example:
    ` + os.Args[0] + ` container docker volumes backup restore-multi --archive /backups/volumes.tar.gz
`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			archivePath := mustutils.Must(cmd.Flags().GetString("archive"))

			mustutils.Must0(dockerutils.RestoreVolumes(ctx, archivePath))
		},
	}

	cmd.Flags().String("archive", "", "Path to the multi volume tar archive file to restore from")

	mustutils.Must0(cmd.MarkFlagRequired("archive"))

	return cmd
}
