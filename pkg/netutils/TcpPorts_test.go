package netutils_test

import (
	"context"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/pkg/netutils"
	"github.com/asciich/asciichgolangpublic/testutils"
)

func TestTcpPortsIsPortOpen(t *testing.T) {
	tests := []struct {
		hostname       string
		portNumber     int
		expectedIsOpen bool
	}{
		{"google.ch", 80, true},
		{"google.ch", 443, true},
		{"google.ch", 442, false},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				isOpen, err := netutils.IsTcpPortOpen(getCtx(), tt.hostname, tt.portNumber)
				require.NoError(t, err)

				require.EqualValues(t, tt.expectedIsOpen, isOpen)
			},
		)
	}
}

func Test_IsPortAvailableForListening(t *testing.T) {
	const testport = 12345

	t.Run("invalid port", func(t *testing.T) {
		isAvailable, err := netutils.IsTcpPortAvailableForListening(getCtx(), 0)
		require.Error(t, err)
		require.False(t, isAvailable)
	})

	t.Run("invalid port", func(t *testing.T) {
		isAvailable, err := netutils.IsTcpPortAvailableForListening(getCtx(), -1)
		require.Error(t, err)
		require.False(t, isAvailable)
	})

	t.Run("open port", func(t *testing.T) {
		isAvailable, err := netutils.IsTcpPortAvailableForListening(getCtx(), testport)
		require.NoError(t, err)
		require.True(t, isAvailable)
	})
}

func Test_WaitPortAvailableForListening(t *testing.T) {
	t.Run("wait", func(t *testing.T) {
		const testport = 12346
		ctx := getCtx()

		isAvailable, err := netutils.IsTcpPortAvailableForListening(ctx, 0)
		require.Error(t, err)
		require.False(t, isAvailable)

		var closed = false
		cWaitOpen := make(chan int, 10)

		go func() {
			// Open port
			ln, err := net.Listen("tcp", ":"+strconv.Itoa(testport))
			require.NoError(t, err)
			logging.LogInfoByCtxf(ctx, "Port '%d' opened for testing.", testport)

			close(cWaitOpen)
			time.Sleep(time.Second * 1)

			err = ln.Close()
			require.NoError(t, err)
			closed = true

			logging.LogInfoByCtxf(ctx, "Port '%d' closed for testing.", testport)
		}()

		<-cWaitOpen

		err = netutils.WaitPortAvailableForListening(ctx, testport)
		require.NoError(t, err)
		require.True(t, closed)
	})

	t.Run("timeout", func(t *testing.T) {
		const testport = 12346
		ctx := getCtx()

		isAvailable, err := netutils.IsTcpPortAvailableForListening(ctx, 0)
		require.Error(t, err)
		require.False(t, isAvailable)

		var closed = false
		cWaitOpen := make(chan int, 10)

		go func() {
			// Open port
			ln, err := net.Listen("tcp", ":"+strconv.Itoa(testport))
			require.NoError(t, err)
			logging.LogInfoByCtxf(ctx, "Port '%d' opened for testing.", testport)

			close(cWaitOpen)
			time.Sleep(time.Second * 1)

			err = ln.Close()
			require.NoError(t, err)
			closed = true

			logging.LogInfoByCtxf(ctx, "Port '%d' closed for testing.", testport)
		}()

		<-cWaitOpen

		ctx, _ = context.WithTimeout(ctx, time.Millisecond * 100)
		err = netutils.WaitPortAvailableForListening(ctx, testport)
		require.Error(t, err)
		require.False(t, closed)
	})
}
