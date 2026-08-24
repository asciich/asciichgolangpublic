package testingcmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/testingcmd/testsuitecmd"
	"os"
)

func NewTestingCmd() *cobra.Command {
	const short = "Testing related commands"

	cmd := &cobra.Command{
		Use:   "testing",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` testing testing`,
	}

	cmd.AddCommand(
		testsuitecmd.NewTestSuiteCmd(),
	)

	return cmd
}
