package dockercmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
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
