package operateheadscalecmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewGetUserIdCmd(options *OperateOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-user-id",
		Short: "Get the user ID of the given user name.",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			if len(args) != 1 {
				logging.LogFatal("Please specify exaclty one user name to get the headscale ID.")
			}

			userName := args[0]

			id := mustutils.Must(options.GetHeadScale(cmd).GetUserId(ctx, userName))
			fmt.Printf("%d\n", id)

			logging.LogGoodByCtxf(ctx, "Got user id '%d' for headscale user '%s'.", id, userName)
		},
	}

	return cmd
}
