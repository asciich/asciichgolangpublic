package netutils

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/contextutils"
)

func TestLookupIpV4(t *testing.T) {
	require.EqualValues(
		t,
		[]string{"80.74.146.168"},
		MustLookupIpV4(contextutils.ContextVerbose(), "asciich.ch"),
	)
}
