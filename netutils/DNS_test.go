package netutils

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/contextutils"
)

func TestDnsLookupIpV4(t *testing.T) {
	require.EqualValues(
		t,
		[]string{"80.74.146.168"},
		MustDnsLookupIpV4(contextutils.ContextVerbose(), "asciich.ch"),
	)
}

func TestDnsReverseLookup(t *testing.T) {
	require.EqualValues(
		t,
		[]string{"ns24.kreativmedia.ch."},
		MustDnsReverseLookup(contextutils.ContextVerbose(), "80.74.146.168"),
	)
}
