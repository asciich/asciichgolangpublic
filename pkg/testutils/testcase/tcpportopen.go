package testcase

import (
	"context"
	"fmt"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/netutils"
	"github.com/asciich/asciichgolangpublic/pkg/testutils/testresults"
	"github.com/asciich/asciichgolangpublic/pkg/testutils/testutilsinterfaces"
)

type TestCaseExecutorTcpPortOpen struct {
	TestCaseExecutorBase
}

func (t *TestCaseExecutorTcpPortOpen) GetName() (string, error) {
	return "tcp_port_open", nil
}

func (t *TestCaseExecutorTcpPortOpen) Run(ctx context.Context) (testutilsinterfaces.TestResult, error) {
	tStart := time.Now()

	name, err := t.GetTestCaseName()
	if err != nil {
		return nil, err
	}

	result := &testresults.TestCaseResult{
		Name: name,
	}

	port, err := t.GetPort()
	if err != nil {
		return nil, err
	}

	host, err := t.GetHost()
	if err != nil {
		return nil, err
	}

	isOpen, err := netutils.IsTcpPortOpen(ctx, host, port)
	if err != nil {
		return nil, err
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
		err := result.SetFailedMessage(
			fmt.Sprintf("The TCP port '%d' on '%s' is not open.", port, host),
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
