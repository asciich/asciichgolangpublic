package backupcmd

import "github.com/spf13/cobra"

func NewBackupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "backup",
		Short: "Backup docker volumes",
	}

	cmd.AddCommand(
		NewCreateCmd(),
		NewCreateMultiCmd(),
		NewRestoreCmd(),
		NewRestoreMultiCmd(),
	)

	return cmd
}
