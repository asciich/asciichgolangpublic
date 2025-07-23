package netutils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/contextutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/netutils"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func TestDnsLookupIpV4(t *testing.T) {
	ips, err := netutils.DnsLookupIpV4(getCtx(), "asciich.ch")
	require.NoError(t, err)
	require.EqualValues(t, []string{"80.74.146.168"}, ips)
}

func TestDnsReverseLookup(t *testing.T) {
	fqdns, err := netutils.DnsReverseLookup(getCtx(), "80.74.146.168")
	require.NoError(t, err)
	require.EqualValues(t, []string{"ns24.kreativmedia.ch."}, fqdns)
}
