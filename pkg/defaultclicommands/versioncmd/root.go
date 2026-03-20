package versioncmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/binaryinfo"
)

func NewVersionCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version information for this binary.",
		Run: func(cmd *cobra.Command, args []string) {
			binaryinfo.PrintInfo()
		},
	}

	return cmd
}
