package messengergeneric_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/messengerutils/messengergeneric"
	"github.com/asciich/asciichgolangpublic/pkg/messengerutils/messengerinterfaces"
)

func Test_GetContentAsString(t *testing.T) {
	t.Run("Not set", func(t *testing.T) {
		msg := &messengergeneric.Message{}

		got, err := msg.GetContentAsString()
		require.Error(t, err)
		require.Empty(t, got)
	})

	t.Run("Hello world", func(t *testing.T) {
		msg := &messengergeneric.Message{
			Message: "hello world",
		}

		got, err := msg.GetContentAsString()
		require.NoError(t, err)
		require.EqualValues(t, "hello world", got)
	})
}

func Test_GetSenderAccountAsString(t *testing.T) {
	t.Run("Not set", func(t *testing.T) {
		msg := &messengergeneric.Message{}

		got, err := msg.GetSenderAccountAsString()
		require.Error(t, err)
		require.Empty(t, got)
	})

	t.Run("Hello world", func(t *testing.T) {
		msg := &messengergeneric.Message{
			SenderAccount: "+123456",
		}

		got, err := msg.GetSenderAccountAsString()
		require.NoError(t, err)
		require.EqualValues(t, "+123456", got)
	})
}

func Test_IsSenderAccount(t *testing.T) {
	t.Run("matches", func(t *testing.T) {
		msg := &messengergeneric.Message{
			SenderAccount: "+123456",
		}

		got, err := msg.IsSenderAccount("+123456")
		require.NoError(t, err)
		require.True(t, got)
	})

	t.Run("mismatch", func(t *testing.T) {
		msg := &messengergeneric.Message{
			SenderAccount: "+123456",
		}

		got, err := msg.IsSenderAccount("+111111")
		require.NoError(t, err)
		require.False(t, got)
	})
}

func Test_GenereicMessageMatchesInterface(t *testing.T) {
	var msg messengerinterfaces.Message

	msg = &messengergeneric.Message{}

	// just to ensure the compiler accepts this test:
	_, err := msg.GetContentAsString()
	require.Error(t, err)
}
