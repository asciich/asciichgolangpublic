package vectordatabasecmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/aicmd/vectordatabasecmd/chromacmd"
)

func NewVectorDatabaseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vector-database",
		Short: "Vector database related commands.",
	}

	cmd.AddCommand(
		chromacmd.NewChromaCmd(),
	)

	return cmd
}
