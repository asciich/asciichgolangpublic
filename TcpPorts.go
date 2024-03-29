package asciichgolangpublic

import (
	"net"
	"strconv"
	"strings"
	"time"
)

type TcpPortsService struct{}

func NewTcpPortsService() (t *TcpPortsService) {
	return new(TcpPortsService)
}

func TcpPorts() (t *TcpPortsService) {
	return NewTcpPortsService()
}

// Check if a TCP port on the given hostnameOrIp with given portNumber is open.
// The evaluation is done by opening a TCP socket and close it again.
func (t *TcpPortsService) IsPortOpen(hostnameOrIp string, port int, verbose bool) (isOpen bool, err error) {
	hostnameOrIp = strings.TrimSpace(hostnameOrIp)
	if hostnameOrIp == "" {
		return false, TracedErrorEmptyString("hostnameOrIp")
	}

	if port <= 0 {
		return false, TracedErrorf("Invalid port number '%d'.", port)
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
			LogInfof("Port '%d' on host '%s' is open.", port, hostnameOrIp)
		} else {
			LogInfof("Port '%d' on host '%s' is NOT open.", port, hostnameOrIp)
		}
	}

	return isOpen, nil
}

func (t *TcpPortsService) MustIsPortOpen(hostnameOrIp string, port int, verbose bool) (isOpen bool) {
	isOpen, err := t.IsPortOpen(hostnameOrIp, port, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return isOpen
}
