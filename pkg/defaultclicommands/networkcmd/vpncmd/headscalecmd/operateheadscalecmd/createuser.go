package operateheadscalecmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewCreateUserCmd(options *OperateOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use: "create-user",
		Short: "Create a headscale user.",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			if len(args) != 1 {
				logging.LogFatal("Please specify exactly one username to create.")
			}

			userName := args[0]

			mustutils.Must0(options.GetHeadScale(cmd).CreateUser(ctx, userName))

			logging.LogGoodByCtxf(ctx, "Created headscale user '%s'.", userName)
		},
	}

	return cmd
}