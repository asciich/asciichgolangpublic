package dockercmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/containercmd/dockercmd/volumescmd"
)

func NewDockerCmd() *cobra.Command {
	const short = "docker related commands."

	cmd := &cobra.Command{
		Use:   "docker",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` container docker docker`,
	}

	cmd.AddCommand(
		NewListContainerNames(),
		volumescmd.NewVolumesCmd(),
	)

	return cmd
}
