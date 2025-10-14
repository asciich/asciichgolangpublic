package ansibleutils

import (
	"context"
	"fmt"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/ansibleutils/ansibleparemeteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexec"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/shellutils/shelllinehandler"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func RunPlaybook(ctx context.Context, options *RunOptions) error {
	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	binPath, err := ansibleparemeteroptions.GetAnsiblePlaybookPath(options)
	if err != nil {
		return err
	}

	playbookPath, err := options.GetPlaybookPath()
	if err != nil {
		return err
	}

	limit, err := options.GetLimit()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Run ansible playbook '%s' using '%s' against hosts '%s' started.", playbookPath, binPath, limit)

	cmd := []string{binPath, playbookPath, "--limit", limit, "--inventory=" + limit + ","}

	if len(options.Tags) > 0 {
		cmd = append(cmd, fmt.Sprintf("--tags=%s", strings.Join(options.Tags, ",")))
	}

	joined, err := shelllinehandler.Join(cmd)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "CLI used to run ansible playbook: %s", joined)

	_, err = commandexecutorexec.RunCommand(
		commandexecutorgeneric.WithLiveOutputOnStdoutIfVerbose(ctx),
		&parameteroptions.RunCommandOptions{
			Command: cmd,
		},
	)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Run ansible playbook '%s' using '%s' against hosts '%s' finished.", playbookPath, binPath, limit)

	return nil
}
