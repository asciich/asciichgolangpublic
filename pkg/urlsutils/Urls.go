package urlsutils

import (
	"github.com/asciich/asciichgolangpublic/pkg/datatypes/stringsutils"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
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
