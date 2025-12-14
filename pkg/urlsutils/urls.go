package urlsutils

import (
	"fmt"
	"net/url"
	"strconv"

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

// GetPath parses the given raw URL string and returns the path component.
// If the URL cannot be parsed, an empty string is returned.
func GetPath(inputUrl string) (string, error) {
	if inputUrl == "" {
		return "", tracederrors.TracedErrorEmptyString("inputUrl")
	}

	u, err := url.Parse(inputUrl)
	if err != nil {
		return "", tracederrors.TracedErrorf("Failed to parse as url: %w", err)
	}

	return u.Path, nil
}

func SetPort(inputUrl string, port int) (string, error) {
	if inputUrl == "" {
		return "", tracederrors.TracedErrorEmptyString("inputUrl")
	}

	if port <= 0 {
		return "", tracederrors.TracedErrorf("Invalid port number: %d", port)
	}

	u, err := url.Parse(inputUrl)
	if err != nil {
		return "", tracederrors.TracedErrorf("Failed to parse as url: %w", err)
	}

	hostname := u.Hostname()
	newHost := hostname + ":" + strconv.Itoa(port)
	u.Host = newHost
	newUrl := u.String()
	return newUrl, nil
}