package exoscaleuserscmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/cloudutils/exoscaleutils/exoscalenativeclient"
	"github.com/asciich/asciichgolangpublic/pkg/cloudutils/exoscaleutils/exoscalenativeclientoo"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewListCmd() *cobra.Command {
	const short = "List all users in the Exoscale account."

	cmd := &cobra.Command{
		Use:   "list",
		Short: short,
		Long: short + `

Usage:
	# Export the credentails as ` + exoscalenativeclient.ENV_VAR_EXOSCALE_API_KEY + `and ` + exoscalenativeclient.ENV_VAR_EXOSCALE_API_SECRET + `
	` + os.Args[0] + ` cloud exoscale iam users list`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			client := mustutils.Must(exoscalenativeclientoo.NewNativeClientFromEnvVars(ctx))

			iam := mustutils.Must(client.IAM())
			users := mustutils.Must(iam.Users())

			for _, u := range mustutils.Must(users.ListUserNames(ctx)) {
				fmt.Println(u)
			}

			logging.LogGoodByCtxf(ctx, "List exoscale users finished.")
		},
	}

	return cmd
}
