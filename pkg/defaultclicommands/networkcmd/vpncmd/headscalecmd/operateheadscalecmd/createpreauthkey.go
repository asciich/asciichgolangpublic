package operateheadscalecmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/datatypes/stringsutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewCreatePreauthKeyCmd(options *OperateOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-preauth-key",
		Short: "Create a preauth key for a the specified headscale user.",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			if len(args) != 1 {
				logging.LogFatal("Please specify exaclty one user name to create a preauth key.")
			}

			userName := args[0]

			preauthKey := mustutils.Must(options.GetHeadScale(cmd).GeneratePreauthKeyForUser(ctx, userName))
			fmt.Print(stringsutils.EnsureEndsWithExactlyOneLineBreak(preauthKey))

			logging.LogGoodByCtxf(ctx, "Generated preauth key for user '%s'.", userName)
		},
	}

	return cmd
}
