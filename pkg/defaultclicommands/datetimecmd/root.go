package datetimecmd

import (
	"github.com/spf13/cobra"
	"os"
)

func NewDateTimeCmd() *cobra.Command {
	const short = "Date or time related commands"

	cmd := &cobra.Command{
		Use:   "datetime",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` datetime datetime`,
	}

	cmd.AddCommand(
		NewPrintRfc822Cmd(),
	)

	return cmd
}
