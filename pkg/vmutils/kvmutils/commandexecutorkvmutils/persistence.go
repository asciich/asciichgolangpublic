package commandexecutorkvmutils

import (
	"context"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/datatypes/stringsutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func (k *CommandExecutrKvmHypervisor) IsVmPersistent(ctx context.Context, vmName string) (isPersistent bool, err error) {
	if vmName == "" {
		return false, tracederrors.TracedErrorEmptyString("vmName")
	}

	output, err := k.RunKvmCommandAndGetStdout(ctx, []string{"dominfo", vmName})
	if err != nil {
		return false, err
	}

	var found bool
	for _, line := range stringsutils.SplitLines(output, true) {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "Persistent:") {
			continue
		}

		value := strings.TrimSpace(strings.TrimPrefix(line, "Persistent:"))

		switch value {
		case "yes":
			isPersistent = true
		case "no":
			isPersistent = false
		default:
			return false, tracederrors.TracedErrorf(
				"Unexpected value '%s' in Persistent line of dominfo output for VM '%s'.", value, vmName,
			)
		}

		found = true
		break
	}

	if !found {
		return false, tracederrors.TracedErrorf(
			"Could not find 'Persistent:' line in dominfo output for VM '%s'.", vmName,
		)
	}

	logging.LogInfoByCtxf(ctx, "VM '%s' persistent: %v.", vmName, isPersistent)

	return isPersistent, nil
}
