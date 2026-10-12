package hostcmd

import (
	"github.com/spf13/cobra"

	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/hostcmd/hostnamecmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/hostcmd/hostreconcilercmd"
)

func NewHostCmd() *cobra.Command {
	const short = "Host management commands"

	cmd := &cobra.Command{
		Use:   "host",
		Short: short,
		Long:  short,
	}

	cmd.AddCommand(
		hostnamecmd.NewHostnameCmd(),
		hostreconcilercmd.NewHostReconcilerCmd(),
	)

	return cmd
}
