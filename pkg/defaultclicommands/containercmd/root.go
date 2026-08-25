package containercmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/containercmd/dockercmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/containercmd/imagescmd"
	"os"
)

func NewContainerCmd() *cobra.Command {
	const short = "Container related commands. Includes image handling but also handling docker."

	cmd := &cobra.Command{
		Use:   "container",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` container container`,
	}

	cmd.AddCommand(
		dockercmd.NewDockerCmd(),
		imagescmd.NewImagesCmd(),
	)

	return cmd

}
