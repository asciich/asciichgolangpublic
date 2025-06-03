package netutils

import (
	"context"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/asciich/asciichgolangpublic/logging"
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
