package signalclirestapiutils

import (
	"context"
	"encoding/json"

	"github.com/asciich/asciichgolangpublic/pkg/httputils"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httpoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/messengerutils/messengerinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func SendMessage(ctx context.Context, apiUrl string, message messengerinterfaces.Message) error {
	if apiUrl == "" {
		return tracederrors.TracedErrorEmptyString("apiUrl")
	}

	if message == nil {
		return tracederrors.TracedErrorNil("message")
	}

	logging.LogChangedByCtxf(ctx, "Send signal message using singal cli rest api started.")

	content, err := message.GetContentAsString()
	if err != nil {
		return err
	}

	senderAccount, err := message.GetSenderAccountAsString()
	if err != nil {
		return err
	}

	recipients, err := message.GetRecipientsAsStringSlice()
	if err != nil {
		return err
	}

	payloadData := struct {
		Message    string   `json:"message"`
		Number     string   `json:"number"`
		Recipients []string `json:"recipients"`
	}{
		Message:    content,
		Number:     senderAccount,
		Recipients: recipients,
	}

	payload, err := json.Marshal(payloadData)
	if err != nil {
		return tracederrors.TracedErrorf("Failed to marshal send message data: %w", err)
	}

	_, err = httputils.SendRequest(ctx,
		&httpoptions.RequestOptions{
			Url:    apiUrl,
			Method: "POST",
			Path:   "/v2/send",
			Header: map[string]string{
				"Content-Type": "application/json",
			},
			Data: payload,
		},
	)
	if err != nil {
		return err
	}

	logging.LogChangedByCtxf(ctx, "Send signal message using singal cli rest api finished.")

	return nil
}

func SendResponse(ctx context.Context, apiUrl string, messageToRespondTo messengerinterfaces.Message, message messengerinterfaces.Message) error {
	logging.LogChangedByCtxf(ctx, "Send response to signal message using singal cli rest api started.")

	if apiUrl == "" {
		return tracederrors.TracedErrorEmptyString("apiUrl")
	}

	if messageToRespondTo == nil {
		return tracederrors.TracedErrorNil("messageToRespondTo")
	}

	if message == nil {
		return tracederrors.TracedErrorNil("message")
	}

	content, err := message.GetContentAsString()
	if err != nil {
		return err
	}

	senderAccount, err := message.GetSenderAccountAsString()
	if err != nil {
		return err
	}

	recipientAccount, err := messageToRespondTo.GetSenderAccountAsString()
	if err != nil {
		return err
	}

	quoteMessage, err := messageToRespondTo.GetContentAsString()
	if err != nil {
		return err
	}

	quoteTimestamp, err := messageToRespondTo.GetTimestampMilliseconds()
	if err != nil {
		return err
	}

	payloadData := struct {
		Message        string   `json:"message"`
		Number         string   `json:"number"`
		Recipients     []string `json:"recipients"`
		QuoteTimestamp int64    `json:"quote_timestamp"`
		QuoteAuthor    string   `json:"quote_author"`
		QuoteMessage   string   `json:"quote_message"`
	}{
		Message:        content,
		Number:         senderAccount,
		Recipients:     []string{recipientAccount},
		QuoteTimestamp: quoteTimestamp,
		QuoteAuthor:    senderAccount,
		QuoteMessage:   quoteMessage,
	}

	payload, err := json.Marshal(payloadData)
	if err != nil {
		return tracederrors.TracedErrorf("Failed to marshal send message data: %w", err)
	}

	_, err = httputils.SendRequest(ctx,
		&httpoptions.RequestOptions{
			Url:    apiUrl,
			Method: "POST",
			Path:   "/v2/send",
			Header: map[string]string{
				"Content-Type": "application/json",
			},
			Data: payload,
		},
	)
	if err != nil {
		return err
	}

	logging.LogChangedByCtxf(ctx, "Send response to signal message using singal cli rest api finished.")

	return nil
}
