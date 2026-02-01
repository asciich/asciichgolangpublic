package exoscalenativeclient_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/cloudutils/exoscaleutils/exoscalenativeclient"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/dnsutils/dnsgeneric"
)

func TestGetDomainUuid(t *testing.T) {
	t.Run("not found", func(t *testing.T) {
		ctx := getCtx()

		client, err := exoscalenativeclient.NewNativeClientFromEnvVars(ctx)
		require.NoError(t, err)

		uuid, err := exoscalenativeclient.GetDomainUuid(ctx, client, "doesnotexists.ch")
		require.Error(t, err)
		require.Empty(t, uuid)
		require.True(t, dnsgeneric.IsErrDnsDomainNotFound(err))
		require.True(t, dnsgeneric.IsErrNotFound(err))
		require.False(t, dnsgeneric.IsErrDnsDomainRecordNotFound(err))
	})
}
