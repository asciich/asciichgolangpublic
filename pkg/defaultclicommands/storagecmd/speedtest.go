package storagecmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/datatypes/bytesutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/storage/s3/storageutils"
)

func NewSpeedTestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "speed-test",
		Short: "Perform a storage speed test by writing and reading the --file of given --size.",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			path, err := cmd.Flags().GetString("file")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			if path == "" {
				logging.LogFatal("Please specify --file")
			}

			size, err := cmd.Flags().GetString("size")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			if size == "" {
				logging.LogFatal("Please specify --size")
			}

			sizeBytes := mustutils.Must(bytesutils.ParseSizeStringAsInt64(size))
			result := mustutils.Must(storageutils.RunSpeedTest(ctx, path, sizeBytes))

			logging.LogGoodByCtx(ctx, mustutils.Must(result.GetResultMessage()))
		},
	}

	cmd.Flags().String("file", "", "Path of the file to write and read. Be aware: This file is overwritten and deleted at the end of the test.")
	cmd.Flags().String("size", "", "File size to test with. Use human readable notation like '1MB' or '1GB'.")

	return cmd
}
