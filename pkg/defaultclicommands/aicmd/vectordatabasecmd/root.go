package vectordatabasecmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/aicmd/vectordatabasecmd/chromacmd"
	"os"
)

func NewVectorDatabaseCmd() *cobra.Command {
	const short = "Vector database related commands."

	cmd := &cobra.Command{
		Use:   "vector-database",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` ai vectordatabase vector-database`,
	}

	cmd.AddCommand(
		chromacmd.NewChromaCmd(),
	)

	return cmd
}
