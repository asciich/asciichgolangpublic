package testcase

import (
	"context"
	"fmt"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/shellutils/shelllinehandler"
	"github.com/asciich/asciichgolangpublic/pkg/testutils/testresults"
	"github.com/asciich/asciichgolangpublic/pkg/testutils/testutilsinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type TestCaseExecutorCommand struct {
	TestCaseExecutorBase
}

func (t *TestCaseExecutorCommand) GetName() (string, error) {
	return "command", nil
}

func (t *TestCaseExecutorCommand) Run(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor) (testutilsinterfaces.TestResult, error) {
	tStart := time.Now()

	name, err := t.GetTestCaseName()
	if err != nil {
		return nil, err
	}

	result := &testresults.TestCaseResult{
		Name: name,
	}

	if commandExecutor == nil {
		return nil, tracederrors.TracedErrorNil("commandExecutor")
	}

	command, err := t.GetCommand()
	if err != nil {
		return nil, err
	}

	splitted, err := shelllinehandler.Split(command)
	if err != nil {
		return nil, err
	}

	output, err := commandExecutor.RunCommand(ctx, &parameteroptions.RunCommandOptions{
		Command:           splitted,
		AllowAllExitCodes: true,
	})
	if err != nil {
		return nil, err
	}

	exitCode, err := output.GetReturnCode()
	if err != nil {
		return nil, err
	}

	tEnd := time.Now()

	if output.IsExitSuccess() {
		err := result.SetSuccessMessage(
			fmt.Sprintf("The test command '%s' was executed successfully.", command),
		)
		if err != nil {
			return nil, err
		}
	} else {
		err := result.SetFailedMessage(
			fmt.Sprintf("The test command '%s' failed and exited with code %d.", command, exitCode),
		)
		if err != nil {
			return nil, err
		}
	}

	err = result.SetTimeStart(&tStart)
	if err != nil {
		return nil, err
	}

	err = result.SetTimeEnd(&tEnd)
	if err != nil {
		return nil, err
	}

	return result, nil
}
