package testsuite

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/testutils/testcase"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func Test_FailsIfNoNameSet(t *testing.T) {
	ctx := getCtx()

	testSuite := &TestSuite{}

	result, err := testSuite.Run(ctx)
	require.Error(t, err)
	require.Nil(t, result)
}

func Test_FailsIfNoTestCasesSet(t *testing.T) {
	ctx := getCtx()

	testSuite := &TestSuite{
		Name: "example test suite",
	}

	result, err := testSuite.Run(ctx)
	require.Error(t, err)
	require.Nil(t, result)
}

func Test_EchoHelloWorldCommand(t *testing.T) {
	ctx := getCtx()
	testSuite := &TestSuite{
		Name: "hello world",
		TestCases: []*testcase.TestCase{
			{
				Name:     "hello world",
				TestType: "command",
				Command:  "echo hello world",
			},
		},
	}

	result, err := testSuite.Run(ctx)
	require.NoError(t, err)

	isPassed, err := result.IsPassed(ctx)
	require.NoError(t, err)
	require.True(t, isPassed)
}
