package testsuitecmd

import (
	"github.com/spf13/cobra"
	"os"
)

func NewTestSuiteCmd() *cobra.Command {
	const short = "Run test suites"

	cmd := &cobra.Command{
		Use:   "test-suite",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` testing testsuite test-suite`,
	}

	cmd.AddCommand(
		NewRunCmd(),
	)

	return cmd
}
