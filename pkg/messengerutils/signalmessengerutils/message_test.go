package signalmessengerutils_test

import (
	"encoding/json"
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

func Test_GetSenderAccount(t *testing.T) {
	msg := &signalmessengerutils.Message{}

	msgData := `
{
    "envelope": {
		"source": "+41711111111",
		"sourceNumber": "+41711111111",
		"sourceUuid": "12abcdef-a1aa-11aa-11aa-1234aaaaaaaa",
		"sourceName": "Reto Hasler",
		"sourceDevice": 1,
		"timestamp": 1779568480623,
		"serverReceivedTimestamp": 1779568482093,
		"serverDeliveredTimestamp": 1779568488710,
		"dataMessage": {
		"timestamp": 1779568480623,
			"message": "Ggghhh",
			"expiresInSeconds": 0,
			"isExpirationUpdate": false,
			"viewOnce": false
		}
	}
}
`

	err := json.Unmarshal([]byte(msgData), msg)
	require.NoError(t, err)

	senderAccount, err := msg.GetSenderAccountAsString()
	require.NoError(t, err)
	require.EqualValues(t, "+41711111111", senderAccount)
}

func Test_GetTimestampMilliseconds(t *testing.T) {
	msg := &signalmessengerutils.Message{}

	msgData := `
{
    "envelope": {
		"source": "+41711111111",
		"sourceNumber": "+41711111111",
		"sourceUuid": "12abcdef-a1aa-11aa-11aa-1234aaaaaaaa",
		"sourceName": "Reto Hasler",
		"sourceDevice": 1,
		"timestamp": 1779568480623,
		"serverReceivedTimestamp": 1779568482093,
		"serverDeliveredTimestamp": 1779568488710,
		"dataMessage": {
		"timestamp": 1779568480623,
			"message": "Ggghhh",
			"expiresInSeconds": 0,
			"isExpirationUpdate": false,
			"viewOnce": false
		}
	}
}
`

	err := json.Unmarshal([]byte(msgData), msg)
	require.NoError(t, err)

	ts, err := msg.GetTimestampMilliseconds()
	require.NoError(t, err)
	require.EqualValues(t, 1779568480623, ts)
}

func Test_IsSenderAccount(t *testing.T) {
	t.Run("matches", func(t *testing.T) {
		msg := &signalmessengerutils.Message{
			Account: "+123456",
		}

		got, err := msg.IsSenderAccount("+123456")
		require.NoError(t, err)
		require.True(t, got)
	})

	t.Run("mismatch", func(t *testing.T) {
		msg := &signalmessengerutils.Message{
			Account: "+123456",
		}

		got, err := msg.IsSenderAccount("+111111")
		require.NoError(t, err)
		require.False(t, got)
	})
}
