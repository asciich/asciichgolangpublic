package commandexecutorfile

import (
	"context"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func GetSha256Sum(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, path string) (string, error) {
	if commandExecutor == nil {
		return "", tracederrors.TracedErrorNil("commandExecutor")
	}

	if path == "" {
		return "", tracederrors.TracedErrorEmptyString("path")
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return "", err
	}

	output, err := commandExecutor.RunCommandAndGetStdoutAsString(ctx, &parameteroptions.RunCommandOptions{
		Command: []string{"sha256sum", path},
	})
	if err != nil {
		return "", err
	}

	sha256sum := strings.Split(output, " ")[0]

	logging.LogInfoByCtxf(ctx, "Sha256sum of '%s' is '%s' on '%s'.", path, sha256sum, hostDescription)

	return sha256sum, nil
}
