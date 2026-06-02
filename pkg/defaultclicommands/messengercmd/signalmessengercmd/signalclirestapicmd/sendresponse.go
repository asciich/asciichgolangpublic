package signalclirestapicmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/messengerutils/messengergeneric"
	"github.com/asciich/asciichgolangpublic/pkg/messengerutils/signalmessengerutils"
	"github.com/asciich/asciichgolangpublic/pkg/messengerutils/signalmessengerutils/signalclirestapiutils"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewSendResponseCmd() *cobra.Command {
	const short = "Send a Signal response."

	// Add new command for sending a response
	cmd := &cobra.Command{
		Use:   "send-response",
		Short: "Send a Signal response.",
		Long: short + `

Usage:
` + os.Args[0] + ` messenger signal signal-cli-rest-api send-response --api-url=https://url-of-singal-cli-rest-api --account-number +4... --message="this is the response" --quote-timestamp="1779570020369" --quote-author="+4.." --quote-message="this is the quoted message"
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

			message, err := cmd.Flags().GetString("message")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			if message == "" {
				logging.LogFatal("Please specify --message")
			}

			quoteTimestamp, err := cmd.Flags().GetInt64("quote-timestamp")
			if err != nil {
				logging.LogGoErrorFatal(err)
			}

			if quoteTimestamp <= 0 {
				logging.LogFatalf("Please specify a valid --quote-timestamp. Got '%d'.", quoteTimestamp)
			}

			quoteAuthor, err := cmd.Flags().GetString("quote-author")
			if err != nil {
				logging.LogGoErrorFatal(err)
			}

			if quoteAuthor == "" {
				logging.LogFatal("Please specify --quote-author")
			}

			quoteMessage, err := cmd.Flags().GetString("quote-message")
			if err != nil {
				logging.LogGoErrorFatal(err)
			}

			if quoteMessage == "" {
				logging.LogFatal("Please specify --quote-message")
			}

			messageToQuote := &signalmessengerutils.Message{
				Envelope: &signalmessengerutils.Envelope{
					SourceNumber: quoteAuthor,
					DataMessage: &signalmessengerutils.DataMessage{
						Message: quoteMessage,
					},
					Timestamp: quoteTimestamp,
				},
			}

			mustutils.Must0(
				signalclirestapiutils.SendResponse(
					ctx,
					apiUrl,
					messageToQuote,
					&messengergeneric.Message{
						Message:       message,
						SenderAccount: accountNumber,
						Recipients:    []string{},
					},
				),
			)

			logging.LogGoodByCtxf(ctx, "Signal response send")
		},
	}

	cmd.Flags().String("api-url", "", "Url the the Signal CLI Rest API server.")
	cmd.Flags().String("account-number", "", "Account number (phone number with +CountryCode, without spaces in between) of the sender account.")
	cmd.Flags().String("message", "", "The message (payload) to send.")
	cmd.Flags().Int64("quote-timestamp", 0, "The timestamp of the message to quote in milliseconds.")
	cmd.Flags().String("quote-author", "", "The account number (phone number with +CountryCode, without spaces in between) of the sender of the quote message.")
	cmd.Flags().String("quote-message", "", "The message to quote")

	return cmd
}
