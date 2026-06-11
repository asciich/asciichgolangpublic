package datetimecmd

import "github.com/spf13/cobra"

func NewDateTimeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "datetime",
		Short: "Date or time related commands",
	}

	cmd.AddCommand(
		NewPrintRfc822Cmd(),
	)

	return cmd
}
