package testsuite

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/testutils/testutilsinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/testutils/testutilsoptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func RunFromFilePath(ctx context.Context, path string, options *testutilsoptions.RunTestSuiteOptions) (testutilsinterfaces.TestResult, error) {
	if path == "" {
		return nil, tracederrors.TracedErrorEmptyString("path")
	}

	logging.LogInfoByCtxf(ctx, "Run test suite from '%s' started.", path)

	testSuite, err := LoadFromFile(ctx, path)
	if err != nil {
		return nil, err
	}

	result, err := testSuite.Run(ctx)
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Run test suite from '%s' finished.", path)

	return result, nil
}
