package errorutils

import (
	"fmt"

	"github.com/asciich/asciichgolangpublic/tracederrors"
)

func GetErrorMessage(e error) string {
	if e == nil {
		return ""
	}

	tracedError, err := tracederrors.GetAsTracedError(e)
	if err != nil {
		return e.Error()
	}

	msg, err := tracedError.GetErrorMessage()
	if err != nil {
		return e.Error()
	}

	return msg
}

func AppendToErrorMessage(e error, appendix string) error {
	if e == nil {
		return fmt.Errorf("%s", appendix)
	}

	tracedError, err := tracederrors.GetAsTracedError(e)
	if err != nil {
		return fmt.Errorf("%w %s", e, appendix)
	}

	msg, err := tracedError.GetErrorMessage()
	if err != nil {
		return fmt.Errorf("%w %s", e, appendix)
	}

	err = tracedError.SetFormattedError(fmt.Errorf("%s %s", msg, appendix))
	if err != nil {
		return fmt.Errorf("%w %s", e, appendix)
	}

	return tracedError
}
