package urlsutils

import (
	"github.com/asciich/asciichgolangpublic/datatypes/stringsutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

func CheckIsUrl(url string) (isUrl bool, err error) {
	isUrl = IsUrl(url)
	if !isUrl {
		return false, tracederrors.TracedErrorf("'%s' is not an URL.", url)
	}

	return isUrl, nil
}

func IsUrl(url string) (isUrl bool) {
	if url == "" {
		return false
	}

	return stringsutils.HasAtLeastOnePrefix(url, []string{
		"https://",
		"http://",
	})
}

func MustCheckIsUrl(url string) (isUrl bool) {
	isUrl, err := CheckIsUrl(url)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return isUrl
}
