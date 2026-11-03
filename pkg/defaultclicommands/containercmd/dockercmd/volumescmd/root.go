package volumescmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/containercmd/dockercmd/volumescmd/backupcmd"
)

func NewVolumesCmd() *cobra.Command {
	const short = "Docker volume related commands."

	cmd := &cobra.Command{
		Use:   "volumes",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` container docker volumes`,
	}

	cmd.AddCommand(
		backupcmd.NewBackupCmd(),

		NewListCmd(),
	)

	return cmd
}
