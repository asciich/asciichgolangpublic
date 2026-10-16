package imagescmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/containercmd/imagescmd/archivecmd"
)

func NewImagesCmd() *cobra.Command {
	const short = "Handle docker images."

	cmd := &cobra.Command{
		Use:   "images",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` container images images`,
	}

	cmd.AddCommand(
		NewCreateSingleFileDockerImage(),

		archivecmd.NewArchiveCmd(),
	)

	return cmd
}
