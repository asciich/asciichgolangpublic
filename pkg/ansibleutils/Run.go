package ansibleutils

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/ansibleutils/ansibleparemeteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexec"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
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

	_, err = commandexecutorexec.RunCommand(
		commandexecutorgeneric.WithLiveOutputOnStdoutIfVerbose(ctx),
		&parameteroptions.RunCommandOptions{
			Command: []string{binPath, playbookPath, "--limit", limit, "--inventory="+limit+","},
		},
	)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Run ansible playbook '%s' using '%s' against hosts '%s' finished.", playbookPath, binPath, limit)

	return nil
}
