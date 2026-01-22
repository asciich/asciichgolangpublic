package htmlutils

import "errors"

var ErrChildNodeNotFound = errors.New("child node not found")
var ErrNoHtmlBodyFound = errors.New("no html body found")
var ErrAttrNotFound = errors.New("attr not found")

func IsErrNoHtmlBodyFound(err error) bool {
	if err == nil {
		return false
	}

	return errors.Is(err, ErrNoHtmlBodyFound)
}

func IsErrAttrNotFound(err error) bool {
	if err == nil {
		return false
	}

	return errors.Is(err, ErrAttrNotFound)
}

func IsErrChildNodeNotFound(err error) bool {
	if err == nil {
		return false
	}

	return errors.Is(err, ErrChildNodeNotFound)
}
