package datetimecmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

func NewPrintRfc822Cmd() *cobra.Command {
	const short = "Print the current date and time in RFC822 format."

	cmd := &cobra.Command{
		Use:   "print-rfc822",
		Short: short,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(time.Now().Format(time.RFC822))
		},
	}

	return cmd
}
