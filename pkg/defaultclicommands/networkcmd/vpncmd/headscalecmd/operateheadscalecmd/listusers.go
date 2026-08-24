package operateheadscalecmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewListUsersCmd(options *OperateOptions) *cobra.Command {
	const short = "List headscale users."

	cmd := &cobra.Command{
		Use:   "list-users",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` network vpn headscale operate list-users`,

		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			for _, user := range mustutils.Must(options.GetHeadScale(ctx, cmd).ListUserNames(ctx)) {
				fmt.Println(user)
			}

			logging.LogGoodByCtxf(ctx, "List headscale users finished")
		},
	}

	return cmd
}
