package netutils

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

// Check if a TCP port on the given hostnameOrIp with given portNumber is open.
// The evaluation is done by opening a TCP socket and close it again.
func IsTcpPortOpen(ctx context.Context, hostnameOrIp string, port int) (isOpen bool, err error) {
	hostnameOrIp = strings.TrimSpace(hostnameOrIp)
	if hostnameOrIp == "" {
		return false, tracederrors.TracedErrorEmptyString("hostnameOrIp")
	}

	if port <= 0 {
		return false, tracederrors.TracedErrorf("Invalid port number '%d'.", port)
	}

	portString := strconv.Itoa(port)

	timeout := time.Second * 1
	connection, err := net.DialTimeout(
		"tcp",
		net.JoinHostPort(hostnameOrIp, portString),
		timeout,
	)
	if err != nil {
		isOpen = false
	} else {
		if connection != nil {
			connection.Close()
			isOpen = true
		} else {
			isOpen = false
		}
	}

	if isOpen {
		logging.LogInfoByCtxf(ctx, "Port '%d' on host '%s' is open.", port, hostnameOrIp)
	} else {
		logging.LogInfoByCtxf(ctx, "Port '%d' on host '%s' is NOT open.", port, hostnameOrIp)
	}

	return isOpen, nil
}

func IsTcpPortAvailableForListening(ctx context.Context, port int) (bool, error) {
	if port <= 0 {
		return false, tracederrors.TracedErrorf("Invalid port number '%d", port)
	}

	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		if opErr, ok := err.(*net.OpError); ok && opErr.Op == "listen" {
			logging.LogInfoByCtxf(ctx, "Port %d is already in use.", port)
			return false, nil
		}
		fmt.Printf("Error checking port %d: %v\n", port, err)
		return false, tracederrors.TracedErrorf("Failed to check for port '%d' available: %w", port, err)
	}
	err = ln.Close()
	if err != nil {
		return false, tracederrors.TracedErrorf("Unable to close port '%d' opened for checking availability", port)
	}

	logging.LogInfoByCtxf(ctx, "Port '%d' is available for listening.", port)
	return true, nil
}

func WaitPortAvailableForListening(ctx context.Context, port int) error {
	if port <= 0 {
		return tracederrors.TracedErrorf("Invalid port number '%d", port)
	}

	logging.LogInfoByCtxf(ctx, "Wait for port '%d' available for listening started.", port)

	for {
		err := ctx.Err()
		if err != nil {
			return tracederrors.TracedErrorf("Wait for port %d avaialbe for listening failed: %w", port, err)
		}

		isAvailable, err := IsTcpPortAvailableForListening(contextutils.WithSilent(ctx), port)
		if err != nil {
			return err
		}

		if isAvailable {
			break
		}

		time.Sleep(time.Millisecond * 50)
	}

	logging.LogInfoByCtxf(ctx, "Wait for port '%d' available for listening finished. The port is availabe for listening.", port)
	return nil
}
