package htmlutils

import "errors"

var ErrNoHtmlBodyFound = errors.New("no html body found")

func IsErrNoHtmlBodyFound(err error) bool {
	if err == nil {
		return false
	}

	return errors.Is(err, ErrNoHtmlBodyFound)
}
