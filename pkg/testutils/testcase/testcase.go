package testcase

import (
	"context"
	"strconv"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/testutils/testutilsinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type TestCase struct {
	Name         string `yaml:"name"`
	TestType     string `yaml:"test_type"`
	Command      string `yaml:"command,omitempty"`
	Description  string `yaml:"description"`
	Port         string `yaml:"port,omitempty"`
	Host         string `yaml:"host,omitempty"`
	Namespace    string `yaml:"namespace,omitempty"`
	Cluster      string `yaml:"cluster,omitempty"`
	ResourceName string `yaml:"resource_name,omitempty"`
	SecretKey    string `yaml:"secret_key" json:"secret_key"`
	TargetHost   string `yaml:"target_host" json:"target_host"`
	TargetUser   string `yaml:"target_user" json:"target_user"`
	TargetPort   int    `yaml:"target_port" json:"target_port"`

	data            any
	commandExecutor commandexecutorinterfaces.CommandExecutor
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

func (t *TestCase) GetNamespace() (string, error) {
	if t.Namespace == "" {
		return "", tracederrors.TracedError("namespace not set")
	}

	return t.Namespace, nil
}

func (t *TestCase) GetCluster() (string, error) {
	if t.Cluster == "" {
		return "", tracederrors.TracedError("cluster not set")
	}

	return t.Cluster, nil
}

func (t *TestCase) GetResourceName() (string, error) {
	if t.ResourceName == "" {
		return "", tracederrors.TracedError("resource_name not set")
	}

	return t.ResourceName, nil
}

func (t *TestCase) SetCommandExecutor(commandExecutor commandexecutorinterfaces.CommandExecutor) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	t.commandExecutor = commandExecutor

	return nil
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

	if t.commandExecutor == nil {
		return nil, tracederrors.TracedError("commandExecutor not set on TestCase. TestSuite must call SetCommandExecutor before Run.")
	}

	result, err := executor.Run(ctx, t.commandExecutor)
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

func (t *TestCase) GetSecretKey() (string, error) {
	if t.SecretKey == "" {
		return "", tracederrors.TracedErrorEmptyString("SecretKey")
	}
	return t.SecretKey, nil
}

func (t *TestCase) GetTargetHost() (string, error) {
	if t.TargetHost == "" {
		return "", tracederrors.TracedErrorEmptyString("TargetHost")
	}
	return t.TargetHost, nil
}

func (t *TestCase) GetTargetUser() (string, error) {
	if t.TargetUser == "" {
		return "", tracederrors.TracedErrorEmptyString("TargetUser")
	}
	return t.TargetUser, nil
}

func (t *TestCase) GetTargetPort() (int, error) {
	if t.TargetPort <= 0 {
		return 0, tracederrors.TracedError("TargetPort not set")
	}
	return t.TargetPort, nil
}
