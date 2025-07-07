package tracederrors_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

func TestTracedErrorIsError(t *testing.T) {
	require := require.New(t)

	var err error = tracederrors.TracedError("example error")
	_, ok := err.(tracederrors.TracedErrorType)
	require.True(ok)
}

func TestTracedErrorIsTracedError(t *testing.T) {
	require := require.New(t)

	var err error = tracederrors.TracedError("example error")
	require.True(errors.Is(err, tracederrors.ErrTracedError))
}

func TestTracedErrorWrap(t *testing.T) {
	require := require.New(t)

	exampleError := errors.New("exampleError")

	var errFmt error = fmt.Errorf("%w", exampleError)
	require.True(errors.Is(errFmt, exampleError))

	var tracedErrorWithWrapFormatted error = tracederrors.TracedErrorf("%w", exampleError)
	require.True(errors.Is(tracedErrorWithWrapFormatted, exampleError))

	var tracedErrorWithWrap error = tracederrors.TracedError(exampleError)
	require.True(errors.Is(tracedErrorWithWrap, exampleError))

	var wrappedAgain error = fmt.Errorf("again: %w", tracedErrorWithWrap)
	require.True(errors.Is(wrappedAgain, exampleError))
	require.True(errors.Is(wrappedAgain, tracederrors.ErrTracedError))

	var wrappedAgain2 error = fmt.Errorf("again2: %w", wrappedAgain)
	require.True(errors.Is(wrappedAgain2, exampleError))
	require.True(errors.Is(wrappedAgain2, tracederrors.ErrTracedError))
}

func testFunctionRaisingError(errorMessage string) (err error) {
	return tracederrors.TracedError(errorMessage)
}

func TestTracedErrorStackTraceInMessage(t *testing.T) {
	tests := []struct {
		testmessage string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				require := require.New(t)

				err := testFunctionRaisingError(tt.testmessage)

				require.Contains(err.Error(), tt.testmessage)
				require.Contains(err.Error(), "testFunctionRaisingError")
			},
		)
	}
}

func TestTracedErrorEmptyString(t *testing.T) {

	tests := []struct {
		stringName string
	}{
		{"varName"},
		{"AnoterVarName"},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				require := require.New(t)

				var err error = tracederrors.TracedErrorEmptyString(tt.stringName)
				require.Contains(err.Error(), "'"+tt.stringName+"' is empty string")
				require.True(tracederrors.IsTracedError(err))
				require.True(tracederrors.IsEmptyStringError(err))
				require.False(tracederrors.IsNilError(err))
				require.False(tracederrors.IsNotImplementedError(err))
			},
		)
	}
}

func TestTracedErrorNil(t *testing.T) {

	tests := []struct {
		stringName string
	}{
		{"varName"},
		{"AnoterVarName"},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				require := require.New(t)

				var err error = tracederrors.TracedErrorNil(tt.stringName)
				require.Contains(err.Error(), "'"+tt.stringName+"' is nil")
				require.True(tracederrors.IsTracedError(err))
				require.True(tracederrors.IsNilError(err))
				require.False(tracederrors.IsEmptyStringError(err))
				require.False(tracederrors.IsNotImplementedError(err))
			},
		)
	}
}

func TestTracedErrorNotImplemented(t *testing.T) {

	tests := []struct {
		stringName string
	}{
		{"varName"},
		{"AnoterVarName"},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				require := require.New(t)

				var err error = tracederrors.TracedErrorNotImplemented()
				require.Contains(err.Error(), "Not implemented")
				require.True(tracederrors.IsTracedError(err))
				require.True(tracederrors.IsNotImplementedError(err))
				require.False(tracederrors.IsNilError(err))
				require.False(tracederrors.IsEmptyStringError(err))
			},
		)
	}
}

func TestTracedErrorGetErrorMessage(t *testing.T) {
	tests := []struct {
		errorMessage         string
		expectedErrorMessage string
	}{
		{"errorMessage", "errorMessage"},
		{"errorMessage2", "errorMessage2"},
		{"errorMessage3", "errorMessage3"},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				tracedError, err := tracederrors.GetAsTracedError(tracederrors.TracedError(tt.errorMessage))
				require.NoError(t, err)

				require.EqualValues(t, tt.expectedErrorMessage, mustutils.Must(tracedError.GetErrorMessage()))
			},
		)
	}
}
