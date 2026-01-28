package publicips_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/publicips"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func TestGetPublicIp(t *testing.T) {
	ip, err := publicips.GetPublicIp(getCtx())
	require.NoError(t, err)
	require.Regexp(t, "^[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}$", ip)
}
