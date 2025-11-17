package urlsutils

import (
	"fmt"
	"net/url"

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

func GetBaseUrl(inputUrl string) (string, error) {
	if inputUrl == "" {
		return "", tracederrors.TracedErrorEmptyString("inputUrl")
	}

	parsedURL, err := url.Parse(inputUrl)
	if err != nil {
		return "", tracederrors.TracedErrorf("Failed to parse URL '%s': %w", inputUrl, err)
	}

	return fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host), nil
}
