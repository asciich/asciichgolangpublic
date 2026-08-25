package testwebservercmd

import (
	"github.com/spf13/cobra"
	"os"
)

func NewTestWebServerCmd() *cobra.Command {
	const short = "A simple testwebserver providing different pages to test and experiment with a webserer."

	cmd := &cobra.Command{
		Use:   "testwebserver",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` http testwebserver testwebserver`,
	}

	cmd.AddCommand(
		NewRunCmd(),
	)

	return cmd
}
