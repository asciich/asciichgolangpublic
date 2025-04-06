package netutils

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func TestDnsLookupIpV4(t *testing.T) {
	require.EqualValues(
		t,
		[]string{"80.74.146.168"},
		MustDnsLookupIpV4(getCtx(), "asciich.ch"),
	)
}

func TestDnsReverseLookup(t *testing.T) {
	require.EqualValues(
		t,
		[]string{"ns24.kreativmedia.ch."},
		MustDnsReverseLookup(getCtx(), "80.74.146.168"),
	)
}
