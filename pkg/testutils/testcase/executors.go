package testcase

import (
	"github.com/asciich/asciichgolangpublic/pkg/testutils/testutilsinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func LoadExecuctors() ([]testutilsinterfaces.TestCaseExecutor, error) {
	return []testutilsinterfaces.TestCaseExecutor{
		&TestCaseExecutorCommand{},
		&TestCaseExecutorTcpPortOpen{},
	}, nil
}

func GetTestCaseExecutorByTestType(testType string, data any) (testutilsinterfaces.TestCaseExecutor, error) {
	if testType == "" {
		return nil, tracederrors.TracedErrorEmptyString("testType")
	}

	executors, err := LoadExecuctors()
	if err != nil {
		return nil, err
	}

	for _, executor := range executors {
		name, err := executor.GetName()
		if err != nil {
			return nil, err
		}

		if name == testType {
			err = executor.SetData(data)
			if err != nil {
				return nil, err
			}

			return executor, nil
		}
	}

	return nil, tracederrors.TracedErrorf("Test executor not found for testType='%s'", testType)
}
