package testutilsinterfaces

import (
	"context"
	"time"
)

type TestResult interface {
	GetName() (string, error)

	GetNFailed(context.Context) (int, error)
	GetNPassed(context.Context) (int, error)
	IsPassed(context.Context) (bool, error)

	LogResult(ctx context.Context) error

	SetTimeStart(*time.Time) error
	SetTimeEnd(*time.Time) error

	GetDuration(context.Context) (time.Duration, error)
}
