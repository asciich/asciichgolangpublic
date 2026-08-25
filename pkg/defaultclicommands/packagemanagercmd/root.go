package packagemanagercmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/packagemanagercmd/yaycmd"
	"os"
)

func NewPackageManagerCmd() *cobra.Command {
	const short = "Packagemanager related commmands"

	cmd := &cobra.Command{
		Use:   "packagemanager",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` packagemanager packagemanager`,
	}

	cmd.AddCommand(
		yaycmd.NewYayCmd(),
	)

	return cmd
}
