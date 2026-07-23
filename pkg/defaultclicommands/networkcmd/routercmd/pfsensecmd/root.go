package pfsensecmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/routerutils/pfsenseutils"
)

func NewPfSenseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pfsense",
		Short: "pfsense related commands.",
	}

	AddSubCommandsAndPersistentFlags(cmd)

	return cmd
}

func AddSubCommands(cmd *cobra.Command) {
	cmd.AddCommand(
		NewGetSystemNameCmd(),
	)
}

func AddPersistentFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().String("url", "", "Url of the pfSense router to manage.")
	cmd.PersistentFlags().String("username", "", "Username for login into pfSense. For the passowrd use the environment variable '"+pfsenseutils.ENV_VAR_NAME_PFSENSE_PASSWORD+"'.")
}

func AddSubCommandsAndPersistentFlags(cmd *cobra.Command) {
	AddSubCommands(cmd)
	AddPersistentFlags(cmd)
}
