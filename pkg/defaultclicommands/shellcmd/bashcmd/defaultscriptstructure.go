package bashcmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/bashutils"
)

func NewDefaultScriptStructureCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "default-script-structure",
		Short: "Print the default bash script structure. Useful as starting point for new bash scripts.",

		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(bashutils.GetDefaultScriptStructure())
		},
	}

	return cmd
}
