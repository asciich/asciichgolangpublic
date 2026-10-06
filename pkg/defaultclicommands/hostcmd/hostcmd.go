package hostcmd

import (
	"github.com/spf13/cobra"
	"os"

	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/hostcmd/hostnamecmd"
)

func NewHostCmd() *cobra.Command {
	const short = "Host management commands"

	cmd := &cobra.Command{
		Use:   "host",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` host <subcommand>`,
	}

	cmd.AddCommand(
		hostnamecmd.NewHostnameCmd(),
	)

	return cmd
}
