package messengergeneric

import (
	"context"
	"slices"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/messengerutils/messengerinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func GetNewestDataMessage(ctx context.Context, messages []messengerinterfaces.Message, options *GetNewestMessageOptions) (messengerinterfaces.Message, error) {
	if len(messages) <= 0 {
		return nil, tracederrors.TracedError(ErrEmptyMessageSlice)
	}

	if options == nil {
		options = &GetNewestMessageOptions{}
	}

	var msg messengerinterfaces.Message

	var newestTs int64

	for _, m := range messages {
		isDataMessage, err := m.IsDataMessage()
		if err != nil {
			return nil, err
		}

		if !isDataMessage {
			continue
		}

		if options.AllowedSenderAccounts != nil {
			senderAccount, err := m.GetSenderAccountAsString()
			if err != nil {
				return nil, err
			}

			if !slices.Contains(options.AllowedSenderAccounts, senderAccount) {
				logging.LogInfoByCtxf(ctx, "Message of '%s' ignored since not in the allowed sender accounts '%v'.", senderAccount, options.AllowedSenderAccounts)
				continue
			}
		}

		ts, err := m.GetTimestampMilliseconds()
		if err != nil {
			return nil, err
		}

		if ts > newestTs {
			newestTs = ts
			msg = m
		}
	}

	if msg == nil {
		return nil, tracederrors.TracedError(ErrNoDataMessageFound)
	}

	return msg, nil
}
