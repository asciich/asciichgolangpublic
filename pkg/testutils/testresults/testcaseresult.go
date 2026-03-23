package testresults

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type TestCaseResult struct {
	TestResultBase

	Name string

	SuccessMessage string
	FailedMessage  string
}

func (t *TestCaseResult) GetName() (string, error) {
	if t.Name == "" {
		return "", tracederrors.TracedErrorEmptyString("Name")
	}

	return t.Name, nil
}

func (t *TestCaseResult) GetNFailed(ctx context.Context) (int, error) {
	if t.FailedMessage != "" {
		return 1, nil
	}

	return 0, nil
}

func (t *TestCaseResult) GetNPassed(ctx context.Context) (int, error) {
	if t.FailedMessage == "" {
		if t.SuccessMessage != "" {
			return 1, nil
		}
	}

	return 0, nil
}

func (t *TestCaseResult) IsPassed(ctx context.Context) (bool, error) {
	if t.FailedMessage == "" {
		if t.SuccessMessage != "" {
			return true, nil
		}
	}

	return false, nil
}

func (t *TestCaseResult) SetSuccessMessage(msg string) error {
	if msg == "" {
		return tracederrors.TracedErrorEmptyString("msg")
	}

	t.SuccessMessage = msg
	return nil
}

func (t *TestCaseResult) SetFailedMessage(msg string) error {
	if msg == "" {
		return tracederrors.TracedErrorEmptyString("msg")
	}

	t.FailedMessage = msg
	return nil
}

func (t *TestCaseResult) LogResult(ctx context.Context) error {
	if t.SuccessMessage == "" && t.FailedMessage == "" {
		return tracederrors.TracedError("Test result invalid. Boths SuccessMessage and FailedMessage are not set")
	}

	if t.SuccessMessage != "" && t.FailedMessage != "" {
		return tracederrors.TracedError("Test result invalid. Boths SuccessMessage and FailedMessage are set")
	}

	name, err := t.GetName()
	if err != nil {
		return err
	}

	duration, err := t.GetDuration(ctx)
	if err != nil {
		return err
	}

	if t.SuccessMessage != "" {
		logging.LogGoodByCtxf(ctx, "TestCase '%s' in %s passed: %s", name, duration, t.SuccessMessage)
	} else {
		logging.LogErrorByCtxf(ctx, "TestCase '%s' in %s failed: %s", name, duration, t.FailedMessage)
	}

	return nil
}
