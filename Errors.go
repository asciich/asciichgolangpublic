package asciichgolangpublic

import "errors"

type ErrorsService struct{}

func Errors() (e *ErrorsService) {
	return new(ErrorsService)
}

func NewErrorsService() (e *ErrorsService) {
	return new(ErrorsService)
}

// Returns true if given error 'err' is a TracedError, false otherwise.
func (e ErrorsService) IsTracedError(err error) (isTracedError bool) {
	return errors.Is(err, ErrTracedError)
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
