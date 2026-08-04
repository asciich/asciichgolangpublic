package netutils_test

import (
	"context"
	"net"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/netutils"
)

// Test_Example_GetNextFreePort demonstrates how to find the next free TCP port.
func Test_Example_GetNextFreePort(t *testing.T) {
	ctx := contextutils.WithVerbose(context.TODO())

	port, err := netutils.GetNextFreePort(ctx, 8080)
	require.NoError(t, err)
	require.GreaterOrEqual(t, port, 8080)
	require.LessOrEqual(t, port, 65535)

	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	require.NoError(t, err)
	defer listener.Close()

	t.Logf("Found free port: %d", port)
}

// Test_Example_GetNextFreePort_InvalidPort tests error handling
func Test_Example_GetNextFreePort_InvalidPort(t *testing.T) {
	ctx := contextutils.WithVerbose(context.TODO())

	_, err := netutils.GetNextFreePort(ctx, 0)
	require.Error(t, err)

	_, err = netutils.GetNextFreePort(ctx, -1)
	require.Error(t, err)

	_, err = netutils.GetNextFreePort(ctx, 70000)
	require.Error(t, err)
}

// Test_Example_GetNextFreePort_ConsecutivePorts tests finding multiple ports
func Test_Example_GetNextFreePort_ConsecutivePorts(t *testing.T) {
	ctx := contextutils.WithVerbose(context.TODO())

	port1, err := netutils.GetNextFreePort(ctx, 9000)
	require.NoError(t, err)

	port2, err := netutils.GetNextFreePort(ctx, port1+1)
	require.NoError(t, err)

	require.NotEqual(t, port1, port2)
	require.Greater(t, port2, port1)

	t.Logf("Found consecutive free ports: %d, %d", port1, port2)
}

// Test_Example_GetNextFreePort_WithOccupiedPort tests with occupied port
func Test_Example_GetNextFreePort_WithOccupiedPort(t *testing.T) {
	ctx := contextutils.WithVerbose(context.TODO())

	listener, err := net.Listen("tcp", ":9500")
	require.NoError(t, err)
	defer listener.Close()

	port, err := netutils.GetNextFreePort(ctx, 9500)
	require.NoError(t, err)
	require.GreaterOrEqual(t, port, 9500)

	testListener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	require.NoError(t, err)
	defer testListener.Close()

	t.Logf("Found free port %d (9500 was occupied)", port)
}
