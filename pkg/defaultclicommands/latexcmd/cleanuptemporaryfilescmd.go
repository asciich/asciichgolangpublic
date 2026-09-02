package latexcmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefilesoo"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
)

func NewCleanupTemporaryFilesCmd() *cobra.Command {
	const short = "Cleanup temporary latex files in given directory"

	cmd := &cobra.Command{
		Use:   "cleanup-temporary-files",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` latex cleanup-temporary-files`,

		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			if len(args) != 1 {
				logging.LogFatal("Please specify exatly one directory to clean up")
			}

			directoryPath := args[0]

			dirToCleanUp := mustutils.Must(nativefilesoo.NewDirectoryByPath(directoryPath))

			mustutils.Must0(dirToCleanUp.DeleteFilesMatching(
				ctx,
				&parameteroptions.ListFileOptions{
					MatchBasenamePattern: []string{".*\\.aux", ".*\\.log"},
				},
			))

			logging.LogGoodByCtxf(ctx, "Temporary latex files from '%s' deleted.", directoryPath)
		},
	}

	return cmd
}
