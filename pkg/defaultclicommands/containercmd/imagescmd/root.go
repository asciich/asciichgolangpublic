package imagescmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/containercmd/imagescmd/archivecmd"
)

func NewImagesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "images",
		Short: "Handle docker images.",
	}

	cmd.AddCommand(
		archivecmd.NewArchiveCmd(),
	)

	return cmd
}
