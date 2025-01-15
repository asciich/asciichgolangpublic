package tracederrors

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTracedErrorIsError(t *testing.T) {
	assert := assert.New(t)

	var err error = TracedError("example error")
	_, ok := err.(TracedErrorType)
	assert.True(ok)
}

func TestTracedErrorIsTracedError(t *testing.T) {
	assert := assert.New(t)

	var err error = TracedError("example error")
	assert.True(errors.Is(err, ErrTracedError))
}

func TestTracedErrorWrap(t *testing.T) {
	assert := assert.New(t)

	exampleError := errors.New("exampleError")

	var errFmt error = fmt.Errorf("%w", exampleError)
	assert.True(errors.Is(errFmt, exampleError))

	var tracedErrorWithWrapFormatted error = TracedErrorf("%w", exampleError)
	assert.True(errors.Is(tracedErrorWithWrapFormatted, exampleError))

	var tracedErrorWithWrap error = TracedError(exampleError)
	assert.True(errors.Is(tracedErrorWithWrap, exampleError))

	var wrappedAgain error = fmt.Errorf("again: %w", tracedErrorWithWrap)
	assert.True(errors.Is(wrappedAgain, exampleError))
	assert.True(errors.Is(wrappedAgain, ErrTracedError))

	var wrappedAgain2 error = fmt.Errorf("again2: %w", wrappedAgain)
	assert.True(errors.Is(wrappedAgain2, exampleError))
	assert.True(errors.Is(wrappedAgain2, ErrTracedError))
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
				assert := assert.New(t)

				err := testFunctionRaisingError(tt.testmessage)

				assert.Contains(err.Error(), tt.testmessage)
				assert.Contains(err.Error(), "testFunctionRaisingError")
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
				assert := assert.New(t)

				var err error = TracedErrorEmptyString(tt.stringName)
				assert.Contains(err.Error(), "'"+tt.stringName+"' is empty string")
				assert.True(IsTracedError(err))
				assert.True(IsEmptyStringError(err))
				assert.False(IsNilError(err))
				assert.False(IsNotImplementedError(err))
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
				assert := assert.New(t)

				var err error = TracedErrorNil(tt.stringName)
				assert.Contains(err.Error(), "'"+tt.stringName+"' is nil")
				assert.True(IsTracedError(err))
				assert.True(IsNilError(err))
				assert.False(IsEmptyStringError(err))
				assert.False(IsNotImplementedError(err))
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
				assert := assert.New(t)

				var err error = TracedErrorNotImplemented()
				assert.Contains(err.Error(), "Not implemented")
				assert.True(IsTracedError(err))
				assert.True(IsNotImplementedError(err))
				assert.False(IsNilError(err))
				assert.False(IsEmptyStringError(err))
			},
		)
	}
}


func TestTracedErrorGetErrorMessage(t *testing.T) {
	tests := []struct {
		errorMessage                string
		expectedErrorMessage  string
	}{
		{"errorMessage", "errorMessage"},
		{"errorMessage2", "errorMessage2"},
		{"errorMessage3", "errorMessage3"},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				assert := assert.New(t)
				
				tracedError := MustGetAsTracedError(TracedError(tt.errorMessage))

				assert.EqualValues(
					tt.expectedErrorMessage,
					tracedError.MustGetErrorMessage(),
				)
			},
		)
	}
}