package testresults

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/testutils/testutilsgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/testutils/testutilsinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type TestResult struct {
	TestResultBase

	TestCaseResults []testutilsinterfaces.TestResult
}

func (t *TestResult) GetNPassed(ctx context.Context) (int, error) {
	if len(t.TestCaseResults) <= 0 {
		return 0, tracederrors.TracedError(testutilsgeneric.ErrNoTestCaseResultsSet)
	}

	var nPassed int
	for _, tcr := range t.TestCaseResults {
		isPassed, err := tcr.IsPassed(contextutils.WithSilent(ctx))
		if err != nil {
			return 0, err
		}

		if isPassed {
			nPassed++
		}
	}

	return nPassed, nil
}

func (t *TestResult) GetNFailed(ctx context.Context) (int, error) {
	if len(t.TestCaseResults) <= 0 {
		return 0, tracederrors.TracedError(testutilsgeneric.ErrNoTestCaseResultsSet)
	}

	var nFailed int
	for _, tcr := range t.TestCaseResults {
		isPassed, err := tcr.IsPassed(contextutils.WithSilent(ctx))
		if err != nil {
			return 0, err
		}

		if !isPassed {
			nFailed++
		}
	}

	return nFailed, nil
}

func (t *TestResult) IsPassed(ctx context.Context) (bool, error) {
	if len(t.TestCaseResults) <= 0 {
		return false, tracederrors.TracedError(testutilsgeneric.ErrNoTestCaseResultsSet)
	}

	for _, tcr := range t.TestCaseResults {
		isPassed, err := tcr.IsPassed(contextutils.WithSilent(ctx))
		if err != nil {
			return false, err
		}

		if !isPassed {
			return false, nil
		}
	}

	return true, nil
}

func (t *TestResult) AddTestCaseResult(testCaseResult testutilsinterfaces.TestResult) error {
	if testCaseResult == nil {
		return tracederrors.TracedErrorNil("testCaseResult")
	}

	t.TestCaseResults = append(t.TestCaseResults, testCaseResult)

	return nil
}

func (t *TestResult) LogResult(ctx context.Context) error {
	name, err := t.GetName()
	if err != nil {
		return err
	}

	nPassed, err := t.GetNPassed(ctx)
	if err != nil {
		return err
	}

	nFailed, err := t.GetNFailed(ctx)
	if err != nil {
		return err
	}

	duration, err := t.GetDuration(ctx)
	if err != nil {
		return err
	}

	if nFailed > 0 {
		logging.LogErrorByCtxf(ctx, "%d out of %d test cases of '%s' failed in %s", nFailed, nPassed+nFailed, name, duration)
	} else {
		logging.LogGoodByCtxf(ctx, "All %d test cases of '%s' passed in %s.", nPassed, name, duration)
	}

	return nil
}
