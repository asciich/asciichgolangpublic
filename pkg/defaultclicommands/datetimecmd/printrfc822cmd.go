package datetimecmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

func NewPrintRfc822Cmd() *cobra.Command {
	const short = "Print the current date and time in RFC822 format."

	cmd := &cobra.Command{
		Use:   "print-rfc822",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` datetime print-rfc822`,

		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(time.Now().Format(time.RFC822))
		},
	}

	return cmd
}
