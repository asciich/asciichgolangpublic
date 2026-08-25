package bashcmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/bashutils"
)

func NewDefaultScriptStructureCmd() *cobra.Command {
	const short = "Print the default bash script structure. Useful as starting point for new bash scripts."

	cmd := &cobra.Command{
		Use:   "default-script-structure",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` shell bash default-script-structure`,

		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(bashutils.GetDefaultScriptStructure())
		},
	}

	return cmd
}
