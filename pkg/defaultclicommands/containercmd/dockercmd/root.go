package dockercmd

import (
	"github.com/spf13/cobra"
	"os"
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
	)

	return cmd
}
