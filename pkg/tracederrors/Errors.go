package tracederrors

import (
	"errors"
)

// Returns true if given error 'err' is a TracedError, false otherwise.
func IsTracedError(err error) (isTracedError bool) {
	return errors.Is(err, ErrTracedError)
}

func AddErrorToUnwrapToTracedError(tracedError error, errorToAdd error) error {
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

func GetAsTracedError(errorToConvert error) (tracedError *TracedErrorType, err error) {
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

func IsEmptyStringError(err error) (isEmptyStringError bool) {
	return errors.Is(err, ErrTracedErrorEmptyString)
}

func IsNilError(err error) (IsNilError bool) {
	return errors.Is(err, ErrTracedErrorNil)
}

func IsNotImplementedError(err error) (isNotImplementedError bool) {
	return errors.Is(err, ErrTracedErrorNotImplemented)
}

func UnwrapRecursive(errorToUnwrap error) (errors []error) {
	errors = []error{}

	if errorToUnwrap == nil {
		return errors
	}

	switch x := errorToUnwrap.(type) {
	case interface{ Unwrap() error }:
		toAdd := x.Unwrap()
		if toAdd == nil {
			errors = append(errors, toAdd)
			errors = append(errors, UnwrapRecursive(toAdd)...)
		}

	case interface{ Unwrap() []error }:
		for _, toAdd := range x.Unwrap() {
			errors = append(errors, toAdd)
			errors = append(errors, UnwrapRecursive(toAdd)...)
		}
	}

	return errors
}
