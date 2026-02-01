package exoscalednscmd

import (
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/cloudutils/exoscaleutils/exoscalenativeclient"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/datetime/durationparser"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewCreateRecordWithPublicIp() *cobra.Command {
	const short = "Create a DNS domain record with the current public IP. Usefull for creating or updating dynamic DNS entries when the public IP changes from time to time."

	cmd := &cobra.Command{
		Use:   "create-record-with-public-ip",
		Short: short,
		Long: short + `
		
To create or update the record "dynamic.example.com" use:
	` + os.Args[0] + ` cloud exoscale dns create-record-with-public-ip --verbose --domain="example.com" --record="dynamic"
	
To create and update the record "dynamic.example.com" every 30seconds use:
	` + os.Args[0] + ` cloud exoscale dns create-record-with-public-ip --verbose --domain="example.com" --record="dynamic" --interval="30seconds"
`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			domain, err := cmd.Flags().GetString("domain")
			if err != nil {
				logging.LogGoErrorFatal(err)
			}

			if domain == "" {
				logging.LogFatal("Please specify --domain.")
			}

			record, err := cmd.Flags().GetString("record")
			if err != nil {
				logging.LogGoErrorFatal(err)
			}

			if record == "" {
				logging.LogFatal("Please specify --record.")
			}

			interval, err := cmd.Flags().GetString("interval")
			if err != nil {
				logging.LogGoErrorFatal(err)
			}

			client := mustutils.Must(exoscalenativeclient.NewNativeClientFromEnvVars(ctx))
			if interval == "" {
				mustutils.Must0(exoscalenativeclient.CreateDnsDomainRecordWithCurrentPublicAddress(ctx, client, domain, record))
			} else {
				sleepDuration := mustutils.Must(durationparser.ToSecondsAsTimeDuration(interval))

				for {
					mustutils.Must0(exoscalenativeclient.CreateDnsDomainRecordWithCurrentPublicAddress(ctx, client, domain, record))
					logging.LogInfoByCtxf(ctx, "Wait '%s'='%s' before next update of the public IP address.", interval, *sleepDuration)
					time.Sleep(*sleepDuration)
				}
			}
			logging.LogGoodByCtxf(ctx, "Create Exoscale DNS domain record '%s' of domain '%s' with current public IP address finished.", record, domain)
		},
	}

	cmd.Flags().String("domain", "", "The domain where the record should be created.")
	cmd.Flags().String("record", "", "Name of the record to create.")
	cmd.Flags().String("interval", "", "If set this command repetatively evaluates the current public IP address and updates the record. Usefull for keeping the record up to date with when the public IP address changes.")

	return cmd
}
