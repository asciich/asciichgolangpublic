package testutilsinterfaces

import "context"

type TestCaseExecutor interface {
	GetName() (string, error)
	Run(ctx context.Context) (TestResult, error)

	// Set the needed data to run the test case.
	// Usually the TestCase struct is passed here.
	SetData(data any) error
}
