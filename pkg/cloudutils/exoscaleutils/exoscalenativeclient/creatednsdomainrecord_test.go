package exoscalenativeclient_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/cloudutils/exoscaleutils/exoscalenativeclient"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/continuousintegration"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/dnsutils/dnsoptions"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func TestCreateAndDeleteDnsDomainRecord(t *testing.T) {
	continuousintegration.SkipInGithubCi(t, "No Exoscale credentials available in Github CI.")

	ctx := getCtx()

	const domainName = "asciich-dev.ch"
	const recordName = "testrecord." + domainName

	client, err := exoscalenativeclient.NewNativeClientFromEnvVars(ctx)
	require.NoError(t, err)

	// Delete the entry to get a proper test setup:
	err = exoscalenativeclient.DeleteDNSDomainRecords(ctx, client, domainName, &dnsoptions.DeleteDnsDomainRecordOptions{
		Name: recordName,
	})
	require.NoError(t, err)

	exists, err := exoscalenativeclient.DnsDomainRecordExists(ctx, client, domainName, &dnsoptions.DnsDomainRecordExistsOptions{Name: recordName})
	require.NoError(t, err)
	require.False(t, exists)

	// Creating the entry for the first time:
	ctxCreate := contextutils.WithChangeIndicator(ctx)
	err = exoscalenativeclient.CreateDnsDomainRecord(ctxCreate, client, domainName, &dnsoptions.CreateDnsDomainRecordOptions{
		Name:        recordName,
		RecordType:  "A",
		IPv4Address: "192.168.11.12",
		TTL:         1234,
	})
	require.NoError(t, err)
	require.True(t, contextutils.IsChanged(ctxCreate))

	exists, err = exoscalenativeclient.DnsDomainRecordExists(ctx, client, domainName, &dnsoptions.DnsDomainRecordExistsOptions{Name: recordName})
	require.NoError(t, err)
	require.True(t, exists)

	// Creating the entry for the second time to check idempotence
	ctxCreate = contextutils.WithChangeIndicator(ctx)
	err = exoscalenativeclient.CreateDnsDomainRecord(ctxCreate, client, domainName, &dnsoptions.CreateDnsDomainRecordOptions{
		Name:        recordName,
		RecordType:  "A",
		IPv4Address: "192.168.11.12",
		TTL:         1234,
	})
	require.NoError(t, err)
	require.False(t, contextutils.IsChanged(ctxCreate))

	exists, err = exoscalenativeclient.DnsDomainRecordExists(ctx, client, domainName, &dnsoptions.DnsDomainRecordExistsOptions{Name: recordName})
	require.NoError(t, err)
	require.True(t, exists)

	// Delete the record again
	ctxDelete := contextutils.WithChangeIndicator(ctx)
	err = exoscalenativeclient.DeleteDNSDomainRecords(ctxDelete, client, domainName, &dnsoptions.DeleteDnsDomainRecordOptions{
		Name: recordName,
	})
	require.NoError(t, err)
	require.True(t, contextutils.IsChanged(ctxDelete))

	exists, err = exoscalenativeclient.DnsDomainRecordExists(ctx, client, domainName, &dnsoptions.DnsDomainRecordExistsOptions{Name: recordName})
	require.NoError(t, err)
	require.False(t, exists)

	// Delete the record again to check idempotence
	ctxDelete = contextutils.WithChangeIndicator(ctx)
	err = exoscalenativeclient.DeleteDNSDomainRecords(ctxDelete, client, domainName, &dnsoptions.DeleteDnsDomainRecordOptions{
		Name: recordName,
	})
	require.NoError(t, err)
	require.False(t, contextutils.IsChanged(ctxDelete))

	exists, err = exoscalenativeclient.DnsDomainRecordExists(ctx, client, domainName, &dnsoptions.DnsDomainRecordExistsOptions{Name: recordName})
	require.NoError(t, err)
	require.False(t, exists)
}

func TestCreateUpdatedDnsDomainRecord(t *testing.T) {
	continuousintegration.SkipInGithubCi(t, "No Exoscale credentials available in Github CI.")

	ctx := getCtx()

	const domainName = "asciich-dev.ch"
	const recordName = "testrecord." + domainName

	client, err := exoscalenativeclient.NewNativeClientFromEnvVars(ctx)
	require.NoError(t, err)

	// Delete the entry to get a proper test setup:
	err = exoscalenativeclient.DeleteDNSDomainRecords(ctx, client, domainName, &dnsoptions.DeleteDnsDomainRecordOptions{
		Name: recordName,
	})
	require.NoError(t, err)

	exists, err := exoscalenativeclient.DnsDomainRecordExists(ctx, client, domainName, &dnsoptions.DnsDomainRecordExistsOptions{Name: recordName})
	require.NoError(t, err)
	require.False(t, exists)

	// Creating the entry for the first time:
	ctxCreate := contextutils.WithChangeIndicator(ctx)
	err = exoscalenativeclient.CreateDnsDomainRecord(ctxCreate, client, domainName, &dnsoptions.CreateDnsDomainRecordOptions{
		Name:        recordName,
		RecordType:  "A",
		IPv4Address: "192.168.11.12",
		TTL:         1234,
	})
	require.NoError(t, err)
	require.True(t, contextutils.IsChanged(ctxCreate))

	content, err := exoscalenativeclient.GetDnsDomainRecordContent(ctx, client, domainName, recordName)
	require.NoError(t, err)
	require.EqualValues(t, "192.168.11.12", content)

	// Creating the entry for the second time to check idempotence
	ctxCreate = contextutils.WithChangeIndicator(ctx)
	err = exoscalenativeclient.CreateDnsDomainRecord(ctxCreate, client, domainName, &dnsoptions.CreateDnsDomainRecordOptions{
		Name:        recordName,
		RecordType:  "A",
		IPv4Address: "192.168.11.12",
		TTL:         1234,
	})
	require.NoError(t, err)
	require.False(t, contextutils.IsChanged(ctxCreate))

	content, err = exoscalenativeclient.GetDnsDomainRecordContent(ctx, client, domainName, recordName)
	require.NoError(t, err)
	require.EqualValues(t, "192.168.11.12", content)

	// Creating the entry for the third time with the IPv4 Adress changed.
	ctxCreate = contextutils.WithChangeIndicator(ctx)
	err = exoscalenativeclient.CreateDnsDomainRecord(ctxCreate, client, domainName, &dnsoptions.CreateDnsDomainRecordOptions{
		Name:        recordName,
		RecordType:  "A",
		IPv4Address: "192.168.11.13",
		TTL:         1234,
	})
	require.NoError(t, err)
	require.True(t, contextutils.IsChanged(ctxCreate))

	content, err = exoscalenativeclient.GetDnsDomainRecordContent(ctx, client, domainName, recordName)
	require.NoError(t, err)
	require.EqualValues(t, "192.168.11.13", content)

	// Delete the record again
	ctxDelete := contextutils.WithChangeIndicator(ctx)
	err = exoscalenativeclient.DeleteDNSDomainRecords(ctxDelete, client, domainName, &dnsoptions.DeleteDnsDomainRecordOptions{
		Name: recordName,
	})
	require.NoError(t, err)
	require.True(t, contextutils.IsChanged(ctxDelete))

	exists, err = exoscalenativeclient.DnsDomainRecordExists(ctx, client, domainName, &dnsoptions.DnsDomainRecordExistsOptions{Name: recordName})
	require.NoError(t, err)
	require.False(t, exists)
}
