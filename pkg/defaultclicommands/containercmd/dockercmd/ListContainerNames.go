package dockercmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewListContainerNames() *cobra.Command {
	const short = "List the names of all found containers."

	cmd := &cobra.Command{
		Use:   "list-container-names",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` container docker list-container-names`,

		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			for _, name := range mustutils.Must(dockerutils.ListContainerNames(ctx)) {
				fmt.Println(name)
			}
		},
	}

	return cmd
}
