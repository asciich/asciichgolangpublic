package signalclirestapicmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/messengerutils/messengeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/messengerutils/signalmessengerutils/signalclirestapiutils"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewSendMessageCmd() *cobra.Command {
	const short = "Send a Signal message."

	cmd := &cobra.Command{
		Use: "send-message",
		Short: short,
		Long: short + `

Usage:
` + os.Args[0] + ` messenger signal signal-cli-rest-api send-message --api-url=https://url-of-singal-cli-rest-api --account-number +4... --recipients +4... [--recipients +4...] --message="hello world"
`,
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

			recipients, err := cmd.Flags().GetStringSlice("recipients")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			if len(recipients) <= 0 {
				logging.LogFatal("Please specify at least one --recipient")
			}

			message, err := cmd.Flags().GetString("message")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			if message == "" {
				logging.LogFatal("Please specify --message")
			}

			mustutils.Must0(
				signalclirestapiutils.SendMessage(
					ctx,
					apiUrl,
					&messengeroptions.SendMessageOptions{
						Message: message,
						SenderAccount: accountNumber,
						Recipints: recipients,
					},
				),
			)

			logging.LogGoodByCtxf(ctx, "Signal message to '%s' sent.", strings.Join(recipients, ", "))
		},
	}


	cmd.Flags().String("api-url", "", "Url the the Signal CLI Rest API server.")
	cmd.Flags().String("account-number", "", "Account number (phone number with +CountryCode, without spaces in between) of the sender account.")
	cmd.Flags().StringSlice("recipients", []string{}, "Recepient numbers.")
	cmd.Flags().String("message", "", "The message (payload) to send.")
	
	return cmd
}