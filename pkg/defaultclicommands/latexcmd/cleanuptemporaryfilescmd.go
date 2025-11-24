package latexcmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
)

func NewCleanupTemporaryFilesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cleanup-temporary-files",
		Short: "Cleanup temporary latex files in given directory",
		Long:  "Cleanup temporary latex files in given directory",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			if len(args) != 1 {
				logging.LogFatal("Please specify exatly one directory to clean up")
			}

			directoryPath := args[0]

			cliLatexCleanupTemporaryFilesCmd(ctx, directoryPath)
		},
	}

	return cmd
}

func cliLatexCleanupTemporaryFilesCmd(ctx context.Context, directoryPath string) {
	dirToCleanUp := mustutils.Must(files.GetLocalDirectoryByPath(directoryPath))

	mustutils.Must0(dirToCleanUp.DeleteFilesMatching(
		ctx,
		&parameteroptions.ListFileOptions{
			MatchBasenamePattern: []string{".*\\.aux", ".*\\.log"},
		},
	))

	logging.LogGoodByCtxf(ctx, "Temporary latex files from '%s' deleted.", directoryPath)
}
