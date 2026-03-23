package testsuite

import (
	"context"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/testutils/testcase"
	"github.com/asciich/asciichgolangpublic/pkg/testutils/testresults"
	"github.com/asciich/asciichgolangpublic/pkg/testutils/testutilsinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"gopkg.in/yaml.v3"
)

type TestSuite struct {
	Name        string              `yaml:"name"`
	Description string              `yaml:"description"`
	TestCases   []*testcase.TestCase `yaml:"test_cases"`
}

func LoadFromFile(ctx context.Context, path string) (testutilsinterfaces.TestSuite, error) {
	if path == "" {
		return nil, tracederrors.TracedErrorEmptyString("path")
	}

	logging.LogInfoByCtxf(ctx, "Load test suite from '%s' started.", path)

	content, err := nativefiles.ReadAsBytes(ctx, path)
	if err != nil {
		return nil, err
	}

	testSuite, err := LoadFromBytes(ctx, content)
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Load test suite from '%s' started.", path)

	return testSuite, nil
}

func LoadFromBytes(ctx context.Context, testSuiteData []byte) (testutilsinterfaces.TestSuite, error) {
	if testSuiteData == nil {
		return nil, tracederrors.TracedErrorNil("testSuiteData")
	}

	testSuite := &TestSuite{}

	err := yaml.Unmarshal(testSuiteData, testSuite)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to unmarshal bytes as test suite: %w", err)
	}

	return testSuite, nil
}

func (t *TestSuite) GetName() (string, error) {
	if t.Name == "" {
		return "", tracederrors.TracedError("name not set")
	}

	return t.Name, nil
}

func (t *TestSuite) Run(ctx context.Context) (testutilsinterfaces.TestResult, error) {
	tStart := time.Now()

	name, err := t.GetName()
	if err != nil {
		return nil, err
	}

	result := &testresults.TestResult{}


	logging.LogInfoByCtxf(ctx, "Run test suite '%s' started.", name)

	if len(t.TestCases) <= 0 {
		return nil, tracederrors.TracedErrorf("TestSuite '%s' has no test cases.", name)
	}

	totalTestCases := len(t.TestCases)

	for i, testCase := range t.TestCases {
		tcName, err := testCase.GetName()
		if err != nil {
			return nil, err
		}

		logging.LogInfoByCtxf(ctx, "Run test case '%s' of test suite '%s' (%d/%d).", tcName, name, i+1, totalTestCases)

		testCaseResult, err := testCase.Run(ctx)
		if err != nil {
			return nil, err
		}

		err = result.AddTestCaseResult(testCaseResult)
		if err != nil {
			return nil, err
		}
	}

	tEnd := time.Now()

	err = result.SetName(name)
	if err != nil {
		return nil, err
	}

	err = result.SetTimeStart(&tStart)
	if err != nil {
		return nil, err
	}

	err = result.SetTimeEnd(&tEnd)
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Run test suite '%s' finished.", name)

	return result, nil
}
