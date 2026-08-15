package testcase

import (
	"github.com/asciich/asciichgolangpublic/pkg/datatypes"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type TestCaseExecutorBase struct {
	data any
}

func (t *TestCaseExecutorBase) SetData(data any) error {
	if data == nil {
		return tracederrors.TracedErrorNil("data")
	}

	t.data = data

	return nil
}

func (t *TestCaseExecutorBase) GetDataAsTestCase() (*TestCase, error) {
	if t.data == nil {
		return nil, tracederrors.TracedError("data not set")
	}

	testCase, ok := t.data.(*TestCase)
	if !ok {
		typeName, _ := datatypes.GetTypeName(t.data)
		return nil, tracederrors.TracedErrorf("data is not a TestCase it's of type '%s'", typeName)
	}

	return testCase, nil
}

func (t *TestCaseExecutorBase) GetTestCaseName() (string, error) {
	testCase, err := t.GetDataAsTestCase()
	if err != nil {
		return "", err
	}

	return testCase.GetName()
}

func (t *TestCaseExecutorBase) GetCommand() (string, error) {
	testCase, err := t.GetDataAsTestCase()
	if err != nil {
		return "", err
	}

	return testCase.GetCommand()
}

func (t *TestCaseExecutorBase) GetHost() (string, error) {
	testCase, err := t.GetDataAsTestCase()
	if err != nil {
		return "", err
	}

	return testCase.GetHost()
}

func (t *TestCaseExecutorBase) GetPort() (int, error) {
	testCase, err := t.GetDataAsTestCase()
	if err != nil {
		return 0, err
	}

	return testCase.GetPort()
}

func (t *TestCaseExecutorBase) GetNamespace() (string, error) {
	testCase, err := t.GetDataAsTestCase()
	if err != nil {
		return "", err
	}

	return testCase.GetNamespace()
}

func (t *TestCaseExecutorBase) GetCluster() (string, error) {
	testCase, err := t.GetDataAsTestCase()
	if err != nil {
		return "", err
	}

	return testCase.GetCluster()
}

func (t *TestCaseExecutorBase) GetResourceName() (string, error) {
	testCase, err := t.GetDataAsTestCase()
	if err != nil {
		return "", err
	}

	return testCase.GetResourceName()
}

func (t *TestCaseExecutorBase) GetSecretKey() (string, error) {
	testCase, err := t.GetDataAsTestCase()
	if err != nil {
		return "", err
	}

	return testCase.GetSecretKey()
}

func (t *TestCaseExecutorBase) GetTargetHost() (string, error) {
	testCase, err := t.GetDataAsTestCase()
	if err != nil {
		return "", err
	}

	return testCase.GetTargetHost()
}

func (t *TestCaseExecutorBase) GetTargetUser() (string, error) {
	testCase, err := t.GetDataAsTestCase()
	if err != nil {
		return "", err
	}

	return testCase.GetTargetUser()
}

func (t *TestCaseExecutorBase) GetTargetPort() (int, error) {
	testCase, err := t.GetDataAsTestCase()
	if err != nil {
		return 0, err
	}

	return testCase.GetTargetPort()
}
