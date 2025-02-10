package tracederrors

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTracedErrorIsError(t *testing.T) {
	require := require.New(t)

	var err error = TracedError("example error")
	_, ok := err.(TracedErrorType)
	require.True(ok)
}

func TestTracedErrorIsTracedError(t *testing.T) {
	require := require.New(t)

	var err error = TracedError("example error")
	require.True(errors.Is(err, ErrTracedError))
}

func TestTracedErrorWrap(t *testing.T) {
	require := require.New(t)

	exampleError := errors.New("exampleError")

	var errFmt error = fmt.Errorf("%w", exampleError)
	require.True(errors.Is(errFmt, exampleError))

	var tracedErrorWithWrapFormatted error = TracedErrorf("%w", exampleError)
	require.True(errors.Is(tracedErrorWithWrapFormatted, exampleError))

	var tracedErrorWithWrap error = TracedError(exampleError)
	require.True(errors.Is(tracedErrorWithWrap, exampleError))

	var wrappedAgain error = fmt.Errorf("again: %w", tracedErrorWithWrap)
	require.True(errors.Is(wrappedAgain, exampleError))
	require.True(errors.Is(wrappedAgain, ErrTracedError))

	var wrappedAgain2 error = fmt.Errorf("again2: %w", wrappedAgain)
	require.True(errors.Is(wrappedAgain2, exampleError))
	require.True(errors.Is(wrappedAgain2, ErrTracedError))
}

func testFunctionRaisingError(errorMessage string) (err error) {
	return TracedError(errorMessage)
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

				var err error = TracedErrorEmptyString(tt.stringName)
				require.Contains(err.Error(), "'"+tt.stringName+"' is empty string")
				require.True(IsTracedError(err))
				require.True(IsEmptyStringError(err))
				require.False(IsNilError(err))
				require.False(IsNotImplementedError(err))
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

				var err error = TracedErrorNil(tt.stringName)
				require.Contains(err.Error(), "'"+tt.stringName+"' is nil")
				require.True(IsTracedError(err))
				require.True(IsNilError(err))
				require.False(IsEmptyStringError(err))
				require.False(IsNotImplementedError(err))
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

				var err error = TracedErrorNotImplemented()
				require.Contains(err.Error(), "Not implemented")
				require.True(IsTracedError(err))
				require.True(IsNotImplementedError(err))
				require.False(IsNilError(err))
				require.False(IsEmptyStringError(err))
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
				require := require.New(t)

				tracedError := MustGetAsTracedError(TracedError(tt.errorMessage))

				require.EqualValues(
					tt.expectedErrorMessage,
					tracedError.MustGetErrorMessage(),
				)
			},
		)
	}
}
