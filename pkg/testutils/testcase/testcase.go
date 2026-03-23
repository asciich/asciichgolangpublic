package testcase

import (
	"context"
	"strconv"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/testutils/testutilsinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type TestCase struct {
	Name        string `yaml:"name"`
	TestType    string `yaml:"test_type"`
	Command     string `yaml:"command,omitempty"`
	Description string `yaml:"description"`
	Port        string `yaml:"port,omitempty"`
	Host        string `yaml:"host,omitempty"`

	data any
}

func (t *TestCase) GetName() (string, error) {
	if t.Name == "" {
		return "", tracederrors.TracedError("name not set")
	}

	return t.Name, nil
}

func (t *TestCase) GetTestType() (string, error) {
	if t.TestType == "" {
		return "", tracederrors.TracedError("test type not set")
	}

	return t.TestType, nil
}

func (t *TestCase) GetCommand() (string, error) {
	if t.Command == "" {
		return "", tracederrors.TracedError("command not set")
	}

	return t.Command, nil
}

func (t *TestCase) GetHost() (string, error) {
	if t.Host == "" {
		return "", tracederrors.TracedError("host not set")
	}

	return t.Host, nil
}

func (t *TestCase) GetPort() (int, error) {
	if t.Port == "" {
		return 0, tracederrors.TracedError("port not set")
	}

	port, err := strconv.Atoi(t.Port)
	if err != nil {
		return 0, tracederrors.TracedErrorf("Failed to convert the given port '%s' to an int", t.Port)
	}

	return port, nil
}

func (t *TestCase) Run(ctx context.Context) (testutilsinterfaces.TestResult, error) {
	name, err := t.GetName()
	if err != nil {
		return nil, err
	}

	testType, err := t.GetTestType()
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Run test case '%s' of type '%s' started.", name, testType)

	executor, err := GetTestCaseExecutorByTestType(testType, t)
	if err != nil {
		return nil, err
	}

	result, err := executor.Run(ctx)
	if err != nil {
		return nil, err
	}

	result.LogResult(ctx)

	logging.LogInfoByCtxf(ctx, "Run test case '%s' of type '%s' finished.", name, testType)

	return result, nil
}

func (t *TestCase) SetData(data any) error {
	if data == nil {
		return tracederrors.TracedErrorNil("data")
	}

	t.data = data

	return nil
}
