package exoscaleuserscmd

import (
	"fmt"
	"os"

	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/cloudutils/exoscaleutils/exoscaleantiveclientoo"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/cloudutils/exoscaleutils/exoscalenativeclient"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/contextutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/logging"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/mustutils"

	"github.com/spf13/cobra"
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

			client := mustutils.Must(exoscaleantiveclientoo.NewNativeClientFromEnvVars(ctx))

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
