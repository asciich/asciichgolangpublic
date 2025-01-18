package urlsutils

import (
	"github.com/asciich/asciichgolangpublic/datatypes/stringsutils"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type UrlsService struct{}

func NewUrlsService() (service *UrlsService) {
	return new(UrlsService)
}

func Urls() (urlService *UrlsService) {
	return NewUrlsService()
}

func (u *UrlsService) CheckIsUrl(url string) (isUrl bool, err error) {
	isUrl = u.IsUrl(url)
	if !isUrl {
		return false, tracederrors.TracedErrorf("'%s' is not an URL.", url)
	}

	return isUrl, nil
}

func (u *UrlsService) IsUrl(url string) (isUrl bool) {
	if url == "" {
		return false
	}

	return stringsutils.HasAtLeastOnePrefix(url, []string{
		"https://",
		"http://",
	})
}

func (u *UrlsService) MustCheckIsUrl(url string) (isUrl bool) {
	isUrl, err := u.CheckIsUrl(url)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return isUrl
}
