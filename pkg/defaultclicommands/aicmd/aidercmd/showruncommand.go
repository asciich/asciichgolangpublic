package aidercmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/aiutils/aiderutils"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewShowRunCommandCmd() *cobra.Command {
	const short = "Shows the full CLI command to run aider in a docker container."

	cmd := &cobra.Command{
		Use:   "show-run-command",
		Short: short,
		Long: short + `

To directly run it use:
  eval $(` + os.Args[0] + ` ai aider show-run-command)
`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(strings.Join(mustutils.Must(aiderutils.GetRunCommand(false)), " "))
		},
	}

	return cmd
}
