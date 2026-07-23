package pfsenseutils

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func Test_ParseDashboardPageAndGetSystemName(t *testing.T) {
	ctx := getCtx()

	content, err := os.ReadFile("./testdata/dashboardpage.html")
	require.NoError(t, err)

	x, err := ParseDashboardPageAndGetSystemName(ctx, content)
	require.NoError(t, err)
	require.EqualValues(t, "router.example.com", x)
}
