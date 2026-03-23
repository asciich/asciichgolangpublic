package testsuitecmd

import "github.com/spf13/cobra"

func NewTestSuiteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test-suite",
		Short: "Run test suites",
	}

	cmd.AddCommand(
		NewRunCmd(),
	)

	return cmd
}
