package asciichgolangpublic

import (
	astrings "github.com/asciich/asciichgolangpublic/datatypes/strings"
	"github.com/asciich/asciichgolangpublic/errors"
	"github.com/asciich/asciichgolangpublic/logging"
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
		return false, errors.TracedErrorf("'%s' is not an URL.", url)
	}

	return isUrl, nil
}

func (u *UrlsService) IsUrl(url string) (isUrl bool) {
	if url == "" {
		return false
	}

	return astrings.HasAtLeastOnePrefix(url, []string{
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
