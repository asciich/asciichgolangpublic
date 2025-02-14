package netutils

import (
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

// Check if a TCP port on the given hostnameOrIp with given portNumber is open.
// The evaluation is done by opening a TCP socket and close it again.
func IsTcpPortOpen(hostnameOrIp string, port int, verbose bool) (isOpen bool, err error) {
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

	if verbose {
		if isOpen {
			logging.LogInfof("Port '%d' on host '%s' is open.", port, hostnameOrIp)
		} else {
			logging.LogInfof("Port '%d' on host '%s' is NOT open.", port, hostnameOrIp)
		}
	}

	return isOpen, nil
}

func MustIsTcpPortOpen(hostnameOrIp string, port int, verbose bool) (isOpen bool) {
	isOpen, err := IsTcpPortOpen(hostnameOrIp, port, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return isOpen
}
