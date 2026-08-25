package versioncmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/binaryinfo"
	"os"
)

func NewVersionCmd() *cobra.Command {
	const short = "Print the version information for this binary."

	var cmd = &cobra.Command{
		Use:   "version",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` version version`,

		Run: func(cmd *cobra.Command, args []string) {
			binaryinfo.PrintInfo()
		},
	}

	return cmd
}
