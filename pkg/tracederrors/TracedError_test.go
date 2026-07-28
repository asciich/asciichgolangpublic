package tracederrors_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func TestTracedErrorIsError(t *testing.T) {
	var err error = tracederrors.TracedError("example error")
	_, ok := err.(tracederrors.TracedErrorType)
	require.True(t, ok)
}

func TestTracedErrorIsTracedError(t *testing.T) {
	var err error = tracederrors.TracedError("example error")
	require.True(t, errors.Is(err, tracederrors.ErrTracedError))
}

func TestTracedErrorWrap(t *testing.T) {
	exampleError := errors.New("exampleError")

	var errFmt error = fmt.Errorf("%w", exampleError)
	require.True(t, errors.Is(errFmt, exampleError))

	var tracedErrorWithWrapFormatted error = tracederrors.TracedErrorf("%w", exampleError)
	require.True(t, errors.Is(tracedErrorWithWrapFormatted, exampleError))

	var tracedErrorWithWrap error = tracederrors.TracedError(exampleError)
	require.True(t, errors.Is(tracedErrorWithWrap, exampleError))

	var wrappedAgain error = fmt.Errorf("again: %w", tracedErrorWithWrap)
	require.True(t, errors.Is(wrappedAgain, exampleError))
	require.True(t, errors.Is(wrappedAgain, tracederrors.ErrTracedError))

	var wrappedAgain2 error = fmt.Errorf("again2: %w", wrappedAgain)
	require.True(t, errors.Is(wrappedAgain2, exampleError))
	require.True(t, errors.Is(wrappedAgain2, tracederrors.ErrTracedError))
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
				err := testFunctionRaisingError(tt.testmessage)

				require.Contains(t, err.Error(), tt.testmessage)
				require.Contains(t, err.Error(), "testFunctionRaisingError")
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
				var err error = tracederrors.TracedErrorEmptyString(tt.stringName)
				require.Contains(t, err.Error(), "'"+tt.stringName+"' is empty string")
				require.True(t, tracederrors.IsTracedError(err))
				require.True(t, tracederrors.IsEmptyStringError(err))
				require.False(t, tracederrors.IsNilError(err))
				require.False(t, tracederrors.IsNotImplementedError(err))
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
				var err error = tracederrors.TracedErrorNil(tt.stringName)
				require.Contains(t, err.Error(), "'"+tt.stringName+"' is nil")
				require.True(t, tracederrors.IsTracedError(err))
				require.True(t, tracederrors.IsNilError(err))
				require.False(t, tracederrors.IsEmptyStringError(err))
				require.False(t, tracederrors.IsNotImplementedError(err))
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
				var err error = tracederrors.TracedErrorNotImplemented()
				require.Contains(t, err.Error(), "Not implemented")
				require.True(t, tracederrors.IsTracedError(err))
				require.True(t, tracederrors.IsNotImplementedError(err))
				require.False(t, tracederrors.IsNilError(err))
				require.False(t, tracederrors.IsEmptyStringError(err))
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

				errMsg, err := tracedError.GetErrorMessage()
				require.NoError(t, err)
				require.EqualValues(t, tt.expectedErrorMessage, errMsg)
			},
		)
	}
}
