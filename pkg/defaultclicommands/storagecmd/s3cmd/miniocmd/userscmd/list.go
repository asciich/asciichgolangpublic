package userscmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/storagecmd/s3cmd/miniocmd/miniocmdoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/storage/s3/nativeminioclient"
)

func NewListUsersCmd(options *miniocmdoptions.MinioCmdOptions) *cobra.Command {
	const short = "List all users."

	cmd := &cobra.Command{
		Use:   "list",
		Short: short,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			adminClient := options.GetAdminClient(ctx, cmd)

			for _, userName := range mustutils.Must(nativeminioclient.ListUserNames(ctx, adminClient)) {
				fmt.Println(userName)
			}

			logging.LogInfoByCtxf(ctx, "List minio users finished.")
		},
	}

	return cmd
}
