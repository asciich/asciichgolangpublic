package userscmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/storagecmd/s3cmd/miniocmd/miniocmdoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/storage/s3/nativeminioclient"
)

func NewDeleteUserCmd(options *miniocmdoptions.MinioCmdOptions) *cobra.Command {
	const short = "Delete a user."

	cmd := &cobra.Command{
		Use:   "delete",
		Short: short,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			adminClient := options.GetAdminClient(ctx, cmd)

			userName := args[0]

			mustutils.Must0(nativeminioclient.DeleteUser(ctx, adminClient, userName))

			logging.LogInfoByCtxf(ctx, "Delete user '%s' finished.", userName)
		},
	}

	return cmd
}
