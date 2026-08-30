package dnsutils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/dnsutils"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func TestDnsLookupIpV4(t *testing.T) {
	ctx, recorder := logging.WithLogRecorder(getCtx())

	ips, err := dnsutils.DnsLookupIpV4(ctx, "asciich.ch")
	require.NoError(t, err)
	require.EqualValues(t, []string{"80.74.146.168"}, ips)

	logOutput := recorder.String()
	require.Contains(t, logOutput, "asciich.ch")
	require.Contains(t, logOutput, "80.74.146.168")
	require.Contains(t, logOutput, "server")
}

func TestDnsReverseLookup(t *testing.T) {
	ctx, recorder := logging.WithLogRecorder(getCtx())

	fqdns, err := dnsutils.DnsReverseLookup(ctx, "80.74.146.168")
	require.NoError(t, err)
	require.EqualValues(t, []string{"ns24.kreativmedia.ch."}, fqdns)

	logOutput := recorder.String()
	require.Contains(t, logOutput, "80.74.146.168")
	require.Contains(t, logOutput, "ns24.kreativmedia.ch.")
	require.Contains(t, logOutput, "server")
}
