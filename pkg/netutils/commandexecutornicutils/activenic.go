package commandexecutornicutils

import (
	"context"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// GetFirstActiveNicName returns the name of the first active network interface
// that is UP and has a global-scope IPv4 address (loopback is excluded).
//
// This is used to auto-detect the interface KubeVip should bind the VIP to when
// no interface is explicitly configured.
func GetFirstActiveNicName(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor) (string, error) {
	if commandExecutor == nil {
		return "", tracederrors.TracedErrorNil("commandExecutor")
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return "", err
	}

	logging.LogInfoByCtxf(ctx, "Get first active NIC name on '%s' started.", hostDescription)

	// "ip -o -4 addr show scope global up" lists only interfaces that are UP and
	// carry a global-scope IPv4 address (loopback is scope host, so it is
	// excluded automatically). "-o" makes ip print one entry per line.
	stdout, err := commandExecutor.RunCommandAndGetStdoutAsString(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{"ip", "-o", "-4", "addr", "show", "scope", "global", "up"},
		},
	)
	if err != nil {
		return "", err
	}

	// Each line looks like:
	//   "2: enp1s0    inet 192.168.122.191/24 metric 1024 ... enp1s0\..."
	// The interface name is the 2nd whitespace-separated field.
	ret := ""
	for _, line := range strings.Split(stdout, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		ret = fields[1]
		break
	}

	if ret == "" {
		return "", tracederrors.TracedErrorf("Unable to auto-detect an active NIC on '%s': no interface with a global-scope IPv4 address found.", hostDescription)
	}

	logging.LogInfoByCtxf(ctx, "Get first active NIC name on '%s' finished. Found NIC '%s'.", hostDescription, ret)

	return ret, nil
}
