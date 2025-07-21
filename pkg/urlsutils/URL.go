package urlsutils

import (
	"strings"

	"github.com/asciich/asciichgolangpublic/datatypes/stringsutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

// Parts of an URL
// ================
//
// Source: https://www.geeksforgeeks.org/components-of-a-url/
//
// Sheme:                  https://
// Subdomain:              www.
// Domain:                 asciich.
// TopLevel domain:        ch/
// Path:                   path/to/file
// Query string separator: ?
// Query string parameter: x=5&y=10
// Fragment                #today

type URL struct {
	urlString string
}

func GetUrlFromString(urlString string) (url *URL, err error) {
	url = NewURL()
	err = url.SetByString(urlString)
	if err != nil {
		return nil, err
	}

	return url, nil
}

func MustGetUrlFromString(urlString string) (url *URL) {
	url, err := GetUrlFromString(urlString)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return url
}

func NewURL() (u *URL) {
	return new(URL)
}

func (u *URL) GetAsString() (urlString string, err error) {
	if len(u.urlString) <= 0 {
		return "", tracederrors.TracedError("urlString not set")
	}

	return u.urlString, nil
}

func (u *URL) GetFqdnAsString() (fqdn string, err error) {
	urlString, err := u.GetAsString()
	if err != nil {
		return "", err
	}

	splitted := strings.SplitN(urlString, "://", 2)
	fqdnAndPath := splitted[len(splitted)-1]

	fqdn = strings.Split(fqdnAndPath, "/")[0]

	if fqdn == "" {
		return "", tracederrors.TracedErrorf(
			"fqdn is empty string after evaluation of urlString='%s'",
			urlString,
		)
	}

	return fqdn, nil
}

func (u *URL) GetFqdnWitShemeAndPathAsString() (fqdnWithSheme string, path string, err error) {
	fqdnWithSheme, err = u.GetFqdnWithShemeAsString()
	if err != nil {
		return "", "", err
	}

	path, err = u.GetPathAsString()
	if err != nil {
		return "", "", err
	}

	return fqdnWithSheme, path, err
}

func (u *URL) GetFqdnWithShemeAsString() (fqdnWithSheme string, err error) {
	sheme, err := u.GetSheme()
	if err != nil {
		return "", err
	}

	fqdn, err := u.GetFqdnAsString()
	if err != nil {
		return "", err
	}

	fqdnWithSheme = sheme + "://" + fqdn

	return fqdnWithSheme, nil
}

func (u *URL) GetPathAsString() (urlPath string, err error) {
	withoutSheme, err := u.GetWithoutShemeAsString()
	if err != nil {
		return "", err
	}
	splitted := strings.Split(withoutSheme, "/")
	if len(splitted) <= 0 {
		return "", tracederrors.TracedError("failed to split 'withoutSheme'")
	}

	pathParts := splitted[1:]
	urlPath = strings.Join(pathParts, "/")

	urlPath = stringsutils.TrimAllSuffix(urlPath, "/")

	return urlPath, nil
}

func (u *URL) GetPathBasename() (basename string, err error) {
	path, err := u.GetPathAsString()
	if err != nil {
		return "", err
	}

	splitted := strings.Split(path, "/")
	if len(splitted) <= 0 {
		return "", tracederrors.TracedErrorf("failed to split '%v'", path)
	}

	basename = splitted[len(splitted)-1]
	return basename, err
}

func (u *URL) GetSheme() (sheme string, err error) {
	urlString, err := u.GetAsString()
	if err != nil {
		return "", err
	}

	splitted := strings.SplitN(urlString, "://", 2)
	if len(splitted) != 2 {
		return "", tracederrors.TracedErrorf(
			"Unable to get sheme from urlString '%s'",
			urlString,
		)
	}

	sheme = splitted[0]
	if sheme == "" {
		return "", tracederrors.TracedError("sheme is empty string after evaluation")
	}

	return sheme, nil
}

func (u *URL) GetUrlString() (urlString string, err error) {
	if u.urlString == "" {
		return "", tracederrors.TracedErrorf("urlString not set")
	}

	return u.urlString, nil
}

func (u *URL) GetWithoutShemeAsString() (urlWithoutSheme string, err error) {
	urlString, err := u.GetAsString()
	if err != nil {
		return "", err
	}

	splitted := strings.Split(urlString, "://")
	urlWithoutSheme = splitted[len(splitted)-1]

	return urlWithoutSheme, nil
}

func (u *URL) MustGetAsString() (urlString string) {
	urlString, err := u.GetAsString()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return urlString
}

func (u *URL) MustGetFqdnAsString() (fqdn string) {
	fqdn, err := u.GetFqdnAsString()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return fqdn
}

func (u *URL) MustGetFqdnWitShemeAndPathAsString() (fqdnWithSheme string, path string) {
	fqdnWithSheme, path, err := u.GetFqdnWitShemeAndPathAsString()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return fqdnWithSheme, path
}

func (u *URL) MustGetFqdnWithShemeAsString() (fqdnWithSheme string) {
	fqdnWithSheme, err := u.GetFqdnWithShemeAsString()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return fqdnWithSheme
}

func (u *URL) MustGetPathAsString() (urlPath string) {
	urlPath, err := u.GetPathAsString()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return urlPath
}

func (u *URL) MustGetPathBasename() (basename string) {
	basename, err := u.GetPathBasename()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return basename
}

func (u *URL) MustGetSheme() (sheme string) {
	sheme, err := u.GetSheme()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return sheme
}

func (u *URL) MustGetUrlString() (urlString string) {
	urlString, err := u.GetUrlString()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return urlString
}

func (u *URL) MustGetWithoutShemeAsString() (urlWithoutSheme string) {
	urlWithoutSheme, err := u.GetWithoutShemeAsString()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return urlWithoutSheme
}

func (u *URL) MustSetByString(urlString string) {
	err := u.SetByString(urlString)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (u *URL) MustSetUrlString(urlString string) {
	err := u.SetUrlString(urlString)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (u *URL) SetByString(urlString string) (err error) {
	urlString = strings.TrimSpace(urlString)
	if len(urlString) <= 0 {
		return tracederrors.TracedError("urlString is empty string")
	}

	u.urlString = urlString

	return nil
}

func (u *URL) SetUrlString(urlString string) (err error) {
	if urlString == "" {
		return tracederrors.TracedErrorf("urlString is empty string")
	}

	u.urlString = urlString

	return nil
}
