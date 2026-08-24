package operateheadscalecmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"os"
)

func NewCreateUserCmd(options *OperateOptions) *cobra.Command {
	const short = "Create a headscale user."

	cmd := &cobra.Command{
		Use:   "create-user",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` network vpn headscale operate create-user`,

		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			if len(args) != 1 {
				logging.LogFatal("Please specify exactly one username to create.")
			}

			userName := args[0]

			mustutils.Must0(options.GetHeadScale(ctx, cmd).CreateUser(ctx, userName))

			logging.LogGoodByCtxf(ctx, "Created headscale user '%s'.", userName)
		},
	}

	return cmd
}
