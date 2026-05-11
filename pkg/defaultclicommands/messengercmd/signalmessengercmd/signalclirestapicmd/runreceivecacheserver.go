package signalclirestapicmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/messengerutils/signalmessengerutils/signalclirestapiutils"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewRunReceiveCacheServerCmd() *cobra.Command {
	const short = "Run a received message cache server. It will poll the Signal CLI rest API and cache the messages so you can receive the history multiple times."

	cmd := &cobra.Command{
		Use:   "run-receive-cache-server",
		Short: short,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			apiUrl, err := cmd.Flags().GetString("api-url")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			if apiUrl == "" {
				logging.LogFatal("Please specify --api-url")
			}

			accountNumber, err := cmd.Flags().GetString("account-number")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			if accountNumber == "" {
				logging.LogFatal("Please specify --account-number")
			}

			mustutils.Must0(
				signalclirestapiutils.RunReceiveCacheServer(ctx, &signalclirestapiutils.ReceiveCacheServerOptions{
					SignalResetClientApiUrl: apiUrl,
					Interval:                "10s",
					CacheSize:               20,
					AccountNumber:           accountNumber,
				}),
			)
		},
	}

	cmd.Flags().String("api-url", "", "Url the the Signal CLI Rest API server.")
	cmd.Flags().String("account-number", "", "Account number (phone number with +CountryCode, without spaces in between) of the receiving account.")

	return cmd
}
