package testingcmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/testingcmd/testsuitecmd"
)

func NewTestingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "testing",
		Short: "Testing related commands",
	}

	cmd.AddCommand(
		testsuitecmd.NewTestSuiteCmd(),
	)

	return cmd
}
