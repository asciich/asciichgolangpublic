package testresults

import (
	"context"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type TestResultBase struct {
	Name      string
	timeStart *time.Time
	timeEnd   *time.Time
}

func (t *TestResult) GetName() (string, error) {
	if t.Name == "" {
		return "", tracederrors.TracedError("Name not set")
	}

	return t.Name, nil
}

func (t *TestResult) SetName(name string) error {
	if name == "" {
		return tracederrors.TracedErrorEmptyString("name")
	}

	t.Name = name

	return nil
}

func (t *TestResultBase) SetTimeStart(tStart *time.Time) error {
	if tStart == nil {
		return tracederrors.TracedErrorNil("tStart")
	}

	t.timeStart = tStart

	return nil
}

func (t *TestResultBase) SetTimeEnd(tEnd *time.Time) error {
	if tEnd == nil {
		return tracederrors.TracedErrorNil("tStart")
	}

	t.timeEnd = tEnd

	return nil
}

func (t *TestResultBase) GetTimeStart() (*time.Time, error) {
	if t.timeStart == nil {
		return nil, tracederrors.TracedError("timeStart not set")
	}

	return t.timeStart, nil
}

func (t *TestResultBase) GetTimeEnd() (*time.Time, error) {
	if t.timeEnd == nil {
		return nil, tracederrors.TracedError("timeEnd not set")
	}

	return t.timeEnd, nil
}

func (t *TestResultBase) GetDuration(ctx context.Context) (time.Duration, error) {
	tStart, err := t.GetTimeStart()
	if err != nil {
		return 0, err
	}

	tEnd, err := t.GetTimeEnd()
	if err != nil {
		return 0, err
	}

	return tEnd.Sub(*tStart), nil
}
