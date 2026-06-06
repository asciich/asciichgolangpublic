package signalclirestapiutils

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/messengerutils/signalmessengerutils"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func Test_appendMessageToCache(t *testing.T) {
	t.Run("nil failes", func(t *testing.T) {
		ctx := getCtx()
		receiveCacheServer := &ReceiveCacheServer{}
		receiveCacheServer.cacheSize = 3

		err := receiveCacheServer.appendMessageToCache(ctx, nil)
		require.Error(t, err)
	})

	t.Run("test removal of old messages when cache is full", func(t *testing.T) {
		// initialize
		ctx := getCtx()
		receiveCacheServer := &ReceiveCacheServer{}
		receiveCacheServer.cacheSize = 3

		// no cached data so far:
		require.Len(t, receiveCacheServer.cache, 0)

		// Append one message:
		err := receiveCacheServer.appendMessageToCache(ctx,
			&signalmessengerutils.Message{
				Envelope: &signalmessengerutils.Envelope{
					DataMessage: &signalmessengerutils.DataMessage{
						Message: "msg1",
					},
				},
			})
		require.NoError(t, err)

		require.Len(t, receiveCacheServer.cache, 1)
		messages := receiveCacheServer.GetMessages()
		require.Len(t, messages, 1)
		content, err := messages[0].GetContentAsString()
		require.NoError(t, err)
		require.EqualValues(t, "msg1", content)

		// Append second message:
		err = receiveCacheServer.appendMessageToCache(ctx,
			&signalmessengerutils.Message{
				Envelope: &signalmessengerutils.Envelope{
					DataMessage: &signalmessengerutils.DataMessage{
						Message: "msg2",
					},
				},
			})
		require.NoError(t, err)

		require.Len(t, receiveCacheServer.cache, 2)
		messages = receiveCacheServer.GetMessages()
		require.Len(t, messages, 2)
		content, err = messages[0].GetContentAsString()
		require.NoError(t, err)
		require.EqualValues(t, "msg1", content)
		content, err = messages[1].GetContentAsString()
		require.NoError(t, err)
		require.EqualValues(t, "msg2", content)

		// Append third message:
		err = receiveCacheServer.appendMessageToCache(ctx,
			&signalmessengerutils.Message{
				Envelope: &signalmessengerutils.Envelope{
					DataMessage: &signalmessengerutils.DataMessage{
						Message: "msg3",
					},
				},
			})
		require.NoError(t, err)

		require.Len(t, receiveCacheServer.cache, 3)
		messages = receiveCacheServer.GetMessages()
		require.Len(t, messages, 3)
		content, err = messages[0].GetContentAsString()
		require.NoError(t, err)
		require.EqualValues(t, "msg1", content)
		content, err = messages[1].GetContentAsString()
		require.NoError(t, err)
		require.EqualValues(t, "msg2", content)
		content, err = messages[2].GetContentAsString()
		require.NoError(t, err)
		require.EqualValues(t, "msg3", content)

		// Append fourth message, this must remove the first one from the cacche:
		err = receiveCacheServer.appendMessageToCache(ctx,
			&signalmessengerutils.Message{
				Envelope: &signalmessengerutils.Envelope{
					DataMessage: &signalmessengerutils.DataMessage{
						Message: "msg4",
					},
				},
			})
		require.NoError(t, err)

		require.Len(t, receiveCacheServer.cache, 3)
		messages = receiveCacheServer.GetMessages()
		require.Len(t, messages, 3)
		content, err = messages[0].GetContentAsString()
		require.NoError(t, err)
		require.EqualValues(t, "msg2", content)
		content, err = messages[1].GetContentAsString()
		require.NoError(t, err)
		require.EqualValues(t, "msg3", content)
		content, err = messages[2].GetContentAsString()
		require.NoError(t, err)
		require.EqualValues(t, "msg4", content)
	})
}
