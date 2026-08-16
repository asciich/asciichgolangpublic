package testcase

import (
	"context"
	"fmt"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/netutils"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/testutils/testresults"
	"github.com/asciich/asciichgolangpublic/pkg/testutils/testutilsinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type TestCaseExecutorTcpPortOpen struct {
	TestCaseExecutorBase
}

func (t *TestCaseExecutorTcpPortOpen) GetName() (string, error) {
	return "tcp_port_open", nil
}

func (t *TestCaseExecutorTcpPortOpen) Run(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor) (testutilsinterfaces.TestResult, error) {
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

	port, err := t.GetPort()
	if err != nil {
		return nil, err
	}

	host, err := t.GetHost()
	if err != nil {
		return nil, err
	}

	// Check if running on localhost or remote
	isLocalhost, err := commandExecutor.IsRunningOnLocalhost()
	if err != nil {
		return nil, err
	}

	var isOpen bool
	if isLocalhost {
		// Use native Go implementation for localhost
		isOpen, err = netutils.IsTcpPortOpen(ctx, host, port)
		if err != nil {
			return nil, err
		}
	} else {
		// Use commandExecutor for remote execution via nc or similar
		output, err := commandExecutor.RunCommand(ctx, &parameteroptions.RunCommandOptions{
			Command:           []string{"nc", "-z", "-w", "5", host, fmt.Sprintf("%d", port)},
			AllowAllExitCodes: true,
		})
		if err != nil {
			return nil, err
		}
		isOpen = output.IsExitSuccess()
	}

	tEnd := time.Now()

	if isOpen {
		err := result.SetSuccessMessage(
			fmt.Sprintf("The TCP port '%d' on '%s' is open.", port, host),
		)
		if err != nil {
			return nil, err
		}
	} else {
		baseMessage := fmt.Sprintf("The TCP port '%d' on '%s' is not open.", port, host)
		failedMessage, err := t.FormatFailedMessage(baseMessage)
		if err != nil {
			return nil, err
		}
		err = result.SetFailedMessage(failedMessage)
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
