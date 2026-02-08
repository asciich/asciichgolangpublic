package commandexecutorheadscale

import (
	"context"
	"strconv"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func GeneratePreauthKeyForUser(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, userName string) (string, error) {
	if commandExecutor == nil {
		return "", tracederrors.TracedErrorNil("commandExecutor")
	}

	if userName == "" {
		return "", tracederrors.TracedErrorEmptyString("userName")
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return "", err
	}

	logging.LogInfoByCtxf(ctx, "Generate preauth key for headscale user '%s' on '%s' started.", userName, hostDescription)

	id, err := GetUserId(ctx, commandExecutor, userName)
	if err != nil {
		return "", err
	}


	preauthKey, err := commandExecutor.RunCommandAndGetStdoutAsString(ctx, &parameteroptions.RunCommandOptions{
		Command: []string{"headscale", "preauthkeys", "create", "--user", strconv.Itoa(id)},
	})
	if err != nil {
		return "", err
	}

	preauthKey = strings.TrimSpace(preauthKey)
	logging.LogChangedByCtxf(ctx, "Generated preauth key for headscale user '%s' on '%s'.", userName, hostDescription)

	logging.LogInfoByCtxf(ctx, "Generate preauth key for headscale user '%s' on '%s' finished.", userName, hostDescription)

	return preauthKey, nil
}
