package exoscalenativeclient_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/cloudutils/exoscaleutils/exoscalenativeclient"
)

func Test_SanitizeNames(t *testing.T) {
	t.Run("both empty", func(t *testing.T) {
		domainName, recordName, err := exoscalenativeclient.SanitizeNames("", "")
		require.Error(t, err)
		require.Empty(t, domainName)
		require.Empty(t, recordName)
	})

	t.Run("domainName empty", func(t *testing.T) {
		domainName, recordName, err := exoscalenativeclient.SanitizeNames("", "example")
		require.Error(t, err)
		require.Empty(t, domainName)
		require.Empty(t, recordName)
	})

	t.Run("recordName empty", func(t *testing.T) {
		domainName, recordName, err := exoscalenativeclient.SanitizeNames("example.com", "")
		require.Error(t, err)
		require.Empty(t, domainName)
		require.Empty(t, recordName)
	})

	t.Run("already sanized", func(t *testing.T) {
		domainName, recordName, err := exoscalenativeclient.SanitizeNames("example.com", "record")
		require.NoError(t, err)
		require.EqualValues(t, "example.com", domainName)
		require.EqualValues(t, "record", recordName)
	})

	t.Run("already record contains domainName as well", func(t *testing.T) {
		domainName, recordName, err := exoscalenativeclient.SanitizeNames("example.com", "record.example.com")
		require.NoError(t, err)
		require.EqualValues(t, "example.com", domainName)
		require.EqualValues(t, "record", recordName)
	})

}
