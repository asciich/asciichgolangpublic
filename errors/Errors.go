package errors

import (
	"errors"
	nativeErrors "errors"
	"log"
)

type ErrorsService struct{}

func Errors() (e *ErrorsService) {
	return new(ErrorsService)
}

func NewErrorsService() (e *ErrorsService) {
	return new(ErrorsService)
}

// Returns true if given error 'err' is a TracedError, false otherwise.
func (e ErrorsService) IsTracedError(err error) (isTracedError bool) {
	return nativeErrors.Is(err, ErrTracedError)
}

func (e *ErrorsService) AddErrorToUnwrapToTracedError(tracedError error, errorToAdd error) error {
	if tracedError == nil {
		return nil
	}

	if errorToAdd == nil {
		return nil
	}

	resultingError, ok := tracedError.(*TracedErrorType)
	if !ok {
		return tracedError
	}

	resultingError.errorsToUnwrap = append(resultingError.errorsToUnwrap, errorToAdd)

	return resultingError
}

func (e *ErrorsService) GetAsTracedError(errorToConvert error) (tracedError *TracedErrorType, err error) {
	if errorToConvert == nil {
		return nil, TracedErrorNil("errorToConvert")
	}

	tracedError, ok := errorToConvert.(*TracedErrorType)
	if !ok {
		tracedErrorNonPointer, ok := errorToConvert.(TracedErrorType)
		if !ok {
			return nil, TracedErrorf("Unable to convert '%v' to TracedError", errorToConvert)
		}

		tracedError = &tracedErrorNonPointer
	}

	return tracedError, nil
}

func (e *ErrorsService) MustGetAsTracedError(errorToConvert error) (tracedError *TracedErrorType) {
	tracedError, err := e.GetAsTracedError(errorToConvert)
	if err != nil {
		log.Panic(err)
	}

	return tracedError
}

func (e ErrorsService) IsEmptyStringError(err error) (isEmptyStringError bool) {
	return errors.Is(err, ErrTracedErrorEmptyString)
}

func (e ErrorsService) IsNilError(err error) (IsNilError bool) {
	return errors.Is(err, ErrTracedErrorNil)
}

func (e ErrorsService) IsNotImplementedError(err error) (isNotImplementedError bool) {
	return errors.Is(err, ErrTracedErrorNotImplemented)
}

func (e ErrorsService) UnwrapRecursive(errorToUnwrap error) (errors []error) {
	errors = []error{}

	if errorToUnwrap == nil {
		return errors
	}

	switch x := errorToUnwrap.(type) {
	case interface{ Unwrap() error }:
		toAdd := x.Unwrap()
		if toAdd == nil {
			errors = append(errors, toAdd)
			errors = append(errors, e.UnwrapRecursive(toAdd)...)
		}

	case interface{ Unwrap() []error }:
		for _, toAdd := range x.Unwrap() {
			errors = append(errors, toAdd)
			errors = append(errors, e.UnwrapRecursive(toAdd)...)
		}
	}

	return errors
}
