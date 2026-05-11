package signalclirestapiutils

import (
	"context"
	"encoding/json"

	"github.com/asciich/asciichgolangpublic/pkg/httputils"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httpoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/messengerutils/signalmessengerutils"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"github.com/asciich/asciichgolangpublic/pkg/urlsutils"
)

func ReceiveMessages(ctx context.Context, apiUrl string, accountNumber string) ([]*signalmessengerutils.Message, error) {
	if apiUrl == "" {
		return nil, tracederrors.TracedErrorEmptyString("apiUrl")
	}

	if accountNumber == "" {
		return nil, tracederrors.TracedErrorEmptyString("accountNumber")
	}

	baseUrl, err := urlsutils.GetSchemeAndHost(apiUrl)
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Receive signal messages rom signal-cli-rest-api '%s' for account number '%s' started.", apiUrl, accountNumber)

	url := baseUrl + "/v1/receive/" + accountNumber

	response, err := httputils.SendRequestAndGetBodyAsBytes(ctx, &httpoptions.RequestOptions{
		Url: url,
	})
	if err != nil {
		return nil, err
	}

	messages := []*signalmessengerutils.Message{}

	err = json.Unmarshal(response, &messages)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Unmarshal response '%s' of received signal messages failed: %w", string(response), err)
	}

	logging.LogInfoByCtxf(ctx, "Receive signal messages rom signal-cli-rest-api '%s' for account number '%s' finished. Received %d messages.", apiUrl, accountNumber, len(messages))

	return messages, nil
}
