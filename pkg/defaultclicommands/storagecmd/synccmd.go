package storagecmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/storage/storageutils"
	"os"
)

func NewSyncCmd() *cobra.Command {
	const short = "Flush OS write cache to the storage. Same as the 'sync' CLI command."

	cmd := &cobra.Command{
		Use:   "sync",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` storage sync`,

		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			mustutils.Must0(storageutils.Sync(ctx))

			logging.LogGoodByCtxf(ctx, "Flushed OS write cache to the storage.")
		},
	}

	return cmd
}
