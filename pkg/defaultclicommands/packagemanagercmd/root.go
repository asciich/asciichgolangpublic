package packagemanagercmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/packagemanagercmd/yaycmd"
)

func NewPackageManagerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "packagemanager",
		Short: "Packagemanager related commmands",
	}

	cmd.AddCommand(
		yaycmd.NewYayCmd(),
	)

	return cmd
}
