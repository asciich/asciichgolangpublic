package exoscalednscmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/cloudutils/exoscaleutils/exoscalenativeclientoo"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewListDomainsCmd() *cobra.Command {
	const short = "List all DNS domains."

	cmd := &cobra.Command{
		Use:   "list-domains",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` cloud exoscale exoscaledns list-domains`,

		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			client := mustutils.Must(exoscalenativeclientoo.NewNativeClientFromEnvVars(ctx))

			dns := mustutils.Must(client.DNS())

			for _, domain := range mustutils.Must(dns.ListDomainNames(ctx)) {
				fmt.Println(domain)
			}

			logging.LogGoodByCtx(ctx, "List DNS domains finished.")
		},
	}

	return cmd
}
