package dockercmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/contextutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/dockerutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/mustutils"
)

func NewListContainerNames() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-container-names",
		Short: "List the names of all found containers.",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			for _, name := range mustutils.Must(dockerutils.ListContainerNames(ctx)) {
				fmt.Println(name)
			}
		},
	}

	return cmd
}
