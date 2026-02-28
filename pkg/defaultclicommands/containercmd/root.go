package containercmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/containercmd/dockercmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/containercmd/imagescmd"
)

func NewContainerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "container",
		Short: "Container related commands. Includes image handling but also handling docker.",
	}

	cmd.AddCommand(
		dockercmd.NewDockerCmd(),
		imagescmd.NewImagesCmd(),
	)

	return cmd

}
