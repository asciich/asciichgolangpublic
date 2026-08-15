package testsuitecmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/testutils/testsuite"
	"github.com/asciich/asciichgolangpublic/pkg/testutils/testutilsoptions"
)

func NewRunCmd() *cobra.Command {
	const short = "Run a test suite."

	cmd := &cobra.Command{
		Use:   "run",
		Short: short,
		Long: short + `

Run one or more test suite files specified as arguments.

Each provided file path is executed sequentially. The test results are logged
after each suite completes. If any test suite fails, execution stops immediately
with a fatal error.

Examples:
  # Run a single test suite file
  testsuite run ./my_tests.yaml

  # Run multiple test suite files
  testsuite run ./suite1.yaml ./suite2.yaml

Arguments:
  At least one file path to a test suite definition file must be provided.

Exit behavior:
  - Exits with a fatal error if no arguments are provided.
  - Exits with a fatal error if any test suite does not pass.
  - Logs a success message if all test suites pass.
`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			if len(args) <= 0 {
				logging.LogFatal("Please specify at least one test suite file.")
			}

			for _, f := range args {
				result := mustutils.Must(testsuite.RunFromFilePath(ctx, f, &testutilsoptions.RunTestSuiteOptions{}))
				mustutils.Must0(result.LogResult(ctx))

				if !mustutils.Must(result.IsPassed(ctx)) {
					logging.LogFatal("Test suite failed.")
				}
			}

			logging.LogGoodByCtxf(ctx, "All tests in '%v' passed.", args)
		},
	}

	return cmd
}
