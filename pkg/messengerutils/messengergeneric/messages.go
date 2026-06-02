package messengergeneric

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/messengerutils/messengerinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func GetNewestDataMessage(ctx context.Context, messages []messengerinterfaces.Message) (messengerinterfaces.Message, error) {
	if len(messages) <= 0 {
		return nil, tracederrors.TracedError(ErrEmptyMessageSlice)
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

		ts, err := m.GetTimestampMilliseconds()
		if err != nil {
			return nil, err
		}

		if ts > newestTs {
			newestTs = ts
			msg = m
		}
	}

	return msg, nil
}
