package signalmessengerutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/messengerutils/messengerinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/messengerutils/signalmessengerutils"
)

func Test_GetMessage(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		var msg messengerinterfaces.Message = &signalmessengerutils.Message{}

		content, err := msg.GetContentAsString()
		require.Error(t, err)
		require.EqualValues(t, "", content)
	})

	t.Run("sent message", func(t *testing.T) {
		var msg messengerinterfaces.Message = &signalmessengerutils.Message{
			Envelope: &signalmessengerutils.Envelope{
				SyncMessage: &signalmessengerutils.SyncMessage{
					SentMessage: &signalmessengerutils.SentMessage{
						Message: "hello world",
					},
				},
			},
		}

		content, err := msg.GetContentAsString()
		require.NoError(t, err)
		require.EqualValues(t, "hello world", content)
	})

	t.Run("received message", func(t *testing.T) {
		var msg messengerinterfaces.Message = &signalmessengerutils.Message{
			Envelope: &signalmessengerutils.Envelope{
				DataMessage: &signalmessengerutils.DataMessage{
					Message: "hello world",
				},
			},
		}

		content, err := msg.GetContentAsString()
		require.NoError(t, err)
		require.EqualValues(t, "hello world", content)
	})
}
