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
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run a test suite.",
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
		},
	}

	return cmd
}
