package storagecmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/storage/s3/storageutils"
)

func NewSyncCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "sync",
		Short: "Flush OS write cache to the storage. Same as the the 'sync' CLI command.",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			mustutils.Must0(storageutils.Sync(ctx))

			logging.LogGoodByCtxf(ctx, "Flushed OS write cache to the storage.")
		},
	}

	return cmd
}