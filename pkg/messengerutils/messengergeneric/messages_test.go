package messengergeneric_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/messengerutils/messengergeneric"
	"github.com/asciich/asciichgolangpublic/pkg/messengerutils/messengerinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/messengerutils/signalmessengerutils"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func Test_GetNewestMessage(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		ctx := getCtx()

		got, err := messengergeneric.GetNewestDataMessage(ctx, []messengerinterfaces.Message{}, &messengergeneric.GetNewestMessageOptions{})
		require.ErrorIs(t, err, messengergeneric.ErrEmptyMessageSlice)
		require.Nil(t, got)
	})

	t.Run("only one entry", func(t *testing.T) {
		ctx := getCtx()

		messages := []messengerinterfaces.Message{
			&signalmessengerutils.Message{
				Envelope: &signalmessengerutils.Envelope{
					Timestamp: 123,
					DataMessage: &signalmessengerutils.DataMessage{
						Message: "hello world",
					},
				},
			},
		}

		got, err := messengergeneric.GetNewestDataMessage(ctx, messages, &messengergeneric.GetNewestMessageOptions{})
		require.NoError(t, err)
		require.NotNil(t, got)

		timestamp, err := got.GetTimestampMilliseconds()
		require.NoError(t, err)
		require.EqualValues(t, timestamp, 123)

		content, err := got.GetContentAsString()
		require.NoError(t, err)

		require.EqualValues(t, "hello world", content)
	})

	t.Run("Two entries already ordered", func(t *testing.T) {
		ctx := getCtx()

		messages := []messengerinterfaces.Message{
			&signalmessengerutils.Message{
				Envelope: &signalmessengerutils.Envelope{
					Timestamp: 123,
					DataMessage: &signalmessengerutils.DataMessage{
						Message: "hello",
					},
				},
			},
			&signalmessengerutils.Message{
				Envelope: &signalmessengerutils.Envelope{
					Timestamp: 124,
					DataMessage: &signalmessengerutils.DataMessage{
						Message: "world",
					},
				},
			},
		}

		got, err := messengergeneric.GetNewestDataMessage(ctx, messages, &messengergeneric.GetNewestMessageOptions{})
		require.NoError(t, err)
		require.NotNil(t, got)

		timestamp, err := got.GetTimestampMilliseconds()
		require.NoError(t, err)
		require.EqualValues(t, timestamp, 124)

		content, err := got.GetContentAsString()
		require.NoError(t, err)

		require.EqualValues(t, "world", content)
	})

	t.Run("Two entries reverse order", func(t *testing.T) {
		ctx := getCtx()

		messages := []messengerinterfaces.Message{
			&signalmessengerutils.Message{
				Envelope: &signalmessengerutils.Envelope{
					Timestamp: 125,
					DataMessage: &signalmessengerutils.DataMessage{
						Message: "hello",
					},
				},
			},
			&signalmessengerutils.Message{
				Envelope: &signalmessengerutils.Envelope{
					Timestamp: 124,
					DataMessage: &signalmessengerutils.DataMessage{
						Message: "world",
					},
				},
			},
		}

		got, err := messengergeneric.GetNewestDataMessage(ctx, messages, &messengergeneric.GetNewestMessageOptions{})
		require.NoError(t, err)
		require.NotNil(t, got)

		timestamp, err := got.GetTimestampMilliseconds()
		require.NoError(t, err)
		require.EqualValues(t, timestamp, 125)

		content, err := got.GetContentAsString()
		require.NoError(t, err)

		require.EqualValues(t, "hello", content)
	})

	t.Run("Non data messages are ignored", func(t *testing.T) {
		ctx := getCtx()

		messages := []messengerinterfaces.Message{
			&signalmessengerutils.Message{
				Envelope: &signalmessengerutils.Envelope{
					Timestamp: 125,
					DataMessage: &signalmessengerutils.DataMessage{
						Message: "hello",
					},
				},
			},
			&signalmessengerutils.Message{
				Envelope: &signalmessengerutils.Envelope{
					Timestamp: 124,
					DataMessage: &signalmessengerutils.DataMessage{
						Message: "world",
					},
				},
			},

			&signalmessengerutils.Message{
				Envelope: &signalmessengerutils.Envelope{
					Timestamp:      126,
					ReceiptMessage: &signalmessengerutils.ReceiptMessage{},
				},
			},
		}

		got, err := messengergeneric.GetNewestDataMessage(ctx, messages, &messengergeneric.GetNewestMessageOptions{})
		require.NoError(t, err)
		require.NotNil(t, got)

		timestamp, err := got.GetTimestampMilliseconds()
		require.NoError(t, err)
		require.EqualValues(t, timestamp, 125)

		content, err := got.GetContentAsString()
		require.NoError(t, err)

		require.EqualValues(t, "hello", content)
	})
}
