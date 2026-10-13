package hostreconciler

import (
	"context"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// runBeforeScript executes every entry of 'beforeScript' on localhost in the given order.
//
// Every entry is executed in its own 'bash -e -o pipefail -c <entry>':
//   - Multi line entries (YAML block scalar '|') stop at the first failing command.
//   - Shell state (e.g. 'cd' or exported variables) is NOT shared between entries.
//
// Execution stops at the first failing entry and its error is returned.
//
// The entries are executed with the privileges of the current process.
//
// Note: The entries are executed on every call. They must be idempotent themselves.
func runBeforeScript(ctx context.Context, beforeScript []string) error {
	err := validateBeforeScript(beforeScript)
	if err != nil {
		return err
	}

	if len(beforeScript) == 0 {
		logging.LogInfoByCtxf(ctx, "No before_script configured. Nothing to run.")
		return nil
	}

	logging.LogInfoByCtxf(ctx, "Run before_script with '%d' entries started.", len(beforeScript))

	executor := commandexecutorexecoo.Exec()

	for i, entry := range beforeScript {
		logging.LogInfoByCtxf(ctx, "Run before_script entry '%d': '%s'", i, entry)

		stdout, err := executor.RunCommandAndGetStdoutAsString(
			ctx,
			&parameteroptions.RunCommandOptions{
				Command: []string{"bash", "-e", "-o", "pipefail", "-c", entry},
			},
		)
		if err != nil {
			return tracederrors.TracedErrorf("before_script entry '%d' failed: '%s': %w", i, entry, err)
		}

		stdout = strings.TrimSpace(stdout)
		if stdout != "" {
			logging.LogInfoByCtxf(ctx, "Output of before_script entry '%d':\n%s", i, stdout)
		}
	}

	logging.LogInfoByCtxf(ctx, "Run before_script with '%d' entries finished.", len(beforeScript))

	return nil
}
