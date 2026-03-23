package testutilsinterfaces

import (
	"context"
)

type TestCase interface {
	GetName() (string, error)
	Run(ctx context.Context) (TestResult, error)
}
