package operateheadscalecmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewListUsersCmd(options *OperateOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-users",
		Short: "List headscale users.",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			for _, user := range mustutils.Must(options.GetHeadScale(cmd).ListUserNames(ctx)) {
				fmt.Println(user)
			}

			logging.LogGoodByCtxf(ctx, "List headscale users finished")
		},
	}

	return cmd
}
