package testutilsinterfaces

import "context"

type TestSuite interface {
	GetName() (string, error)

	Run(ctx context.Context) (TestResult, error)
}
