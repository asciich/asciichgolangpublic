package loggingcmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/loggingcmd/loggingexamplescmd"
	"os"
)

func NewLoggingCmd() *cobra.Command {
	const short = "logging related commands."

	cmd := &cobra.Command{
		Use:   "logging",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` logging logging`,
	}

	cmd.AddCommand(loggingexamplescmd.NewLoggingExamplesCmd())

	return cmd
}
