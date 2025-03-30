package loggingcmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/defaultclicommands/loggingcmd/loggingexamplescmd"
)

func NewLoggingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "logging",
		Short: "logging related commands.",
	}

	cmd.AddCommand(loggingexamplescmd.NewLoggingExamplesCmd())

	return cmd
}