package volumescmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewListFilesInVolumeCmd() *cobra.Command {
	const short = "List all files inside a docker volume."

	cmd := &cobra.Command{
		Use:   "list-files-in-volume",
		Short: short,
		Long: short + `

A temporary container is started to list all files inside the given docker volume.
Hidden files are included and the output is sorted alphabetically.

Usage:
    ` + os.Args[0] + ` container docker volumes list-files-in-volume --volume <volume-name>

Usage example:
    ` + os.Args[0] + ` container docker volumes list-files-in-volume --volume my-volume
`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			volumeName := mustutils.Must(cmd.Flags().GetString("volume"))

			for _, fileName := range mustutils.Must(dockerutils.ListVolumeFileNames(ctx, volumeName)) {
				fmt.Println(fileName)
			}
		},
	}

	cmd.Flags().String("volume", "", "Name of the docker volume to list the files from")

	mustutils.Must0(cmd.MarkFlagRequired("volume"))

	return cmd
}
