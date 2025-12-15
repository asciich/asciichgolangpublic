package testwebservercmd

import "github.com/spf13/cobra"

func NewTestWebServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "testwebserver",
		Short: "A simple testwebserver providing different pages to test and experiment with a webserer.",
	}

	cmd.AddCommand(
		NewRunCmd(),
	)

	return cmd
}
