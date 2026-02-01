package exoscalenativeclient_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/cloudutils/exoscaleutils/exoscalenativeclient"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/continuousintegration"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/dnsutils/dnsoptions"
)

func Test_DeleteDmsDomainRecor(t *testing.T) {
	continuousintegration.SkipInGithubCi(t, "No credentails for Exoscale available in Github.")

	const domainName = "asciich-dev.ch"

	t.Run("Noting to delete", func(t *testing.T) {
		ctx := getCtx()

		client, err := exoscalenativeclient.NewNativeClientFromEnvVars(ctx)
		require.NoError(t, err)

		deleteCtx := contextutils.WithChangeIndicator(ctx)

		err = exoscalenativeclient.DeleteDNSDomainRecords(deleteCtx, client, domainName, &dnsoptions.DeleteDnsDomainRecordOptions{
			Name: "does-not-exist",
		})
		require.NoError(t, err)
		require.False(t, contextutils.IsChanged(deleteCtx))
	})

	t.Run("one to delete delete", func(t *testing.T) {
		ctx := getCtx()

		const recordName = "one-to-delete"

		client, err := exoscalenativeclient.NewNativeClientFromEnvVars(ctx)
		require.NoError(t, err)

		deleteCtx := contextutils.WithChangeIndicator(ctx)

		// Delete existing records to start with a proper test setup.
		err = exoscalenativeclient.DeleteDNSDomainRecords(deleteCtx, client, domainName, &dnsoptions.DeleteDnsDomainRecordOptions{
			Name: recordName,
		})
		require.NoError(t, err)

		exists, err := exoscalenativeclient.DnsDomainRecordExists(ctx, client, domainName, &dnsoptions.DnsDomainRecordExistsOptions{
			Name: recordName,
		})
		require.NoError(t, err)
		require.False(t, exists)

		// Create one record to delete:
		err = exoscalenativeclient.CreateDnsDomainRecord(ctx, client, domainName, &dnsoptions.CreateDnsDomainRecordOptions{
			RecordType:  "A",
			IPv4Address: "192.168.11.14",
			Name:        recordName,
			TTL:         60,
		})
		require.NoError(t, err)

		exists, err = exoscalenativeclient.DnsDomainRecordExists(ctx, client, domainName, &dnsoptions.DnsDomainRecordExistsOptions{
			Name: recordName,
		})
		require.NoError(t, err)
		require.True(t, exists)

		// Delete record
		err = exoscalenativeclient.DeleteDNSDomainRecords(deleteCtx, client, domainName, &dnsoptions.DeleteDnsDomainRecordOptions{
			Name: recordName,
		})
		require.NoError(t, err)
		require.True(t, contextutils.IsChanged(deleteCtx))

		exists, err = exoscalenativeclient.DnsDomainRecordExists(ctx, client, domainName, &dnsoptions.DnsDomainRecordExistsOptions{
			Name: recordName,
		})
		require.NoError(t, err)
		require.False(t, exists)
	})

	t.Run("three to delete delete", func(t *testing.T) {
		ctx := getCtx()

		const recordName = "three-to-delete"

		client, err := exoscalenativeclient.NewNativeClientFromEnvVars(ctx)
		require.NoError(t, err)

		deleteCtx := contextutils.WithChangeIndicator(ctx)

		// Delete existing records to start with a proper test setup.
		err = exoscalenativeclient.DeleteDNSDomainRecords(deleteCtx, client, domainName, &dnsoptions.DeleteDnsDomainRecordOptions{
			Name: recordName,
		})
		require.NoError(t, err)

		exists, err := exoscalenativeclient.DnsDomainRecordExists(ctx, client, domainName, &dnsoptions.DnsDomainRecordExistsOptions{
			Name: recordName,
		})
		require.NoError(t, err)
		require.False(t, exists)

		// Create one record to delete:
		for i := range 3 {
			err = exoscalenativeclient.CreateDnsDomainRecord(ctx, client, domainName, &dnsoptions.CreateDnsDomainRecordOptions{
				RecordType:           "A",
				IPv4Address:          "192.168.11.1" + strconv.Itoa(i),
				Name:                 recordName,
				TTL:                  60,
				AllowMultipleEntries: true,
			})
		}
		require.NoError(t, err)

		exists, err = exoscalenativeclient.DnsDomainRecordExists(ctx, client, domainName, &dnsoptions.DnsDomainRecordExistsOptions{
			Name: recordName,
		})
		require.NoError(t, err)
		require.True(t, exists)

		list, err := exoscalenativeclient.ListDomainRecords(ctx, client, domainName, &dnsoptions.ListDnsDomainRecordOptions{
			Name: recordName,
		})
		require.NoError(t, err)
		require.Len(t, list, 3)

		// Delete record
		err = exoscalenativeclient.DeleteDNSDomainRecords(deleteCtx, client, domainName, &dnsoptions.DeleteDnsDomainRecordOptions{
			Name: recordName,
		})
		require.NoError(t, err)
		require.True(t, contextutils.IsChanged(deleteCtx))

		exists, err = exoscalenativeclient.DnsDomainRecordExists(ctx, client, domainName, &dnsoptions.DnsDomainRecordExistsOptions{
			Name: recordName,
		})
		require.NoError(t, err)
		require.False(t, exists)
	})

}
