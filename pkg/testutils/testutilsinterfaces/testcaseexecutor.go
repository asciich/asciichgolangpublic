package testutilsinterfaces

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
)

type TestCaseExecutor interface {
	GetName() (string, error)
	Run(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor) (TestResult, error)

	// Set the needed data to run the test case.
	// Usually the TestCase struct is passed here.
	SetData(data any) error
}
