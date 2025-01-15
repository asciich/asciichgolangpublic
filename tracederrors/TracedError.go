package tracederrors

import (
	"errors"
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/asciich/asciichgolangpublic/datatypes"
)

var ErrTracedError = errors.New("asciichgolangpublic TracedError base")
var ErrTracedErrorEmptyString = errors.New("asciichgolangpublic TracedError empty string")
var ErrTracedErrorNil = errors.New("asciichgolangpublic TracedError nil")
var ErrTracedErrorNotImplemented = errors.New("asciichgolangpublic TracedError not implemented")

type TracedErrorType struct {
	formattedError error
	functionCalls  []string
	errorsToUnwrap []error
}

// Create a new error with given error or error message.
// TracedErrors extends the error message by a human readable stack trace.
// Error wrapping using '%w' in format string is supported.
func TracedErrorf(formatString string, args ...interface{}) (tracedError error) {
	return TracedError(
		fmt.Errorf(formatString, args...),
	)
}

// Create a new error with given error or error message.
// TracedErrors extends the error message by a human readable stack trace.
func TracedError(errorMessageOrError interface{}) (tracedError error) {
	var internalError error = nil

	switch errorToAdd := errorMessageOrError.(type) {
	case error:
		internalError = errorToAdd
	default:
		internalError = fmt.Errorf("%v", errorToAdd)
	}

	tracedErrorToReturn := TracedErrorType{
		formattedError: internalError,
	}

	programCounters := make([]uintptr, 128)

	const skipNewTracedError = 1
	const skipRuntimeCallers = 1
	const skipFrames = skipNewTracedError + skipRuntimeCallers

	n_callers := runtime.Callers(skipFrames, programCounters)
	frames := runtime.CallersFrames(programCounters[:n_callers])
	more := true
	for more {
		var stackFrameInfo runtime.Frame
		stackFrameInfo, more = frames.Next()

		if stackFrameInfo.Function == "asciichgolangpublic.TracedErrorf" {
			continue
		}

		if stackFrameInfo.Function == "asciichgolangpublic.TracedError" {
			continue
		}

		if stackFrameInfo.Function == "testing.tRunner" {
			break
		}

		if stackFrameInfo.Function == "runtime.main" {
			break
		}

		stackInfoString := fmt.Sprintf("%s (%s:%d)", stackFrameInfo.Function, stackFrameInfo.File, stackFrameInfo.Line)

		tracedErrorToReturn.functionCalls = append([]string{stackInfoString}, tracedErrorToReturn.functionCalls...)
	}

	return tracedErrorToReturn
}

func NewTracedErrorType() (t *TracedErrorType) {
	return new(TracedErrorType)
}

func TracedErrorEmptyString(stringVarName string, errorToUnwrap ...error) (tracedError error) {
	var toReturn TracedErrorType = TracedErrorf("'%s' is empty string", stringVarName).(TracedErrorType)

	toReturn.errorsToUnwrap = append(toReturn.errorsToUnwrap, ErrTracedErrorEmptyString)
	toReturn.errorsToUnwrap = append(toReturn.errorsToUnwrap, errorToUnwrap...)

	return toReturn
}

func TracedErrorNil(nilVarName string) (tracedError error) {
	var toReturn TracedErrorType = TracedErrorf("'%s' is nil", nilVarName).(TracedErrorType)
	toReturn.errorsToUnwrap = append(toReturn.errorsToUnwrap, ErrTracedErrorNil)
	return toReturn
}

func TracedErrorNilf(formatString string, args ...interface{}) (tracedError error) {
	message := fmt.Sprintf(formatString, args...)
	return TracedErrorNil(message)
}

func TracedErrorNotImplemented() (tracedError error) {
	var toReturn TracedErrorType = TracedError("Not implemented").(TracedErrorType)
	toReturn.errorsToUnwrap = append(toReturn.errorsToUnwrap, ErrTracedErrorNotImplemented)
	return toReturn
}

func (t *TracedErrorType) GetErrorMessage() (errorMessage string, err error) {
	formattedError, err := t.GetFormattedError()
	if err != nil {
		return "", err
	}

	errorMessage = formattedError.Error()

	return errorMessage, nil
}

func (t *TracedErrorType) GetErrorsToUnwrap() (errorsToUnwrap []error, err error) {
	if t.errorsToUnwrap == nil {
		return nil, TracedErrorf("errorsToUnwrap not set")
	}

	if len(t.errorsToUnwrap) <= 0 {
		return nil, TracedErrorf("errorsToUnwrap has no elements")
	}

	return t.errorsToUnwrap, nil
}

func (t *TracedErrorType) GetFormattedError() (formattedError error, err error) {

	return t.formattedError, nil
}

func (t *TracedErrorType) GetFunctionCalls() (functionCalls []string, err error) {
	if t.functionCalls == nil {
		return nil, TracedErrorf("functionCalls not set")
	}

	if len(t.functionCalls) <= 0 {
		return nil, TracedErrorf("functionCalls has no elements")
	}

	return t.functionCalls, nil
}

func (t *TracedErrorType) MustGetErrorMessage() (errorMessage string) {
	errorMessage, err := t.GetErrorMessage()
	if err != nil {
		log.Panic(err)
	}

	return errorMessage
}

func (t *TracedErrorType) MustGetErrorsToUnwrap() (errorsToUnwrap []error) {
	errorsToUnwrap, err := t.GetErrorsToUnwrap()
	if err != nil {
		log.Panic(err)
	}

	return errorsToUnwrap
}

func (t *TracedErrorType) MustGetFormattedError() (formattedError error) {
	formattedError, err := t.GetFormattedError()
	if err != nil {
		log.Panic(err)
	}

	return formattedError
}

func (t *TracedErrorType) MustGetFunctionCalls() (functionCalls []string) {
	functionCalls, err := t.GetFunctionCalls()
	if err != nil {
		log.Panic(err)
	}

	return functionCalls
}

func (t *TracedErrorType) MustSetErrorsToUnwrap(errorsToUnwrap []error) {
	err := t.SetErrorsToUnwrap(errorsToUnwrap)
	if err != nil {
		log.Panic(err)
	}
}

func (t *TracedErrorType) MustSetFormattedError(formattedError error) {
	err := t.SetFormattedError(formattedError)
	if err != nil {
		log.Panic(err)
	}
}

func (t *TracedErrorType) MustSetFunctionCalls(functionCalls []string) {
	err := t.SetFunctionCalls(functionCalls)
	if err != nil {
		log.Panic(err)
	}
}

func (t *TracedErrorType) SetErrorsToUnwrap(errorsToUnwrap []error) (err error) {
	if errorsToUnwrap == nil {
		return TracedErrorf("errorsToUnwrap is nil")
	}

	if len(errorsToUnwrap) <= 0 {
		return TracedErrorf("errorsToUnwrap has no elements")
	}

	t.errorsToUnwrap = errorsToUnwrap

	return nil
}

func (t *TracedErrorType) SetFormattedError(formattedError error) (err error) {
	t.formattedError = formattedError

	return nil
}

func (t *TracedErrorType) SetFunctionCalls(functionCalls []string) (err error) {
	if functionCalls == nil {
		return TracedErrorf("functionCalls is nil")
	}

	if len(functionCalls) <= 0 {
		return TracedErrorf("functionCalls has no elements")
	}

	t.functionCalls = functionCalls

	return nil
}

func (t TracedErrorType) Error() (errorMessage string) {
	errorMessage = ""

	if len(t.functionCalls) > 0 {
		errorMessage += strings.Join(t.functionCalls, "\n")
		errorMessage += "\n"
		errorMessage += "\n"
	}

	allErrors := UnwrapRecursive(t)
	for _, unwrapped := range allErrors {
		unwrapType, err := datatypes.GetTypeName(unwrapped)
		if err != nil {
			continue
		}
		errorMessage += fmt.Sprintf(
			"This error unwraps to type '%s'\n",
			unwrapType,
		)
	}
	errorMessage += "\n"

	errorMessage += t.formattedError.Error()

	return errorMessage
}

func (t TracedErrorType) Unwrap() (errors []error) {
	errors = []error{
		ErrTracedError, // Every traced error is automatically a ErrTraccedError.
	}

	errors = append(errors, t.formattedError)

	if t.errorsToUnwrap != nil {
		errors = append(errors, t.errorsToUnwrap...)
	}

	switch x := t.formattedError.(type) {
	case interface{ Unwrap() error }:
		errors = append(errors, x.Unwrap())
	case interface{ Unwrap() []error }:
		errors = append(errors, x.Unwrap()...)
	}

	return errors
}
