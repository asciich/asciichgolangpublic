package dockercmd

import "github.com/spf13/cobra"

func NewDockerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "docker",
		Short: "docker related commands.",
	}

	cmd.AddCommand(
		NewListContainerNames(),
	)

	return cmd
}
