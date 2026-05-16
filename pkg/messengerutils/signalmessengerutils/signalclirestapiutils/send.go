package signalclirestapiutils

import (
	"context"
	"encoding/json"

	"github.com/asciich/asciichgolangpublic/pkg/httputils"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httpoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/messengerutils/messengeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func SendMessage(ctx context.Context, apiUrl string, options *messengeroptions.SendMessageOptions) error {
	if apiUrl == "" {
		return tracederrors.TracedErrorEmptyString("apiUrl")
	}

	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	message, err := options.GetMessage()
	if err != nil {
		return err
	}

	senderAccount, err := options.GetSenderAccount()
	if err != nil {
		return err
	}

	recipients, err := options.GetRecipients()
	if err != nil {
		return err
	}

	payloadData := struct {
		Message    string   `json:"message"`
		Number     string   `json:"number"`
		Recipients []string `json:"recipients"`
	}{
		Message:    message,
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

	logging.LogChangedByCtxf(ctx, "Sent signal message using singal cli rest api.")

	return nil
}
